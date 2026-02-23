# 🤖 AI Context: Documentation Maintenance

**Directory**: `docs/`  
**Your Role**: Keep specifications in sync with code, maintain architecture documentation  
**File Format**: Markdown (.md)  

**📖 Read First**: `../AI_INSTRUCTIONS.md` (project overview)

---

## 🎯 What You're Maintaining

**Living Documentation** for Tokkatot 2.0 System

Purpose:
- Single source of truth for system design
- Guide for developers implementing features
- Reference for architectural decisions
- Record of farmer-centric requirements

**Golden Rule**: If code exists but spec doesn't, spec is incomplete. If spec exists but code doesn't match, code is wrong.

---

## 📁 Documentation Structure

```
docs/
├── README.md                    Navigation hub START HERE
├── ARCHITECTURE.md              System design, data flows, user journeys
├── TECH_STACK.md                Technology choices (Go, Vue.js, PostgreSQL)
├── AUTOMATION_USE_CASES.md      Real farmer scenarios (schedules) ⭐ NEW
│
├── guides/
│   └── SETUP.md                 Complete setup guide (PostgreSQL, Go, build)
│
├── implementation/              Component specs (read before coding)
│   ├── API.md                   66 REST endpoints, WebSocket
│   ├── DATABASE.md              PostgreSQL schema (10 tables)
│   ├── FRONTEND.md              Vue.js 3 PWA, farmer accessibility
│   ├── AI_SERVICE.md            PyTorch disease detection, FastAPI
│   ├── EMBEDDED.md              ESP32 firmware, MQTT
│   └── SECURITY.md              JWT auth, registration keys
│
└── troubleshooting/             Problem solving
    ├── DATABASE.md              Connection issues, schema fixes
    └── API_TESTING.md           Test endpoints, debug tokens
```

---

## ✅ When to Update Documentation

### After API changes
- **File to update**: `implementation/API.md`
- **What to add**: New endpoint spec, request/response examples, error codes
- **Example**: Added 7 schedule endpoints → Update API.md with full spec

### After database schema changes
- **File to update**: `implementation/DATABASE.md`
- **What to add**: New table CREATE statement, field explanations, indexes, relationships
- **Recent example**: Added `action_duration` and `action_sequence` fields to schedules table

### After adding new features
- **File to update**: `AUTOMATION_USE_CASES.md` (if schedule-related), `ARCHITECTURE.md` (if major feature)
- **What to add**: User story, real-world scenario, JSON examples, farmer benefits
- **Recent example**: Multi-step sequences for pulse feeding → Added to AUTOMATION_USE_CASES.md

### After architecture changes
- **File to update**: `ARCHITECTURE.md`, `TECH_STACK.md`
- **What to add**: Data flow diagrams, component diagrams, design decisions
- **Example**: Changed from microservices to monolith → Update ARCHITECTURE.md

---

## 🚨 MANDATORY: Update Docs After Building Features

**CRITICAL REQUIREMENT**: Documentation is NOT optional. When you complete significant work, you MUST update documentation immediately.

### Why This Matters

**Without documentation updates**:
- ❌ Future AI sessions rediscover features by reading code (slow, error-prone)
- ❌ Knowledge is lost between sessions (no institutional memory)
- ❌ New developers can't understand system without digging through implementation
- ❌ Features become "forgotten" and get re-implemented differently

**With documentation updates**:
- ✅ Future AI sessions read specs and understand instantly
- ✅ Knowledge persists across sessions (institutional memory)
- ✅ New developers read docs and start contributing immediately
- ✅ Features are discoverable and consistent

### When You MUST Update Docs

**Update documentation when you complete**:
- ✅ New feature (e.g., pulse feeding automation, disease detection API)
- ✅ Database schema change (new table, new fields like `action_sequence`)
- ✅ New API endpoint or modified endpoint behavior
- ✅ New automation pattern (schedule types, sensor triggers)
- ✅ New UI component (schedule builder, device control panel)
- ✅ Architecture decision (SQLite fallback, JWT auth flow)
- ✅ Integration work (Go → FastAPI → PyTorch, MQTT protocol)

**Don't update for**:
- ❌ Minor bug fixes (unless they reveal missing documentation)
- ❌ Code refactoring without functional changes
- ❌ Variable/function renames
- ❌ Comment additions, formatting/linting changes

