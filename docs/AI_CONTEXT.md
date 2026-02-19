# 🤖 AI Context: Documentation & Specifications

**Directory**: `docs/`  
**Your Role**: Keep all specifications in sync with code, maintain architecture, document requirements  
**File Format**: Markdown (.md)  

---

## 🎯 What You're Maintaining

**Living Documentation** for Tokkatot 2.0 System

Purpose:
- Single source of truth for system design
- Guide for developers implementing features
- Reference for architectural decisions
- Record of farmer-centric requirements

Golden Rule: **If code exists but spec doesn't, spec is incomplete. If spec exists but code doesn't match, code is wrong.**

---

## 📁 File Hierarchy (Read in this order)

### Core Specifications (00_, 01_, 02_)

**Start here** - Read these before any implementation:

1. **`00_SPECIFICATIONS_INDEX.md`**
   - Overview of Tokkatot project
   - How to navigate other docs
   - Quick reference of all 14 spec files
   - Farmer-centric principles

2. **`01_SPECIFICATIONS_ARCHITECTURE.md`**
   - System design (3-tier: Client/API/Data)
   - Microservices descriptions
   - Data flow patterns (5 patterns: device control, sensors, schedules, sync, AI)
   - Integration points with Go/Python/JavaScript/C

3. **`02_SPECIFICATIONS_REQUIREMENTS.md`**
   - Functional requirements (FR1-FR2)
   - Non-functional requirements (performance, security, compatibility)
   - User stories for farmers
   - Acceptance criteria

### Implementation Guides (IG_*)

**How to build components** - Read these per component:

- **`IG_SPECIFICATIONS_API.md`** (66 endpoints)
  - Authentication (8 endpoints)
  - User management (5)
  - Farm management (8)
  - Device management (10)
  - Device control (8)
  - Scheduling (7)
  - Monitoring & alerts (8)
  - Analytics & reporting (5)
  - AI endpoints (3)
  - WebSocket, error handling, rate limiting

- **`IG_SPECIFICATIONS_DATABASE.md`**
  - PostgreSQL schema (13 tables)
  - Table relationships, constraints
  - Data types, indexes
  - Migration strategy

- **`IG_SPECIFICATIONS_SECURITY.md`**
  - Authentication flow (Email/Phone JWT)
  - Authorization (3-role system)
  - Password hashing, token expiry
  - Input validation, SQL injection prevention
  - No MFA for farmers (optional for admins)

- **`IG_SPECIFICATIONS_FRONTEND.md`**
  - Vue.js 3 UI specifications
  - Page layouts (dashboard, disease detection, profile, settings)
  - Accessibility (48px+ fonts, WCAG AAA, Khmer/English)
  - Mobile-first responsiveness
  - WebSocket integration

- **`IG_SPECIFICATIONS_EMBEDDED.md`**
  - ESP32 firmware architecture
  - GPIO pinout, MQTT topics
  - Sensor reading intervals
  - OTA update mechanism
  - Local Raspberry Pi agent

- **`IG_SPECIFICATIONS_AI_SERVICE.md`**
  - PyTorch ensemble model (99% accuracy)
  - 3 FastAPI endpoints
  - Disease classes, input/output formats
  - Performance metrics
  - Docker deployment

- **`IG_TOKKATOT_2.0_FARMER_CENTRIC_SPECIFICATIONS.md`**
  - Design principles for elderly farmers
  - Accessibility standards
  - Multilingual requirements
  - Device setup process (team-only)

### Operational Guides (OG_*)

**How to run/deploy** - Operational reference:

- **`OG_SPECIFICATIONS_DEPLOYMENT.md`** - Infrastructure, Docker, CI/CD
- **`OG_SPECIFICATIONS_TECHNOLOGY_STACK.md`** - Why Go/Vue/PyTorch/PostgreSQL
- **`OG_PROJECT_TIMELINE.md`** - 27-35 week schedule
- **`OG_TEAM_STRUCTURE.md`** - Roles, responsibilities, team size
- **`OG_RISK_MANAGEMENT.md`** - Risks and mitigation strategies

---

## 📝 When to Update Documentation

### Update Immediately

✅ **Add new endpoint** → Document in `IG_SPECIFICATIONS_API.md`
- Endpoint path, method, parameters, response format
- Example request/response (JSON)
- Error cases and HTTP status codes
- Which role can access (Owner/Manager/Viewer)

✅ **Change database schema** → Document in `IG_SPECIFICATIONS_DATABASE.md`
- New table or new column
- Data types, constraints, indexes
- Foreign key relationships
- Migration strategy

✅ **Change architecture** → Document in `01_SPECIFICATIONS_ARCHITECTURE.md`
- Add microservice diagram
- Update data flow pattern
- Document integration points
- Update component descriptions

