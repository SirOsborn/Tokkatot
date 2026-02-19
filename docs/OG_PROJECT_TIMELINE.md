# Tokkatot 2.0: Project Timeline & Milestones

**Document Version**: 2.0  
**Last Updated**: February 2026  
**Status**: Planning Phase  

**Total Duration**: 27-35 weeks (6-8 months)  
**Start Date**:  End of February 2026 
**Target Launch**: Q2/Q3 2026  

---

## Project Phases

### Phase 1: Planning & Architecture (Weeks 1-2)
**Duration**: 2 weeks  
**Team**: Tech Lead, Architects (Neath, Heng)  

**Activities**:
- Team kickoff and assignment
- Technology stack finalization
- Development environment setup
- Database schema finalization
- API contract definition
- Infrastructure planning (DigitalOcean setup)

**Deliverables**:
- Development environment ready
- Database schema approved
- API endpoints documented (Swagger)
- Infrastructure provisioned (staging)
- CI/CD pipeline created

**Success Criteria**:
- All team members have development environment
- First "hello world" API deployed to staging
- Database migrations working

---

### Phase 2: Backend Development (Weeks 3-10)
**Duration**: 8 weeks  
**Team**: Backend (Heng), DevOps, Secondary Dev  

**Weeks 3-4: Core Services**
- Authentication service (JWT, login, registration)
- User management and RBAC
- Database connection pooling
- Error handling framework
- Logging infrastructure

**Weeks 5-6: Device Services**
- Device registration and management
- Device communication handlers (MQTT receive)
- Device state tracking
- Command queue implementation
- Health monitoring

**Weeks 7-8: Advanced Features**
- Schedule engine (time-based, duration-based, condition-based)
- Real-time WebSocket server
- Batch operations
- Export functionality
- Notification engine

**Weeks 9-10: Integration & Polish**
- Third-party integrations (email, SMS)
- Rate limiting and throttling
- Comprehensive error handling
- Performance optimization

**Deliverables**:
- 30+ API endpoints implemented
- Database fully operational
- Unit tests (70%+ coverage)
- Staging environment fully functional
- API documentation complete

**Success Criteria**:
- All APIs return proper responses
- Authentication required and enforced
- No data loss on restart
- API response time < 200ms (95th percentile)

---

### Phase 3: Frontend Development (Weeks 7-14)
**Duration**: 8 weeks (overlaps with backend)  
**Team**: Frontend Developer (Rika), Secondary Dev  

**Weeks 7-8: UI Framework & Responsive Design**
- Vue.js 3 project setup
- TailwindCSS configuration
- Component library creation
- Responsive breakpoints

**Weeks 9-10: Core Pages**
- Home/Dashboard page
- Control page (device management)
- Authentication (login/signup)
- Settings page

**Weeks 11-12: Feature Pages**
- Schedules page (create, edit, list)
- Analytics/History page (charts)
- Notifications and alerts
- User management page

**Weeks 13-14: Polish & Optimization**
- Accessibility audit (WCAG AAA)
- Performance optimization (< 150KB bundle)
- Offline support (Service Workers)
- Dark mode (optional)
- Responsive testing on real devices

**Deliverables**:
- Complete web application (Vue.js)
- All planned pages functional
- Bundle size < 150KB
- Passes accessibility audit
- Works offline with cached data

**Success Criteria**:
- Application loads in < 2 seconds on 4G
- All buttons are 48x48px minimum
- 7:1 contrast ratio on all text
- Works on phones with 1-2GB RAM

---

### Phase 4: Embedded System Development (Weeks 8-14)
**Duration**: 7 weeks  
**Team**: Embedded (Neath)  

**Weeks 8-9: Firmware Foundation**
- ESP-IDF setup and configuration
- GPIO and relay control module
- MQTT client integration
- WiFi management

**Weeks 10-11: Sensor and Device Control**
- DHT22 sensor reading
- ADC implementation
- Sensor averaging and filtering
- PWM control (if needed)

**Weeks 12-13: Advanced Features**
- Local schedule execution (offline mode)
- OTA (Over-The-Air) firmware updates
- NVS (Non-Volatile Storage) for configs
- Status indicators (LEDs)
- Watchdog and recovery

**Weeks 14: Testing & Documentation**
- Unit testing on hardware
- Integration testing
- Field testing in test farm
- Documentation

