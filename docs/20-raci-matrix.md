# 20 — RACI Matrix per Phase

Who is **R**esponsible, **A**ccountable, **C**onsulted, **I**nformed for each major activity across the engagement. Resolves the most common project failure mode: an activity that everyone thinks someone else is doing.

> **RACI conventions**: There is exactly **one A (Accountable)** per row — the single person/team on the hook. **R (Responsible)** is who does the work (can be multiple). **C (Consulted)** is two-way input. **I (Informed)** is one-way notification. If a cell would have both R and A, mark it **R/A**.

---

## 1. Stakeholders defined

### Consultant team (from [doc 18](18-team-roles-and-skills.md))

| Code | Role |
|---|---|
| **CSA** | Consultant Solution Architect |
| **CLE** | Consultant Lead Engineer |
| **CE** | Consultant Engineer(s) |
| **CDO** | Consultant DevOps Engineer |
| **CSE** | Consultant Security Engineer |
| **CPM** | Consultant Project Manager |

### Client team

| Code | Role |
|---|---|
| **EA** | Client Enterprise / Solution Architect |
| **APO** | Client API Product Owner |
| **NET** | Client Network Team |
| **SEC** | Client Security / CISO Office |
| **OPS** | Client Operations / Platform Team |
| **PKI** | Client PKI / Certificate Team |
| **IAM** | Client Identity & Access Management Team |
| **LEG** | Client Legal / DPO / Compliance |
| **BIZ** | Client Business / API Consumer Group |
| **CPM-C** | Client Project Manager |

### External

| Code | Role |
|---|---|
| **MS** | MuleSoft Account / Support |
| **PT** | Pen Test Firm |

---

## 2. Phase 1 — Discovery & Design

| Activity | CSA | CLE | CE | CDO | CSE | CPM | EA | APO | NET | SEC | OPS | IAM | LEG | CPM-C | MS |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| Kickoff meeting | C | C | — | C | C | R/A | C | C | C | C | C | C | I | C | I |
| Stakeholder workshops (logistics + facilitation) | R | C | — | — | C | R/A | I | I | I | I | I | I | — | C | — |
| Current-state assessment | R/A | C | — | C | C | I | C | I | C | C | C | C | — | I | I |
| Functional + NFR requirements gathering | R | C | — | — | C | R/A | C | R | I | C | C | C | I | C | — |
| Architecture envisioning | R/A | C | — | — | C | I | C | I | C | C | C | C | — | I | C |
| High-level architecture document | R/A | C | — | — | C | I | C | I | I | I | I | I | I | I | — |
| Detailed design — gateway + Redis | R/A | R | — | C | C | I | C | I | C | C | C | I | — | I | C |
| Detailed design — security + identity | C | C | — | — | R/A | I | C | I | I | C | I | C | C | I | — |
| Detailed design — network + connectivity | C | C | — | C | — | I | C | I | R/A | I | I | I | — | I | — |
| ARB submission + review | R | C | — | — | C | R/A | A | C | C | C | C | C | C | C | — |
| Design sign-off + baselining | R | C | — | — | C | R/A | A | C | C | C | C | C | C | C | I |
| Phase 1 status reporting (weekly) | C | C | — | — | — | R/A | I | I | I | I | I | I | — | C | — |

**Key reads:**
- Architecture decisions: Consultant Architect drives, Client Enterprise Architect approves at ARB
- Network design: Network Team is the A — Consultant proposes, Network Team accepts/modifies
- Identity / security design: Consultant Security Engineer drives in partnership with IAM + Security teams

---

## 3. Phase 2 — Foundation Setup

### 3.1 Infrastructure & networking

| Activity | CSA | CLE | CE | CDO | CSE | CPM | EA | APO | NET | SEC | OPS | PKI | IAM | LEG | CPM-C | MS |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| VM provisioning (per env) | I | C | — | R | — | I | I | — | C | — | R/A | — | — | — | I | — |
| FW rules + DNS + LB setup | C | C | — | R | C | I | I | — | R/A | C | C | — | — | — | I | — |
| TLS certificates (internal + public CA) | C | C | — | R | C | — | I | — | I | C | I | R/A | — | — | I | — |
| HashiCorp Vault deployment + integration | C | R | — | R/A | C | I | I | — | C | C | C | — | — | — | I | — |
| Active Directory / IdP wiring | C | R | — | C | R | I | I | — | I | C | I | — | R/A | — | I | — |
| Smoke-test infrastructure stack per env | C | R | C | R/A | C | I | I | — | C | C | C | — | C | — | I | — |

