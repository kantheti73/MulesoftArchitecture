# 16 — Implementation & Operations Estimate

Consulting-grade effort estimate for delivering the Omni Gateway platform described in docs 01–15 and operating it in steady state. Hours only — pricing is left to commercial discussion. Numbers are intentionally **realistic-and-defensible** rather than aggressive in either direction.

---

## 1. Basis & scope

| | |
|---|---|
| **Target architecture** | On-prem Omni Gateway (Connected mode), 4 environments (DEV / QA / Acceptance / Prod), 5M calls/day Prod target, citizen-data workload |
| **Component coverage** | Omni Gateway (per [doc 09](09-onprem-install.md)) · Redis Sentinel ([doc 10](10-redis-cache.md)) · Anypoint Platform config · IdP federation ([doc 03](03-identity.md)) · CI/CD via Anypoint CLI + GitHub Actions ([doc 04](04-cicd.md)) · Observability into Splunk + Datadog ([doc 05](05-observability.md)) |
| **Estimating approach** | Bottom-up per task, validated against industry benchmarks for comparable Anypoint Gateway engagements |
| **Buffer methodology** | Each task includes realistic execution time + reasonable rework/coordination overhead. No separate "buffer" line — buffer is embedded in task estimates |
| **Currency** | Hours only |

---

## 2. Estimate summary (executive view)

| Phase | Duration | Effort (hrs) |
|---|---|---|
| **Phase 1** — Discovery & Design | 6 weeks | **~180** |
| **Phase 2** — Foundation Setup | 12 weeks | **~960** |
| **Phase 3** — First Wave API Onboarding (~10 APIs assumed) | 10 weeks (parallel to Phase 2 tail) | **~1,050** |
| **Phase 4** — Hypercare & Knowledge Transfer | 8 weeks | **~310** |
| **Initial delivery total (Phases 1–4)** | ~6 months elapsed | **~2,500 hrs** |
| **Per-API onboarding** (after first wave) | varies (see §7) | **40–250 hrs per API depending on complexity** |
| **Steady-state operations** (after hypercare) | recurring monthly | **~130 hrs / month** |
| **Annual recurring** (audits, drills, upgrades) | yearly | **~280 hrs / year** in addition to monthly |

For year-2 budget planning at a typical pace of ~20 new APIs/year + steady-state ops + annual recurring, expect **~3,000 hrs/year ongoing engagement**.

---

## 3. Phase 1 — Discovery & Design

| Task | Hrs | Notes |
|---|---|---|
| Stakeholder workshops (3 sessions) | 24 | Architect + BA; one per stakeholder group (security, network, API teams) |
| Current-state assessment | 24 | Existing API estate, IdP config, network topology |
| Requirements gathering & documentation | 32 | Functional + non-functional, SLAs, compliance scope |
| Architecture envisioning sessions | 24 | Joint with client architecture team |
| High-level architecture document | 24 | Even with this repo as a starting point, client-specific tailoring |
| Detailed design — gateway + Redis | 24 | |
| Detailed design — security + identity | 16 | OAuth flows, mTLS, IdP integration plan |
| ARB submission & 2 review cycles | 24 | Prep + presentation + iteration |
| Design sign-off & baselining | 8 | |
| Phase 1 PM coordination | (incl. above) | Project manager allocated 20% |
| **Subtotal** | **200** | |
| Less: leverage existing repo docs | **(-20)** | Honest discount for docs we can adapt |
| **Phase 1 total** | **~180** | |

---

## 4. Phase 2 — Foundation Setup

### 4.1 Infrastructure setup (per environment)

| Task | DEV | QA | UAT | Prod |
|---|---:|---:|---:|---:|
| VM provisioning coordination | 8 | 16 | 12 | 16 |
| Network: FW rules + DNS + LB setup | 16 | 24 | 24 | 32 |
| TLS certificates (internal CA + public CA) | 8 | 16 | 16 | 24 |
| HashiCorp Vault integration | 16 | 24 | 16 | 24 |
| Active Directory / IdP wiring | 8 | 16 | 12 | 16 |
| Smoke test infrastructure stack | 8 | 16 | 12 | 16 |
| **Per-env subtotal** | **64** | **112** | **92** | **128** |

