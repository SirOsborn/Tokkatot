# Tokkatot 2.0: Risk Management Plan

**Document Version**: 2.0  
**Last Updated**: February 2026  
**Status**: Active Planning

---

## Risk Management Process

### Risk Assessment Matrix

```
        Likelihood
        │
   High │ Medium │ Critical
        ├────────┼────────
   Med  │ Low    │ Medium
        ├────────┼────────
   Low  │ Low    │ Low
        └────────┼────────
              Impact
```

**Risk Score** = Likelihood × Impact  
- **Critical** (15+): Immediate action required
- **High** (10-15): Senior team involvement required
- **Medium** (5-10): Team lead tracking required
- **Low** (<5): Monitor, escalate if changes

---

## Identified Risks

### Risk 1: Technology Selection Challenges

**Description**: Go vs Node.js decision might impact team productivity  
**Likelihood**: Medium (60%)  
**Impact**: High (8/10) - Would require rewrite if wrong  
**Score**: 4.8 (Medium)  

**Symptoms**:
- Go learning curve steep for some team members
- Node.js would be faster to develop initially

**Mitigation**:
1. **Proactive**: Conduct technology sprint (Week 1)
   - Heng validates Go choice with PoC
   - Compare productivity vs Node.js
   - Document decision rationale
2. **Reactive**: If Go problematic by Week 5
   - Migrate to Express.js (2-3 week effort)
   - Prioritize critical APIs first

**Contingency Plan**: Have Express.js skeleton ready Week 3  
**Owner**: Heng (Backend Lead)  
**Status**: Green (assessment in progress)

---

### Risk 2: API Contract Misalignment

**Description**: Frontend and Backend developers define different API contracts  
**Likelihood**: Medium (50%)  
**Impact**: High (8/10) - Would cause integration delays  
**Score**: 4.0 (Medium)  

**Symptoms**:
- Frontend expects different response format than API provides
- Missing or extra API fields
- Wrong HTTP status codes

**Mitigation**:
1. **Proactive**: Share OpenAPI spec before Week 4
   - Heng creates detailed API spec Week 2
   - Rika reviews and approves Week 3
   - Mock API created for frontend testing
2. **Weekly**: API review meetings (every Tuesday)
   - Frontend shows what they're building
   - Backend confirms API matches

**Contingency Plan**: 
- Implement API adapter layer (1-2 days to fix)

**Owner**: Heng + Rika (joint responsibility)  
**Status**: Green (planning to prevent)

---

### Risk 3: Database Performance Issues

**Description**: Database queries might become slow with scale  
**Likelihood**: Low (30%)  
**Impact**: High (8/10) - Platform becomes unusable  
**Score**: 2.4 (Low)  

**Symptoms**:
- API responses > 500ms
- Dashboard charts don't render
- Schedule execution delays
- Database CPU near max

**Mitigation**:
1. **Proactive**: Load testing Week 15
   - Simulate 1000 concurrent users
   - Identify slow queries
   - Add indexes and optimize
2. **Preventive**: Code review standards
   - Database queries reviewed before merge
   - N+1 query detection
   - Index requirements documented

**Contingency Plan**:
- Add read replicas (1-2 weeks deployment)
- Cache layer optimization
- Query optimization expert consultation

**Owner**: Heng + Raingsey (joint)  
**Status**: Green (load testing planned)

---

### Risk 4: Firmware Reliability Issues

**Description**: ESP32 firmware crashes or doesn't handle edge cases  
**Likelihood**: Medium (50%)  
**Impact**: High (9/10) - Farm can't operate without reliable devices  
**Score**: 4.5 (Medium)  

**Symptoms**:
- Devices go offline unexpectedly
- Commands don't execute
- WiFi connection drops frequently
- OTA updates fail silently

**Mitigation**:
1. **Proactive**: Extended testing Week 13-14
   - Field test in real farm (72+ hours)
   - Test network interruptions
   - Test OTA update flow
   - Test against command flood
2. **Preventive**: Code practices
   - Watchdog timer enabled (restart every 10s without reset)
   - Connection retry with backoff
   - Graceful degradation (local mode on cloud failure)
   - Comprehensive logging

**Contingency Plan**:
- Distribute known-good firmware version (1 week)
- Develop firmware hotfix process
- Local fallback guarantees

**Owner**: Neath (Embedded Lead)  
**Status**: Green (testing planned)

---

### Risk 5: Frontend Performance on Low-End Devices

