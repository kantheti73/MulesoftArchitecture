# 17 — Engagement Models per Phase

How to structure the commercial relationship for each phase of the work described in [doc 16](16-consulting-estimate.md). The wrong engagement model on a phase produces predictable failure modes: T&M on a vague Discovery becomes endless; Fixed-Bid on first-wave APIs becomes a margin-eroding scope fight; per-API tiering on Hypercare creates the wrong incentives.

The recommendation: **hybrid** — different model per phase, matched to that phase's uncertainty profile.

---

## 1. The three engagement models

| Model | Who carries scope risk | Who carries execution risk | When it fits |
|---|---|---|---|
| **Time & Materials (T&M)** | **Client** | Consultant | Scope is genuinely uncertain at signing |
| **Fixed-Bid** | Consultant | Consultant | Scope is well-understood and stable |
| **Retainer / Managed Service** | Shared (by SLA) | Consultant | Steady-state operations with predictable shape |

A fourth structure — **T&M with cap** (a.k.a. "not-to-exceed") — is sometimes treated as a separate model. In practice it's T&M with a ceiling, and behaves as T&M from a risk standpoint up to the cap, then like fixed-bid above it. Useful for short well-bounded sub-phases.

---

## 2. Recommended model per phase

| Phase | Recommended model | Why |
|---|---|---|
| **Phase 1 — Discovery & Design** | **T&M with cap** (e.g., 200 hrs cap) | Real discovery has unknown unknowns; pricing risk on client during exploration, but capped so commitment is bounded |
| **Phase 2 — Foundation Setup** | **Fixed-Bid** | Scope is well-understood (install, config, CI/CD, security baseline). Consultant is the expert; clear deliverables make this fair fixed-bid territory. |
| **Phase 3 — First-Wave APIs** | **T&M per-API at tiered rates** | Each API is a discrete unit; tier picker from [doc 16 §7](16-consulting-estimate.md#7-per-api-onboarding-model-the-unit-economics) sets the price per unit. Client knows the unit price upfront; consultant doesn't carry "the 8th API was secretly Tier 3" risk |
| **Phase 4 — Hypercare** | **Fixed-Bid** (8-week window) | Time-boxed, well-understood. Consultant manages hour allocation; client gets predictable cost |
| **Steady-state Operations** | **Monthly Retainer** | Predictable run rate; covers ~130 hrs/month baseline. New APIs billed separately at tier rates |
| **Annual Recurring** (audits, drills, upgrades) | **Quoted per-event** | These are scheduled but each has bounded scope; quote when triggered |

---

## 3. Why per-phase mixing beats a single model

### "Just do everything T&M"

- **Pro**: zero consultant risk; flexibility to handle whatever emerges
- **Con**: client sees an open-ended commitment; budget approvals stall; CFO pushback
- **Real outcome**: client demands weekly burn reports and aggressive cap negotiations; relationship becomes adversarial on every Friday timesheet review

### "Just do everything Fixed-Bid"

- **Pro**: client gets cost certainty; easy approval
- **Con**: consultant carries all scope risk; either pads heavily (lose to lower bidder) or under-prices (lose margin during execution)
- **Real outcome**: phase 1 ambiguity causes consultant to spend $X discovering scope that wasn't quoted; change orders multiply; client perceives nickel-and-diming

### "Just do retainer"

- **Pro**: simplest commercial relationship; ongoing predictable revenue
- **Con**: front-loaded delivery work doesn't fit a steady run rate; either consultant under-resources Phase 2 (delay) or client over-pays during steady state
- **Real outcome**: project under-resources or over-staffs; either way someone loses

### The hybrid

Each phase gets the model that aligns the incentives correctly for that phase. Client carries scope risk where scope is genuinely unknown; consultant carries it where they're the expert. Per-API tiering aligns incentives around the actual unit of work.

---

## 4. Detailed model-per-phase rationale

### 4.1 Phase 1 — T&M with cap

**Capped at**: 200 hrs (slightly above the ~180 hrs estimate; small buffer)

**Why T&M**:
- Discovery surfaces unknowns by design
- Workshop facilitation expands or contracts based on stakeholder availability
- Architecture iterations depend on ARB feedback you can't predict

**Why capped**:
- Bounds client commitment
- Forces consultant to ration time wisely
- Creates trigger point for renegotiation if cap is approached (vs silent burn)

**What to include in the SOW**:
- Cap hours explicitly
- Burn-rate reporting weekly
- Re-baselining process at 75% of cap (have the conversation early, not at the cap)
- Hours-not-used **do** carry forward into Phase 2 (optional but builds goodwill)

### 4.2 Phase 2 — Fixed-Bid

**Bid amount**: ~960 hrs at agreed blended rate (with consultant's margin embedded)

**Why fixed-bid**:
- Install + config work is well-bounded
- Consultant has done this many times; their estimate is more accurate than client's
- Client gets cost certainty for the capital portion of the project (easier to approve)

**Critical SOW protections for the consultant**:
- **Assumption-based exclusions** (cite [doc 16 §11](16-consulting-estimate.md#11-estimating-assumptions) — "client provides X; if not, change order")
- **Change-order triggers** explicitly defined (assumption violation → change order; not "scope discussion")
- **Delivery acceptance criteria** explicit (what does "done" look like for each workstream)
- **Hold-time tracking** (if client team blocks consultant, clock stops on the consultant side; or hold time gets billed separately at agreed rate)

### 4.3 Phase 3 — T&M per-API at tiered rates

**Pricing structure**: per-API at one of three published tier rates ([doc 16 §7](16-consulting-estimate.md#7-per-api-onboarding-model-the-unit-economics))

**Why T&M per unit**:
- API count is known; complexity per API is variable
- Tier rates are published; client sees unit economics upfront
- Tier classification done jointly at intake (no surprise reclassifications mid-build)
- Consultant doesn't carry "Tier 3 surprise" risk; client doesn't carry "everything is suddenly Tier 3" risk

**Tier classification process** (write into the SOW):
1. Client submits API intake form (existing OAS spec or functional brief)
2. Consultant proposes tier within 3 business days
3. Tier mismatch escalation: joint architect review, decision within 2 days
4. Locked tier = locked price; if scope changes mid-build, change order applies

### 4.4 Phase 4 — Fixed-Bid

**Bid amount**: ~310 hrs over 8-week window

**Why fixed-bid**:
- Time-boxed (8 weeks is the time-box, not the hours)
- Hours allocation is the consultant's optimization problem
- Client gets clear cost; consultant manages on-call/incident-response economically

**SOW protections**:
- Defined hypercare exit criteria (what does "ready for handover" mean)
- 24/7 vs business-hours coverage explicitly chosen (24/7 is a premium add)
- Severity-defined response SLAs (Sev 1 = N min, Sev 2 = N hrs)
- Hypercare extension trigger if exit criteria not met (with renegotiation, not free extension)

### 4.5 Steady-state — Monthly retainer

**Retainer size**: ~130 hrs/month

**Why retainer**:
- Run rate is predictable from ops profile
- Aligns with how client procurement budgets ongoing services
- Creates a stable, ongoing relationship

**Critical retainer mechanics**:
- **Unused hours**: roll over up to N months (3 is common), then expire (avoids consultant becoming bank)
- **Overage**: billed at agreed T&M rate per hour (don't absorb silently)
- **New-API onboarding**: explicitly OUT of retainer; tier rates per Phase 3 apply
- **Annual recurring activities**: explicitly OUT of retainer; quoted per event
- **Quarterly review**: retainer size revisited based on actual run-rate trends (right-sizing both directions)

### 4.6 Annual recurring — Quoted per-event

| Event | Approach |
|---|---|
| Annual security review / pen test | Quote 8 weeks before; fixed-bid based on prior years + delta |
| Anypoint major upgrade | Quote when MuleSoft announces target version; fixed-bid for the upgrade window |
| Architecture review board annual walkthrough | Fixed-bid (well-bounded scope) |
| Capacity planning review | Fixed-bid (analytical) |
| Full-scale DR drill | Fixed-bid (well-bounded scope) |
| Compliance attestation support | T&M with cap (auditor-driven scope) |

---

## 5. Risk allocation matrix (which model carries which risk)

| Risk | T&M | Fixed-Bid | Retainer |
|---|---|---|---|
| Scope grows | Client absorbs | Consultant absorbs (or change order) | Up to retainer cap, consultant absorbs |
| Scope shrinks | Client benefits (lower spend) | Consultant retains (margin opportunity) | Hours roll over per retainer terms |
| Schedule slips for client reasons | Client absorbs hold time | Consultant typically absorbs unless explicit hold-time clause | Consultant absorbs |
| Schedule slips for consultant reasons | Consultant absorbs | Consultant absorbs heavily (margin erosion) | Consultant absorbs |
| Technology surprise (e.g. MuleSoft breaking change) | Client absorbs additional hours | **Change-order trigger** | Within reason, retainer absorbs; major changes = change order |
| Resource turnover on consultant side | Consultant absorbs (rebadging cost) | Consultant absorbs | Consultant absorbs |
| Pen test surfaces architectural changes | Client absorbs new hours | Change-order trigger | Out of retainer scope; quoted separately |

---

## 6. Pricing-risk dials inside each model

Mechanisms to fine-tune the risk balance within a chosen model:

### Within a T&M engagement
- **Cap** (a.k.a. not-to-exceed): converts to fixed-bid above cap
- **Floor** (minimum commitment): ensures consultant doesn't staff and then lose hours
- **Burn-rate alerts**: trigger conversation at 50/75/90% of cap
- **Phased caps**: separate cap per sub-phase

### Within a Fixed-Bid engagement
- **Assumption-based exclusions** (cite [doc 16 §11](16-consulting-estimate.md#11-estimating-assumptions))
- **Change-order triggers**: explicit list of conditions that trigger a change order (assumption violation, scope addition, schedule compression)
- **Acceptance criteria**: define "done" so completion is unambiguous
- **Hold-time clause**: client-caused delays accrue billable hold time
- **Payment milestones**: 30% / 30% / 30% / 10% (kickoff / mid / delivery / acceptance)

### Within a Retainer engagement
- **Unused-hour rollover** (capped, e.g. 3 months)
- **Overage billing** (agreed T&M rate)
- **Quarterly right-sizing review**
- **SLA-based credits**: missed response SLA = credit toward next month
- **Exclusion list**: what's NOT in the retainer (new APIs, annual events)

---

## 7. Common engagement anti-patterns

| Anti-pattern | Why bad | What to do instead |
|---|---|---|
| Fixed-bid on Phase 1 Discovery | Consultant has to pad heavily for unknown unknowns OR lose money | T&M with cap |
| T&M on Phase 4 Hypercare without time-box | Becomes "extended hypercare forever"; client unhappy with bill | Fixed-bid for the time-box; explicit exit criteria |
| Per-API tiering on Phase 1–2 work | Architects don't think in "APIs"; foundational work doesn't fit per-API math | Per-API only for Phase 3+ |
| Retainer with no exclusions | New APIs and annual events silently consume retainer hours; consultant under-resourced for ongoing ops | Explicit exclusion list in retainer SOW |
| Retainer with unlimited rollover | Consultant becomes a "credit bank"; client treats hours as savings; consultant resourcing problem | Cap rollover at 2–3 months max |
| Fixed-bid without change-order definition | Every scope question becomes a margin fight | Define change-order triggers explicitly |
| Single model across all phases | Inevitably mis-aligns incentives on at least one phase | Hybrid per [§2](#2-recommended-model-per-phase) |
| Pricing locked at SOW for multi-year engagement | Inflation + tech change erodes margin | Annual rate-card review clause |

---

## 8. SOW structure that supports the hybrid model

A single master SOW with separate work orders per phase:

```
Master Services Agreement (MSA)
  · Governance (steering committee, escalation, dispute resolution)
  · Rate card (referenced by each work order)
  · Change-order process
  · IP / confidentiality / data protection
  · Termination terms

Work Order #1 — Phase 1 (T&M with cap)
  · 200-hour cap
  · Weekly burn-rate reporting
  · Re-baseline at 75% cap

Work Order #2 — Phase 2 (Fixed-Bid)
  · ~960 hrs at blended rate
  · Deliverable acceptance criteria
  · Assumption-based exclusions list
  · Change-order triggers
  · Milestone payment schedule

Work Order #3 — Phase 3 (T&M per-API)
  · Tier rates published
  · Tier classification process
  · Tier lock = price lock

Work Order #4 — Phase 4 (Fixed-Bid for hypercare window)
  · 8-week time-box
  · Exit criteria
  · Severity-based SLAs

Work Order #5 — Steady-state (Monthly retainer)
  · ~130 hrs/month
  · Unused hour rollover (capped)
  · Overage rate
  · Exclusion list

Work Order #6+ — Annual recurring (quoted per event)
```

The MSA establishes the umbrella; work orders carry the per-phase commercial terms. This pattern lets you adjust each phase's model without renegotiating the master agreement.

---

## 9. Recommendation summary

Use **all five models in combination**, one per phase:

| Phase | Model | Risk balance |
|---|---|---|
| 1 — Discovery | T&M with cap | Client bears scope risk to cap; consultant bears it above |
| 2 — Foundation | Fixed-Bid | Consultant bears scope risk (their expertise) |
| 3 — First-wave APIs | T&M per-API tiered | Shared via tier classification |
| 4 — Hypercare | Fixed-Bid | Consultant bears within time-box |
| 5 — Steady-state | Retainer | Predictable; shared via SLA |
| 6 — Annual recurring | Quoted per-event | Quoted as scope materializes |

This is the structure most experienced integration consultancies use because it survives contact with reality.

---

## Related

- [16 — Implementation & Operations Estimate](16-consulting-estimate.md) — the hours these models price
- [18 — Team Roles & Skills](18-team-roles-and-skills.md) — staffing model these engagements deliver against
- [19 — Pricing Model](19-pricing-model.md) — converts hours + models into rate-card pricing
- [20 — RACI Matrix](20-raci-matrix.md) — clarifies who owns what in each phase