**Infrastructure total: ~396 hrs** (across all 4 envs)

### 4.2 Omni Gateway + Redis installation (per environment)

| Task | DEV | QA | UAT | Prod |
|---|---:|---:|---:|---:|
| Omni Gateway install + registration | 8 | 12 | 12 | 16 |
| Redis Sentinel cluster install + TLS/AUTH config | 8 | 16 | 12 | 16 |
| Connected-mode validation | 4 | 8 | 8 | 8 |
| Smoke test (end-to-end) | 4 | 8 | 8 | 16 |
| DR drill (failover validation) | — | 16 | 8 | 16 |
| **Per-env subtotal** | **24** | **60** | **48** | **72** |

**Install total: ~204 hrs**

### 4.3 Anypoint Platform configuration (one-time)

| Task | Hrs |
|---|---:|
| Business groups / environment hierarchy setup | 16 |
| Connected Apps per environment (4) | 16 |
| Custom domains + DLB configuration | 24 |
| Anypoint Monitoring dashboards | 24 |
| Audit log forwarding to SIEM | 16 |
| Access Management (RBAC, SSO) | 16 |
| **Subtotal** | **112** |

### 4.4 CI/CD pipeline setup (one-time)

| Task | Hrs |
|---|---:|
| Repo structure + branching strategy | 16 |
| OpenAPI lint pipeline (Spectral) | 16 |
| Anypoint CLI integration scripts | 24 |
| Publish-to-Exchange automation | 16 |
| Policy YAML application + drift cleanup | 32 |
| Smoke test automation in CI | 24 |
| Environment promotion gates | 16 |
| Rollback procedures + scripts | 16 |
| Documentation of pipeline | 16 |
| **Subtotal** | **176** |

### 4.5 Security posture (one-time)

| Task | Hrs |
|---|---:|
| IdP federation setup + testing | 32 |
| OAuth flows (Client Credentials, Auth Code+PKCE, mTLS, JWT validation) | 40 |
| mTLS configuration + cert rotation automation | 32 |
| Threat protection policy bundle (JSON threat, schema validate, IP allowlists) | 24 |
| WAF / Anypoint Security Edge tuning | 16 |
| Pen test coordination (with security firm) | 16 |
| Pen test remediation | 40 |
| Security review pack (CISO / DPO sign-off prep) | 24 |
| DPA / BAA paperwork support | 16 |
| **Subtotal** | **240** |

### 4.6 Observability foundation (one-time)

| Task | Hrs |
|---|---:|
| Splunk HEC ingest setup + indexes | 24 |
| Datadog agent deployment + dashboards | 24 |
| W3C `traceparent` propagation validation | 16 |
| Alert routing (PagerDuty / Teams / email) | 16 |
| Runbook authoring (per-alert) | 24 |
| **Subtotal** | **104** |

### Phase 2 total

| Workstream | Hrs |
|---|---:|
| Infrastructure (4 envs) | 396 |
| Omni Gateway + Redis installation (4 envs) | 204 |
| Anypoint Platform configuration | 112 |
| CI/CD pipeline | 176 |
| Security posture | 240 |
| Observability foundation | 104 |
| Less: cross-workstream coordination already counted | -272 |
| **Phase 2 total** | **~960** |

The deduction reflects that some tasks (e.g. test execution) overlap across workstreams and shouldn't be double-counted.

---

## 5. Phase 3 — First wave API onboarding

Assumes a representative first wave of **10 APIs** balanced across complexity tiers:

| Complexity | Count | Per-API hrs (see §7) | Subtotal |
|---|---:|---:|---:|
| Simple (passthrough) | 3 | 50 | 150 |
| Standard (policies + light transform) | 5 | 100 | 500 |
| Complex (multi-backend / transformation / partner SLA) | 2 | 200 | 400 |
| **Subtotal** | **10** | | **1,050** |