### 3.2 Omni Gateway + Redis installation

| Activity | CSA | CLE | CE | CDO | CSE | CPM | EA | NET | SEC | OPS | PKI | IAM | CPM-C | MS |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| Omni Gateway install + registration (per env) | C | R/A | R | R | C | I | I | C | I | C | I | C | I | C |
| Redis Sentinel cluster install + TLS/AUTH | C | R | C | R/A | C | I | I | C | C | C | R | — | I | — |
| Connected-mode validation against Anypoint | C | R/A | R | C | C | I | I | I | I | I | — | — | I | C |
| End-to-end smoke test (per env) | C | R/A | R | C | C | I | I | I | I | C | — | C | I | — |
| DR drill (failover validation) | C | R/A | R | R | C | I | C | C | C | C | — | — | C | — |

### 3.3 Anypoint Platform configuration

| Activity | CSA | CLE | CE | CDO | CSE | CPM | EA | APO | OPS | IAM | CPM-C | MS |
|---|---|---|---|---|---|---|---|---|---|---|---|---|
| Business groups / environment hierarchy | R/A | R | — | — | C | I | C | C | I | C | I | C |
| Connected Apps per env | C | R/A | R | — | C | — | I | I | I | C | I | C |
| Custom domains + DLB config | C | R | C | R/A | — | I | I | I | C | — | I | — |
| Anypoint Monitoring dashboards | C | R | R/A | C | — | I | I | C | C | — | I | — |
| Audit log forwarding to SIEM | C | R | C | R/A | C | I | I | — | C | — | I | — |
| Anypoint Access Management (RBAC + SSO) | C | R | — | — | R/A | I | I | — | I | C | I | — |

### 3.4 CI/CD pipeline setup

| Activity | CSA | CLE | CE | CDO | CSE | CPM | OPS | CPM-C |
|---|---|---|---|---|---|---|---|---|
| Repo structure + branching strategy | C | R/A | R | C | — | I | C | I |
| OpenAPI lint pipeline (Spectral) | C | R | R/A | C | — | I | C | I |
| Anypoint CLI integration scripts | C | R/A | R | C | — | I | — | I |
| Publish-to-Exchange automation | C | R/A | R | — | — | I | — | I |
| Policy YAML application + drift cleanup | R | R/A | R | C | — | I | — | I |
| Smoke test automation in CI | C | R | R/A | C | — | I | — | I |
| Environment promotion gates | C | R/A | C | C | C | I | C | I |
| Rollback procedures + scripts | C | R/A | R | R | — | I | C | I |
| Pipeline documentation | C | R | R/A | C | — | I | C | I |

### 3.5 Security posture

| Activity | CSA | CLE | CSE | CPM | EA | SEC | IAM | LEG | PT | CPM-C |
|---|---|---|---|---|---|---|---|---|---|---|
| IdP federation setup + testing | C | R | R/A | I | C | C | R | — | — | I |
| OAuth flows (multiple grant types) | C | R | R/A | I | C | C | R | — | — | I |
| mTLS configuration + cert rotation automation | C | R | R/A | I | I | C | — | — | — | I |
| Threat protection policy bundle | R | R | R/A | I | C | C | — | — | — | I |
| WAF / Anypoint Security Edge tuning | C | R | R/A | I | I | C | — | — | — | I |
| Pen test scope definition + coordination | C | I | R/A | C | I | R | — | C | C | I |
| Pen test execution | I | I | C | I | I | C | — | — | R/A | I |
| Pen test remediation | C | R | R/A | I | I | C | C | — | C | I |
| Security review pack (CISO sign-off) | C | I | R/A | C | I | A | — | C | — | I |
| DPA / BAA paperwork support | I | I | R | C | I | C | — | R/A | — | I |

