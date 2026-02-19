# Tokkatot 2.0: Complete Technical Specifications Index

**Version**: 2.0 (Production Release)  
**Release Date**: Q2 2026  
**Status**: Development Phase  
**Last Updated**: February 2026

---

## 📋 Documentation Structure

This documentation is organized into three layers for clarity:
1. **Core Specifications (00_)** - Start here, read in order
2. **Implementation Guides (IG_)** - Technical implementation details
3. **Operational Guides (OG_)** - Team processes and infrastructure

Each file contains detailed requirements and can be read independently or as part of the complete specification.

### Quick Navigation Guide

**🎯 START HERE (Read in Order):**
1. 👉 **[01_SPECIFICATIONS_ARCHITECTURE.md](01_SPECIFICATIONS_ARCHITECTURE.md)** - Overall system design, data flow, and architecture overview
2. 👉 **[02_SPECIFICATIONS_REQUIREMENTS.md](02_SPECIFICATIONS_REQUIREMENTS.md)** - Functional and non-functional requirements (farmer-centric approach)

**🏗️ IMPLEMENTATION GUIDES (IG_*):**
- **[IG_SPECIFICATIONS_DATABASE.md](IG_SPECIFICATIONS_DATABASE.md)** - Database schema, tables, normalization, relationships
- **[IG_SPECIFICATIONS_API.md](IG_SPECIFICATIONS_API.md)** - Backend API endpoints (63 total), simplified for farmer usage
- **[IG_SPECIFICATIONS_FRONTEND.md](IG_SPECIFICATIONS_FRONTEND.md)** - UI/UX for low-literacy farmers, large fonts, high contrast, Khmer/English
- **[IG_SPECIFICATIONS_EMBEDDED.md](IG_SPECIFICATIONS_EMBEDDED.md)** - ESP32 device firmware, Raspberry Pi agent, device setup
- **[IG_SPECIFICATIONS_SECURITY.md](IG_SPECIFICATIONS_SECURITY.md)** - Authentication, encryption, simplified role system
- **[IG_TOKKATOT_2.0_FARMER_CENTRIC_SPECIFICATIONS.md](IG_TOKKATOT_2.0_FARMER_CENTRIC_SPECIFICATIONS.md)** - Farmer accessibility, multilingual, device setup by team

**⚙️ OPERATIONAL GUIDES (OG_*):**
- **[OG_SPECIFICATIONS_TECHNOLOGY_STACK.md](OG_SPECIFICATIONS_TECHNOLOGY_STACK.md)** - Technology selections, version requirements
- **[OG_SPECIFICATIONS_DEPLOYMENT.md](OG_SPECIFICATIONS_DEPLOYMENT.md)** - Cloud infrastructure, Docker, CI/CD pipelines
- **[OG_TEAM_STRUCTURE.md](OG_TEAM_STRUCTURE.md)** - Team responsibilities, role assignments, handoff procedures
- **[OG_PROJECT_TIMELINE.md](OG_PROJECT_TIMELINE.md)** - Development phases, milestones, timeline
- **[OG_RISK_MANAGEMENT.md](OG_RISK_MANAGEMENT.md)** - Risk analysis, mitigation strategies, contingency plans

---

## 🎯 Reading Paths for Different Roles

### 👨‍💼 Project Manager / Team Lead
1. Read: 02_SPECIFICATIONS_REQUIREMENTS.md (understand what will be built)
2. Read: OG_PROJECT_TIMELINE.md (understand schedule)
3. Read: OG_TEAM_STRUCTURE.md (understand roles)
4. Read: OG_RISK_MANAGEMENT.md (understand obstacles)

### 🏗️ System Architect / Tech Lead
1. Read: 01_SPECIFICATIONS_ARCHITECTURE.md (understand system design)
2. Read: OG_SPECIFICATIONS_TECHNOLOGY_STACK.md (understand tech choices)
3. Read: OG_SPECIFICATIONS_DEPLOYMENT.md (understand infrastructure)
4. Read: All IG_SPECIFICATIONS_*.md files (understand implementation details)

