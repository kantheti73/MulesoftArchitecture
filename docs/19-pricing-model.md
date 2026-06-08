# 19 — Pricing Model

How to translate the hours from [doc 16](16-consulting-estimate.md) into priced engagements using the role mix from [doc 18](18-team-roles-and-skills.md) and the commercial structure from [doc 17](17-engagement-models.md).

> **This doc establishes the structure**. Actual rates are intentionally left as **placeholders** in the accompanying spreadsheet template — fill in based on your firm's rate card and the engagement's commercial terms before sharing externally.

---

## 1. What this doc does (and doesn't)

| In scope | Out of scope |
|---|---|
| Rate-card structure (role × seniority × geography) | Specific dollar/euro/GBP amounts |
| Calculation method (hours × blended rate) | Your firm's internal cost model |
| Discount strategy (volume / multi-year / retainer) | Specific discount percentages |
| Margin guidance (% built in per role) | Your firm's margin policy |
| Risk pricing (markup for uncertain phases) | Insurance / liability pricing |
| Travel / expenses model | Per-diem rates for your geography |
| T&M ceiling structures | Tax / regulatory adjustments |
| SLA-based credits + incentives | Specific credit amounts |

The companion file [`pricing/pricing-template.xlsx`](../pricing/pricing-template.xlsx) is the working spreadsheet. The formulas are wired; you fill in the rates and hours per engagement.

---

## 2. Rate-card structure

The rate card has three dimensions:

| Dimension | Typical values |
|---|---|
| **Role** | Solution Architect · Integration Lead · Integration Engineer · DevOps Engineer · Security Engineer · Project Manager · Technical Writer |
| **Seniority within role** | Principal / Senior / Mid / Junior |
| **Geography** | Onshore (US/UK/EU) · Nearshore (e.g. EU-Eastern / LATAM for US) · Offshore (India / Philippines / etc.) |

That's **7 roles × 4 seniorities × 3 geographies = 84 cells maximum**, though most engagements use a subset. The spreadsheet ships with the full grid; fill in only the cells you'll actually use.

### Why three geographies

Mixed-shore delivery is common in consulting. Onshore architects + nearshore leads + offshore engineers is a typical cost-optimized mix. The rate card has to express that. Even pure-onshore engagements benefit from having the structure visible.

---

## 3. Calculation method

For each phase:

```
Phase cost = Σ (role × seniority × geography × hours × rate)
```

The spreadsheet does this with named ranges and SUMPRODUCT — you change a rate or an hour count and totals update everywhere.

### Cell layout

```
Sheet 1: Rate Card
  rows = role × seniority
  cols = geography
  cells = $rate/hr

Sheet 2: Hours by Phase
  rows = phase
  cols = role × seniority × geography (or a flattened view)
  cells = hours

Sheet 3: Calculated Cost
  pulls Hours × Rate Card → priced
  totals per phase + grand total

Sheet 4: Discounts & Adjustments
  volume discount tiers
  multi-year discount
  retainer discount

Sheet 5: Risk + Margin
  risk markup per phase
  margin % per role

Sheet 6: Summary (client-facing)
  cleaned-up totals per phase
  no internal cost data shown
```

---

## 4. Discount strategy (typical structures)

| Discount type | Trigger | Typical depth | Purpose |
|---|---|---|---|
| **Volume / multi-year** | Total engagement value > threshold | 5–15% off list | Reward commitment; secure ongoing book of business |
| **Retainer commitment** | 12-month retainer signed | 8–12% off T&M rate for retainer hours | Reward predictability; reduces sales overhead |
| **Multi-engagement / framework** | Master agreement across multiple projects | 10–20% off | Strategic relationship pricing |
| **Pre-payment** | Quarterly or annual pre-payment | 3–8% off | Cash-flow benefit to consultant |
| **Reference / case-study** | Client agrees to be a case study | 2–5% off | Marketing value |
| **Beta / first-of-its-kind** | First engagement using new methodology | 5–15% off | Risk-sharing for novel work |

**Critical**: discounts should be **named** in the SOW (e.g., "12-month retainer commitment discount: 10%"). Unnamed discounts get forgotten and re-litigated at renewal.

---

## 5. Margin guidance (% to build in)

Different roles carry different margin profiles:

| Role | Typical margin band | Notes |
|---|---|---|
| Solution Architect | High (40–60%) | Scarce skill; high client-value perception |
| Integration Lead | Mid-high (35–50%) | Strong technical depth; team multiplier |
| Integration Engineer | Mid (25–40%) | More commoditized; competitive market |
| DevOps Engineer | Mid-high (35–50%) | Scarce skill set; production-critical |
| Security Engineer | High (40–55%) | Specialized; risk-laden role |
| Project Manager | Mid (25–40%) | Easily compared; price-sensitive |
| Technical Writer | Lower (20–35%) | Easy to find; often de-prioritized |

These bands are illustrative — your firm's policy may differ. The spreadsheet has a margin column per role you can set.