### 3.6 Observability foundation

| Activity | CSA | CLE | CE | CDO | CPM | OPS | CPM-C |
|---|---|---|---|---|---|---|---|
| Splunk HEC ingest setup + indexes | C | R | C | R/A | I | C | I |
| Datadog agent deployment + dashboards | C | R | R | R/A | I | C | I |
| W3C `traceparent` propagation validation | R | R/A | R | C | I | I | I |
| Alert routing (PagerDuty / Teams) | C | R | C | R/A | I | C | I |
| Runbook authoring (per-alert) | C | R | R | R/A | I | C | I |

---

## 4. Phase 3 — First-Wave API Onboarding (per-API)

This table applies once per API. R/A may shift slightly based on tier (Tier 3 brings the Architect closer).

| Activity | CSA | CLE | CE | CDO | CSE | CPM | EA | APO | OPS | IAM | BIZ | CPM-C |
|---|---|---|---|---|---|---|---|---|---|---|---|---|
| API intake + tier classification | R | C | — | — | C | I | C | R/A | — | — | C | I |
| Discovery + requirements (with API team) | C | C | R | — | — | I | C | R/A | — | — | R | I |
| OpenAPI spec authoring / review | C | R | R/A | — | — | I | C | C | — | — | C | I |
| Policy bundle config | C | R/A | R | — | C | — | I | I | — | — | — | I |
| Transformation logic (Function / Logic App) | C | R/A | R | C | — | I | I | C | C | — | C | I |
| Backend integration + auth wiring | C | R | R/A | C | C | I | I | C | C | C | C | I |
| IdP scope / role configuration | — | R | C | — | C | I | — | C | — | R/A | — | I |
| Test cases (positive + negative + idempotency) | C | R | R/A | — | — | I | — | C | — | — | C | I |
| Performance smoke test | C | R/A | R | C | — | I | — | C | C | — | — | I |
| Documentation + portal listing | I | C | R | — | — | I | — | R/A | I | — | C | I |
| Promotion through environments | C | R/A | R | R | C | I | I | C | C | — | — | I |
| Production go-live sign-off | C | C | — | C | C | I | A | R | C | — | C | A |

**Per-tier deviations:**
- **Tier 1 (Simple)**: CE handles most of it; CLE reviews
- **Tier 2 (Standard)**: same shape as above (the default)
- **Tier 3 (Complex)**: CSA is C on most rows, R on architecture-affecting decisions; CLE more hands-on

---

## 5. Phase 4 — Hypercare & Knowledge Transfer

| Activity | CSA | CLE | CE | CDO | CSE | CPM | EA | OPS | SEC | CPM-C |
|---|---|---|---|---|---|---|---|---|---|---|
| Daily ops standup (joint) | C | R | R | R | — | R/A | — | R | — | R |
| Severity 1 incident lead | C | R/A | R | R | C | I | I | C | C | I |
| Severity 2 incident lead | I | R/A | R | R | I | I | — | C | — | I |
| Severity 3 incident triage | I | C | R/A | C | — | I | — | R | — | I |
| Post-incident review (per Sev 1/2) | C | R | C | C | C | R/A | I | C | C | I |
| Operations runbook authoring | C | R | R | R/A | C | I | I | R | — | I |
| Architecture handover documentation | R/A | R | — | C | C | I | C | C | — | I |
| Training sessions delivery (3 sessions) | C | R/A | R | R | C | I | C | R | — | I |
| Knowledge-transfer shadowing | C | R | R | R | — | I | I | R/A | — | I |
| Hypercare exit-criteria validation | C | R | — | C | C | R/A | C | A | — | A |

**Critical:** the Client Ops team is the **A** on knowledge-transfer shadowing — they're the ones who have to operate the platform after the consultant leaves. They have to demonstrate readiness before exit.

---

## 6. Steady-state operations (monthly)