### 💻 Backend Developer
1. Read: 02_SPECIFICATIONS_REQUIREMENTS.md (understand scope)
2. Read: IG_SPECIFICATIONS_API.md (understand endpoints and simplified role system)
3. Read: IG_SPECIFICATIONS_DATABASE.md (understand data structure)
4. Read: IG_SPECIFICATIONS_SECURITY.md (understand auth/encryption)
5. Read: OG_SPECIFICATIONS_DEPLOYMENT.md (understand deployment)

### 🎨 Frontend Developer
1. Read: 02_SPECIFICATIONS_REQUIREMENTS.md (understand scope)
2. Read: IG_SPECIFICATIONS_FRONTEND.md (understand UI/UX for farmers)
3. Read: IG_TOKKATOT_2.0_FARMER_CENTRIC_SPECIFICATIONS.md (understand farmer accessibility needs)
4. Read: IG_SPECIFICATIONS_API.md (understand backend contracts)
5. Read: IG_SPECIFICATIONS_SECURITY.md (understand authentication)

### 🔧 Embedded Systems Developer
1. Read: IG_SPECIFICATIONS_EMBEDDED.md (understand firmware architecture)
2. Read: 02_SPECIFICATIONS_REQUIREMENTS.md (understand device requirements)
3. Read: IG_SPECIFICATIONS_API.md (understand cloud communication)
4. Read: IG_TOKKATOT_2.0_FARMER_CENTRIC_SPECIFICATIONS.md (understand device setup by team, OTA updates)

### 🔐 DevOps / Infrastructure Engineer
1. Read: OG_SPECIFICATIONS_DEPLOYMENT.md (understand infrastructure)
2. Read: IG_SPECIFICATIONS_SECURITY.md (understand security requirements)
3. Read: OG_SPECIFICATIONS_TECHNOLOGY_STACK.md (understand tech requirements)
4. Read: 01_SPECIFICATIONS_ARCHITECTURE.md (understand system overview)

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
00_SPECIFICATIONS_INDEX (this file)
├── 01_SPECIFICATIONS_ARCHITECTURE
│   ├── 02_SPECIFICATIONS_REQUIREMENTS
│   ├── IG_SPECIFICATIONS_API
│   └── Integration diagrams reference all specs
├── 02_SPECIFICATIONS_REQUIREMENTS
│   └── All IG_* and OG_* specs implement these requirements
├── IG_SPECIFICATIONS_DATABASE
│   └── Used by: IG_SPECIFICATIONS_API, Backend, IG_SPECIFICATIONS_SECURITY
├── IG_SPECIFICATIONS_API
│   ├── Uses: IG_SPECIFICATIONS_DATABASE schema, IG_SPECIFICATIONS_SECURITY auth
│   └── Referenced by: IG_SPECIFICATIONS_FRONTEND, OG_SPECIFICATIONS_DEPLOYMENT
├── IG_SPECIFICATIONS_FRONTEND
│   ├── Uses: IG_SPECIFICATIONS_API endpoints, IG_SPECIFICATIONS_SECURITY auth
│   ├── IG_TOKKATOT_2.0_FARMER_CENTRIC_SPECIFICATIONS (farmer accessibility)
│   └── Large fonts, high contrast, simple navigation (48px+, WCAG AAA)
├── IG_SPECIFICATIONS_EMBEDDED
│   ├── Uses: IG_SPECIFICATIONS_API endpoints, IG_SPECIFICATIONS_SECURITY
│   ├── Device setup by Tokkatot team only
│   └── OTA update requirements
├── IG_SPECIFICATIONS_SECURITY
│   ├── Simplified role system for farmers
│   └── Used by: IG_SPECIFICATIONS_API, IG_SPECIFICATIONS_FRONTEND, IG_SPECIFICATIONS_EMBEDDED
├── IG_TOKKATOT_2.0_FARMER_CENTRIC_SPECIFICATIONS
│   ├── Phone/Email registration
│   ├── Farmer accessibility features
│   └── Multilingual (Khmer + English)
├── OG_SPECIFICATIONS_DEPLOYMENT
│   ├── Uses: All IG_*/02_* specs
│   └── Infrastructure guide
├── OG_SPECIFICATIONS_TECHNOLOGY_STACK
│   └── Prerequisites for all other specs
└── Supporting operational docs:
    ├── OG_TEAM_STRUCTURE
    ├── OG_PROJECT_TIMELINE
    └── OG_RISK_MANAGEMENT
