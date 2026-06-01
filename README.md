# MulesoftArchitecture

Architecture and reference materials for MuleSoft-based integration designs — focused on the **API Gateway role**, with the heavy integration work (orchestration, transformation, EDI) delegated to an on-prem Microsoft stack.

## Docs

| Doc | What it covers |
|---|---|
| [01 — API Gateway Architecture](docs/01-api-gateway-architecture.md) | Scope, product choice (Anypoint Flex Gateway, Connected mode), internal+external dual-listener topology, SaaS vs on-prem trade-offs, private connectivity (Private Space + Transit Gateway + Direct Connect — no public internet), sizing for 100K calls/day |
| [02 — Policies](docs/02-policies.md) | Policy bundles per listener (external/internal/partner), recommended order, SLA tiers, custom policies, anti-patterns |
| [03 — Identity](docs/03-identity.md) | IdP choice (Entra/Okta/Ping/Cognito), OAuth Client Credentials for partners, mTLS for internal services, Auth Code+PKCE for users, JWKS caching, backend identity propagation |
| [04 — CI/CD](docs/04-cicd.md) | Spec → policy → publish pipeline; GitHub Actions + Anypoint CLI; dev/stg auto + prd manual gate; drift cleanup; rollback |
| [05 — Observability](docs/05-observability.md) | SLO set, metrics catalog, dashboards, access-log schema (PII-safe), W3C tracing end-to-end, Splunk HEC + Datadog wiring, alert matrix with runbooks |

## Conventions

- **API specs:** OAS 3.x in Anypoint Exchange, published from this Git repo via CI.
- **Naming:** `<domain>-<purpose>-api` (e.g. `orders-public-api`, `inventory-internal-api`).
- **Environments:** `dev`, `stg`, `prd` — separate Anypoint business groups or environments per tenancy.
- **Diagrams:** Mermaid (renders inline on GitHub). PNG renders to follow.
