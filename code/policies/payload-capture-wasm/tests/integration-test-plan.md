# Integration test plan — payload-capture-wasm

The unit tests in `src/lib.rs` (`cargo test`) cover the pure functions
(base64 encoding, JSON redaction, sampling math). Those don't need the
gateway. Integration testing requires the actual WASM filter loaded into
a Flex/Omni Gateway instance.

## Test environments

| Env | What gets tested | When |
|---|---|---|
| **Local Flex Gateway in Docker** | Smoke test the WASM filter loads, parses config, dispatches HTTP calls to a mock Service Bus endpoint | Every PR |
| **DEV Anypoint env** | End-to-end: real Anypoint policy attachment, real Service Bus DEV queue, real test API | Per merge to main |
| **QA env** | Load test: 1,200 TPS for 30 min with sample_rate_percent=10; verify gateway latency unchanged + queue depth manageable | Per release |
| **Prod** | Canary: enable on 1 low-volume API at 1% sample, watch for 24h, then promote | Release window |

## Local smoke test (Docker)

```powershell
# 1. Build the WASM filter
cargo build --release --target wasm32-wasi
Copy-Item target/wasm32-wasi/release/payload_capture_wasm.wasm ./payload-capture-wasm.wasm

# 2. Run Flex Gateway in Docker with the filter loaded
docker run --rm -it `
  -e LICENSE_KEY=$env:ANYPOINT_LICENSE_KEY `
  -p 8081:8081 -p 9999:9999 `
  -v "${PWD}:/policies" `
  mulesoft/flex-gateway:1.4 `
  start --connected=false `
       --config /policies/test-flex-config.yaml

# 3. Spin up a mock Service Bus endpoint (a simple HTTPS listener)
#    that 201's everything POSTed to it
# In a separate terminal:
docker run --rm -p 8888:80 `
  -v "${PWD}/tests/mock-sb-responses:/srv" `
  mockserver/mockserver

# 4. Hit the gateway and confirm the mock SB received the capture
curl -X POST http://localhost:8081/test/api/orders `
  -H "Content-Type: application/json" `
  -H "X-Correlation-ID: smoke-test-001" `
  -d '{"customer":"acme","items":[{"sku":"123","qty":2}]}'

# 5. Verify mock SB log shows the captured payload
docker logs <mock-server-container-id>
```

## End-to-end DEV test

```powershell
# After deploy-policy.ps1 has uploaded the policy to Exchange and you've
# attached it to a DEV API with appropriate config

# 1. Provision the DEV Service Bus queue
az servicebus queue create `
  --namespace-name ssp-bus-dev `
  --name api-payload-capture `
  --resource-group ssp-rg

# 2. Generate a SAS token (Send-only) and put in Anypoint Secure Property
$sasKey = az servicebus queue authorization-rule keys list `
  --namespace-name ssp-bus-dev `
  --queue-name api-payload-capture `
  --name SenderOnly `
  --query primaryKey -o tsv
# (Compute the SAS Authorization header value — see Microsoft Service Bus
# REST docs for the SAS header format)

# 3. Hit a real DEV API a few times
for ($i = 1; $i -le 10; $i++) {
  curl -X POST https://api-dev.yourco.com/orders `
    -H "Authorization: Bearer $devToken" `
    -H "Content-Type: application/json" `
    -d "{\"orderId\":\"e2e-test-$i\",\"customer\":\"test\"}"
}

# 4. Read messages from the queue and verify
az servicebus queue show `
  --namespace-name ssp-bus-dev `
  --name api-payload-capture `
  --query "countDetails"

# Or with Service Bus Explorer / az-cli message peek
```

## Load test (QA)

Goal: confirm the WASM filter doesn't degrade gateway latency under load.

```
Target: 1,200 TPS sustained for 30 minutes (= our worst-case spike from doc 09 §6.2)
Sample rate: 10%
Expected:
  - Gateway p99 latency increase < 3 ms vs baseline (without policy)
  - Service Bus queue depth peaks < 1,000 (drained by consumer)
  - No 5xx from gateway
  - No DLQ messages from SB (publish-side errors)
```

Use k6 / Apache JMeter / Vegeta. Capture before/after percentiles.

## Sampling verification

Goal: confirm the actual sample rate matches the configured rate.

```
Send 10,000 requests with sample_rate_percent=10
Expect: 800-1,200 messages on the queue (uniform random ~10%)
```

## PII redaction verification

Goal: confirm denied fields are blanked.

```
Configure redact_body_fields: ["/ssn","/dob"]
Send: {"name":"Jane","ssn":"123-45-6789","dob":"1980-01-01","amount":100}
Expect on queue:
  {"name":"Jane","ssn":"[REDACTED]","dob":"[REDACTED]","amount":100}
```

## Failure-tolerance verification

Goal: confirm SB outage doesn't break the gateway request path.

```
1. Make SB queue unreachable (block FW, or change cluster config to invalid host)
2. Send 100 gateway requests
3. Verify: all 100 succeed at the gateway (no 5xx); WARN logs in
   gateway about failed dispatch_http_call; no DLQ messages (because no
   message ever reached SB)
4. Restore SB
5. Verify subsequent requests resume capture
```

## What to do if a test fails

| Symptom | Likely cause | Fix |
|---|---|---|
| Filter fails to load on gateway start | Cargo build target mismatch (`wasm32-wasip1` vs `wasm32-wasi`) | Verify Cargo build target matches Omni's expected ABI |
| Filter loads but config parse error in logs | Manifest YAML → JSON conversion drops a field | Check anypoint-cli upload payload; verify `configurationData` field names match the `PolicyConfig` struct |
| Captures work but never reach SB | Upstream cluster name mismatch | Ensure `sb_http_cluster_name` matches what's declared in Flex Gateway upstream config |
| SB returns 401 | SAS Authorization header malformed or expired | Regenerate SAS; ensure the header is the complete `SharedAccessSignature sr=...&sig=...` value |
| Gateway latency spike under load | Body buffering exceeded memory | Lower `max_payload_bytes`; check filter is releasing buffers post-publish |
| Some requests captured, some not | Sampling working as designed | If 100% needed, set `sample_rate_percent: 100` |
| Capture payloads contain PII you didn't expect | Redaction missing fields | Add to `redact_body_fields`; verify JSON pointer paths are correct |
