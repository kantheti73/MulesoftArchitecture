# `code/` — Custom Implementations

This is where actual implementation artifacts live (as opposed to architecture docs in `docs/` and commercial templates in `pricing/`). The repo started architecture-first; this folder is for production-bound code as use cases require it.

## Current contents

| Path | Purpose | Status |
|---|---|---|
| [`policies/payload-capture-wasm/`](policies/payload-capture-wasm/) | Custom Omni Gateway policy that captures every API call payload + headers and publishes to an Azure Service Bus queue (fire-and-forget, async, with PII redaction + sampling) | Reference implementation — needs build + integration test |

## Why a `code/` folder

The architecture docs (01–21) deliberately don't carry implementation code — most of the design work doesn't require it. But certain capabilities have no useful documentation without an actual implementation:

- Custom WASM gateway policies (no SaaS substitute exists for many use cases)
- Saga reference implementations for Durable Functions / Logic Apps
- Sample backend stubs for end-to-end testing
- One-off operational scripts

These belong here. Architecture remains in `docs/`.

## Conventions

- One subdirectory per capability under the right category (`policies/`, `orchestrators/`, `samples/`, etc.)
- Each subdir has its own `README.md` with build / deploy / operate instructions
- Production code includes its own tests
- Languages are chosen pragmatically per capability (Rust for WASM policies, C#/Python for Functions, etc.) — no monorepo language constraint

## Anti-patterns

- **Don't put architecture explanation here** — that's `docs/`'s job. The README for a `code/` subdir should be operationally focused (how to build, deploy, debug), not strategy.
- **Don't carry custom code you don't actually need** — every line is maintenance burden ([doc 18 §3 anti-pattern](../docs/18-team-roles-and-skills.md)). Reach for built-in policies / managed services first; custom code is a last resort.
- **Don't skip the README in a code subdir** — without it, the code is unmaintainable in 6 months when no one remembers why it exists.

## Related

- [`docs/02-policies.md`](../docs/02-policies.md) — policy strategy at the gateway tier
- [`docs/11-azure-service-bus-integration.md`](../docs/11-azure-service-bus-integration.md) — Service Bus integration patterns (the payload-capture policy consumes this)
- [`docs/07-data-protection.md`](../docs/07-data-protection.md) — PII implications of capturing payloads (critical for the payload-capture use case)