### Which Files You MUST Update

**Implementation specs** (`docs/implementation/*.md`) - Update when code changes:

| File | Update When | What to Add |
|------|-------------|-------------|
| `API.md` | New/modified endpoint | Full request/response examples, error codes |
| `DATABASE.md` | Schema changed | CREATE statement, explain all fields |
| `FRONTEND.md` | New UI component | Vue.js code, user flow, screenshots |
| `AI_SERVICE.md` | Model/API changed | PyTorch architecture, FastAPI endpoints |
| `EMBEDDED.md` | Firmware behavior changed | C code, MQTT topics, GPIO pins |

**AI knowledge files** - Update to teach future AI sessions:

| File | Update When | What to Add |
|------|-------------|-------------|
| `middleware/AI_CONTEXT.md` | Added Go pattern | Code examples, function signatures, gotchas |
| `frontend/AI_CONTEXT.md` | Added Vue.js component | Component structure, API calls, styling |
| `ai-service/AI_CONTEXT.md` | Changed model/preprocessing | PyTorch code, data flow, performance notes |
| `embedded/AI_CONTEXT.md` | Added MQTT topic/GPIO | C code, protocol spec, safety checks |
| `AI_INSTRUCTIONS.md` | Added major system concept | High-level overview, business context |

**Use case docs** - Update when solving farmer problems:

| File | Update When | What to Add |
|------|-------------|-------------|
| `AUTOMATION_USE_CASES.md` | Solved real farmer scenario | Complete scenario, JSON example, benefits |

### Update Timing (Goldilocks Rule)