✅ **Add feature requirement** → Document in `02_SPECIFICATIONS_REQUIREMENTS.md`
- Feature description (FR#)
- Acceptance criteria
- User story perspective
- Non-functional implications

✅ **Fix farm-facing issue** → Update `IG_TOKKATOT_2.0_FARMER_CENTRIC_SPECIFICATIONS.md`
- Accessibility considerations
- Farmer education
- Multilingual implications

---

## 🔄 Specification Template Examples

### API Endpoint Template

```markdown
#### 64. Get AI Service Health
\`\`\`
GET /ai/health
Authorization: Bearer {access_token}

Response (200):
{
  "status": "healthy",
  "model_loaded": true,
  "device": "cuda",
  "timestamp": "2026-02-19T10:30:00Z"
}

Requires: user role >= Viewer
Data Source: AI Service (FastAPI)
\`\`\`
```

### Database Schema Template

```markdown
#### predictions table
- **id** (UUID, Primary Key)
- **user_id** (UUID, Foreign Key → users)
- **device_id** (UUID, Foreign Key → devices, nullable)
- **disease** (VARCHAR(50): Coccidiosis|Newcastle|Salmonella|Healthy)
- **confidence** (FLOAT: 0.0-1.0)
- **image_hash** (VARCHAR(64): SHA256 hash)
- **created_at** (TIMESTAMP)

**Indexes**:
- INDEX ON (user_id, created_at DESC)
- INDEX ON (farm_id, created_at DESC)

**Constraints**:
- confidence >= 0 AND confidence <= 1
```

### Requirement Template

```markdown
### FR38: Disease Detection Prediction (AI Service)

**User Story**: As a farmer, I want to upload a chicken feces image and get instant AI-powered disease prediction so I can take treatment action quickly.

**Acceptance Criteria**:
- [ ] Farmer can upload PNG/JPEG image via web app
- [ ] Image size max 5MB, auto-rejects larger
- [ ] Prediction returns within 3 seconds (CPU) or <500ms (GPU)
- [ ] Result shows disease name + confidence percentage
- [ ] Result shows treatment recommendations (step-by-step)
- [ ] Result saved to database for history
- [ ] If confidence < 50%, return "uncertain" with guidance to retake

**Non-Functional**:
- Accuracy: 99% (ensemble voting)
- Language: Khmer + English
- Accessibility: Text readable on 48px font
```

---

## ✅ Sync Checklist (Before Submitting PR)

When making code changes, verify documentation is current:

- [ ] Core spec files reference updated components (`01_SPECIFICATIONS_ARCHITECTURE.md`)
- [ ] API spec matches code endpoints (`IG_SPECIFICATIONS_API.md`)
- [ ] Database schema matches code migrations (`IG_SPECIFICATIONS_DATABASE.md`)
- [ ] Requirements still match code behavior (`02_SPECIFICATIONS_REQUIREMENTS.md`)
- [ ] If Python code changed, AI spec updated (`IG_SPECIFICATIONS_AI_SERVICE.md`)
- [ ] If frontend changed, UI spec updated (`IG_SPECIFICATIONS_FRONTEND.md`)
- [ ] If embedded code changed, firmware spec updated (`IG_SPECIFICATIONS_EMBEDDED.md`)
- [ ] README.md still accurate for setup instructions

---

## 📊 DocStrings in Code

**Python (FastAPI)**:
```python
async def predict(file: UploadFile = File(...)) -> PredictionResponse:
    """
    Predict disease from fecal image (simple response).
    
    SPEC: See IG_SPECIFICATIONS_API.md Endpoint 65
    
    Args:
        file: PNG/JPEG image, max 5MB
    
    Returns:
        PredictionResponse with disease, confidence, recommendation
    
    Raises:
        HTTPException(400): Invalid image format/size
        HTTPException(503): Model not loaded
    """
```

**Go**:
```go
// AuthLogin handles POST /auth/login
// SPEC: See IG_SPECIFICATIONS_API.md Endpoint 1
//
// Accepts: {"email": "...", "password": "..."}
// Returns: {"access_token": "...", "refresh_token": "..."}
func AuthLogin(w http.ResponseWriter, r *http.Request) {
```

**JavaScript**:
```javascript
/**
 * Predict disease from image
 * SPEC: See IG_SPECIFICATIONS_API.md Endpoint 65
 * 
 * @param {File} imageFile - PNG/JPEG image, max 5MB
 * @returns {Promise<Object>} {disease, confidence, recommendation}
 */
async function predictDisease(imageFile) {
```

---

## 🔗 Cross-References

**Always link between specs**:

In `IG_SPECIFICATIONS_API.md`:
```markdown
For database schema, see [IG_SPECIFICATIONS_DATABASE.md](IG_SPECIFICATIONS_DATABASE.md)
For security, see [IG_SPECIFICATIONS_SECURITY.md](IG_SPECIFICATIONS_SECURITY.md)
```

In code comments:
```go
// SPEC: IG_SPECIFICATIONS_API.md Endpoint 65 (Predict Disease)
// SPEC: IG_SPECIFICATIONS_DATABASE.md table "predictions"
```

---

## 📚 Markdown Best Practices

### Document Structure

```markdown
# Main Title (H1)

## Section (H2)

### Subsection (H3)

#### Detail (H4)

- Bullet point
- Another point
  - Indent for sub-point

| Column 1 | Column 2 |
|----------|----------|
| Data     | Data     |

\`\`\`code
code block
\`\`\`
```

### Emphasis

```markdown
**Bold** - Important concepts
*Italic* - Emphasis
`Code` - Function names, variables
> Quote - Important callouts
```

---

## 🆘 Common Spec Issues & Fixes

### Issue: Endpoint documented but not implemented
```
ERROR: API spec says POST /farms/{id}/devices but code missing
```
**Fix**: Either implement the endpoint immediately or remove from spec

### Issue: Database table exists but not documented
```
ERROR: Code has "prediction_logs" table but IG_SPECIFICATIONS_DATABASE.md missing
```
**Fix**: Add full table documentation with schema

### Issue: Spec says 3-second response time but code takes 10 seconds
```
ERROR: Performance spec violated
```
**Fix**: Either update code to meet spec OR update spec with new target + justification

### Issue: UI spec says 48px buttons but code has 32px
```
ERROR: Accessibility spec violated
```
**Fix**: Update UI to match accessibility spec

---

## 📈 Documentation Structure

**Each spec file should have**:

1. **Title & Metadata**
   - Document version
   - Last updated date
   - Status (Draft/Final/Deprecated)

2. **Overview** (0.5-1 page)
   - What this component does
   - Key features
   - Context in larger system

3. **Details** (5-10 pages typical)
   - Specifications
   - Schemas
   - Endpoints
   - Formats

4. **Cross-References**
   - Links to related specs
   - Links to code files
   - Links to other documentation

5. **Examples**
   - JSON request/response
   - Code examples
   - Usage patterns

---

## 🎯 Your Responsibilities

### You Own Documentation For:
- Keeping specs current with code changes
- Updating examples when outputs change
- Maintaining cross-references (not broken links)
- Ensuring farmer-centric language
- Recording architectural decisions

### You Should Ask Before:
- Major architectural changes (spec implications) → Ask tech lead
- Database schema changes (migration complexity) → Ask DBA
- Breaking API changes → Ask product owner
- Deleting/deprecating endpoints → Discuss with team

---

## 📞 When Something is Unclear

1. **First**: Check if spec is unclear (ambiguous wording)
   - If yes: Update spec to be clearer
   - If no: Code implementation is wrong

2. **Second**: Check if spec is outdated
   - If yes: Update spec to match code or vice versa
   - If no: Proceed with implementation

3. **Third**: Ask the team
   - Spec gap → Product owner
   - Implementation gap → Relevant developer
   - Architecture gap → Tech lead

---

## ✨ Example: Adding New Endpoint

**Scenario**: Add `POST /api/farms/{id}/photos` to upload farm photos

**Steps**:

1. **Document in IG_SPECIFICATIONS_API.md**:
   ```markdown
   #### 67. Upload Farm Photo
   \`\`\`
   POST /farms/{farm_id}/photos
   Authorization: Bearer {access_token}
   
   Request:
   - photo: Image file (PNG/JPEG, max 10MB)
   - title: String (optional)
   
   Response (201):
   {
     "photo_id": "uuid",
     "url": "https://cdn.tokkatot.local/photos/...",
     "created_at": "2026-02-19T10:30:00Z"
   }
   \`\`\`
   ```

2. **Update IG_SPECIFICATIONS_DATABASE.md**:
   ```markdown
   #### photos table
   - id (UUID, Primary Key)
   - farm_id (FK farms)
   - title (VARCHAR)
   - url (VARCHAR)
   - created_at (TIMESTAMP)
   ```

3. **Update 02_SPECIFICATIONS_REQUIREMENTS.md**:
   ```markdown
   ### FR40: Farm Photo Upload
   
   Farmer can upload photos of farm/equipment for documentation
   ```

4. **Implement in Go** (middleware):
   ```go
   // SPEC: IG_SPECIFICATIONS_API.md Endpoint 67
   func UploadFarmPhoto(w http.ResponseWriter, r *http.Request) {
   ```

5. **Implement database**:
   ```go
   // Add to IG_SPECIFICATIONS_DATABASE.md migration
   CREATE TABLE photos (...)
   ```

---

**Documentation is code. Treat it with same care. Outdated docs = team confusion = bugs.** 🎯

Now go keep those specs current! 📚✨
