# GitHub Pull Request Code Review Checklist

**For GitHub Copilot & AI Code Reviewers**

This checklist ensures all PRs to Tokkatot maintain quality, security, and consistency.

---

## 🔍 Security Checks (CRITICAL)

- [ ] **No secrets leaked**
  - No API keys, passwords, or database credentials in code
  - All `.env` values are used via `os.getenv()` or Go `os.Getenv()`
  - Check: Are there any hardcoded connection strings?

- [ ] **No model files committed**
  - No `*.pth`, `*.h5`, `*.pkl`, `*.bin` files in diff
  - Check `/ai-service/outputs/` - should only show `.gitignore` changes
  - Model files should be in `.gitignore`

- [ ] **Error messages safe**
  - No internal paths exposed (e.g., `/var/www/app/models/...`)
  - No database column names in error messages
  - Generic "Internal Server Error" for unexpected failures
  - Check API error responses

- [ ] **Input validation present**
  - File uploads validated (size, type, extension)
  - JSON payloads validated against schema
  - String length checks (no buffer overflows)
  - No SQL injection risks (parameterized queries only)

- [ ] **Authentication/Authorization**
  - JWT token validation on protected endpoints
  - Role checks (Owner/Manager/Viewer) enforced
  - User can only access their own data
  - Device operations restricted to Tokkatot team only

---

## 📚 Documentation Checks

- [ ] **Specs are in sync** (CRITICAL - see AI_INSTRUCTIONS.md "Documentation Update Protocol")
  - API endpoints documented in `docs/implementation/API.md` ✓
  - Database schema matches `docs/implementation/DATABASE.md` ✓
  - Architecture changes reflected in `docs/ARCHITECTURE.md` ✓
  - Requirements updated in `docs/02_SPECIFICATIONS_REQUIREMENTS.md` ✓
  - Automation use cases in `docs/AUTOMATION_USE_CASES.md` (if schedule-related) ✓

- [ ] **AI Context files updated** (for significant changes only)
  - Component `AI_CONTEXT.md` updated if new patterns added (middleware/, frontend/, ai-service/, embedded/)
  - `AI_INSTRUCTIONS.md` updated if major system concept added
  - Examples tested and working (JSON valid, code compiles)
  - Cross-references added between related docs

- [ ] **Code is documented**
  - Functions have comments explaining purpose
  - Complex business logic explained (especially around "farmer-centric" logic)
  - Error messages are user-friendly
  - Inline comments for non-obvious code

- [ ] **PR description is clear**
  - What problem does this solve? (farmer problem preferred)
  - How was it tested?
  - Any breaking changes?
  - Link to related issue (if any)
  - Which docs were updated? (list files)

---

## 🏗️ Architecture & Design

- [ ] **Follows tech stack**
  - Go for API (not Node/Python)
  - Vue.js 3 for frontend (not React/Angular)
  - SQLite for embedded (not PostgreSQL on device)
  - PyTorch for AI (not TensorFlow)

- [ ] **Database changes**
  - Schema migration included (if any)
  - Changes reflected in `docs/implementation/DATABASE.md` (with examples)
  - API.md updated with new fields in request/response examples
  - No breaking changes without deprecation period
  - Foreign key relationships maintained
  - Indexes added for new query patterns

- [ ] **API design**
  - Endpoints follow REST conventions
  - HTTP status codes correct (200, 201, 400, 401, 403, 404, 500)
  - Pagination implemented for list endpoints (if applicable)
  - Rate limiting headers present (if applicable)

- [ ] **Farmer-centric design**
  - UI uses 48px+ fonts/buttons (accessibility)
  - Features require ≤2 clicks to complete (simplicity)
  - No overwhelming choice/options (max 5 options shown)
  - Khmer/English language toggle supported

---

## ✅ Code Quality

- [ ] **No hardcoded values**
  - Database URLs use env vars
  - API endpoints configurable
  - Feature flags for new features
  - Timeouts & retries configured

- [ ] **Error handling comprehensive**
  - Try-catch blocks around risky operations
  - Network failures handled gracefully
  - Database connection failures handled
  - File I/O errors handled

- [ ] **Type safety**
  - Go: No `interface{}` unless necessary
  - Python: Type hints on all functions
  - JavaScript: PropTypes or JSDoc type comments

- [ ] **Performance acceptable**
  - API endpoints respond < 1 second
  - AI predictions complete < 3 seconds (CPU)
  - Frontend pages load < 2 seconds (4G network)
  - No N+1 queries in database code