---

## 6. Risk pricing (markup for uncertain phases)

Different phases carry different execution risk. Add a risk markup:

| Phase | Risk markup typical |
|---|---|
| Phase 1 — Discovery (T&M with cap) | 0–5% (T&M passes risk to client; minimal markup) |
| Phase 2 — Foundation (Fixed-Bid) | 10–20% (consultant bears scope risk) |
| Phase 3 — First-wave APIs (T&M tiered) | 5–10% (per-API tier shares risk) |
| Phase 4 — Hypercare (Fixed-Bid) | 8–15% (consultant manages incident burn) |
| Steady-state (Retainer) | 0–5% (predictable run rate) |

Risk markup is applied **to fixed-bid phases especially**. It's not separately itemized in the client price — it's embedded in the rate or hour count.

### When to raise risk markup

| Trigger | Markup adjustment |
|---|---|
| New client, no prior relationship | +5% |
| Aggressive timeline (client wants 4 months instead of 6) | +5–10% |
| Multiple novel components (e.g. first time using Omni Gateway + first time on this client's network) | +5–10% |
| Compliance scope expansion (HIPAA / FedRAMP) | +10–20% |
| Geographic / time-zone challenges in team | +5% |
| Known difficult client (history of scope changes) | +10–15% |
| Strong prior relationship + repeat business | -5% (give them a "loyalty" discount on the markup) |

---

## 7. Travel & expenses

Two common models:

| Model | How it works | When to use |
|---|---|---|
| **T&E pass-through** | Actual travel + lodging + per-diem billed at cost | Default for most engagements |
| **T&E baked-in** | Travel cost estimated and included in fixed-bid; consultant manages within budget | Short, well-defined engagements; clients who hate expense reports |

For US engagements, typical T&E:
- Flights, hotels, ground transport at actual
- Per-diem at GSA rates (or client's policy if stricter)
- 1 trip per phase milestone is typical (kickoff, mid-phase review, close-out)
- Steady-state: quarterly visits typical

Build into the SOW: **explicit travel budget** OR **monthly cap on billable T&E**.

---

## 8. T&M ceiling structures

For T&M phases (Phase 1, Phase 3), the SOW needs cap mechanics:

| Mechanism | How it works |
|---|---|
| **Hard cap** | Above the cap, consultant works free OR work stops | Cleanest for client; consultant carries overage risk |
| **Soft cap with renegotiation trigger** | At cap, work pauses; renegotiation happens before continuing | Most realistic — neither side wants surprise stoppage |
| **Phased caps** | Each sub-phase has its own cap | Best for long T&M engagements |
| **Burn-rate alerts** | At 50% / 75% / 90% of cap, formal notice + status review | Required for any cap structure |

**Recommended language for SOW**: "Soft cap at 200 hours. At 150 hours, parties will meet to review burn rate and either confirm cap or renegotiate. Work continues only after explicit written go-ahead."

---

## 9. SLA-based credits & incentives

For retainer + steady-state operations, you can structure incentives:

| Type | How it works | Effect |
|---|---|---|
| **Penalty credit** | Missed SLA (e.g. Sev 1 response > 30 min) = N% credit toward next month | Aligns consultant on operational quality |
| **Performance bonus** | Hit all SLAs for 12 months = bonus payment | Costly to client; rare; only for premium engagements |
| **Innovation budget** | N hours/year of consultant time for client-chosen innovation work | Builds relationship; consultant gets to flex |
| **Knowledge-transfer credits** | N hours/year reserved for client team training | Pre-paid; client can spend or lose |

Most engagements use penalty-credit (downside) without a bonus (upside). That's market norm.

---

## 10. Sample priced scenarios

Three plausible engagement profiles using the structure. **Hour estimates only**; rates are placeholders in the spreadsheet.

### Scenario A — Conservative (cost-sensitive client)

- Mixed-shore team (onshore architect + nearshore lead + offshore engineers)
- Lower per-API tier mix (5 Tier 1 + 4 Tier 2 + 1 Tier 3 in first wave)
- 12-month retainer commitment → discount applied
- Total Year 1: ~3,600 hrs

### Scenario B — Standard

- Onshore architect + onshore lead + nearshore engineers
- Even tier mix (3 Tier 1 + 5 Tier 2 + 2 Tier 3)
- 12-month retainer + pre-payment quarterly → discount applied
- Total Year 1: ~4,050 hrs

### Scenario C — Premium (large enterprise, compliance-heavy)

- All-onshore team
- Heavy Tier 3 mix (1 Tier 1 + 4 Tier 2 + 5 Tier 3)
- HIPAA scope → additional security hours + risk markup
- 24/7 on-call premium
- Total Year 1: ~5,200 hrs

Use the spreadsheet to model your own scenarios.

---

## 11. Common pricing anti-patterns

| Anti-pattern | Why bad |
|---|---|
| Pricing without role mix specified | Hours look fine, but team composition is over- or under-senior; margin or quality suffers |
| Single rate for the whole team | Either undersells architect time or overcharges for engineer time; opaque to client |
| No risk markup on fixed-bid phases | First scope surprise eats margin |
| Bundled discount with no named driver | At renewal, client expects discount to continue without remembering why |
| Travel as a vague "estimated $X/month" line | Either over-billed (client annoyed) or under-billed (consultant absorbing) — be explicit |
| T&M with no cap | Client gets nervous and demands weekly justification; relationship strain |
| Retainer with no overage definition | Either consultant absorbs spikes (margin hit) or surprise bill (relationship hit) |
| Discount applied as a percentage with no expiration | Discount becomes the new list price within 12 months |
| Margin opaque per role (one blended margin number) | Can't optimize; can't explain price moves |
| Pricing locked for multi-year with no annual review clause | Inflation + skill-market shifts erode margin silently |

---

## 12. SOW pricing language patterns

Examples of pricing language that holds up:

### T&M with cap

> "Consultant will deliver Phase 1 — Discovery & Design on a Time & Materials basis at the rates set forth in Exhibit A. Total fees will not exceed two hundred (200) hours without written agreement of both parties. At one hundred fifty (150) hours of consumption, parties shall meet to review status, scope, and remaining work. Either party may request renegotiation of the cap at that meeting."

### Fixed-bid with change-order trigger

> "Consultant will deliver Phase 2 — Foundation Setup for a fixed fee of [amount], based on the scope, assumptions, and exclusions in Sections X, Y, Z of this Statement of Work. Any of the following events shall trigger a Change Order in accordance with Section [N] of the Master Services Agreement: [list of assumption violations]."

### Per-API tiered

> "Each new API onboarded under Phase 3 will be classified jointly into one of three tiers (defined in Exhibit B). The fee for each API shall be the tier rate published in Exhibit A. Tier classification shall be agreed in writing prior to commencement of work on that API; reclassification mid-build shall be a Change Order."

### Retainer with rollover

> "Consultant will provide ongoing operations support on a monthly retainer basis equal to one hundred thirty (130) hours per calendar month. Hours not consumed in a given month shall roll forward for up to two (2) successive months; hours not consumed within three (3) months of accrual shall expire. Consumption exceeding the monthly retainer shall be billed at the agreed Time & Materials rate."

### SLA penalty credit

> "Failure to meet the Severity 1 response SLA of thirty (30) minutes in any calendar month shall result in a credit equal to ten percent (10%) of that month's retainer fee, applied to the following month's invoice."

---

## 13. Operating the spreadsheet template

The spreadsheet [`pricing/pricing-template.xlsx`](../pricing/pricing-template.xlsx) has 6 sheets:

| Sheet | Purpose | What you edit |
|---|---|---|
| `01_RateCard` | Per-role × per-geography hourly rates | All rate cells (`<RATE>` placeholders) |
| `02_HoursByPhase` | Hours allocated per role per phase | Adjust hours per engagement scope |
| `03_CalculatedCost` | Auto-computed cost per phase | Don't edit — formula-driven |
| `04_Discounts` | Discount tiers + applied discount % | Set discount % for this engagement |
| `05_RiskMargin` | Risk markup per phase + margin per role | Set per engagement risk posture |
| `06_Summary` | Clean client-facing summary | Don't edit — pulls from other sheets |

**Workflow**:
1. Open the template
2. Fill in `01_RateCard` for the geographies / roles you'll use
3. Adjust `02_HoursByPhase` if the engagement deviates from the standard estimate in doc 16
4. Set discounts in `04_Discounts`
5. Set risk markup in `05_RiskMargin`
6. The Summary sheet shows the final number — share that with the client; **don't share the cost/margin sheets externally**

---

## 14. Internal-only vs client-facing

| Sheet | Share with client? | Why |
|---|---|---|
| `01_RateCard` | Sometimes (engagement-specific rate exhibit) | OK if your firm publishes rate cards |
| `02_HoursByPhase` | **Optional** — useful for transparency | Some firms share, some don't |
| `03_CalculatedCost` | **Internal only** | Reveals internal cost build-up |
| `04_Discounts` | **Internal only** | Reveals discount strategy |
| `05_RiskMargin` | **Internal only** | Reveals margin policy |
| `06_Summary` | **Yes** | Client-facing total + phase breakdown |

The template ships with a clear color code (yellow = edit; gray = formula; red = internal-only). Respect the coloring when sharing.

---

## Related

- [16 — Implementation & Operations Estimate](16-consulting-estimate.md) — the hours this pricing applies to
- [17 — Engagement Models](17-engagement-models.md) — the commercial structures this pricing is shaped to
- [18 — Team Roles & Skills](18-team-roles-and-skills.md) — the roles this rate card prices
- [20 — RACI Matrix](20-raci-matrix.md) — what each role does for the price quoted
- [`pricing/pricing-template.xlsx`](../pricing/pricing-template.xlsx) — the working spreadsheet