**Deliverables**:
- Complete ESP32 firmware (production-ready)
- OTA update mechanism working
- Offline mode functional
- Device documentation

**Success Criteria**:
- Firmware compiles with no warnings
- All GPIO modules functional
- MQTT communication reliable
- OTA updates work with automatic rollback

---

### Phase 5: Integration & Testing (Weeks 13-18)
**Duration**: 6 weeks  

**Weeks 13-15: System Integration**
- Backend ↔ Frontend API integration
- Frontend ↔ Device communication (via backend)
- Real-time updates (WebSocket)
- Offline-online sync testing

**Integration Tests**:
- User creation to device control
- Schedule creation to device execution
- Data flow (device → analytics page)
- Command execution end-to-end
- Error scenarios and recovery

**Weeks 16-17: Performance & Security**
- Load testing (1000 concurrent users target)
- Security penetration testing
- Database performance tuning
- API optimization
- Infrastructure auto-scaling configuration

**Weeks 18: UAT (User Acceptance Testing)**
- Testing with farm operators (1-2 farms)
- Feedback collection
- Bug fixes

**Deliverables**:
- Integrated system (all components working together)
- Load test report (1000 concurrent users)
- Security audit report (completed)
- Bug list prioritized and tracked

**Success Criteria**:
- No critical bugs
- System handles 1000 concurrent users
- API response time < 200ms under load
- No data loss scenarios
- Security audit passed

---

### Phase 6: AI/ML Integration (Weeks 14-18)
**Duration**: 5 weeks (overlaps with integration)  
**Team**: AI/ML Specialist (with team support)  

**Optional Phase**: Chicken disease detection  

**Activities**:
- Train/validate disease detection model
- API endpoint for image upload/analysis
- Model deployment to production
- Performance optimization
- Documentation

**Deliverables**:
- Functional disease detection API
- Model accuracy > 85%
- Response time < 5 seconds

---

### Phase 7: Staging & Pre-Production (Weeks 19-20)
**Duration**: 2 weeks  

**Activities**:
- Deploy to staging environment (full production replica)
- Smoke testing
- Migration testing (data from v1.0)
- Backup/restore procedures
- Disaster recovery drill
- Performance baseline establishment
- Monitoring dashboards setup

**Deliverables**:
- Staging environment identical to production
- Migration procedures tested
- Runbooks documented
- 100% monitoring coverage

**Success Criteria**:
- Staging identical to production
- Migration completes successfully
- Monitoring shows all systems healthy

---

### Phase 8: Production Deployment (Weeks 21-22)
**Duration**: 2 weeks  

**Week 21: Blue-Green Deployment**
- Deploy v2.0 to parallel infrastructure
- Run both v1.0 and v2.0 simultaneously
- Test in production with internal users
- Verify no critical issues

**Week 22: Launch**
- Production cutover (v1.0 → v2.0)
- Data migration
- User communication
- 24/7 monitoring and support
- Hotfix team on standby

**Deliverables**:
- v2.0 in production
- All users migrated
- V1.0 as fallback (72 hours)

**Success Criteria**:
- Zero data loss
- < 1 hour downtime
- All users can login
- APIs responding < 200ms
- Error rate < 0.1%

---

### Phase 9: Post-Launch Support (Weeks 23-26)
**Duration**: 4 weeks  

**Activities**:
- Monitor production metrics (real-time dashboards)
- Handle user issues and support tickets
- Performance optimization (if needed)
- Bug fixes and patches
- User training and onboarding

**Deliverables**:
- Production system running stably
- User feedback collected
- Optimization recommendations

**Success Criteria**:
- Uptime > 99%
- Error rate < 0.05%
- User satisfaction > 80%

---

### Phase 10: Post-Release Improvements (Weeks 27-35)
**Duration**: 9 weeks  

**Planned Improvements**:
- Additional device types and sensors
- Mobile app (React Native, optional)
- Weather API integration
- Multi-farm dashboard
- Advanced analytics
- Compliance certifications (if needed)

**Ongoing Activities**:
- Feature development for v2.1
- Security updates
- Performance improvements
- User feedback implementation

---

## Timeline Gantt Chart

