# 14 — Redis Assumption Register

The Redis design in [doc 10](10-redis-cache.md) rests on a stack of assumptions about traffic, workload shape, MuleSoft behavior, and your environment. This doc enumerates them explicitly so they can be reviewed, challenged, and either confirmed or replaced before they harden into "facts" that no one remembers were originally guesses.

For each assumption: **the assumption · the basis · what breaks if it's wrong · how to verify**.

---

## 1. Working set assumptions

### A1. Total Redis steady-state working set is ~500 MB – 1 GB at 5M/day

| | |
|---|---|
| **Basis** | Rate-limit counters (~150 KB) + runtime config cache (~5–20 MB) + request data cache (~200–500 MB) + idempotency cache (~50–100 MB) |
| **What breaks if wrong** | Memory pressure → eviction of legitimately-needed entries → policy misses → wrong behavior under load |
| **Verify** | Run QA at full Prod traffic for 7 days; sample `INFO memory` every 5 min; measure 95th-percentile working set |
| **Hidden assumption inside this** | Response cache TTLs are short (≤ 60 s); long-lived cached responses would inflate the request-data cache significantly |

### A2. Idempotency cache is ~50–100 MB for 24h window

| | |
|---|---|
| **Basis** | At 5M/day with ~30% writes (idempotent), ~1.5M idempotent keys / day × 24h × ~50 bytes per key |
| **What breaks if wrong** | If 80% of traffic is writes (high-volume event ingestion via [doc 11](11-azure-service-bus-integration.md)), the cache could be 250 MB+ |
| **Verify** | Profile actual request mix per API in QA; tune idempotency window per API tier (24h overkill for many internal APIs) |

### A3. ~1,000 active rate-limit counter keys

| | |
|---|---|
| **Basis** | 4 SLA tiers × ~250 distinct clients |
| **What breaks if wrong** | Onboarding to 5,000+ partners would 5× key count — still trivial in absolute terms but worth knowing |
| **Verify** | Pull partner list from IdP / Anypoint client registry; multiply by SLA tier count |

---

## 2. Performance assumptions

### A4. Each gateway request makes ~2–3 Redis ops

| | |
|---|---|
| **Basis** | Authorization cache check + rate-limit counter increment + (for writes) idempotency check |
| **What breaks if wrong** | If custom policies add more Redis lookups (e.g., per-request enrichment from cache), peak ops could double |
| **Verify** | Profile in QA at sustained load; instrument Redis exporter `commands_processed_total` over wall-clock time |
| **Implied** | Peak Redis ops: ~2,000–2,500/sec at 1,200 TPS gateway spike |

### A5. Redis round-trip latency stays under 2 ms

| | |
|---|---|
| **Basis** | Same DC, same subnet, TLS over Gigabit LAN |
| **What breaks if wrong** | Cross-DC Redis adds 5–10 ms × 2-3 ops = ~30 ms per gateway request — material at 600+ TPS |
| **Verify** | Run `redis-cli --latency` from a gateway replica to its primary; alarm if p99 > 5 ms |
| **Critical:** | This is the assumption most likely to silently degrade — network changes erode latency budgets without anyone noticing until a load test |

### A6. Single Redis master per DC handles the load

| | |
|---|---|
| **Basis** | A single Redis 7.x process handles 100k+ ops/sec for simple GET/SET on modern hardware; our ~2,500 ops/sec is < 3% of capacity |
| **What breaks if wrong** | At 10M+/day with heavier policy chains, you may need Redis Cluster sharding or scale-up |
| **Verify** | CPU on Redis master in QA load test should be < 30%; if higher, profile slow commands |

---

## 3. Topology assumptions

### A7. Sentinel 3-node quorum per DC (not Redis Cluster sharding)

| | |
|---|---|
| **Basis** | Our working set fits comfortably in one node; sharding adds complexity without throughput benefit at this scale |
| **What breaks if wrong** | If working set exceeds Redis node memory in future, need to move to Cluster — different connection logic, different ops model |
| **Verify** | Memory headroom check every quarter; trigger Cluster discussion at 50% memory utilization |

### A8. Per-DC active Redis Sentinel + custom sync daemon (revised 2026-06-09)

