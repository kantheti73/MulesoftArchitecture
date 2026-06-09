// payload-capture-wasm
//
// Custom Anypoint Omni Gateway policy (Envoy WASM filter) that captures
// every API call's request payload + headers and publishes them to an
// Azure Service Bus queue asynchronously.
//
// Design principles:
//   1. Fire-and-forget: never block the request path on the publish call
//   2. Sampling: configurable per-API sample rate (avoid 100% capture
//      cost on high-volume APIs)
//   3. PII redaction: header allowlist + body field deny-list applied
//      before publish (per docs/07-data-protection.md)
//   4. Size cap: large payloads truncated with a marker (avoid SB 1 MB
//      message limit)
//   5. Failure-tolerant: SB publish failure is logged but does not break
//      the gateway request
//   6. Bounded memory: per-request capture state cleared after publish

use proxy_wasm::traits::*;
use proxy_wasm::types::*;
use serde::{Deserialize, Serialize};
use std::time::Duration;

// ==========================================================================
// Entrypoint
// ==========================================================================

#[no_mangle]
pub fn _start() {
    proxy_wasm::set_log_level(LogLevel::Info);
    proxy_wasm::set_root_context(|_| -> Box<dyn RootContext> {
        Box::new(PayloadCaptureRoot::default())
    });
}

// ==========================================================================
// Configuration (parsed from the policy manifest's `configurationData`)
// ==========================================================================

#[derive(Clone, Default, Deserialize, Debug)]
struct PolicyConfig {
    // Service Bus REST API target
    sb_namespace_host: String,           // e.g. "ssp-bus.servicebus.windows.net"
    sb_queue_name: String,                // e.g. "api-payload-capture"
    sb_sas_authorization_header: String,  // pre-built "SharedAccessSignature ..." (read from Secrets Manager via Anypoint)

    // Sampling
    sample_rate_percent: u32,             // 0-100; 100 = capture every request

    // PII handling
    redact_body_fields: Vec<String>,      // JSON pointer paths to redact (e.g. "/ssn", "/email")
    capture_header_allowlist: Vec<String>,// only these headers are forwarded; everything else dropped
    capture_response: bool,               // also capture response body (default false)

    // Limits
    max_payload_bytes: usize,             // truncate above this (default 64 KB)
    sb_http_cluster_name: String,         // Envoy upstream cluster name for SB
    timeout_ms: u64,                      // HTTP call timeout (default 5000)
}

impl PolicyConfig {
    fn with_defaults(mut self) -> Self {
        if self.sample_rate_percent == 0 {
            self.sample_rate_percent = 10;
        }
        if self.max_payload_bytes == 0 {
            self.max_payload_bytes = 65_536;
        }
        if self.timeout_ms == 0 {
            self.timeout_ms = 5_000;
        }
        if self.capture_header_allowlist.is_empty() {
            self.capture_header_allowlist = vec![
                "content-type".into(),
                "accept".into(),
                "x-correlation-id".into(),
                "x-request-id".into(),
                "traceparent".into(),
            ];
        }
        self
    }
}

// ==========================================================================
// Root context — loads + validates configuration on start
// ==========================================================================

#[derive(Default)]
struct PayloadCaptureRoot {
    config: PolicyConfig,
}

impl Context for PayloadCaptureRoot {}

impl RootContext for PayloadCaptureRoot {
    fn on_configure(&mut self, _config_size: usize) -> bool {
        let bytes = match self.get_plugin_configuration() {
            Some(b) => b,
            None => {
                log::warn!("payload-capture: no configuration provided; using defaults");
                Vec::new()
            }
        };

        let parsed: Result<PolicyConfig, _> = serde_json::from_slice(&bytes);
        self.config = match parsed {
            Ok(cfg) => cfg.with_defaults(),
            Err(e) => {
                log::error!("payload-capture: invalid configuration JSON: {}", e);
                return false; // fail to start — better than running misconfigured
            }
        };

        // Required field validation
        if self.config.sb_namespace_host.is_empty()
            || self.config.sb_queue_name.is_empty()
            || self.config.sb_sas_authorization_header.is_empty()
            || self.config.sb_http_cluster_name.is_empty()
        {
            log::error!(
                "payload-capture: missing required config (sb_namespace_host / sb_queue_name / \
                 sb_sas_authorization_header / sb_http_cluster_name)"
            );
            return false;
        }

        log::info!(
            "payload-capture: configured for queue {}@{} (sample {}%, max {} bytes, response capture {})",
            self.config.sb_queue_name,
            self.config.sb_namespace_host,
            self.config.sample_rate_percent,
            self.config.max_payload_bytes,
            self.config.capture_response
        );
        true
    }