```

---

## ✅ Document Checklist

Before starting development, ensure:

- [ ] **All roles have read** their assigned specification documents (use reading paths above)
- [ ] **Farmer-centric approach understood**: Phone/Email registration, simplified roles, device setup by team
- [ ] **Tech Lead approved** the technology stack choices
- [ ] **All APIs reviewed** with simplified role system (not complex RBAC)
- [ ] **Database schema reviewed** with simplified role system
- [ ] **Security requirements approved** with farmer-friendly authentication
- [ ] **Frontend wireframes reviewed** with farmer accessibility: 48px+ fonts, WCAG AAA contrast, Khmer/English
- [ ] **Embedded device setup process** designed for team installation only
- [ ] **Deployment procedures tested** in staging environment
- [ ] **Team members assigned** per OG_TEAM_STRUCTURE.md
- [ ] **Timeline milestones scheduled** per OG_PROJECT_TIMELINE.md
- [ ] **Risk mitigation plans acknowledged** per OG_RISK_MANAGEMENT.md

---

## 🌾 Farmer-Centric Design Principles

**For Elderly Farmers with Low Digital Literacy in Cambodia:**

✅ **Registration**: Email OR phone number (not both required)  
✅ **Device Setup**: Tokkatot team installs and configures (farmers don't manage)  
✅ **Roles**: Simplified - just Owner, Manager, Viewer (no complex permissions)  
✅ **Language**: Khmer primary, English secondary (seamless toggle)  
✅ **UI**: 48px+ fonts, WCAG AAA contrast, 5-click max to any feature  
✅ **Performance**: < 2 second load on 4G, < 150KB bundle  
✅ **Offline**: Works without internet (Raspberry Pi local fallback)  
✅ **Support**: Phone support in Khmer 7am-8pm Cambodia time

---

## 📞 Questions? References

**Documentation Organization:**
- **00_** - Core specifications (read in numbered order)
- **IG_** - Implementation Guides (how to build components)
- **OG_** - Operational Guides (how to run the project)

**Structure:**
- All specifications follow: Overview → Requirements → Technical Details → Integration Points → Constraints/Targets
- No code examples provided (implementation is developer responsibility)
- All specifications are technology/language agnostic where possible
- Each file can be read independently or as part of the complete specification

**File Format:**
- Markdown (.md) for readability and version control
- Organized hierarchically with clear sections
- Cross-referenced between related documents
- Updated incrementally as requirements evolve

**Common Questions:**
- **"How do I get started?"** → Start with 01_SPECIFICATIONS_ARCHITECTURE.md
- **"Where do I implement [feature X]?"** → Check 02_SPECIFICATIONS_REQUIREMENTS.md for functional req, then navigate to specific IG_* file
- **"What are the constraints for my component?"** → Find your component IG_* file, check Constraints/Performance Targets section
- **"How does [X] connect to [Y]?"** → Check Integration Points section in relevant IG_* files or 01_SPECIFICATIONS_ARCHITECTURE.md
- **"How do I simplify for farmers?"** → Read IG_TOKKATOT_2.0_FARMER_CENTRIC_SPECIFICATIONS.md

---

## 📝 Version History

| Version | Date | Changes |
|---------|------|---------|
| 2.0-FarmerCentric | Feb 2026 | Simplified architecture for low-literacy farmers |
| | | Renamed files: 00_, IG_, OG_ prefixes for clarity |
| | | Phone/Email registration support |
| | | Device setup by Tokkatot team only |
| | | Simplified role system (Owner/Manager/Viewer) |
| 2.0 | Feb 2026 | Initial production specification suite |
| | | Reorganized from single document into modular files |
| | | Removed Option 3 (self-hosted option) |

---

**Last Updated**: February 18, 2026  
**Status**: Ready for Development Team Review