- [ ] **Logging & debugging**
  - Structured logging with timestamps
  - Request IDs for tracing
  - Useful debug logs (not verbose spam)
  - Error logs include stack traces

---

## 🧪 Testing

- [ ] **Tests included**
  - Unit tests for business logic
  - Integration tests for API endpoints
  - Tests for edge cases (empty input, timeouts, etc)
  - Tests pass locally (`go test`, `python -m pytest`, etc)

- [ ] **No regressions**
  - Existing tests still pass
  - No breaking changes to APIs
  - Database migrations are reversible
  - Backward compatibility maintained

---

## 📦 Deployment Readiness

- [ ] **Docker builds successfully**
  - `docker build` completes without errors
  - Image size reasonable (< 1GB for AI service)
  - Non-root user used in Dockerfile
  - Health checks configured

- [ ] **Environment configuration**
  - All config via `.env` or environment variables
  - `.env.example` updated with new variables
  - Secrets not in code

- [ ] **Git compliance**
  - `.gitignore` updated if new files added
  - No accidental binary files committed
  - Commit messages clear and descriptive

---

## 🔄 Integration Points

- [ ] **Component integration correct**
  - API responses match frontend expectations
  - Device commands reach ESP32 correctly
  - AI predictions stored in database
  - WebSocket updates broadcast correctly

- [ ] **Cross-service communication**
  - Go API calls AI service correctly
  - AI service returns expected format
  - Database transactions consistent
  - MQTT messages follow protocol spec

---

## 🚀 Pre-Merge Checklist

- [ ] Approved by at least one human reviewer
- [ ] All checks pass (linting, tests, builds)
- [ ] Branch is up-to-date with main
- [ ] No merge conflicts
- [ ] Commit history is clean (squash if needed)
- [ ] PR description complete
- [ ] Related issues linked/closed

---

## 📝 Common Issues to Watch For

**Red Flags** (Request changes):
- ❌ Secrets found in code
- ❌ Model files in git
- ❌ Specification not updated
- ❌ No tests for new features
- ❌ SQL injection vulnerability
- ❌ Hardcoded passwords

**Warnings** (Request review):
- ⚠️ Large file added (> 10MB) - may be unnecessary
- ⚠️ Database migration without rollback plan
- ⚠️ Breaking API change without version bump
- ⚠️ Performance critical code not benchmarked
- ⚠️ Accessibility not tested (font sizes, colors, etc)

**Good to improve** (Suggest):
- 💡 Add more comments
- 💡 Simplify complex logic
- 💡 Extract function to reduce duplication
- 💡 Add type hints
- 💡 Update README with usage example

---

## 🎯 Decision Matrix

| Issue | Severity | Action |
|-------|----------|--------|
| Secrets in code | 🔴 CRITICAL | Block PR, require changes, audit git history |
| Model files in git | 🔴 CRITICAL | Block PR, remove from history |
| No tests | 🟠 HIGH | Request changes |
| Spec not updated | 🟠 HIGH | Request changes |
| SQL injection | 🔴 CRITICAL | Block PR immediately |
| Farmer UX issue | 🟠 HIGH | Request changes or discussion |
| Code style | 🟡 LOW | Approve, suggest improvements |
| Comments missing | 🟡 LOW | Approve, suggest improvements |

---

## 💬 Suggestion Format

**Good feedback**:
```
✅ Love how you handled the error case here. Consider also handling the timeout case:

```go
// Add context timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```

Would make it more robust.
```

**Poor feedback**:
```
❌ This code is bad.
```

Always explain why and suggest alternatives!

---

## 🤖 For AI Code Reviewers

When using AI to review PRs:

1. **Read the entire PR context first** - Don't just flag issues, understand intent
2. **Cross-reference specs** - Check that code matches IG_SPECIFICATIONS_*.md files
3. **Test locally if possible** - Try to build/run the code
4. **Example realistic scenarios** - "What if farmer is offline?" "What if image upload fails?"
5. **Be constructive** - Explain why, not just what
6. **Know your limits** - If unsure about security, flag for human review

---

## ✨ Approval Template

```markdown
✅ **APPROVED** (or 🔄 **REQUEST CHANGES**, or ❓ **COMMENT**)

**Summary**: 
[What does this PR do?]

**Strengths**:
- ✅ Good error handling
- ✅ Tests included
- ✅ Spec updated

**Suggestions**:
- 💡 Consider adding X
- 💡 Performance: Y could be optimized

**Concerns**:
- ⚠️ Need clarification on Z

**Final**: Ready to merge! 🚀
```

---

*Last Updated: February 19, 2026*  
*For questions: See AI_INSTRUCTIONS.md*