| | |
|---|---|
| **Basis** | OSS Redis can't safely do active/active multi-master across DCs natively. Three viable shapes: (a) per-DC scope with SLA wording change, (b) active/passive, or (c) per-DC active + selective custom sync daemon for idempotency keys + rate-limit counter deltas only. We're choosing (c) per [doc 22](22-redis-cross-dc-sync.md) to avoid Redis Enterprise license. |
| **What breaks if wrong** | Sync daemon failure → idempotency keys diverge between DCs (duplicate-write risk on cross-DC retry); rate-limit counter drift up to 2× during partition |
| **Verify** | Sync daemon SLI: 1–2s lag for idempotency keys, 5–10s for counters. Quarterly partition drill validates auto-heal behavior. Document drift bounds in partner SLA. |
| **Hidden assumption inside this** | Eventually consistent semantics are acceptable for our workload. Strong cross-DC consistency would still require Redis Enterprise. |

### A9. Same Redis cluster supports both Connected-mode runtime config cache AND distributed policy state

| | |
|---|---|
| **Basis** | MuleSoft documentation describes a single shared Redis store |
| **What breaks if wrong** | If MuleSoft requires separate stores per concern, infrastructure footprint doubles |
| **Verify** | Confirm with MuleSoft account team — is one Redis enough, or are multiple required? |

---

## 4. Security assumptions

### A10. TLS in transit enforced (port 6380, plaintext 6379 disabled)