```
Phase 1: Planning & Arch        [██] Weeks 1-2
Phase 2: Backend Dev            [████████] Weeks 3-10
Phase 3: Frontend Dev           [████████] Weeks 7-14 (overlaps)
Phase 4: Embedded Dev           [███████] Weeks 8-14 (overlaps)
Phase 5: Integration & Test     [██████] Weeks 13-18 (overlaps)
Phase 6: AI/ML (optional)       [█████] Weeks 14-18 (overlaps)
Phase 7: Staging & Pre-Prod     [██] Weeks 19-20
Phase 8: Production Deploy      [██] Weeks 21-22
Phase 9: Post-Launch Support    [████] Weeks 23-26
Phase 10: Improvements          [█████████] Weeks 27-35

        Month 1         Month 2        Month 3        Month 4       Month 5       Month 6
    |----------|    |----------|    |----------|    |----------|    |----------|    |----------|
    1          5   10          15   20          25   30          35
```

---

## Resource Allocation

### Team Capacity Planning

**Phase 1 (Weeks 1-2)**:
- Tech Lead: 100% (daily standup, architecture reviews)
- Backend Lead (Heng): 50% (environment setup)
- Frontend Lead (Rika): 50% (figma design)
- Embedded (Neath): 50% (hardware planning)

**Phase 2-3 (Weeks 3-14)**:
- Backend: 100% (2 devs)
- Frontend: 100% (1 dev)
- Embedded: 100% (1 dev)
- DevOps: 50% (infrastructure support)
- QA: 50% (test planning)

**Phase 5 (Weeks 13-18)**:
- QA: 100% (integration testing, bug tracking)
- All devs: 80% (bug fixes, optimization)

**Phase 7-8 (Weeks 19-22)**:
- DevOps: 100% (deployment)
- All devs: 50% (standby for hotfixes)

**Phase 9 (Weeks 23-26)**:
- Support: 100% (user issues)
- Dev: 50% (critical bugs)

---

## Key Milestones & Gates

| Milestone | Week | Criteria |
|-----------|------|----------|
| **Kickoff Approved** | 1 | Team assembled, plan reviewed |
| **Dev Env Ready** | 2 | All devs can build and deploy |
| **Authentication Done** | 4 | Login/JWT working, tested |
| **Device API Working** | 6 | Devices can receive commands |
| **Frontend MVP** | 10 | Home, Control, Settings pages done |
| **End-to-End Test** | 14 | Device → API → Frontend working |
| **Performance OK** | 16 | Load test passed, target metrics met |
| **Security Audit** | 17 | Passed with no critical issues |
| **Staging Ready** | 20 | Production replica ready |
| **Data Migration OK** | 21 | V1 data successfully migrated |
| **Go-Live** | 22 | Production deployment complete |
| **Stability** | 26 | 4 weeks stable operation |

---

## Risk Timeline

### High-Risk Periods

**Weeks 6-8**: Backend and Frontend diverge in API contracts
- Mitigation: Weekly API review meetings

**Weeks 13-15**: Integration issues discovered late
- Mitigation: Early integration testing, API mocking

**Week 21-22**: Production cutover risks
- Mitigation: Blue-green deployment, rollback plan

**Week 26-35**: Feature scope creep
- Mitigation: Strict feature freeze, v2.1 planning

---

## Buffer & Contingency

**Planned Timeline**: 27 weeks (ideal scenario)  
**Realistic Timeline**: 30-35 weeks  
**Contingency Buffer**: 2-3 weeks for unexpected issues  

**Scenarios Extending Timeline**:
- Major security issues (2-3 weeks)
- Database performance problems (1-2 weeks)
- Scope creep (5-10 weeks)
- Key team member departure (3-5 weeks)

---

## Success Metrics

**On-Time Completion**:
- Within 35-week target
- All planned features delivered
- No critical bugs

**Quality Metrics**:
- > 70% unit test coverage
- Zero critical security issues
- Uptime > 99% in production
- User satisfaction > 80%

**Performance Metrics**:
- API response < 200ms (95th percentile)
- Bundle size < 150KB
- Load time < 2 seconds on 4G
- Support < 2 hours response time

---

**Version History**
| Version | Date | Changes |
|---------|------|---------|
| 2.0 | Feb 2026 | Initial project timeline |

**Related Documents**
- TEAM_STRUCTURE.md
- RISK_MANAGEMENT.md
- SPECIFICATIONS_REQUIREMENTS.md
