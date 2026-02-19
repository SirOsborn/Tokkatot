# Tokkatot 2.0: Team Structure & Responsibilities

**Document Version**: 2.0  
**Last Updated**: February 2026  
**Status**: Organization Plan

---

## Team Overview

**Total Team Size**: 6-8 people  
**Duration**: 6-8 months (with scaling after launch)  
**Reporting Structure**: Project Manager ← Team Lead ← Individual Contributors  

---

## Team Members & Roles

### Engineering Leadership

#### 1. Tech Lead / Project Manager
**Name**: [Lead Name - TBD]  
**Role**: Overall project leadership and execution  

**Responsibilities**:
- Project planning and timeline management
- Resource allocation
- Daily standup facilitation
- Stakeholder communication
- Risk tracking and mitigation
- Decision arbitration (technical/non-technical)
- Production deployment oversight

**Time Commitment**: 100% (full-time)  
**Key Skills**: Project management, full-stack development, decision-making  

**Success Criteria**:
- Project delivered on-time or within 2-week buffer
- Team morale high
- Zero critical escalations

---

### Backend Engineering

#### 2. Backend Lead - API & Services
**Name**: Heng  
**Role**: Backend architecture and development  

**Responsibilities**:
- API design and endpoint development
- Database schema and migrations
- Service architecture (Auth, Device, Schedule, Data services)
- Performance optimization
- Integration with DevOps/Deployment
- Code review for backend PRs
- Technology stack evaluation

**Time Commitment**: 100% (full-time, main contributor weeks 3-18)  
**Key Skills**: Go programming, API design, database optimization  
**Reports To**: Tech Lead  

**Deliverables**:
- Week 4: Authentication service complete
- Week 6: Device management complete
- Week 8: Schedule engine complete
- Week 14: Full backend integration tested

**Success Criteria**:
- 30+ API endpoints delivered
- API test coverage > 70%
- Load test passed (1000 concurrent)
- Zero data loss scenarios

---

#### 3. Backend Developer (Secondary)
**Name**: [TBD]  
**Role**: Supporting backend development  

**Responsibilities**:
- API development (assigned endpoints)
- Unit testing
- Bug fixes and debugging
- Database work (migrations, optimization)
- Supporting Heng with feature development
- Code review participation

**Time Commitment**: 100% (weeks 3-18)  
**Reports To**: Heng (Backend Lead)  

**Success Criteria**:
- Assigned features complete on-time
- Code review quality high
- Contributes to test coverage

---

### Frontend Engineering

#### 4. Frontend Lead - UI/UX Developer
**Name**: Rika  
**Role**: Frontend architecture and development  

**Responsibilities**:
- Vue.js 3 application architecture
- UI/UX design and implementation
- Responsive design (mobile-first)
- Accessibility audit and implementation (WCAG AAA)
- Performance optimization (bundle size)
- Component library creation
- Offline support (Service Workers)
- Code review for frontend PRs

**Time Commitment**: 100% (full-time, weeks 7-18)  
**Key Skills**: Vue.js, CSS/TailwindCSS, UI/UX, Accessibility  
**Reports To**: Tech Lead  

**Deliverables**:
- Week 10: Home & Control pages complete
- Week 12: Schedules & Analytics complete
- Week 14: Accessibility audit passed
- Week 14: Bundle size < 150KB

**Success Criteria**:
- Web UI fully functional
- Passes WCAG AAA accessibility audit
- All 48x48px buttons
- Loads in < 2s on 4G
- Works on 1-2GB RAM devices

---

### Embedded Systems

#### 5. Embedded Systems Lead - Firmware
**Name**: Neath  
**Role**: ESP32 firmware development  

**Responsibilities**:
- Firmware architecture design
- GPIO and relay control module
- Sensor integration (DHT22, ADC)
- MQTT communication implementation
- OTA (Over-The-Air) firmware updates
- Local scheduling (offline mode)
- Watchdog and recovery mechanisms
- Embedded code review
- Hardware testing and validation

**Time Commitment**: 100% (full-time, weeks 8-18)  
**Key Skills**: C/C++, ESP-IDF, Embedded Systems, MQTT  
**Hardware**: ESP32 development boards, relay modules, sensors  
**Reports To**: Tech Lead  

