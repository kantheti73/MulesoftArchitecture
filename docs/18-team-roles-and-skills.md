# 18 — Team Roles & Skills

Concrete role descriptions for the consulting team that delivers the engagement in [doc 16](16-consulting-estimate.md) under the engagement model in [doc 17](17-engagement-models.md). Each role has a defined seniority, skill profile, day-to-day responsibilities, and allocation across phases — so SOW staffing tables are explicit, not aspirational.

---

## 1. Team composition at a glance

| Role | Seniority | Year-1 allocation | Year-2 allocation |
|---|---|---|---|
| **Solution Architect** | Principal / 12+ yrs | ~0.3 FTE | ~0.15 FTE |
| **Integration Lead Engineer** | Senior / 8+ yrs | ~0.6 FTE | ~0.4 FTE |
| **Integration Engineer × 2** | Mid / 3–5 yrs | ~1.5 FTE | ~1.0 FTE |
| **DevOps / SRE Engineer** | Senior / 5+ yrs | ~0.4 FTE | ~0.2 FTE |
| **Security Engineer** | Senior / 5+ yrs | ~0.2 FTE (Phase 2 heavy) | ~0.1 FTE |
| **Technical Project Manager** | Senior / 5+ yrs | ~0.25 FTE | ~0.15 FTE |
| **Technical Writer** | Mid / 3+ yrs (optional) | ~0.1 FTE (Phase 4 heavy) | ~0.05 FTE |
| **Total** | — | **~3.35 FTE** | **~2.05 FTE** |

The team **ramps in for Phase 2**, peaks during **Phase 3 first-wave APIs**, **ramps down through Phase 4**, and settles into a leaner steady-state team.

---

## 2. Solution Architect

### Role profile

**Alternate titles**: Principal Architect · Integration Architect · Enterprise Architect

**Years of experience**: **12+** with at least **5 years on Anypoint Platform / API Gateway architectures** specifically.

### Core skills (must-have)

| Skill | Depth |
|---|---|
| Anypoint Platform architecture (Flex / Omni Gateway, API Manager, Exchange) | Expert |
| API design (REST, OpenAPI 3.x, gRPC) | Expert |
| OAuth 2.0 / OIDC / JWT / mTLS | Expert — can debug JWT validation issues from a packet trace |
| Distributed systems (CAP, eventual consistency, idempotency, saga patterns) | Expert |
| Identity federation (Entra ID, Okta, PingFederate, custom OIDC) | Strong |
| Hybrid networking (ExpressRoute, Direct Connect, VPN, NAT, FW design) | Strong |
| Cloud platforms (Azure or AWS deep + the other one familiar) | Strong |
| Microsoft integration stack adjacency (BizTalk, Logic Apps, Service Bus, Azure Functions) | Working knowledge |
| Architecture-decision-record (ADR) writing + ARB facilitation | Strong |

### Nice-to-have

- MuleSoft Certified Platform Architect or equivalent (MCPA)
- TOGAF / Open Group certification
- Public-cloud architect cert (AWS SA Pro / Azure Solutions Architect Expert)
- CISSP or equivalent security cert (for citizen-data engagements)

### Day-to-day in this engagement

- Owns architecture decisions and ADRs
- Faces off with client's enterprise architect / ARB
- Final sign-off on policy bundles, identity flows, security posture
- Drives the design docs (01–15 in this repo); keeps them current
- First responder on architecture-level production incidents
- Quarterly architecture review with steering committee

### Phase allocation

| Phase | Allocation |
|---|---|
| 1 — Discovery | 0.6 FTE (deeply involved) |
| 2 — Foundation | 0.4 FTE |
| 3 — First-wave APIs | 0.2 FTE (consulting only, not implementing) |
| 4 — Hypercare | 0.2 FTE |
| Steady-state | 0.1 FTE (sustaining + new-API tier review) |

### Anti-patterns

