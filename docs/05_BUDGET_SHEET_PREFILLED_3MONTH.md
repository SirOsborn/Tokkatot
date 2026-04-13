# Tokkatot Technical Budget Sheet (Prefilled, 3-Month Pilot)

This prefilled sheet is technical spending only (no man-day/labor).

Assumptions used:
- One shared cloud backend for all pilot farms.
- 2 standard farm kits + 1 AI farm kit.
- AI image retention enabled (object storage budget included).
- Cloudflare paid add-ons are budgeted as small optional reserves.

---

## 1) Prefilled Line Items

| Category | Subcategory | Item | Qty | Unit Cost (USD) | Duration (Months) | Subtotal (USD) | Priority | Phase | Owner | Notes |
|---|---|---|---:|---:|---:|---:|---|---|---|---|
| Cloud Infra | Compute | AWS EC2 instance | 1 | 12.00 | 3 | 36.00 | Must-have | P1 |  | Shared pilot app server |
| Cloud Infra | Storage | EBS gp3 volume | 1 | 6.00 | 3 | 18.00 | Must-have | P1 |  | App + DB disk |
| Cloud Infra | Network | Public IPv4 | 1 | 3.60 | 3 | 10.80 | Must-have | P1 |  | Static public IP |
| Cloud Infra | Backup | EBS snapshots | 1 | 4.00 | 3 | 12.00 | Must-have | P1 |  | Backup reserve |
| Cloud Infra | Data Transfer | Outbound bandwidth reserve | 1 | 8.00 | 3 | 24.00 | Must-have | P1 |  | Traffic headroom |
| Cloudflare | Edge Security | WAF/rate-limit paid add-on (if used) | 1 | 5.00 | 3 | 15.00 | Nice-to-have | P1 |  | Optional reserve |
| Cloudflare | Website | Pages/Workers paid tier (if needed) | 1 | 5.00 | 3 | 15.00 | Nice-to-have | P1 |  | Optional reserve |
| Data Platform | Object Storage | AI image storage bucket (S3 or R2) | 1 | 10.00 | 3 | 30.00 | Must-have | P3 |  | Private image storage |
| Data Platform | Requests | Object storage operation reserve | 1 | 5.00 | 3 | 15.00 | Must-have | P3 |  | GET/PUT request reserve |
| Data Platform | Lifecycle | Tiering/archive reserve | 1 | 4.00 | 3 | 12.00 | Nice-to-have | P4 |  | For long retention optimization |
| Software | Landing Page | Final UI/content/assets completion | 1 | 80.00 | 1 | 80.00 | Must-have | P1 |  | Technical spend only |
| Software | Web App | QA, bug bash, pilot hardening tools | 1 | 25.00 | 3 | 75.00 | Must-have | P2 |  | Tooling reserve |
| Software | Security | Testing tools/scanners/logging tools | 1 | 20.00 | 3 | 60.00 | Nice-to-have | P2 |  | Tooling reserve |
| AI Service | Integration | AI API deployment and integration infra | 1 | 30.00 | 3 | 90.00 | Must-have | P3 |  | AI path readiness |
| AI Service | Data Ops | Labeling/storage utilities | 1 | 15.00 | 3 | 45.00 | Nice-to-have | P4 |  | Dataset organization |
| Hardware | Standard Kit | Gateway kit for standard farm | 2 | 300.00 | 1 | 600.00 | Must-have | P2 |  | Pi4 + sensors + relays |
| Hardware | AI Kit | Gateway kit for AI farm | 1 | 330.00 | 1 | 330.00 | Must-have | P3 |  | Pi5 + AI HAT+ + camera |
| Hardware | Spare Parts | Sensors/relays/SD cards/cables reserve | 1 | 170.00 | 1 | 170.00 | Must-have | P2 |  | Failure replacement |
| Hardware | Enclosure & Power | Cases, PSU, cooling, weather protection | 3 | 33.33 | 1 | 99.99 | Must-have | P2 |  | Per deployment set |
| R&D Testing | Bench Testing | Test consumables and replacements | 1 | 120.00 | 1 | 120.00 | Must-have | P2 |  | Hardware bench verification |
| R&D Testing | Field Testing | On-site test support materials | 1 | 130.00 | 1 | 130.00 | Must-have | P3 |  | Pilot validation materials |
| R&D Testing | Validation | Calibration and reliability tests | 1 | 100.00 | 1 | 100.00 | Must-have | P3 |  | Stability checks |
| Risk Buffer | Contingency | Technical contingency reserve (12%) | 1 | 250.53 | 1 | 250.53 | Must-have | P1-P5 |  | 12% of pre-contingency subtotal |

---

## 2) Summary

| Summary Metric | Value (USD) |
|---|---:|
| Cloud Infra Subtotal | 100.80 |
| Cloudflare/Web Subtotal | 30.00 |
| Data Platform Subtotal | 57.00 |
| Software Subtotal | 215.00 |
| AI Service Subtotal | 135.00 |
| Hardware Subtotal | 1199.99 |
| R&D Testing Subtotal | 350.00 |
| Technical Subtotal (Before Contingency) | 2087.79 |
| Contingency (12%) | 250.53 |
| Grand Total (3-Month Technical Budget) | 2338.32 |

---

## 3) Funding Fit Snapshot

- If using startup activities funding pool of `5948 USD`, this prefilled technical plan uses about `39.31%`.
- Remaining budget can cover non-technical categories (logistics, marketing, workshops, etc.) per program rules.

---

## 4) Adjustment Rules

- If hardware quotes increase, protect security and core reliability first; reduce only nice-to-have tooling.
- If AI image volume grows fast, increase object storage and lifecycle reserve first.
- If cloud load increases, upgrade compute before cutting backups.

---

**Template Version:** Prefilled v1.0  
**Last Updated:** 2026-04-12