**⏱️ TOO FAST** (don't do this):
- ❌ Updating after every single function (creates noise)
- ❌ Updating after variable renames (not significant)

**⏱️ TOO SLOW** (don't do this):
- ❌ Waiting weeks for "perfect time" (knowledge gets lost)
- ❌ Building 10 features before documenting any (overwhelming)

**⏱️ JUST RIGHT** (do this):
- ✅ Update after 30-60 minutes of significant work
- ✅ Update when feature is complete (end of session)
- ✅ Update before switching to different component
- ✅ Update after 3-5 related changes (consolidate)

### Mandatory Update Checklist

**Before ending your session, verify**:

```markdown
[ ] Did I change database schema? → MUST update docs/implementation/DATABASE.md
[ ] Did I add/modify API endpoints? → MUST update docs/implementation/API.md  
[ ] Did I add UI components? → MUST update docs/implementation/FRONTEND.md
[ ] Did I change firmware behavior? → MUST update docs/implementation/EMBEDDED.md
[ ] Did I solve a farmer problem? → MUST update docs/AUTOMATION_USE_CASES.md
[ ] Did I add reusable patterns? → MUST update component AI_CONTEXT.md files
[ ] Did I test all examples? → Verify JSON/SQL/code compiles and runs
[ ] Did I add cross-references? → Link related docs together
```

**This is how we maintain institutional knowledge - it's not optional!**

### Example: action_sequence Feature Documentation

**What was built**: Multi-step automation for pulse feeding (ON 30s, pause 10s, repeat)

**Documentation updates REQUIRED** (completed same session):
1. ✅ `docs/implementation/DATABASE.md` - Added `action_sequence JSONB` field spec with schema
2. ✅ `docs/implementation/API.md` - Updated 4 schedule endpoints with field examples
3. ✅ `docs/implementation/FRONTEND.md` - Added Action Sequence Builder UI (300+ lines)
4. ✅ `docs/implementation/EMBEDDED.md` - Added ESP32 execution code (200+ lines)
5. ✅ `docs/AUTOMATION_USE_CASES.md` - Created 500+ line guide with 3 farmer scenarios
6. ✅ `middleware/AI_CONTEXT.md` - Added schedule automation section with Go code
7. ✅ `AI_INSTRUCTIONS.md` - Added automation & schedules overview section

**Result**: Future AI sessions immediately understand this feature exists, how it works, why farmers need it, and how to build on it.

**If not documented**: Future AI sessions would never discover this feature or would reinvent it differently, wasting time and creating inconsistency.

---

## 📁 Documentation Structure
- **Trigger**: Completed multi-step schedule automation
- **Updates made**:
  1. ✅ `DATABASE.md` - Added action_sequence field schema
  2. ✅ `API.md` - Added field to schedule endpoints
  3. ✅ `AUTOMATION_USE_CASES.md` - Created 500+ line guide with farmer scenarios
  4. ✅ `middleware/AI_CONTEXT.md` - Added schedule automation section
  5. ✅ `AI_INSTRUCTIONS.md` - Added automation & schedules section
- **Result**: Future AI knows about pulse feeding, conveyor belt sequences

**Example 2**: AI Documentation Consolidation (Feb 2026)
- **Trigger**: User noticed overlapping AI docs, requested cleanup
- **Updates made**:
  1. ✅ Rewrote all 5 `AI_CONTEXT.md` files (removed duplication)
  2. ✅ Added cross-references to `docs/` folder
  3. ✅ Updated `AI_INSTRUCTIONS.md` with AI context file index
  4. ✅ Updated `docs/AI_CONTEXT.md` with self-documentation guidance
- **Result**: Clear doc hierarchy, no overlap, future AI knows where to find info

---

## 🔄 Documentation Update Workflow

1. **Code change made** (e.g., added `action_sequence` field to schedules)
2. **Identify affected docs**:
   - Database schema change → `implementation/DATABASE.md`
   - New API field → `implementation/API.md`
   - Farmer use case → `AUTOMATION_USE_CASES.md`
3. **Update all affected files** (don't skip any)
4. **Add real-world examples** (JSON, cron expressions, SQL queries)
5. **Cross-reference** (e.g., "See AUTOMATION_USE_CASES.md for examples")

---

## 📝 Documentation Standards

### Writing Style
- **Farmer-first language**: "Turn feeder ON at 6AM for 15 minutes" (not "Execute scheduled job")
- **Concrete examples**: Show actual JSON, not abstract schemas
- **Visual hierarchy**: Use tables, lists, code blocks generously
- **Cross-references**: Link to related docs

### Code Examples
- **Include full context**: Not just the new field, but the whole request/response
- **Show realistic data**: Use farmer names (Sokha), phone numbers (012345678), real cron expressions
- **Comment complex logic**: Explain WHY, not just WHAT

### Recent Updates (Feb 2026)
- ✅ Created `AUTOMATION_USE_CASES.md` - 500+ lines of real farmer scenarios
- ✅ Updated `DATABASE.md` - schedules table with new fields
- ✅ Updated `README.md` - added AUTOMATION_USE_CASES.md to navigation
- ✅ Updated `AI_INSTRUCTIONS.md` - added schedule automation section

---

## 🚫 What NOT to Include

**Don't duplicate**:
- ❌ Business model details (already in `/AI_INSTRUCTIONS.md`)
- ❌ Tech stack rationale (already in `TECH_STACK.md`)
- ❌ Setup instructions (already in `guides/SETUP.md`)
- ❌ Component-specific patterns (already in `middleware/AI_CONTEXT.md`, `frontend/AI_CONTEXT.md`, etc.)

**Keep docs/ focused on**:
- ✅ System architecture
- ✅ API specifications
- ✅ Database schemas
- ✅ User journeys
- ✅ Real farmer use cases

---

## 📚 Component AI Context Files

**These exist in each service folder** (don't duplicate in docs/):
- `middleware/AI_CONTEXT.md` - Go API specifics, file structure, common patterns
- `frontend/AI_CONTEXT.md` - Vue.js components, UI patterns, accessibility
- `ai-service/AI_CONTEXT.md` - PyTorch model, FastAPI endpoints, training
- `embedded/AI_CONTEXT.md` - ESP32 firmware, MQTT protocol, sensor drivers

**docs/AI_CONTEXT.md** (this file) only covers documentation maintenance, not implementation.

---

## ⚠️ Critical Rules

1. **Update docs IMMEDIATELY after code changes** - don't let them drift
2. **Test examples before documenting** - all JSON should be valid and tested
3. **Keep examples real** - use farmer scenarios, not abstract theoretical cases
4. **Cross-reference extensively** - help readers navigate between related docs
5. **Version history** - add "Updated: [date]" at top when making changes

**End of docs/AI_CONTEXT.md**