**Description**: Web app is too slow or doesn't work on 1-2GB RAM phones  
**Likelihood**: Medium (70%)  
**Impact**: High (9/10) - Users can't use the app  
**Score**: 6.3 (Medium-High)  

**Symptoms**:
- Bundle size > 200KB
- App takes > 5 seconds to load
- Buttons feel unresponsive
- Charts don't render
- Crashes on old devices

**Mitigation**:
1. **Proactive**: Real device testing Week 12
   - Test on actual 1-2GB RAM phones (Neath's phone?)
   - Test on slow 4G networks (throttle to 1 Mbps)
   - Test zoom to 200% magnification
2. **Preventive**: Performance budgets
   - Bundle size: < 150KB
   - Initial load: < 2 seconds
   - Page transitions: < 300ms
   - Code splitting and lazy loading
3. **Optimization**: Ongoing
   - Remove unused dependencies
   - CSS/JS tree-shaking
   - Image optimization (WebP)
   - Service worker caching

**Contingency Plan**:
- Switch to lighter framework (Svelte, etc) - 2 week rewrite
- Reduce feature set - defer to v2.1
- Progressive enhancement - basic mode for slow devices

**Owner**: Rika (Frontend Lead)  
**Status**: Yellow (proactive testing needed)

---

### Risk 6: Team Member Departure

**Description**: Key team member (especially Neath or Heng) leaves the project  
**Likelihood**: Low-Medium (25%)  
**Impact**: High (9/10) - Critical knowledge lost  
**Score**: 2.25 (Low)  

**Symptoms**:
- Team member gives notice
- Critical code only one person understands
- Undocumented decisions

**Mitigation**:
1. **Proactive**: Knowledge documentation
   - Architecture decision records (ADRs) added to repo
   - Code comments and docstrings
   - Weekly pair programming with backup person
2. **Preventive**: Cross-training
   - Pair programming 2 hours/week (Heng ↔ Secondary Dev)
   - Pair programming 1 hour/week (Neath ↔ Secondary Embedded)
3. **Retention**: Team morale
   - Competitive compensation
   - Career development
   - Recognition and bonuses

**Contingency Plan**:
- Hire replacement (2-3 weeks recruitment)
- Contract consultant (1 week onboarding)
- Redistribute work to team

**Owner**: Tech Lead  
**Status**: Green (prevention in place)

---

### Risk 7: Cloud Environment Outage

**Description**: DigitalOcean has major outage affecting production  
**Likelihood**: Very Low (5%)  
**Impact**: Critical (10/10) - Platform completely down  
**Score**: 0.5 (Very Low)  

**Symptoms**:
- All services unreachable
- Databases cannot connect
- Users cannot login or control devices

**Mitigation**:
1. **Preventive**: High availability
   - Multi-region setup (optional, future)
   - Database backups (daily, 30-day retention)
   - Infrastructure-as-Code (can rebuild in 4 hours)
   - Staged database replicas
2. **Detection**: Monitoring
   - Prometheus/Grafana 24/7 monitoring
   - Alert on service downtime
   - Pagerduty escalation

**Contingency Plan**:
- Rebuild infrastructure (4 hours)
- Restore from most recent backup (< 1 hour data loss)
- Switch to AWS (24-48 hours, manual setup)
- Local fallback on Raspberry Pi (farmers can still control devices)

**Owner**: Raingsey (DevOps)  
**Status**: Green (disaster recovery plan ready)

---

### Risk 8: Security Vulnerabilities Discovered

**Description**: Critical security vulnerability found (XSS, SQL injection, etc)  
**Likelihood**: Medium (40%)  
**Impact**: Critical (10/10) - User data compromised  
**Score**: 4.0 (Medium)  

**Symptoms**:
- Security researcher reports vulnerability
- Vulnerability discovered during audit
- Exploit found in wild

**Mitigation**:
1. **Proactive**: Security practices
   - OWASP Top 10 review Week 16
   - Dependency vulnerability scanning (Snyk, etc)
   - Penetration testing (professional, Week 17)
   - Code review security checklist
2. **Preventive**: Secure coding
   - Input validation on all endpoints
   - Parameterized queries (no SQL injection)
   - Output encoding (no XSS)
   - CSRF token protection
3. **Detection**: Monitoring
   - Log suspicious activities
   - Alert on multiple failed login attempts
   - Monitor API error rates

**Contingency Plan**:
- Security hotfix process (< 4 hours to patch + deploy)
- Communication plan (notify affected users)
- Incident postmortem

**Owner**: Tech Lead + Security specialist  
**Status**: Yellow (audit scheduled Week 17)

---

### Risk 9: Scope Creep

**Description**: Project requirements expand beyond original scope  
**Likelihood**: High (80%)  
**Impact**: Medium (6/10) - Timeline slips 2-4 weeks  
**Score**: 4.8 (Medium)  

**Symptoms**:
- New feature requests during development
- "While we're at it" modifications
- Stakeholder feature additions
- Timeline slippage

**Mitigation**:
1. **Preventive**: Scope management
   - Locked feature set before Week 4
   - "Nice-to-have" moved to v2.1 backlog
   - Change control process
   - Weekly scope review with stakeholders
2. **Reactive**: Feature triage
   - Critical: implement immediately
   - Important: defer to v2.1
   - Nice-to-have: definitely v2.1
   - Prioritize speed over perfection

**Contingency Plan**:
- Reduce scope (defer less critical features)
- Add team members (short-term consultants)
- Extend timeline (notify stakeholders early)

**Owner**: Tech Lead  
**Status**: Yellow (scope lock Week 4)

---

### Risk 10: Data Migration Issues

**Description**: Data migration from v1.0 to v2.0 loses data or creates corruption  
**Likelihood**: Low (25%)  
**Impact**: Critical (10/10) - Farm historical data lost  
**Score**: 2.5 (Low)  

**Symptoms**:
- Some records fail to migrate
- Data type conversion errors
- Foreign key constraint violations
- Timestamp timezone issues

**Mitigation**:
1. **Proactive**: Migration planning
   - Schema mapping created (Week 18)
   - Migration scripts tested repeatedly (Week 19-21)
   - Dry-run on production copy (Week 21)
   - Validation checks after migration
2. **Preventive**: Fallback plan
   - Keep v1.0 running in parallel (48-72 hours)
   - Full database backup before migration
   - Rollback procedure documented

**Contingency Plan**:
- Restore from backup (< 1 hour)
- Retry migration with fixes
- Manual data entry as last resort (unlikely)

**Owner**: Heng + Raingsey  
**Status**: Green (migration plan being created)

---

## Risk Tracking & Reporting

### Risk Register (Active Tracking)

| # | Risk | Status | Owner | Review |
|---|------|--------|-------|--------|
| 1 | Tech Selection | Green | Heng | Weekly |
| 2 | API Alignment | Green | Heng/Rika | Weekly |
| 3 | DB Performance | Green | Heng | Bi-weekly |
| 4 | Firmware Issues | Green | Neath | Weekly |
| 5 | Frontend Perf | Yellow | Rika | Weekly |
| 6 | Team Departure | Green | Tech Lead | Monthly |
| 7 | Cloud Outage | Green | Raingsey | Monthly |
| 8 | Security Issues | Yellow | Tech Lead | Weekly |
| 9 | Scope Creep | Yellow | Tech Lead | Weekly |
| 10 | Migration | Green | Heng | Bi-weekly |

### Risk Review Meetings

**Weekly**: Standup (all team)
- Quick status on active risks
- New risks identified

**Bi-Weekly**: Risk Review (Tech Lead + Team Leads)
- Deep dive on Yellow/Red risks
- Mitigation execution status
- Escalation if needed

**Monthly**: Executive Review (if stakeholders exist)
- Overall project health
- Timeline/budget impact
- Major risk changes

### Risk Status Colors

🟢 **Green**: Low risk, under control, on mitigation plan  
🟡 **Yellow**: Medium risk, requires active attention  
🔴 **Red**: High risk, immediate action needed  

---

## Escalation Triggers

**Auto-Escalate to Tech Lead if**:
- Risk score increases by 2+ points
- Mitigation plan not executed by deadline
- New high-risk item identified
- Status change to Yellow/Red

**Auto-Escalate to CTO/Stakeholders if**:
- Timeline extends beyond 35 weeks
- Critical security issue
- Major funding/scope impact
- Team member departure

---

## Lessons Learned (Post-Project)

**Post-Launch Review** (Week 26):
- Which risks materialized? Why?
- Which mitigations were effective?
- Which risks were overblown?
- What new risks were missed?

**Continuous Improvement**:
- Update risk register for future projects
- Share learnings with team
- Refine estimation and planning

---

**Version History**
| Version | Date | Changes |
|---------|------|---------|
| 2.0 | Feb 2026 | Initial risk management plan |

**Related Documents**
- PROJECT_TIMELINE.md (timeline visibility)
- TEAM_STRUCTURE.md (team resources)
- SPECIFICATIONS_REQUIREMENTS.md (scope definition)