Onboarding runs in parallel with Phase 2's final weeks once the platform is stable in QA.

---

## 6. Phase 4 — Hypercare & knowledge transfer

| Task | Hrs |
|---|---:|
| 8-week hypercare on-call (consultant rotation) | 160 |
| Daily ops standups (joint) | 20 |
| Issue triage and resolution buffer | 60 |
| Operations runbook authoring | 32 |
| Architecture handover documentation | 16 |
| 3× training sessions (prep + delivery) | 24 |
| **Phase 4 total** | **~310** |

---

## 7. Per-API onboarding model (the unit economics)

Three tiers cover most realistic APIs. Use these for future-wave budgeting.

### Tier 1 — Simple passthrough API

Single backend, standard policies, no transformation.

| Task | Hrs |
|---|---:|
| Discovery + spec gathering | 6 |
| OpenAPI spec authoring / review | 12 |
| Standard policy bundle attach | 6 |
| IdP scope / role configuration | 4 |
| Smoke + contract tests | 8 |
| Backend integration validation | 6 |
| Documentation + portal listing | 4 |
| Promotion through environments (DEV → QA → UAT → Prod) | 4 |
| **Tier 1 total** | **~50** |

### Tier 2 — Standard API with policies + light transform

Most production APIs. Single backend, custom policy chain, light header/payload transformation, full test coverage.

| Task | Hrs |
|---|---:|
| Discovery + requirements (sometimes spans 2 sessions) | 12 |
| OpenAPI spec authoring | 16 |
| Policy bundle (custom rate limits, scopes, claims) | 12 |
| Transformation logic (Function or Logic App snippet) | 16 |
| Backend integration + auth wiring | 12 |
| Test cases (positive + negative + idempotency) | 16 |
| Performance smoke test | 6 |
| Documentation + portal | 6 |
| Promotion through environments | 4 |
| **Tier 2 total** | **~100** |

### Tier 3 — Complex API (composite / orchestrated / partner-facing)

Multi-backend composite, custom orchestrator (Function / Logic App), partner-tier SLA, EDI mapping, or content-based routing.

| Task | Hrs |
|---|---:|
| Discovery + design workshops (often 3+ sessions) | 24 |
| OpenAPI spec + canonical event schema | 24 |
| Orchestrator implementation (Function / Logic App) | 40 |
| Multiple backend integrations | 24 |
| Custom policy + tier-based routing | 16 |
| Idempotency + saga / compensation logic | 16 |
| Comprehensive test cases (incl. load + chaos) | 24 |
| Partner onboarding documentation + SLA wording | 12 |
| Performance + DR test | 12 |
| Promotion + UAT signoff (often partner involvement) | 8 |
| **Tier 3 total** | **~200** |

### Per-API tier picker

| If the API… | Tier |
|---|---|
| Just proxies one backend with standard auth | 1 |
| Has 1–2 backend hops + transformation that fits in a small Function | 2 |
| Spans 3+ backends, has a workflow, or has partner-specific SLAs | 3 |
| Involves EDI translation or AS2 | 3 (add Logic Apps EIP effort separately) |

---

## 8. Steady-state operations & maintenance (monthly)

After Phase 4 hypercare, the engagement transitions to a run-rate model.

| Activity | Hrs / month |
|---|---:|
| Platform health monitoring + reporting | 16 |
| OS + Omni Gateway + Redis patching | 24 |
| Certificate rotation (averaged from quarterly) | 8 |
| Bug triage + fix execution | 32 |
| Policy tuning + drift cleanup | 16 |
| DR drill (averaged from quarterly) | 8 |
| Documentation maintenance | 8 |
| Monthly status report + steering committee | 8 |
| On-call coverage premium (averaged) | 8 |
| **Steady-state subtotal** | **~128 hrs / month** |

This is roughly **0.75 FTE** of consulting effort baseline. New APIs are billed separately per §7.

---

## 9. Annual recurring activities (in addition to monthly)

