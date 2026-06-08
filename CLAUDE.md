# MulesoftArchitecture — Working Notes for Claude

Reference architecture for an on-prem Anypoint **Omni Gateway** (formerly Flex Gateway) deployment. Documentation-heavy project — no application code yet.

## Read first

- [`README.md`](README.md) — full doc index
- [`docs/01-api-gateway-architecture.md`](docs/01-api-gateway-architecture.md) — overall architecture
- [`docs/16-consulting-estimate.md`](docs/16-consulting-estimate.md) — effort estimate (hours only)
- [`docs/20-raci-matrix.md`](docs/20-raci-matrix.md) — who's R/A/C/I across phases

## Target environment

- **5M calls/day Prod** (NOT the legacy 100K/day baseline — that was the original assumption, updated 2026-06-06; all docs + decks reflect 5M)
- **4 environments**: DEV / QA / Acceptance / Prod
- **QA mirrors Prod** for accurate load testing (doubles gateway compute spend; "ephemeral QA" called out as the alternative)
- **Citizen-data workload** — PII handling discipline mandatory; see [`docs/07-data-protection.md`](docs/07-data-protection.md)
- **Omni Gateway in Connected mode** — **Redis is REQUIRED** infrastructure (not optional; corrected in [`docs/10-redis-cache.md`](docs/10-redis-cache.md) after fact-check against MuleSoft docs)

## Naming convention

- **"Omni Gateway"** is the current MuleSoft product name; **"Flex Gateway"** is the same product (pre-rebrand)
- Both names appear in docs 01–09 (written before the discovery); docs 10+ use "Omni Gateway"
- They are interchangeable — don't try to disambiguate them as different products

## Diagrams

- **Mermaid** for inline doc diagrams (renders on GitHub natively)
- **Python `diagrams` lib** for PNG architecture images — sources in `presentations/diagrams/`, PNGs in `presentations/img/`
- **Python `python-pptx`** for the 2 PowerPoint decks in `presentations/`
- Regenerate diagrams: `python presentations/diagrams/<file>.py` (requires Graphviz + diagrams pip pkg)
- Regenerate decks: `python presentations/build_pptx.py` and `python presentations/build_onprem_pptx.py`
- Regenerate pricing template: `python pricing/build_pricing_template.py`

## Doc conventions

- Numbered docs (`01-`...`20-`); each ~1000–2000 words
- Sequence diagrams for flows; component-responsibility tables over prose; honest framing of anti-patterns
- Cross-link liberally between docs

## Cross-references between sizing/cost docs

- `docs/09 §6.5` — canonical per-env sizing matrix
- `docs/01 §6` — mirrors the sizing for the SaaS path
- `docs/05 §12` — observability cost across all 4 envs at 5M/day
- `docs/16` — hours estimate
- `docs/19` — pricing model
- `pricing/pricing-template.xlsx` — fill-in workbook

## Commercial pack (docs 16–20 + `pricing/`)

- **Engagement model**: hybrid per phase (T&M cap → Fixed-Bid → tiered T&M per-API → Fixed-Bid hypercare → retainer)
- **Per-API tiers**: 50 / 100 / 200 hrs (Simple / Standard / Complex)
- **Pricing template** uses openpyxl-generated xlsx with formulas wired across 6 sheets

## Git / push conventions

- Branch: `main`
- Push via `gh` CLI (already auth'd; uses `gho_` browser-issued token in OS keyring)
- **Never** use the PAT pasted in chat on 2026-05-18 — it's revoked
- Commit messages: full explanatory body (multi-line); always include rationale not just the what

## What this project does NOT have yet

- No Terraform / Ansible / IaC code
- No working CI/CD pipelines (only documented in `docs/04`)
- No actual deployed environment
- No reference implementation code (only conceptual)

This is **architecture documentation + commercial pack**, designed to be the basis for a future implementation engagement. The deliverables are SOW-ready, not Day-1-deployable.