**Deliverables**:
- Week 9: GPIO control working
- Week 11: Sensors and MQTT working
- Week 13: OTA updates functional
- Week 14: Field testing complete

**Success Criteria**:
- Firmware compiles without warnings
- All device controls functional
- MQTT reliable under poor connectivity
- OTA updates work with automatic rollback
- Field testing passes

---

### Infrastructure & DevOps

#### 6. DevOps Engineer / Infrastructure
**Name**: Raingsey  
**Role**: Infrastructure setup and CI/CD  

**Responsibilities**:
- DigitalOcean infrastructure provisioning
- Kubernetes cluster setup (or Docker Compose)
- CI/CD pipeline (GitHub Actions)
- Database setup and management
- Monitoring and alerting (Prometheus/Grafana)
- Security hardening
- Deployment procedures
- Disaster recovery planning
- Infrastructure-as-Code (Terraform)

**Time Commitment**: 50% (weeks 1-4), 100% (weeks 19-26)  
**Key Skills**: DevOps, Kubernetes/Docker, CI/CD, Linux, Infrastructure  
**Reports To**: Tech Lead  

**Deliverables**:
- Week 2: Staging environment ready
- Week 5: CI/CD pipeline working
- Week 10: Monitoring dashboards set up
- Week 19: Production infrastructure ready

**Success Criteria**:
- Infrastructure as Code complete
- CI/CD deployments fully automated
- 99.5% uptime achieved
- All monitoring/alerting in place

---

### Quality Assurance

#### 7. QA Lead / Test Automation
**Name**: [TBD]  
**Role**: Testing coordination and automation  

**Responsibilities**:
- Test planning and strategy
- Test case creation
- Automated test framework setup
- Integration testing coordination
- Performance testing (load test, stress test)
- Security testing coordination
- Bug tracking and prioritization
- Release candidate verification

**Time Commitment**: 50% (weeks 13-20), 100% (weeks 21-26)  
**Key Skills**: QA automation, testing tools, performance testing  
**Reports To**: Tech Lead  

**Deliverables**:
- Week 14: Integration test suite ready
- Week 16: Load test completed (1000 concurrent)
- Week 17: Security audit passed
- Week 22: Production sign-off

**Success Criteria**:
- No critical bugs at launch
- 1000 concurrent user load test passed
- Security audit passed

---

### Support Roles

#### 8. AI/ML Specialist (Optional)
**Name**: [TBD]  
**Role**: Disease detection model and integration  

**Responsibilities**:
- Train disease detection model
- API endpoint development
- Model optimization for production
- Integration with backend

**Time Commitment**: 50% (weeks 14-18)  
**Reports To**: Heng (Backend Lead)  

**Deliverables**:
- Week 18: Disease detection API working
- Model accuracy > 85%

---

## Communication Structure

### Standup Meetings

**Daily Standup**: 15 minutes (9:00 AM)
- All team members
- Format: What I did, what I'm doing, blockers
- Location: Zoom/in-person

**Weekly Planning**: 1 hour (Monday 10:00 AM)
- Tech Lead + Team Leads
- Review previous week
- Plan upcoming week
- Adjust timeline if needed

**Bi-Weekly Review**: 1.5 hours (every other Friday)
- Entire team
- Demo of completed work
- Stakeholder feedback
- Retrospective (what went well, improve)

### Communication Channels

**Real-time Chat**: Slack
- `#tokkatot-general` - announcements
- `#tokkatot-backend` - backend discussions
- `#tokkatot-frontend` - frontend discussions
- `#tokkatot-embedded` - embedded systems
- `#tokkatot-bugs` - bug reports
- `#tokkatot-urgent` - critical issues

**Email**: For formal communications
- Decisions, approvals
- Escalations
- External communications

**GitHub Issues**: Technical tracking
- Bug reports
- Feature requests
- Code review discussions

---

## Decision-Making & Escalation

### Decision Authority