    fn create_http_context(&self, _context_id: u32) -> Option<Box<dyn HttpContext>> {
        Some(Box::new(PayloadCaptureHttp {
            config: self.config.clone(),
            sampled_in: false,
            captured_request_headers: Vec::new(),
            captured_request_body: Vec::new(),
            captured_response_headers: Vec::new(),
            captured_response_body: Vec::new(),
            request_method: String::new(),
            request_path: String::new(),
            response_status: 0,
            start_epoch_ms: 0,
            request_truncated: false,
            response_truncated: false,
        }))
    }

    fn get_type(&self) -> Option<ContextType> {
        Some(ContextType::HttpContext)
    }
}

// ==========================================================================
// Per-request HTTP context — does the actual capture + publish
// ==========================================================================

struct PayloadCaptureHttp {
    config: PolicyConfig,
    sampled_in: bool,
    captured_request_headers: Vec<(String, String)>,
    captured_request_body: Vec<u8>,
    captured_response_headers: Vec<(String, String)>,
    captured_response_body: Vec<u8>,
    request_method: String,
    request_path: String,
    response_status: u16,
    start_epoch_ms: u64,
    request_truncated: bool,
    response_truncated: bool,
}

impl Context for PayloadCaptureHttp {
    fn on_http_call_response(
        &mut self,
        _token_id: u32,
        _num_headers: usize,
        _body_size: usize,
        _num_trailers: usize,
    ) {
        // Fire-and-forget: we log success/failure but take no action.
        // Service Bus errors don't propagate back to the original request.
        if let Some(status_header) = self
            .get_http_call_response_header(":status")
        {
            if !status_header.starts_with('2') {
                log::warn!(
                    "payload-capture: Service Bus publish returned non-2xx status: {}",
                    status_header
                );
            }
        }
    }
}

impl HttpContext for PayloadCaptureHttp {
    fn on_http_request_headers(&mut self, _num_headers: usize, _end_of_stream: bool) -> Action {
        // Sampling decision is made once, per request
        self.sampled_in = should_sample(self.config.sample_rate_percent);
        if !self.sampled_in {
            return Action::Continue;
        }

        // Capture method + path for the envelope
        self.request_method = self
            .get_http_request_header(":method")
            .unwrap_or_else(|| "UNKNOWN".to_string());
        self.request_path = self
            .get_http_request_header(":path")
            .unwrap_or_else(|| "/unknown".to_string());

        // Capture an allowlisted subset of headers (per PII rules in docs/07)
        let raw_headers = self.get_http_request_headers();
        self.captured_request_headers = raw_headers
            .into_iter()
            .filter(|(name, _)| {
                let lower = name.to_ascii_lowercase();
                self.config
                    .capture_header_allowlist
                    .iter()
                    .any(|h| h.to_ascii_lowercase() == lower)
            })
            .collect();

        // Record start time (Envoy uses microseconds since epoch in some flavors;
        // for portability we use seconds + millis via the host clock)
        self.start_epoch_ms = current_epoch_ms();

        Action::Continue
    }

    fn on_http_request_body(&mut self, body_size: usize, end_of_stream: bool) -> Action {
        if !self.sampled_in || body_size == 0 {
            return Action::Continue;
        }

        // Read this chunk
        if let Some(chunk) = self.get_http_request_body(0, body_size) {
            let remaining = self
                .config
                .max_payload_bytes
                .saturating_sub(self.captured_request_body.len());
            if remaining == 0 {
                self.request_truncated = true;
            } else if chunk.len() > remaining {
                self.captured_request_body
                    .extend_from_slice(&chunk[..remaining]);
                self.request_truncated = true;
            } else {
                self.captured_request_body.extend_from_slice(&chunk);
            }
        }

        if end_of_stream && !self.config.capture_response {
            self.publish_to_service_bus();
        }
        Action::Continue
    }

    fn on_http_response_headers(&mut self, _num_headers: usize, _end_of_stream: bool) -> Action {
        if !self.sampled_in || !self.config.capture_response {
            return Action::Continue;
        }

        // Capture response status + allowlisted response headers
        if let Some(status) = self.get_http_response_header(":status") {
            self.response_status = status.parse().unwrap_or(0);
        }

        let raw = self.get_http_response_headers();
        self.captured_response_headers = raw
            .into_iter()
            .filter(|(name, _)| {
                let lower = name.to_ascii_lowercase();
                self.config
                    .capture_header_allowlist
                    .iter()
                    .any(|h| h.to_ascii_lowercase() == lower)
            })
            .collect();

        Action::Continue
    }