| Activity | CSA | CLE | CE | CDO | CSE | CPM | EA | OPS | SEC | CPM-C |
|---|---|---|---|---|---|---|---|---|---|---|
| Platform health monitoring + reporting | I | C | C | R | — | R/A | I | C | I | I |
| OS + Omni Gateway + Redis patching | I | C | — | R/A | — | I | I | R | — | I |
| Certificate rotation (quarterly avg) | I | C | — | R | C | I | I | R/A | I | I |
| Bug triage + fix execution | C | R | R/A | C | — | I | I | C | — | I |
| Policy tuning + drift cleanup | C | R/A | R | C | — | I | I | C | — | I |
| DR drill (quarterly) | C | R/A | R | R | C | I | I | C | — | I |
| Documentation maintenance | C | R | R/A | R | — | I | I | R | — | I |
| Monthly status report + steering committee | C | C | — | C | C | R/A | I | C | I | I |
| On-call coverage premium | — | R | R | R/A | — | I | I | C | — | I |
| New API onboarding requests intake | C | R/A | C | — | — | I | C | C | — | C |

**Pattern**: in steady-state, **Client Ops** moves into the R role on routine operations (patching, monitoring, cert rotation) with consultants providing backup, escalation, and expertise. This is the **goal of knowledge transfer** — if consultants are still R/A on routine ops in month 12, KT didn't work.

---

## 7. Annual recurring activities

| Activity | CSA | CLE | CE | CDO | CSE | CPM | EA | OPS | SEC | LEG | CPM-C | PT | MS |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| Annual security review + pen test | C | C | C | C | R/A | I | C | C | R | I | C | R | — |
| Anypoint major-version upgrade | R | R/A | R | R | C | I | C | C | C | — | C | — | C |
| ARB annual walkthrough | R/A | C | — | — | C | I | A | C | C | I | C | — | — |
| Capacity planning review | R/A | R | C | R | — | I | C | R | — | — | C | — | C |
| Full-scale DR drill | C | R | R | R | C | I | C | R/A | C | — | C | — | — |
| Compliance attestation audit support | C | C | — | C | R | I | C | C | R/A | R | C | — | — |

---

## 8. Common RACI smells to watch for

| Smell | Why bad | Fix |
|---|---|---|
| Multiple A's on one row | Diffused accountability; nobody owns it when it fails | Pick exactly one A; others become R or C |
| No A on a row | Activity orphaned; nothing happens | Always assign an A, even for low-priority work |
| Same person is A on > 30% of rows for a phase | Bottleneck waiting to happen | Distribute A among the team based on actual ownership |
| Client is "I" on something that affects their ops | Surprise when it goes live | Promote to "C" minimum so they have input |
| Consultant is "A" on long-term ops in steady-state | KT didn't transfer accountability | Re-baseline RACI at month 3/6 of steady-state to push A toward client |
| RACI matrix never reviewed after kickoff | Roles drift; arguments about "I thought you were doing that" | Quarterly RACI review at steering committee |

---

## 9. RACI as a SOW artifact

The RACI matrix should be an **appendix to the SOW**, not a separate document. This means:

- Changes to the RACI require change-control (formal amendment to the SOW)
- The matrix locks scope and ownership simultaneously
- New activities require an explicit RACI row before they start
- Dispute resolution: the SOW's RACI is the first reference for "who was supposed to do this?"

---

## 10. How to use this matrix

1. **At SOW signature**: walk through the matrix with the client; capture corrections to ownership assumptions before commencement
2. **At each phase kickoff**: review the rows for that phase with the relevant teams; confirm A's are still right
3. **At weekly status**: any incident or scope question routes back to the matrix as the first authority on ownership
4. **At quarterly review**: re-baseline the matrix; promote client team into R/A roles as KT progresses
5. **At engagement close**: the final RACI shows the steady-state ownership distribution — useful artifact for the client's ongoing operations

---

## Related

- [16 — Implementation & Operations Estimate](16-consulting-estimate.md) — the activities this matrix assigns
- [17 — Engagement Models](17-engagement-models.md) — the commercial structure these roles deliver under
- [18 — Team Roles & Skills](18-team-roles-and-skills.md) — the roles referenced in this matrix
- [19 — Pricing Model](19-pricing-model.md) — the pricing applied to the consultant roles in this matrix