- **Architect who codes daily**: dilutes architecture quality and creates single-point-of-failure on code; the architect should review, not write production code outside of POCs
- **Architect 100% on one client**: stale perspective; ideally split across 2 engagements so cross-pollination happens

---

## 3. Integration Lead Engineer

### Role profile

**Alternate titles**: Senior MuleSoft Engineer · Tech Lead · API Platform Lead

**Years of experience**: **8+** with **4+ on Anypoint Platform** and **production CI/CD experience**.

### Core skills (must-have)

| Skill | Depth |
|---|---|
| Anypoint Platform hands-on (Omni Gateway install, API Manager config, Exchange publishing) | Expert |
| Anypoint CLI scripting + automation | Expert |
| CI/CD pipelines (GitHub Actions / Azure DevOps / Jenkins / Tekton) | Expert |
| OAuth flows + IdP integration (concrete debugging skills) | Expert |
| Code review at scale (reviewing engineers' policy bundles, OAS specs, IaC) | Strong |
| Linux system administration (RHEL / Ubuntu) | Strong |
| Container fundamentals (Docker; K8s if relevant) | Strong |
| HashiCorp Vault integration | Strong |
| Network troubleshooting at HTTP/TLS level (tcpdump, curl, openssl) | Strong |
| Mentoring + technical interview skills | Strong |

### Nice-to-have

- MuleSoft Certified Developer / Integration Specialist
- Kubernetes administrator cert (CKA) if K8s deployment path
- Redis Certified Developer
- DataWeave proficient (for transformation work in Logic Apps / Functions even though not in Omni)

### Day-to-day in this engagement

- Runs the Phase 2 installation: hands on keyboard for Omni Gateway / Redis / Vault setup
- Builds CI/CD pipelines that engineers then use day-to-day
- Code-reviews every PR from engineers (policy bundles, OAS, transformation code)
- Owns the "is this ready to merge to main" gate
- On-call rotation lead during hypercare
- Mentors mid-level engineers

### Phase allocation

| Phase | Allocation |
|---|---|
| 1 — Discovery | 0.2 FTE (sanity-checks design feasibility) |
| 2 — Foundation | 0.8 FTE (heavy hands-on) |
| 3 — First-wave APIs | 0.6 FTE (reviews + complex APIs) |
| 4 — Hypercare | 0.6 FTE (incident lead) |
| Steady-state | 0.3 FTE (sustaining + escalations) |

### Anti-patterns

- **Lead engineer who won't review** (or rubber-stamps PRs): code quality erodes; mid-level engineers stop growing
- **Lead engineer who tries to write all the production code themselves**: bottleneck; team velocity drops to one person's pace
- **Promoting from "good mid-level engineer" without leadership/mentoring assessment**: technical strength ≠ leadership readiness

---

## 4. Integration Engineer (mid-level, × 2)

### Role profile

**Alternate titles**: MuleSoft Developer · API Engineer · Integration Developer

**Years of experience**: **3–5** with **2+ on Anypoint Platform** or equivalent API platform.

### Core skills (must-have)

| Skill | Depth |
|---|---|
| Anypoint hands-on (API Manager, policy attachment, Exchange publishing) | Strong |
| OpenAPI 3.x authoring | Strong |
| OAuth 2.0 concepts; can implement Client Credentials and Auth Code flows | Strong |
| One backend programming language (Python / Node.js / C# / Java) | Strong |
| Git, branching, PR workflow | Strong |
| Writing/running unit + integration tests | Strong |
| Reading network captures (Wireshark / browser dev tools) | Working |
| Linux command line | Working |

### Nice-to-have

- MuleSoft Certified Developer (associate or platform)
- Azure Fundamentals or AWS Cloud Practitioner
- Postman / contract-testing experience
- Some experience with Azure Functions / Logic Apps (for orchestrator work)

### Day-to-day in this engagement

- Builds API onboarding work (Tier 1 / Tier 2 APIs are the engineer's bread and butter)
- Writes OAS specs in collaboration with client API teams
- Configures policy bundles
- Writes tests (unit + contract + smoke)
- Triages bugs in steady-state
- Light on-call rotation (Sev 2/3 first responder during hypercare; Sev 1 escalates to lead)

### Phase allocation

| Phase | Allocation per engineer |
|---|---|
| 1 — Discovery | 0 FTE (not engaged yet) |
| 2 — Foundation | 0.5 FTE (shadowing + assisting lead) |
| 3 — First-wave APIs | 1.0 FTE (primary onboarding hands) |
| 4 — Hypercare | 0.6 FTE |
| Steady-state | 0.3 FTE (sustaining + ongoing API onboarding) |

**Two engineers typical** — gives parallelism on API onboarding and resilience to absence.

### Anti-patterns

- **Engineer with no test discipline**: produces APIs that pass dev but fail in QA; lead engineer's review burden balloons
- **Engineer assigned a Tier 3 API solo**: under-equipped; should pair with lead engineer for complex APIs
- **Engineer not rotating through different APIs**: becomes the only person who knows that API; bus-factor risk

---

## 5. DevOps / SRE Engineer

### Role profile

**Alternate titles**: Platform Engineer · SRE · Infrastructure Engineer

**Years of experience**: **5+** with **3+ years operating production platforms** (any major stack).

### Core skills (must-have)

| Skill | Depth |
|---|---|
| Linux administration (RHEL / Ubuntu / CoreOS at production scale) | Expert |
| IaC (Terraform / Ansible / both) | Strong |
| CI/CD (GitHub Actions / Azure DevOps / Jenkins) | Strong |
| Observability stack (Splunk / Datadog / Prometheus / Grafana / OpenTelemetry) | Strong |
| Networking (TCP/IP, DNS, TLS, FW rules, load balancers) | Strong |
| Containers (Docker; K8s if relevant) | Strong |
| HashiCorp Vault administration | Working |
| Incident response (PagerDuty workflows, post-incident reviews) | Strong |
| Scripting (Bash + Python or PowerShell) | Strong |

### Nice-to-have

- HashiCorp Certified Vault Operations Professional
- Kubernetes CKA / CKS
- Splunk / Datadog admin certs
- F5 BIG-IP administration experience (LTM / GTM) — high value for the on-prem deployment

### Day-to-day in this engagement

- Owns infrastructure provisioning and configuration
- Builds and maintains the observability stack
- Owns Redis Sentinel cluster operations
- Patches OS, Anypoint runtime, Redis
- Manages on-call rotation tooling (PagerDuty / Opsgenie)
- Owns DR drills (quarterly execution)
- Authors runbooks

### Phase allocation

| Phase | Allocation |
|---|---|
| 1 — Discovery | 0.1 FTE (network + infra design input) |
| 2 — Foundation | 0.6 FTE (heavy hands-on) |
| 3 — First-wave APIs | 0.2 FTE (supports environment promotion) |
| 4 — Hypercare | 0.4 FTE (on-call + incident response) |
| Steady-state | 0.2 FTE (ongoing patching + monitoring + drills) |

### Anti-patterns

- **DevOps engineer who treats networking as someone else's problem**: leaves the platform with broken FW rules / DNS / certs that nobody understands
- **DevOps engineer with no on-call experience**: doesn't appreciate why runbooks need to be readable at 3 AM
- **Single DevOps engineer for the whole engagement**: bus factor; client takeover at handover becomes painful

---

## 6. Security Engineer

### Role profile

**Alternate titles**: Application Security Engineer · API Security Engineer · AppSec Consultant

**Years of experience**: **5+** with **3+ years on API / identity security specifically**.

### Core skills (must-have)

| Skill | Depth |
|---|---|
| OAuth 2.0 / OIDC / JWT / SAML (deep) | Expert |
| TLS configuration and debugging (openssl, ciphers, cert chains) | Expert |
| Threat modeling (STRIDE / PASTA) | Strong |
| API security (OWASP API Top 10) | Strong |
| Pen test scope definition + remediation triage | Strong |
| Compliance frameworks (SOC 2, ISO 27001; HIPAA / FedRAMP if in scope) | Strong |
| IdP configuration (Entra / Okta / Ping) | Strong |
| Vault + secret management practices | Strong |
| Static / dynamic application scanning | Working |

### Nice-to-have

- CISSP / CSSLP
- OSCP (offensive security background helps with threat modeling)
- AWS Certified Security – Specialty / Azure SC-100

### Day-to-day in this engagement

- Owns the security posture work in Phase 2 §4.5 of doc 16
- Designs OAuth flows + IdP integration; works with IAM team
- Authors threat models per API and per architectural component
- Coordinates pen test with third-party firm
- Triages pen test findings; assigns remediation
- Owns DPA / BAA paperwork support
- Annual security review lead

### Phase allocation

| Phase | Allocation |
|---|---|
| 1 — Discovery | 0.1 FTE (threat model envisioning) |
| 2 — Foundation | 0.5 FTE (heavy in security posture workstream) |
| 3 — First-wave APIs | 0.2 FTE (per-API security review) |
| 4 — Hypercare | 0.1 FTE (incident response if security-related) |
| Steady-state | 0.1 FTE (ongoing reviews; annual heavy) |

### Anti-patterns

- **Security engineer who blocks but doesn't enable**: just says "no" without offering compliant alternatives; team routes around them
- **Security engineer engaged late in Phase 2**: rework cost; security must be in Phase 1
- **Security engineer who only does docs, not config**: ends up with paper security that doesn't match runtime reality

---

## 7. Technical Project Manager

### Role profile

**Alternate titles**: Delivery Manager · Engagement Manager · Scrum Master (for agile-flavored engagements)

**Years of experience**: **5+** with **technical-platform delivery background** (not generic PM).

### Core skills (must-have)

| Skill | Depth |
|---|---|
| Multi-stakeholder coordination (client + consultant + vendor) | Expert |
| Reading enough technical content to ask the right clarifying questions | Strong |
| Risk management + escalation discipline | Strong |
| Standard PM tooling (Jira / Azure Boards / similar) | Strong |
| Agile + Waterfall hybrid (most enterprise engagements are hybrid) | Strong |
| Steering committee facilitation | Strong |
| Burn-rate tracking + financial reporting back to consulting firm | Strong |

### Nice-to-have

- PMP / Prince2 / Scrum Master cert
- Background as a developer or analyst (not just career PM)

### Day-to-day in this engagement

- Single point of contact for client PM
- Maintains the project plan; updates weekly
- Runs weekly steering committee
- Tracks burn rate and risks
- Coordinates change orders when assumptions break
- Owns the engagement P&L (works with consulting firm leadership)
- Facilitates Phase gate reviews

### Phase allocation

| Phase | Allocation |
|---|---|
| 1 — Discovery | 0.2 FTE |
| 2 — Foundation | 0.3 FTE |
| 3 — First-wave APIs | 0.3 FTE |
| 4 — Hypercare | 0.2 FTE |
| Steady-state | 0.1 FTE |

### Anti-patterns

- **PM with no technical background**: can't ask the right questions; becomes a status-meeting scheduler
- **PM running multiple engagements simultaneously**: time-sliced attention misses risks
- **PM whose primary job is shielding the architect from the client**: relationship breakdown waiting to happen

---

## 8. Technical Writer (optional but recommended)

### Role profile

**Alternate titles**: Documentation Engineer · Content Designer

**Years of experience**: **3+** with **technical writing for API / platform products**.

### Core skills

| Skill | Depth |
|---|---|
| Writing API documentation (OpenAPI, developer portals) | Strong |
| Writing operational runbooks | Strong |
| Markdown + diagram authoring (Mermaid / draw.io) | Strong |
| Version-controlled docs (Git workflow) | Strong |
| Working from engineer interviews + reading code | Strong |

### Day-to-day in this engagement

- Authors runbooks during Phase 2 (in collaboration with DevOps + Lead Eng)
- Writes API onboarding documentation as part of each Tier 1–3 onboarding
- Maintains the developer portal content
- Polishes architect's design docs for external readability

### Why "optional"

In smaller engagements the architect + lead engineer write docs themselves. In larger engagements a dedicated writer dramatically improves the quality of the deliverable docs at relatively low cost. Recommend including unless the client explicitly declines.

### Phase allocation

| Phase | Allocation |
|---|---|
| 2 — Foundation | 0.2 FTE (runbook authoring) |
| 3 — First-wave APIs | 0.1 FTE (per-API docs) |
| 4 — Hypercare | 0.2 FTE (operational handover docs) |
| Steady-state | 0.05 FTE (docs maintenance) |

---

## 9. Team-shape decision rules

| If… | Adjust team like this |
|---|---|
| Engagement is purely greenfield | Lower Architect (no current state to assess) |
| Heavy EDI / B2B scope | Add Logic Apps EIP specialist (~0.3 FTE Phase 2/3) |
| Multiple regions (>1 Private Space) | Add DevOps capacity (~+0.2 FTE) |
| Custom WASM policies required | Add Rust/Go developer (~0.3 FTE during policy development) |
| Client team is mature ops org | Lower DevOps allocation (~-0.2 FTE) |
| Client team has no MuleSoft experience | Add Engineer capacity for KT (~+0.3 FTE Phase 4) |
| 24/7 on-call required post-hypercare | Add on-call rotation (~+0.4 FTE distributed across engineers + DevOps) |
| Compliance scope = HIPAA / FedRAMP | Security Engineer becomes ~0.4 FTE Year 1 instead of 0.2 |
| Tier 3 APIs > 30% of first wave | Add Engineer (~+0.5 FTE Phase 3) |

---

## 10. Skill-progression / career path

For consulting-firm operations: this engagement is a **growth platform** for staff. Realistic progression:

| From → To | Years | Trigger |
|---|---|---|
| Junior Engineer → Integration Engineer | 1–2 | Independently delivers a Tier 1 API end-to-end |
| Integration Engineer → Integration Lead | 3–4 | Owns code reviews; mentors 2+ engineers; runs a CI/CD pipeline build |
| Integration Lead → Solution Architect | 4–6 | Owns architecture decisions across multiple engagements; faces off with client architects independently |
| Engineer → DevOps Engineer | varies | Strong infra / scripting interest + production on-call experience |
| Engineer → Security Engineer | varies | Deep interest in OAuth / threat modeling + cert pursuit |

This shapes who you put on the engagement: **don't staff all senior**, and **don't staff all junior**. The mix is the team.

---

## 11. Anti-patterns at the team-composition level

| Anti-pattern | Why bad |
|---|---|
| Architect-only team in Phase 2 | Implementation slows; architect bottlenecked on hands-on work |
| All-senior team | Margin disaster; over-engineered solution; client overpays |
| All-junior team | Quality + judgment issues; client unhappy; senior burnout from constant rework |
| No PM | Status reporting + risk tracking falls through; client perception suffers |
| Single individual in any role with no backup | Bus factor = 1; vacation / illness = engagement stops |
| Onshore-only or offshore-only with no rationale | Should be cost / time-zone / skill optimized, not dogmatic |
| Team that never sees the client face-to-face | Loses architectural empathy with client context; design becomes academic |

---

## Related

- [16 — Implementation & Operations Estimate](16-consulting-estimate.md) — the hours these roles are sized to deliver
- [17 — Engagement Models](17-engagement-models.md) — the commercial models these roles staff against
- [19 — Pricing Model](19-pricing-model.md) — the rate card structure these roles map to
- [20 — RACI Matrix](20-raci-matrix.md) — what each role is responsible / accountable for per phase
