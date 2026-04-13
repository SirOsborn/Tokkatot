# Tokkatot Technical Budget Sheet Template (3-Month Pilot)

Use this sheet for technical spending only (no man-day/labor).

---

## 1) How to Use

1. Fill `Qty`, `Unit Cost (USD)`, and `Duration (Months)`.
2. Compute `Subtotal = Qty x Unit Cost x Duration` (or one-time lump sum with `Duration=1`).
3. Mark `Priority` as `Must-have` or `Nice-to-have`.
4. Keep 10-15% contingency in the final total.

---

## 2) Line-Item Budget Table

| Category | Subcategory | Item | Qty | Unit Cost (USD) | Duration (Months) | Subtotal (USD) | Priority | Phase | Owner | Notes |
|---|---|---|---:|---:|---:|---:|---|---|---|---|
| Cloud Infra | Compute | AWS EC2 instance | 1 |  | 3 |  | Must-have | P1 |  |  |
| Cloud Infra | Storage | EBS gp3 volume | 1 |  | 3 |  | Must-have | P1 |  |  |
| Cloud Infra | Network | Public IPv4 | 1 |  | 3 |  | Must-have | P1 |  |  |
| Cloud Infra | Backup | EBS snapshots | 1 |  | 3 |  | Must-have | P1 |  |  |
| Cloud Infra | Data Transfer | Outbound bandwidth reserve | 1 |  | 3 |  | Must-have | P1 |  |  |
| Cloudflare | Edge Security | WAF/rate-limit paid add-on (if used) | 1 |  | 3 |  | Nice-to-have | P1 |  |  |
| Cloudflare | Website | Pages/Workers paid tier (if needed) | 1 |  | 3 |  | Nice-to-have | P1 |  |  |
| Data Platform | Object Storage | AI image storage bucket (S3 or R2) | 1 |  | 3 |  | Must-have | P3 |  |  |
| Data Platform | Requests | Object storage operation reserve | 1 |  | 3 |  | Must-have | P3 |  |  |
| Data Platform | Lifecycle | Tiering/archive reserve | 1 |  | 3 |  | Nice-to-have | P4 |  |  |
| Software | Landing Page | Final UI/content/assets completion | 1 |  | 1 |  | Must-have | P1 |  | Technical spend only |
| Software | Web App | QA, bug bash, pilot hardening tools | 1 |  | 3 |  | Must-have | P2 |  |  |
| Software | Security | Testing tools/scanners/logging tools | 1 |  | 3 |  | Nice-to-have | P2 |  |  |
| AI Service | Integration | AI API deployment and integration infra | 1 |  | 3 |  | Must-have | P3 |  |  |
| AI Service | Data Ops | Labeling/storage utilities | 1 |  | 3 |  | Nice-to-have | P4 |  |  |
| Hardware | Standard Kit | Gateway kit for standard farm | 2 |  | 1 |  | Must-have | P2 |  | Pi4 + sensors + relays |
| Hardware | AI Kit | Gateway kit for AI farm | 1 |  | 1 |  | Must-have | P3 |  | Pi5 + AI HAT+ + camera |
| Hardware | Spare Parts | Sensors/relays/SD cards/cables reserve | 1 |  | 1 |  | Must-have | P2 |  |  |
| Hardware | Enclosure & Power | Cases, PSU, cooling, weather protection | 3 |  | 1 |  | Must-have | P2 |  |  |
| R&D Testing | Bench Testing | Test consumables and replacements | 1 |  | 3 |  | Must-have | P2 |  |  |
| R&D Testing | Field Testing | On-site test support materials | 1 |  | 3 |  | Must-have | P3 |  | Excluding logistics if separate |
| R&D Testing | Validation | Calibration and reliability tests | 1 |  | 3 |  | Must-have | P3 |  |  |
| Risk Buffer | Contingency | Technical contingency reserve (10-15%) | 1 |  | 1 |  | Must-have | P1-P5 |  | From subtotal |

---

## 3) Summary Block

| Summary Metric | Value (USD) |
|---|---:|
| Cloud Infra Subtotal |  |
| Cloudflare/Web Subtotal |  |
| Data Platform Subtotal |  |
| Software Subtotal |  |
| AI Service Subtotal |  |
| Hardware Subtotal |  |
| R&D Testing Subtotal |  |
| Technical Subtotal (Before Contingency) |  |
| Contingency (10-15%) |  |
| Grand Total (3-Month Technical Budget) |  |

---

## 4) Phase Mapping

- `P1`: Foundation hardening (cloud, DNS, security baseline)
- `P2`: Multi-tenant safety + standard farm rollout
- `P3`: AI farm integration + field validation
- `P4`: Data lifecycle + optimization
- `P5`: Stabilization and handover

---

## 5) Budget Guardrails

- Prefer monthly billing during pilot.
- Avoid long-term commitments before KPI validation.
- Keep all tenant data isolation and security controls as non-negotiable `Must-have`.
- Keep all AI image storage private and tracked with metadata.

---

**Template Version:** 1.0  
**Last Updated:** 2026-04-10