| Activity | Hrs / year |
|---|---:|
| Annual security review + pen test coordination | 80 |
| Anypoint platform major-version upgrade | 80 |
| Architecture review board annual walkthrough | 24 |
| Capacity planning review | 24 |
| Disaster-recovery full-scale drill (Prod-validated) | 40 |
| Annual compliance attestation support (SOC 2 / DPA renewal) | 32 |
| **Annual recurring subtotal** | **~280 hrs / year** |

---

## 10. Year-1 vs Year-2+ rollup

| Bucket | Year 1 (one-time + run-rate from month 7) | Year 2 ongoing |
|---|---:|---:|
| Phases 1–4 (initial delivery) | 2,500 | — |
| Steady-state ops (6 months in Year 1, 12 months Year 2) | 768 | 1,536 |
| New API onboarding (assume 5 in Year 1 post-go-live, 20 in Year 2) | ~500 | ~2,000 |
| Annual recurring | 280 | 280 |
| **Year total** | **~4,050 hrs** | **~3,820 hrs** |

Realistic resource shape for Year 1: **architect + lead engineer + 2 engineers + part-time DevOps + part-time PM**, ramping in/out across phases. Year 2 trims toward **lead engineer + 1–2 engineers + part-time architect + part-time PM**.

---

## 11. Estimating assumptions

These are explicit so the estimate can be defended (and renegotiated cleanly if they change):