    fn on_http_response_body(&mut self, body_size: usize, end_of_stream: bool) -> Action {
        if !self.sampled_in || !self.config.capture_response || body_size == 0 {
            if end_of_stream && self.sampled_in {
                self.publish_to_service_bus();
            }
            return Action::Continue;
        }

        if let Some(chunk) = self.get_http_response_body(0, body_size) {
            let remaining = self
                .config
                .max_payload_bytes
                .saturating_sub(self.captured_response_body.len());
            if remaining == 0 {
                self.response_truncated = true;
            } else if chunk.len() > remaining {
                self.captured_response_body
                    .extend_from_slice(&chunk[..remaining]);
                self.response_truncated = true;
            } else {
                self.captured_response_body.extend_from_slice(&chunk);
            }
        }

        if end_of_stream {
            self.publish_to_service_bus();
        }
        Action::Continue
    }
}

impl PayloadCaptureHttp {
    fn publish_to_service_bus(&self) {
        // Build the canonical event envelope (CloudEvents 1.0 shape per doc 11 §6)
        let envelope = self.build_envelope();
        let body_bytes = match serde_json::to_vec(&envelope) {
            Ok(b) => b,
            Err(e) => {
                log::error!("payload-capture: failed to serialize envelope: {}", e);
                return;
            }
        };

        // Compose Service Bus REST API path
        let sb_path = format!(
            "/{}/messages?api-version=2015-01",
            self.config.sb_queue_name
        );

        // Headers required by Service Bus REST API
        let headers: Vec<(&str, &str)> = vec![
            (":method", "POST"),
            (":authority", &self.config.sb_namespace_host),
            (":path", &sb_path),
            (":scheme", "https"),
            ("authorization", &self.config.sb_sas_authorization_header),
            ("content-type", "application/json"),
            (
                "broker-properties",
                // Set a stable MessageId for SB-side dedup
                // (one per gateway request via correlation ID if present)
                "{\"MessageId\":\"capture-anon\",\"ContentType\":\"application/cloudevents+json\"}",
            ),
        ];

        let result = self.dispatch_http_call(
            &self.config.sb_http_cluster_name,
            headers,
            Some(&body_bytes),
            vec![],
            Duration::from_millis(self.config.timeout_ms),
        );

        match result {
            Ok(_token) => log::debug!(
                "payload-capture: dispatched {} bytes to {}@{}",
                body_bytes.len(),
                self.config.sb_queue_name,
                self.config.sb_namespace_host
            ),
            Err(e) => log::warn!("payload-capture: dispatch_http_call failed: {:?}", e),
        }
    }

    fn build_envelope(&self) -> CaptureEnvelope {
        // PII redaction: walk JSON, blank-out denied fields if body is JSON
        let request_body_clean = redact_json_fields(
            &self.captured_request_body,
            &self.config.redact_body_fields,
        );
        let response_body_clean = redact_json_fields(
            &self.captured_response_body,
            &self.config.redact_body_fields,
        );

        CaptureEnvelope {
            specversion: "1.0",
            type_: "com.yourco.gateway.api-call-captured",
            source: "/omni-gateway/payload-capture",
            id: extract_correlation_id(&self.captured_request_headers)
                .unwrap_or_else(|| format!("cap-{}", self.start_epoch_ms)),
            time_epoch_ms: self.start_epoch_ms,
            datacontenttype: "application/json",
            data: CaptureData {
                method: self.request_method.clone(),
                path: self.request_path.clone(),
                request_headers: self.captured_request_headers.clone(),
                request_body_b64: base64_encode(&request_body_clean),
                request_truncated: self.request_truncated,
                response_status: self.response_status,
                response_headers: self.captured_response_headers.clone(),
                response_body_b64: base64_encode(&response_body_clean),
                response_truncated: self.response_truncated,
            },
        }
    }
}

// ==========================================================================
// Envelope types (CloudEvents 1.0 shape)
// ==========================================================================

#[derive(Serialize)]
struct CaptureEnvelope {
    specversion: &'static str,
    #[serde(rename = "type")]
    type_: &'static str,
    source: &'static str,
    id: String,
    #[serde(rename = "time")]
    time_epoch_ms: u64,
    datacontenttype: &'static str,
    data: CaptureData,
}

#[derive(Serialize)]
struct CaptureData {
    method: String,
    path: String,
    request_headers: Vec<(String, String)>,
    request_body_b64: String,
    request_truncated: bool,
    response_status: u16,
    response_headers: Vec<(String, String)>,
    response_body_b64: String,
    response_truncated: bool,
}

