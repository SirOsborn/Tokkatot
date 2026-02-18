# Tokkatot 2.0: Complete Technical Specifications Index

**Version**: 2.0 (Production Release)  
**Release Date**: Q2 2026  
**Status**: Development Phase  
**Last Updated**: February 2026

---

## 📋 Documentation Structure

This documentation has been organized into focused specification files for clarity and ease of navigation. Each file contains detailed requirements and specifications for a specific topic.

### Quick Navigation Guide

**START HERE:**
- 👉 **[SPECIFICATIONS_ARCHITECTURE.md](SPECIFICATIONS_ARCHITECTURE.md)** - Overall system design, data flow, and architecture overview
- 👉 **[SPECIFICATIONS_REQUIREMENTS.md](SPECIFICATIONS_REQUIREMENTS.md)** - Functional and non-functional requirements

**IMPLEMENTATION GUIDES:**
- **[SPECIFICATIONS_DATABASE.md](SPECIFICATIONS_DATABASE.md)** - Database schema, tables, normalization, relationships
- **[SPECIFICATIONS_API.md](SPECIFICATIONS_API.md)** - Backend API endpoints (30+), request/response formats, authentication
- **[SPECIFICATIONS_FRONTEND.md](SPECIFICATIONS_FRONTEND.md)** - UI/UX specifications, responsive design, accessibility, performance targets
- **[SPECIFICATIONS_EMBEDDED.md](SPECIFICATIONS_EMBEDDED.md)** - ESP32 firmware architecture, local Raspberry Pi agent, device protocols
- **[SPECIFICATIONS_DEPLOYMENT.md](SPECIFICATIONS_DEPLOYMENT.md)** - Cloud infrastructure, Docker, CI/CD, deployment procedures
- **[SPECIFICATIONS_SECURITY.md](SPECIFICATIONS_SECURITY.md)** - Authentication, encryption, RBAC, security architecture

**OPERATIONAL GUIDES:**
- **[SPECIFICATIONS_DATA_LOGGING.md](SPECIFICATIONS_DATA_LOGGING.md)** - Logging strategy, data retention, monitoring, analytics
- **[SPECIFICATIONS_TECHNOLOGY_STACK.md](SPECIFICATIONS_TECHNOLOGY_STACK.md)** - Technology selections, rationale, version requirements
- **[TEAM_STRUCTURE.md](TEAM_STRUCTURE.md)** - Team responsibilities, role assignments, handoff procedures
- **[PROJECT_TIMELINE.md](PROJECT_TIMELINE.md)** - Development phases, milestones, timeline (27-35 weeks)
- **[RISK_MANAGEMENT.md](RISK_MANAGEMENT.md)** - Risk analysis, mitigation strategies, contingency plans

**FARMER-CENTRIC FEATURES:**
- **[TOKKATOT_2.0_FARMER_CENTRIC_ADDITIONS.md](TOKKATOT_2.0_FARMER_CENTRIC_ADDITIONS.md)** - Remote updates, sync strategy, performance optimization for low-end phones, accessibility for elderly farmers, multilingual support

---

## 🎯 Reading Paths for Different Roles

### 👨‍💼 Project Manager / Team Lead
1. Read: SPECIFICATIONS_REQUIREMENTS.md (understand what will be built)
2. Read: PROJECT_TIMELINE.md (understand schedule)
3. Read: TEAM_STRUCTURE.md (understand roles)
4. Read: RISK_MANAGEMENT.md (understand obstacles)

### 🏗️ System Architect / Tech Lead
1. Read: SPECIFICATIONS_ARCHITECTURE.md (understand system design)
2. Read: SPECIFICATIONS_TECHNOLOGY_STACK.md (understand tech choices)
3. Read: SPECIFICATIONS_DEPLOYMENT.md (understand infrastructure)
4. Read: All other SPECIFICATIONS_*.md files (understand details)

### 💻 Backend Developer
1. Read: SPECIFICATIONS_REQUIREMENTS.md (understand scope)
2. Read: SPECIFICATIONS_API.md (understand endpoints and contracts)
3. Read: SPECIFICATIONS_DATABASE.md (understand data structure)
4. Read: SPECIFICATIONS_SECURITY.md (understand auth/encryption)
5. Read: SPECIFICATIONS_DEPLOYMENT.md (understand deployment)

### 🎨 Frontend Developer
1. Read: SPECIFICATIONS_REQUIREMENTS.md (understand scope)
2. Read: SPECIFICATIONS_FRONTEND.md (understand UI/UX specs)
3. Read: TOKKATOT_2.0_FARMER_CENTRIC_ADDITIONS.md (understand farmer needs, accessibility)
4. Read: SPECIFICATIONS_API.md (understand backend contracts)
5. Read: SPECIFICATIONS_SECURITY.md (understand authentication)

### 🔧 Embedded Systems Developer (Neath)
1. Read: SPECIFICATIONS_EMBEDDED.md (understand firmware architecture)
2. Read: SPECIFICATIONS_REQUIREMENTS.md (understand device requirements)
3. Read: SPECIFICATIONS_API.md (understand cloud communication)
4. Read: TOKKATOT_2.0_FARMER_CENTRIC_ADDITIONS.md (understand OTA update requirements)

### 🔐 DevOps / Infrastructure Engineer
1. Read: SPECIFICATIONS_DEPLOYMENT.md (understand infrastructure)
2. Read: SPECIFICATIONS_SECURITY.md (understand security requirements)
3. Read: SPECIFICATIONS_DATA_LOGGING.md (understand monitoring)
4. Read: SPECIFICATIONS_TECHNOLOGY_STACK.md (understand tech requirements)

---

## 📌 Key Project Information

