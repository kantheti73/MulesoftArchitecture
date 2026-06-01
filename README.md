# MulesoftArchitecture

Architecture and reference materials for MuleSoft-based integration designs — focused on the **API Gateway role**, with the heavy integration work (orchestration, transformation, EDI) delegated to an on-prem Microsoft stack.

## Docs

| Doc | What it covers |
|---|---|
| [01 — API Gateway Architecture](docs/01-api-gateway-architecture.md) | Scope, product choice (Anypoint Flex Gateway, Connected mode), internal+external dual-listener topology, SaaS vs on-prem trade-offs, private connectivity (Private Space + Transit Gateway + Direct Connect — no public internet), sizing for 100K calls/day |

## Conventions

- **API specs:** OAS 3.x in Anypoint Exchange, published from this Git repo via CI.
- **Naming:** `<domain>-<purpose>-api` (e.g. `orders-public-api`, `inventory-internal-api`).
- **Environments:** `dev`, `stg`, `prd` — separate Anypoint business groups or environments per tenancy.
- **Diagrams:** Mermaid (renders inline on GitHub). PNG renders to follow.