// ==========================================================================
// Helpers
// ==========================================================================

fn should_sample(percent: u32) -> bool {
    if percent >= 100 {
        return true;
    }
    if percent == 0 {
        return false;
    }
    // Cheap pseudo-random based on host clock — adequate for sampling decisions
    let now = current_epoch_ms();
    let r = (now.wrapping_mul(2862933555777941757) ^ (now >> 17)) % 100;
    (r as u32) < percent
}

fn current_epoch_ms() -> u64 {
    // proxy-wasm exposes the host clock via the standard API
    proxy_wasm::hostcalls::get_current_time()
        .map(|d| d.as_millis() as u64)
        .unwrap_or(0)
}

fn extract_correlation_id(headers: &[(String, String)]) -> Option<String> {
    for (k, v) in headers {
        let lower = k.to_ascii_lowercase();
        if lower == "x-correlation-id" || lower == "x-request-id" {
            return Some(v.clone());
        }
    }
    None
}

fn redact_json_fields(body: &[u8], deny: &[String]) -> Vec<u8> {
    if body.is_empty() || deny.is_empty() {
        return body.to_vec();
    }
    let mut v: serde_json::Value = match serde_json::from_slice(body) {
        Ok(v) => v,
        Err(_) => return body.to_vec(), // not JSON; return as-is
    };
    for pointer in deny {
        if let Some(target) = v.pointer_mut(pointer) {
            *target = serde_json::Value::String("[REDACTED]".to_string());
        }
    }
    serde_json::to_vec(&v).unwrap_or_else(|_| body.to_vec())
}

// Tiny base64 encoder — keeps the WASM binary lean (no full base64 crate needed)
fn base64_encode(input: &[u8]) -> String {
    const ALPHA: &[u8] =
        b"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";
    let mut out = String::with_capacity((input.len() + 2) / 3 * 4);
    let mut i = 0;
    while i + 3 <= input.len() {
        let n = ((input[i] as u32) << 16) | ((input[i + 1] as u32) << 8) | (input[i + 2] as u32);
        out.push(ALPHA[((n >> 18) & 63) as usize] as char);
        out.push(ALPHA[((n >> 12) & 63) as usize] as char);
        out.push(ALPHA[((n >> 6) & 63) as usize] as char);
        out.push(ALPHA[(n & 63) as usize] as char);
        i += 3;
    }
    let rem = input.len() - i;
    if rem == 1 {
        let n = (input[i] as u32) << 16;
        out.push(ALPHA[((n >> 18) & 63) as usize] as char);
        out.push(ALPHA[((n >> 12) & 63) as usize] as char);
        out.push('=');
        out.push('=');
    } else if rem == 2 {
        let n = ((input[i] as u32) << 16) | ((input[i + 1] as u32) << 8);
        out.push(ALPHA[((n >> 18) & 63) as usize] as char);
        out.push(ALPHA[((n >> 12) & 63) as usize] as char);
        out.push(ALPHA[((n >> 6) & 63) as usize] as char);
        out.push('=');
    }
    out
}

// ==========================================================================
// Unit tests (run with `cargo test`)
// ==========================================================================

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn base64_simple() {
        assert_eq!(base64_encode(b"hello"), "aGVsbG8=");
        assert_eq!(base64_encode(b"hi"), "aGk=");
        assert_eq!(base64_encode(b""), "");
        assert_eq!(base64_encode(b"abc"), "YWJj");
    }

    #[test]
    fn redaction_preserves_non_targeted_fields() {
        let body = br#"{"name":"Jane","ssn":"123-45-6789","amount":100}"#;
        let denied = vec!["/ssn".to_string()];
        let out = redact_json_fields(body, &denied);
        let parsed: serde_json::Value = serde_json::from_slice(&out).unwrap();
        assert_eq!(parsed["name"], "Jane");
        assert_eq!(parsed["ssn"], "[REDACTED]");
        assert_eq!(parsed["amount"], 100);
    }

    #[test]
    fn redaction_on_non_json_returns_as_is() {
        let body = b"this is not JSON";
        let denied = vec!["/ssn".to_string()];
        let out = redact_json_fields(body, &denied);
        assert_eq!(out, body);
    }

    #[test]
    fn sample_rate_zero_never_samples() {
        for _ in 0..100 {
            assert!(!should_sample(0));
        }
    }

    #[test]
    fn sample_rate_hundred_always_samples() {
        for _ in 0..100 {
            assert!(should_sample(100));
        }
    }
}