| | |
|---|---|
| **Basis** | Citizen-data workload; cached request data is a PII surface (per [doc 07](07-data-protection.md) + [doc 10 §4](10-redis-cache.md#4-securing-the-redis-instance-per-mulesoft-guidance)) |
| **What breaks if wrong** | AUTH password sniffable on the wire; cached request data sniffable |
| **Verify** | Network capture in QA — confirm only port 6380 reachable; 6379 refused |

### A11. Redis AUTH password rotated quarterly via Vault

| | |
|---|---|
| **Basis** | Standard credential lifecycle policy |
| **What breaks if wrong** | Stale credentials accumulating in ops scripts; harder to revoke on personnel changes |
| **Verify** | Vault rotation policy active; gateway service principal can read the secret at startup |

### A12. Dangerous commands (FLUSHALL, CONFIG, KEYS, DEBUG, FLUSHDB) renamed/disabled

| | |
|---|---|
| **Basis** | Single command can wipe shared state across all gateway replicas instantly; MuleSoft security guidance |
| **What breaks if wrong** | Operator typo or compromised credential = production outage |
| **Verify** | `redis-cli ... FLUSHALL` from a gateway host should return `ERR unknown command` |

### A13. Redis nodes are dedicated (not shared with other workloads)

| | |
|---|---|
| **Basis** | Noisy-neighbor isolation; PII blast-radius containment |
| **What breaks if wrong** | Other tenant's blowup → your gateway state goes with it; shared cache means broader compliance attestation scope |
| **Verify** | Inventory check on the Redis VMs — confirm no other services co-located |

### A14. mTLS between gateway and Redis (not just AUTH password)

| | |
|---|---|
| **Basis** | Defense in depth; standard for PII-bearing internal services |
| **What breaks if wrong** | Compromised internal network = AUTH password is the only barrier |
| **Verify** | Cert presented during TLS handshake; `tls-auth-clients yes` in `redis.conf` |
| **Caveat** | If your PKI can't issue Redis-specific certs in time, this becomes "AUTH password + network controls" instead — document the trade |

---

## 5. Network assumptions

### A15. Redis nodes are co-located with gateway in the same DC subnet

| | |
|---|---|
| **Basis** | MuleSoft guidance: "Deploy the Redis as close as possible to the Omni Gateway deployment" |
| **What breaks if wrong** | Cross-subnet FW hops add latency + new failure modes |
| **Verify** | Network team confirms Redis VLAN + gateway VLAN are adjacent; FW logs show 0 east-west drops |

### A16. Per-DC firewall allows gateway SG/CIDR → Redis 6380 only

| | |
|---|---|
| **Basis** | Least-privilege; operator laptops use jumpbox not direct |
| **What breaks if wrong** | Direct access from anywhere creates audit gaps and exfiltration paths |
| **Verify** | FW rule review; confirmed during pen test |

---

## 6. Availability assumptions

### A17. Sentinel auto-failover within 5–10 seconds

| | |
|---|---|
| **Basis** | Default `down-after-milliseconds 5000` + reasonable network latency |
| **What breaks if wrong** | Longer failover = gateway requests in flight see Redis errors → fall back to per-replica counters → brief rate-limit relaxation |
| **Verify** | Chaos drill: kill master in QA, time client recovery from gateway logs |

### A18. `fallbackOnError: local` is enabled in Omni Gateway config

| | |
|---|---|
| **Basis** | Redis blip should not equal production outage |
| **What breaks if wrong** | Redis maintenance window = customer-visible downtime |
| **Verify** | Confirm field name + behavior with MuleSoft account team — Omni Gateway naming has shifted between minor versions |
| **Critical** | **This is the single most important Redis-config switch to get right** |

### A19. It's acceptable for rate-limit counters to drift up to 2× during cross-DC partition (revised 2026-06-09)

| | |
|---|---|
| **Basis** | With per-DC Redis + sync daemon (A8 revised), partitions cause each DC to enforce local limits independently; effective global limit can briefly reach 2× the configured value. On reconnect, counters re-baseline. |
| **What breaks if wrong** | If partner SLAs require strict global limit enforcement at all times (incl. during partition), the custom sync daemon design is insufficient — Redis Enterprise required |
| **Verify** | Legal/contract review of partner SLA wording before go-live; document the bounded-drift semantics explicitly |
| **Note** | This is a softer assumption than the prior "reset during cutover" — drift is bounded and auto-heals, vs the harder break of full counter reset. Most partners accept it once explained. |

---

## 7. Data / PII assumptions

### A20. Request data cache contains transient elements only (truncated, short TTL)

| | |
|---|---|
| **Basis** | MuleSoft's caching is policy-driven; long TTLs are atypical for PII-bearing APIs |
| **What breaks if wrong** | Long-lived cached responses become a data-at-rest exposure surface; doc 07 controls insufficient |
| **Verify** | Audit per-API cache policies; document max TTL per API; alarm on cache TTL > 5 min for any PII-bearing route |

### A21. No persistence enabled (`appendonly no`, no RDB snapshots)

| | |
|---|---|
| **Basis** | State is ephemeral; persistence creates data-at-rest exposure with no recovery benefit |
| **What breaks if wrong** | If MuleSoft requires RDB snapshots for any reason, the disk-encryption + backup-isolation story expands |
| **Verify** | `redis-cli CONFIG GET save` returns empty; `appendonly` is `no` |

### A22. Redis is in the same compliance boundary as the gateway runtime

| | |
|---|---|
| **Basis** | DPA / BAA scope covers cache infrastructure same as compute |
| **What breaks if wrong** | If Redis is in a "shared services" zone with different attestation, you have a compliance gap |
| **Verify** | Compliance / DPO sign-off on Redis attestation scope before Prod go-live |

---

## 8. Operational assumptions

### A23. Redis is a Tier-1 service with the same monitoring as the gateway

| | |
|---|---|
| **Basis** | Redis outage = gateway degradation (or outage, if `fallbackOnError` isn't right) |
| **What breaks if wrong** | Silent Redis problems → cascading gateway issues nobody attributes to Redis |
| **Verify** | Redis on the same on-call rotation, dashboards, alerting routes as gateway |

### A24. Dedicated Redis cluster per environment

| | |
|---|---|
| **Basis** | Test traffic can't pollute Prod counters; QA load tests don't affect Prod |
| **What breaks if wrong** | Shared Redis = QA load test crashes Prod rate limiting |
| **Verify** | Inventory: confirm each env has its own Sentinel cluster |
| **Cost note** | Per-env Redis adds up (see [doc 10 §5 per-env table](10-redis-cache.md#per-environment-redis-sizing-matches-the-4-env-matrix-in-doc-09-65)) |

### A25. Patching cadence aligns with gateway VMs (monthly OS patching, quarterly Redis minor)

| | |
|---|---|
| **Basis** | Standard enterprise change windows |
| **What breaks if wrong** | Drift between Redis nodes → quorum disruption during partial patching |
| **Verify** | Ansible / patching pipeline groups gateway + Redis together |

---

## 9. License / cost assumptions

### A26. OSS Redis + custom sync daemon is sufficient for our DR posture (revised 2026-06-09)

| | |
|---|---|
| **Basis** | Per A8 revised, we deploy per-DC OSS Sentinel + custom sync daemon ([doc 22](22-redis-cross-dc-sync.md)) for selective cross-DC replication. Avoids Redis Enterprise license. |
| **What breaks if wrong** | If sync daemon proves unmaintainable in production OR partner SLA contractually requires strong cross-DC consistency, fall back to Redis Enterprise. Carry ~$15-25K/yr+ as budget contingency. |
| **Verify** | Quarterly review of sync daemon SLIs (lag, drop rate, partition recovery time); annual review of "should we have just paid for Enterprise" question |
| **Cost shape** | One-time build ~276 hrs; ongoing maintenance ~76 hrs/yr (per doc 22 §7). Breaks even vs Enterprise list pricing at ~Year 2-3 depending on consultant rates. |

### A27. Redis runs on existing VM capacity (no incremental hardware cost)

| | |
|---|---|
| **Basis** | 16 Redis nodes × 2 vCPU = 32 vCPU; absorbable in most enterprise capacity pools |
| **What breaks if wrong** | If hardware purchase required, ~$30K incremental for 16 small VMs |
| **Verify** | Capacity planning with infrastructure team |

---

## 10. MuleSoft product assumptions

### A28. Omni Gateway 1.4+ supports the Sentinel-mode configuration block as written

| | |
|---|---|
| **Basis** | Public MuleSoft docs reference Sentinel + Cluster support |
| **What breaks if wrong** | Field names may differ; some versions only support direct Redis URL, not Sentinel discovery |
| **Verify** | **Test against your actual Omni Gateway version in DEV before assuming.** Field names have shifted between minor releases. |
| **Critical** | This is the most version-specific assumption in the entire design |

### A29. Anypoint Control Plane repopulates runtime config to Redis on cold start

| | |
|---|---|
| **Basis** | Connected-mode design — Anypoint is authoritative |
| **What breaks if wrong** | Redis cold start without Anypoint reachability = gateway can't accept traffic |
| **Verify** | Cold-start drill: shut down Redis cluster in QA, restart, time until gateway resumes accepting traffic. Expect < 60 s. |

### A30. MuleSoft doesn't change the Connected-mode Redis dependency in a future Omni Gateway version

| | |
|---|---|
| **Basis** | Wishful thinking — product evolution can change requirements |
| **What breaks if wrong** | Major architecture pivot; your Redis infrastructure becomes mandatory or unnecessary depending on direction |
| **Verify** | Track MuleSoft roadmap announcements; reassess each Omni Gateway major version |

---

## Reviewer checklist

Print this. Each assumption gets owner + status + last verified date.

```
Assumption | Owner | Status | Last verified
-----------|-------|--------|---------------
A1  ...    |       |        |
A2  ...    |       |        |
...
A30 ...    |       |        |
```

Review cadence:
- Architecture review board: quarterly walkthrough of A1–A30
- Pre-release: re-verify A1, A4, A5, A18, A20, A28, A29 (the "drift-prone" set)
- Post-incident: review the relevant assumption block to see whether the incident invalidated one we'd thought safe

---

## Highest-risk assumptions (the short list)

If you're short on review time, focus on these five — they're the ones most likely to silently degrade or be wrong out of the gate:

1. **A5** — Redis latency stays under 2 ms (silent network-change risk)
2. **A8** — Active/passive Redis is acceptable for DR (contract risk)
3. **A18** — `fallbackOnError: local` is correctly configured (config-naming version drift)
4. **A20** — Request data cache is transient, not long-lived (PII exposure risk)
5. **A28** — Omni Gateway config block matches your actual version (version drift)

---

## Related

- [10 — Redis Shared Storage](10-redis-cache.md) — the design these assumptions support
- [07 — Data Protection](07-data-protection.md) — the citizen-data lens that drives the security assumptions (A10–A14, A20–A22)
- [09 — On-Prem Install Guide](09-onprem-install.md) — environment + ops assumptions (A23–A27)