| # | Assumption | If wrong → |
|---|---|---|
| 1 | Client provides VMs / network capacity per [doc 09 §6.5](09-onprem-install.md#65-environment-matrix--dev--qa--acceptance--prod) sizing | Add ~80 hrs to Phase 2 for hardware procurement coordination |
| 2 | Client owns ExpressRoute / Direct Connect provisioning | Add ~120 hrs for network procurement + carrier coordination |
| 3 | Client owns HashiCorp Vault (or has Azure Key Vault standardized) | Add ~80 hrs for Vault HA design + install |
| 4 | IdP (Entra ID / Okta / Ping) already operational; we integrate | Add ~120 hrs if IdP itself needs to be deployed |
| 5 | First-wave APIs already have either OAS specs OR clear functional definitions | Add 8–16 hrs per API for forensic spec authoring |
| 6 | Pen test is conducted by a third-party security firm (we coordinate, not execute) | Add ~80 hrs if we lead pen testing |
| 7 | Backend services accept gateway-injected identity headers without backend code changes | Add ~24 hrs per API for backend auth coupling |
| 8 | Anypoint Titanium subscription includes the connectors / policies we use | Add ~40 hrs if custom WASM policies are needed |
| 9 | Client provides one technical PM as a single point of contact | Add ~80 hrs/year for cross-team coordination overhead |
| 10 | Compliance scope: standard SOC 2 + state PII laws (per [doc 07](07-data-protection.md)) — no HIPAA / FedRAMP High in scope | Add ~120 hrs for BAA + compliance attestation if HIPAA in scope; ~400 hrs if FedRAMP path required |
| 11 | All work performed remotely; on-site visits limited to 2 per phase | Add travel hours per requested visit (~16 hrs/visit including travel) |
| 12 | Standard business hours; emergency on-call after hypercare is best-effort within steady-state hrs | Add ~40 hrs/month for 24/7 on-call rotation |

---

## 12. Out of scope (explicit exclusions)

To prevent scope creep — these are NOT included in the estimate unless added separately:

- Production hardware procurement / purchase
- Network circuit (ExpressRoute / Direct Connect) procurement or carrier billing
- IdP product license or initial deployment (assumed operational)
- HashiCorp Vault license or initial deployment
- Splunk / Datadog license fees
- Anypoint Platform subscription costs
- Pen test execution by third-party firms
- Application-side code changes to backend services (BizTalk pipelines, Logic Apps, .NET worker code)
- Data migration / historical data transformation
- 24/7 on-call coverage post-hypercare (best-effort only in steady-state baseline)
- Compliance attestation audits performed by client's auditors (consulting support during audit is in scope; audit itself is not)
- Logic Apps Enterprise Integration Pack / Integration Account configuration for EDI workflows (separate engagement — typically ~200–400 hrs for a starter set)
- Custom WASM policy development (assumes built-in Omni Gateway policies suffice)
- Custom Datadog / Splunk dashboard development beyond the baseline set
- End-user training programs beyond technical knowledge transfer to client ops team

---

## 13. Risks that could materially change the estimate

| Risk | Impact | Mitigation |
|---|---|---|
| Anypoint version gap between current docs and the GA version at install time | + 80–160 hrs Phase 2 rework | Confirm Anypoint major version at SOW signature |
| Network team change-control queue is slow (4+ week lead) | Schedule extension; no hours impact | Start network coordination at Phase 1 kickoff |
| Pen test surfaces architectural changes (not just config fixes) | + 80–240 hrs Phase 2 | Schedule pen test early in Phase 2 to surface issues sooner |
| First-wave APIs are systematically more complex than assumed (e.g. 8 of 10 are Tier 3) | + 600 hrs Phase 3 | Per-API tier review at end of Phase 1; renegotiate before Phase 3 |
| Compliance scope expands mid-project (e.g. HIPAA added) | + 120–400 hrs | Lock compliance scope in SOW |
| Client ops team isn't ready for handover at end of hypercare | Extended hypercare → + 40 hrs/week | Joint readiness criteria at hypercare start |
| BizTalk / backend system changes required to accept gateway auth model | + 16–40 hrs per affected backend | Surface in Phase 1 design |
| MuleSoft makes a breaking change in Connected-mode Redis or policy format mid-engagement | + 80–200 hrs reactive | Track MuleSoft roadmap; reassess at major releases |

---

## 14. Phase-by-phase resource shape (illustrative)

| Phase | Architect | Lead Eng | Engineer | DevOps | PM |
|---|---:|---:|---:|---:|---:|
| Phase 1 (Discovery & Design) | 0.6 FTE | 0.2 FTE | — | 0.1 FTE | 0.2 FTE |
| Phase 2 (Foundation) | 0.4 FTE | 0.8 FTE | 1.0 FTE | 0.6 FTE | 0.3 FTE |
| Phase 3 (First wave APIs) | 0.2 FTE | 0.6 FTE | 2.0 FTE | 0.2 FTE | 0.3 FTE |
| Phase 4 (Hypercare) | 0.2 FTE | 0.6 FTE | 0.6 FTE | 0.4 FTE | 0.2 FTE |
| Steady-state ops | 0.1 FTE | 0.3 FTE | 0.3 FTE | 0.2 FTE | 0.1 FTE |

These are rough load profiles for capacity planning — actual SOW will allocate by week.

---

## 15. How to use this estimate

1. **Take Phase 1–4 (~2,500 hrs) as the initial commitment.** That gets the platform live with the first 10 APIs.
2. **Use §7 (per-API tiers) for future API budgeting.** Each new API is budgeted independently — no surprise scope creep.
3. **Use §8 + §9 (~1,820 hrs/year) for ongoing operations budget.** Predictable run rate.
4. **Review §11 (assumptions) and §12 (exclusions) explicitly in the SOW.** Most estimate disputes trace back to one of those being misaligned at signature.
5. **Revisit §13 (risks) quarterly.** Material changes trigger an estimate update, not a margin write-off.

---

## Related

- [09 — On-Prem Install Guide](09-onprem-install.md) — the Phase 2 install work this estimate covers
- [02 — Policies](02-policies.md) + [03 — Identity](03-identity.md) — the security posture work in Phase 2 §4.5
- [04 — CI/CD](04-cicd.md) — the pipeline scope in Phase 2 §4.4
- [05 — Observability](05-observability.md) — the observability foundation in Phase 2 §4.6
- [14 — Redis Assumption Register](14-redis-assumptions.md) — pattern this doc follows for assumptions discipline
- [15 — Orchestration](15-orchestration.md) — the orchestration complexity that drives Tier 3 per-API estimates