| Decision Type | Authority | Approval Required |
|---|---|---|
| **Technical Architecture** | Tech Lead + Team Leads | CTO (if exists) |
| **Feature Scope** | Tech Lead | Project Manager |
| **Timeline Changes** | Tech Lead | Stakeholders |
| **Infrastructure Changes** | DevOps + Tech Lead | Tech Lead approval |
| **Security Issues** | Tech Lead + DevOps | CTO + Legal (if sensitive) |
| **Production Changes** | Tech Lead + DevOps | Tech Lead sign-off |
| **Hiring/Roles** | Tech Lead | HR/Management |

### Escalation Path

**Level 1**: Team member → Team Lead (same discipline)  
**Level 2**: Team Lead → Tech Lead (cross-team)  
**Level 3**: Tech Lead → Project Manager/CTO (strategic)  
**Emergency**: Direct to Tech Lead (any level)  

---

## Knowledge & Code Ownership

### Code Ownership

| Component | Owner | Backup |
|---|---|---|
| Authentication | Heng | Backend Dev |
| Device Service | Heng | Backend Dev |
| Schedule Engine | Backend Dev | Heng |
| API Gateway | Heng | DevOps |
| Real-time (WebSocket) | Backend Dev | Heng |
| Frontend Core | Rika | Frontend Dev (if hired) |
| Embedded Firmware | Neath | [TBD] |
| Infrastructure | Raingsey | Tech Lead |
| Database | Heng + Raingsey | Tech Lead |

### Knowledge Transfer

- **Pair Programming**: 1-2 hours per week (cross-functional)
- **Documentation**: Code comments, README files
- **Design Documents**: Architecture decisions
- **Video Tutorials**: How-to guides for complex systems

---

## Onboarding & Training

### New Team Member Onboarding (1 week)

**Day 1**:
- Welcome & team introduction
- Development environment setup
- GitHub access, communication tools
- Project overview

**Day 2-3**:
- Technology stack deep dive
- Code repository walkthrough
- Setting up local development
- Running tests and deploying locally

**Day 4**:
- Pair programming with team lead
- First small task/bug fix
- Code review process

**Day 5**:
- Independent task completion
- Team collaboration
- Retrospective

---

## Performance & Recognition

### Key Performance Indicators (KPIs)

**All Team Members**:
- Deliver assigned features on-time (90%+ target)
- Code quality (> 70% test coverage, < 5 critical bugs/month)
- Responsiveness (< 2 hour response to blockers)
- Collaboration (active participation, helping others)

**Team Leads**:
- Team productivity and morale
- Feature delivery quality
- Code review turnaround (< 4 hours)
- Risk identification and mitigation

**Tech Lead**:
- Overall project delivery (on-time, budget)
- Team morale and retention
- Quality metrics maintained
- Stakeholder satisfaction

### Recognition & Rewards

- **Monthly**: Highlight top contributor in standup
- **Quarterly**: Performance review and bonuses
- **Launch**: Team celebration event
- **Post-Launch**: Career development planning

---

## Scaling Plan (Post-Launch)

### Month 1-3 (Launch Phase)
- Current team: 6-8 people
- Focus: Stabilization, critical bugs, user support

### Month 4-6 (Growth Phase)
- Add: 1 Frontend Dev, 1 Backend Dev (optional)
- Focus: Feature development, performance improvement
- Hiring plan: Start recruiting for v2.1 features

### Month 7+ (Scale Phase)
- Add: 2 Backend Devs, 1 Frontend Dev, 1 QA Dev
- Structure: Sub-teams by feature area
- Focus: Multi-region support, advanced features

---

## Contact Information Template

```
Name: [Name]
Email: [email@example.com]
Phone: [+855-XX-XXX-XXXX]
Role: [Role]
Time Zone: [UTC+7 (Cambodia)]
Availability: Mon-Fri 8 AM - 6 PM (flexible)
```

---

**Version History**
| Version | Date | Changes |
|---------|------|---------|
| 2.0 | Feb 2026 | Initial team structure spec |

**Related Documents**
- PROJECT_TIMELINE.md (project schedule)
- RISK_MANAGEMENT.md (risk tracking)
- SPECIFICATIONS_REQUIREMENTS.md (deliverables)