### Project Context
- **Client**: Smart Chicken Farming (local Cambodia implementation)
- **Target Users**: Elderly farmers with low digital literacy, using 1-2GB RAM smartphones, 4G networks
- **Geographic Region**: Cambodia (Khmer primary language, English secondary)
- **Production Target**: Q2 2026

### System Overview
- **Architecture**: 3-tier (Client / Backend API / Data Layer) + Edge Computing
- **Cloud Provider**: DigitalOcean (chosen for cost and simplicity)
- **Local Fallback**: Offline mode on Raspberry Pi 4B with MQTT queue
- **Devices**: ESP32-based controllers for water, feeder, light, fan, heater, conveyor systems
- **Real-Time Communication**: MQTT for devices, Socket.io for web app

### Critical Success Factors
✅ **Reliability**: App never crashes, works with or without internet  
✅ **Performance**: < 2 second load time on 4G networks, < 150KB app bundle  
✅ **Usability**: Large fonts (48px), high contrast (WCAG AAA), simple navigation  
✅ **Synchronization**: 3-level conflict resolution between app and devices  
✅ **Remote Updates**: Firmware patches without visiting farms (OTA)  
✅ **Multilingual**: Khmer + English with seamless toggle  

---

## 📖 How to Use This Documentation

### For Reading
1. **Start with Architecture** - Understand "what" and "how" at system level
2. **Read Requirements** - Understand "what needs to be built"
3. **Read Specific Specs** - Dive into technology/component details
4. **Cross-reference** - Each file contains links to related specifications

### For Implementation
1. **Identify your component** - Find relevant spec file above
2. **Read the specification** - Understand requirements and constraints
3. **Review constraints** - Check performance targets, security requirements, compatibility
4. **Check integration points** - See how your component connects to others
5. **Reference related specs** - Understand dependencies and contracts

### For Updates
When requirements change:
1. Find the relevant specification file
2. Update only that file (don't update the index unless creating new spec)
3. Increment version in that file's header
4. Add changelog entry at bottom of file
5. Notify team leads of changes via SPECIFICATIONS_INDEX.md update message

---

## 🔗 File Dependency Map

```
SPECIFICATIONS_INDEX (this file)
├── SPECIFICATIONS_ARCHITECTURE
│   ├── SPECIFICATIONS_REQUIREMENTS
│   ├── SPECIFICATIONS_API
│   └── Integration diagrams reference all specs
├── SPECIFICATIONS_REQUIREMENTS
│   └── All other specs implement these requirements
├── SPECIFICATIONS_DATABASE
│   └── Used by: API, Backend, Security
├── SPECIFICATIONS_API
│   ├── Uses: Database schema, Security auth
│   └── Referenced by: Frontend, Deployment
├── SPECIFICATIONS_FRONTEND
│   ├── Uses: API endpoints, Security auth
│   └── Farmer-centric additions
├── SPECIFICATIONS_EMBEDDED
│   ├── Uses: API endpoints, Security
│   └── Deployment requirements
├── SPECIFICATIONS_DEPLOYMENT
│   ├── Uses: All frontend/backend/embedded specs
│   └── Deployment infrastructure guide
├── SPECIFICATIONS_SECURITY
│   ├── Used by: API, Frontend, Embedded
│   └── Deployment requirements
├── SPECIFICATIONS_DATA_LOGGING
│   └── Uses: Database design
├── SPECIFICATIONS_TECHNOLOGY_STACK
│   └── Prerequisite for all other specs
└── Supporting docs:
    ├── TEAM_STRUCTURE
    ├── PROJECT_TIMELINE
    ├── RISK_MANAGEMENT
    └── TOKKATOT_2.0_FARMER_CENTRIC_ADDITIONS
```

---

## ✅ Document Checklist

Before starting development, ensure:

- [ ] **All roles have read** their assigned specification documents
- [ ] **Tech Lead approved** the technology stack choices
- [ ] **All APIs reviewed** and endpoints documented
- [ ] **Database schema reviewed** and normalized
- [ ] **Security requirements approved** and audit scheduled
- [ ] **Frontend wireframes reviewed** with farmer accessibility requirements
- [ ] **Embedded firmware architecture reviewed** by Neath
- [ ] **Deployment procedures tested** in staging environment
- [ ] **Team members assigned** per TEAM_STRUCTURE.md
- [ ] **Timeline milestones scheduled** per PROJECT_TIMELINE.md
- [ ] **Risk mitigation plans acknowledged** per RISK_MANAGEMENT.md

---

## 📞 Questions? References

**Structure:**
- All specifications follow: Overview → Requirements → Technical Details → Integration Points → Constraints/Targets
- No code examples provided (implementation is developer responsibility)
- All specifications are technology/language agnostic
- Each file can be read independently or as part of the overall specification

**File Format:**
- Markdown (.md) for readability and version control
- Organized hierarchically with clear sections
- Cross-referenced between related documents
- Updated incrementally as requirements evolve

**Questions about:**
- **"How do I get started?"** → Start with SPECIFICATIONS_ARCHITECTURE.md
- **"Where do I implement [feature X]?"** → Check SPECIFICATIONS_REQUIREMENTS.md, then navigate to specific section
- **"What are the constraints for my component?"** → Find your component spec file, check Constraints/Performance Targets section
- **"How does [X] connect to [Y]?"** → Check Integration Points section in relevant files or Architecture file

---

## 📝 Version History

| Version | Date | Changes |
|---------|------|---------|
| 2.0 | Feb 2026 | Initial production specification suite |
| | | Reorganized from single 3K+ line doc into modular files |
| | | Removed Option 3 (self-hosted option) |
| | | Added farmer-centric specifications |

---

**Last Updated**: February 18, 2026  
**Status**: Ready for Development Team Review
