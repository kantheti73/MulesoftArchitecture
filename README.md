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
| [06 — Azure-Resident Private Space (US Central)](docs/06-azure-private-space.md) | Azure-hosted data plane in centralus, ExpressRoute Private Peering, VNet peering, vWAN options, AWS-vs-Azure gotchas, what to verify with MuleSoft sales |
| [07 — Data Protection (Citizen Data on SaaS)](docs/07-data-protection.md) | 12 data-bleed vectors with mitigations, cloud-specific exposure surfaces, compliance landscape (HIPAA/GDPR/CCPA/FedRAMP), shared responsibility matrix, mandatory controls bundle, memory-dump policy with MuleSoft, residual risks, pre-go-live security checklist |
| [08 — Flex Gateway Deep-Dive](docs/08-flex-gateway.md) | Product-level honest evaluation: Envoy core, performance numbers, high points (architecture/perf/ops/cost/cloud-native), low points (maturity, capability gaps, custom WASM burden, lock-in, license floor), comparison vs Kong/Apigee/AWS/Azure/Mule runtime, verification checklist before committing |
| [09 — On-Prem Install Guide](docs/09-onprem-install.md) | Step-by-step install for both deployment paths: standalone Linux VM (RPM/Deb + systemd hardening) and Kubernetes (Helm + production values + HPA + PDB + NetworkPolicy). Side-by-side VM-vs-K8s trade-offs. Hardware sizing curve from 100K → 50M calls/day. Disaster recovery with 4 posture options + real-time config replication mechanics for both Connected mode (Anypoint pushes) and Local mode (GitOps + Ansible). F5 GSLB failover runbook. 18-item pre-install checklist. |
| [10 — Redis Shared Storage (Connected Mode)](docs/10-redis-cache.md) | **Required infrastructure** for Omni/Flex Gateway in Connected mode (per [MuleSoft security best practices](https://docs.mulesoft.com/gateway/latest/flex-security-best-practices#secure-redis-shared-storage)). Caches runtime configurations + request data + distributed policy state. Local mode does not use Redis. Note on Flex → Omni Gateway rebrand. SaaS path: MuleSoft provides Redis (verify provisioning). On-prem path: Redis Sentinel 3-node per DC (or Redis Enterprise for active/active). Sizing, install with TLS+AUTH+rename-command hardening, DR behavior, anti-patterns, PII implications. |

## Presentations

| File | What it covers |
|---|---|
| [presentations/MuleSoft-APIGateway-Architecture.pptx](presentations/MuleSoft-APIGateway-Architecture.pptx) | **SaaS-centric 10-slide deck**: title · synopsis · architecture (embedded diagram) · product choice · Flex Gateway vs Full Mule Runtime · network components to provision · capacity planning · Azure vs AWS Private Spaces · on-prem deployment · risks & next steps |
| [presentations/MuleSoft-APIGateway-OnPrem.pptx](presentations/MuleSoft-APIGateway-OnPrem.pptx) | **On-prem-only 10-slide deck**: title · why on-prem · on-prem architecture (embedded diagram) · Connected vs Local mode · VM vs Kubernetes install · hardware sizing curve · components to provision · DR postures with RTO/RPO · config replication mechanics · operational ownership shift + risks + next steps |
| [presentations/build_pptx.py](presentations/build_pptx.py) | Generator for the SaaS-centric deck |
| [presentations/build_onprem_pptx.py](presentations/build_onprem_pptx.py) | Generator for the on-prem-only deck |
| [presentations/diagrams/mulesoft_architecture.py](presentations/diagrams/mulesoft_architecture.py) | Generates the SaaS architecture PNG |
| [presentations/diagrams/onprem_architecture.py](presentations/diagrams/onprem_architecture.py) | Generates the on-prem architecture PNG (multi-DC active/active with F5 GSLB, Ansible AWX, Vault) |

## Conventions

- **API specs:** OAS 3.x in Anypoint Exchange, published from this Git repo via CI.
- **Naming:** `<domain>-<purpose>-api` (e.g. `orders-public-api`, `inventory-internal-api`).
- **Environments:** `dev`, `stg`, `prd` — separate Anypoint business groups or environments per tenancy.
- **Diagrams:** Mermaid (renders inline on GitHub). PNG renders to follow.
