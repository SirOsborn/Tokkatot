# Tokkatot 2.0: Security Architecture Specification

**Document Version**: 2.0  
**Last Updated**: February 2026  
**Status**: Final Specification

---

## Overview

This document specifies security requirements, authentication, authorization, encryption, and compliance measures for Tokkatot 2.0.

---

## Authentication & Authorization

### Authentication Methods

**Primary: JWT (JSON Web Tokens)**
- Standard: RFC 7519
- Algorithm: HS256 (HMAC with SHA-256)
- Expiration: 24 hours per token
- Refresh: Automatic on login
- Storage: `localStorage` (client-side; not HttpOnly cookie)

**Token Structure**:
```json
{
  "sub": "user-uuid",
  "email": "farmer@example.com",
  "phone": null,
  "farm_id": "farm-uuid",
  "role": "farmer",
  "iat": 1677000000,
  "exp": 1677086400
}
```
> `email` and `phone` are nullable — only the field used at signup/login is populated.

**Password Requirements**:
- Minimum 8 characters
- Must contain: uppercase, lowercase, number, symbol
- Cannot reuse last 3 passwords
- Expiration: 90 days (prompt for change)
- Failure lockout: 5 failed attempts = 15 minute lockout

### Authorization (2-Role Farm System)

**2 Farm Roles** (farmer-centric — no Admin, no Manager, no Keeper):

**Farmer**
- Full farm management (settings, delete farm)
- Manage all members (invite `farmer` or `viewer`, change role, remove)
- Full device control and scheduling
- View all analytics and audit logs
- Multiple farmers can share the same farm with equal full access

**Viewer**
- Read-only access to all monitoring data and alerts
- Can **acknowledge alerts** (maintenance/worker role for large farms)
- Cannot control devices, cannot change farm settings
- No user management

**Tokkatot System Staff** (not a `farm_users` role — system-level access only):
- Manage registration keys (`registration_keys` table + `generate_reg_key.ps1`)
- Control JWT configuration (`.env` `JWT_SECRET`, `REG_KEY`)
- System-level bypasses for onboarding and support
- Full visibility across all farms for support/diagnostic purposes
- Access managed via server environment and `registration_keys` table, not via farm invitation

**Permission Matrix**:

| Resource | Farmer | Viewer |
|----------|--------|--------|
| User Management | ✓ | ✗ |
| Device Control | ✓ | ✗ |
| Create Schedule | ✓ | ✗ |
| View Analytics | ✓ | ✓ |
| Farm Settings | ✓ | ✗ |
| Audit Logs | ✓ | ✗ |
| Acknowledge Alerts | ✓ | ✓ |

### Multi-Factor Authentication (MFA)

**Not implemented for farmers** — MFA is an optional future feature for Tokkatot staff/admin only:
- Time-based One-Time Password (TOTP)
- SMS codes (optional fallback)
- Backup codes for account recovery

---

## Encryption

### Data in Transit

**HTTPS/TLS**:
- Minimum: TLS 1.3
- Certificate: Let's Encrypt (auto-renewed)
- Cipher Suites: Only modern, secure suites
  - TLS_AES_256_GCM_SHA384
  - TLS_CHACHA20_POLY1305_SHA256
  - TLS_AES_128_GCM_SHA256

**Certificate Pinning** (Mobile Apps):
- Pin API certificate to app bundle
- Prevent man-in-the-middle attacks
- Update mechanism during app updates

**MQTT Communication**:
- TLS 1.3 over port 8883 (external devices)
- Plain TCP 1883 (internal/local devices)

### Data at Rest

**Database Encryption**:
- PostgreSQL: transparent data encryption (TDE)
- All sensitive fields encrypted:
  - Email addresses
  - Phone numbers
  - API keys
  - Device credentials
- Encryption key: AWS KMS or DigitalOcean managed key

**Storage Encryption**:
- S3/Spaces: Server-side encryption (SSE-S3)
- Firmware files: encrypted with public key
- Configuration backups: AES-256 encrypted
- Key rotation: Annual or on compromise

**Backups**:
- Encrypted at rest in S3
- Encryption key stored separately
- Access restricted to ops team

---

## API Security

### Rate Limiting

**Per-Client Limits**:
- Authenticated endpoint: 1000 requests/minute per user
- Public endpoint: 100 requests/minute per IP
- Burst allowance: 50 requests/second
- Penalty: 429 (Too Many Requests) response

**Endpoint-Specific Limits**:
```
POST /api/auth/login       → 5 attempts/15 min per user (phone/id)
POST /api/devices/commands → 100 commands/min per device
GET /api/data/history      → 50 queries/hour per user
```

**Implementation**: Redis-backed token bucket algorithm

### Input Validation

**All inputs validated**:
- SQL injection prevention: parameterized queries
- XSS prevention: input sanitization
- CSRF protection: SameSite cookies, CSRF tokens
- Path traversal: whitelist allowed paths
- File upload: whitelist extensions, scan for malware

**Example Validation**:
```
Device Name:
  - Type: string
  - Length: 1-100 characters
  - Pattern: alphanumeric + spaces/hyphens only
  - Reject: special characters, scripts

Time-Based Duration:
  - Type: integer
  - Range: 1-3600 seconds (1 second to 1 hour)
  - Reject: negative, zero, > 1 hour
```

### API Versioning

**Current Version**: v1 (available in 2.0)  
**Deprecation**: Old versions supported for 12 months  
**Headers**: Accept: application/vnd.tokkatot.v1+json  

---

## Device Security

### Device Authentication

**Device Certificate**:
- Self-signed certificate on each ESP32
- Certificate includes device ID
- Verified by RPi agent before accepting commands
- No TOFU (Trust On First Use) - pre-register devices

**Device Credentials**:
- MQTT username/password: unique per device
- Credentials stored in device EEPROM (protected)
- Credentials rotated every 6 months
- Compromised credentials: immediate revocation

### Device-to-Cloud Communication

**MQTT over TLS**:
- Certificate verification enabled
- Device verifies server certificate
- Certificate pinning on device firmware (optional)
- Topic-based access control:
  ```
  device/{id}/command        → device can receive
  device/{id}/status         → device can publish
  device-admin/{id}/firmware → firmware OTA only
  ```

**Secure Firmware Updates**:
- Firmware signed with private key
- Device verifies signature before installing
- Rollback protection (previous version kept)
- Update restricted to farmer role
- Staged rollout (test on single device first)

---

## API Endpoint Security

### Authentication Required

**Endpoints requiring JWT**:
- All `/api/devices/*` endpoints
- All `/api/schedules/*` endpoints
- All `/api/data/*` endpoints
- All `/api/users/*` endpoints
- All `/api/farms/*` endpoints

**No Authentication Required**:
- POST /api/auth/register
- POST /api/auth/login
- GET /api/health
- GET /api/version

### Response Security

**No Sensitive Data in Responses**:
- Never return: password hashes, API keys, tokens
- Return: user ID, email (for own account only)
- Error messages: generic ("Invalid credentials" not "Email not found")
- Stack traces: only in development, never production

**Security Headers**:
```
Strict-Transport-Security: max-age=31536000; includeSubDomains
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Content-Security-Policy: default-src 'self'
Referrer-Policy: strict-origin-when-cross-origin
```

---

## Compliance & Standards

### GDPR Compliance

**User Rights**:
- **Right to Access**: /api/users/me/export (download all data)
- **Right to Delete**: /api/users/me/delete (anonymize account)
- **Right to Portability**: Export data in standard format
- **Breach Notification**: 72 hours to authorities

**Data Processing**:
- Data Processing Agreement (DPA) with users
- Legitimate interest assessed for each processing
- Privacy impact assessment (PIA) completed
- Consent explicit before processing (opt-in)

### CCPA Compliance (California)

**Consumer Rights**:
- Know what data is collected
- Delete personal data
- Opt-out of data selling
- Non-discrimination for exercising rights

**Implementation**:
- Privacy Policy updated
- Opt-in consent required
- Data retention policy documented

### OWASP Top 10 Mitigation

| OWASP Risk | Mitigation |
|-----------|-----------|
| Injection | Parameterized queries, input validation |
| Broken Auth | JWT + password hashing, MFA option |
| XSS | Input/output encoding, CSP headers |
| XXE | Disable XML parsing, strict schema |
| Broken Access | RBAC, proper permission checks |
| Security Misconfiguration | Infrastructure-as-code, security checklist |
| Sensitive Data | Encryption at rest/transit, PII handling |
| XML External | Disable external entities |
| Broken Func Auth | Role-based permissions, audit logging |
| Using Components | Dependency scanning, automated updates |

---

## Incident Response

### Security Incident Categories

**Severity Levels**:
- **Critical**: Data breach, service unavailable > 1 hour
- **High**: Unauthorized access, partial data exposure
- **Medium**: Failed security control, suspicious activity
- **Low**: Policy violation, minor logging issue

### Response Procedures

**Upon Detection**:
1. **Confirm**: Verify incident is real (not false alarm)
2. **Assess**: Determine scope, affected systems, severity
3. **Contain**: Isolate affected systems, prevent spread
4. **Communicate**: Notify affected users (within 24 hours)
5. **Investigate**: Root cause analysis, gather evidence
6. **Remediate**: Fix vulnerability, patch systems
7. **Recover**: Restore service, validate functionality
8. **Document**: Record all actions for review
9. **Improve**: Update security measures to prevent recurrence

### Communication Plan

**Notification Timeline**:
- Immediate: Notify security team (internal Slack)
- Within 1 hour: Notify affected farmers via app notification (message log)
- Within 24 hours: Publish incident report
- Within 72 hours: GDPR notification to authorities (if required)

---

## Monitoring & Logging

### Security Event Logging

**Events Logged**:
- Failed login attempts
- Successful login (timestamp, IP, device)
- API authentication failures
- Permission denied events
- Data export/download
- Account changes (phone number, password, role)
- Administrative actions

**Log Format** (JSON):
```json
{
  "timestamp": "2026-02-18T14:30:00Z",
  "event_type": "authentication_failure",
  "severity": "warning",
  "user_id": "uuid",
  "identifier": "farmer@example.com",
  "ip_address": "203.0.113.45",
  "reason": "invalid_password",
  "attempt_number": 3,
  "farm_id": "farm-uuid"
}
```

### Real-Time Alerts

**Triggered on**:
- Multiple failed login attempts (5+)
- Access from new IP address
- Permission denied > 10/minute
- API error rate > 5%
- Unauthorized data access attempt

**Alerting Channel**: In-app notification (message log), Slack for security team (internal)

---

## Vulnerability Management

### Dependency Scanning

**Tools**: Snyk, Dependabot, npm audit  
**Frequency**: On every commit, daily automated scan  
**Response**: Critical patches within 24 hours  

**Dependency Restrictions**:
- No packages with known CVEs
- Max age of dependencies: 1 year (unless security issue)
- Automated updates: minor + patch only
- Manual updates: major versions

### Penetration Testing

**Frequency**: Yearly (external), quarterly (internal)  
**Scope**: API endpoints, web app, device communication  
**Remediation**: Critical/High within 2 weeks, Medium within month  

---

## Physical Security (On-Farm)

### Device Protection

**ESP32 Devices**:
- Mounted in weatherproof enclosure
- Anti-tamper seal (optional)
- Serial number recorded
- GPS tracking (optional for valuable deployments)

**Raspberry Pi Agent**:
- Secured in locked cabinet
- Physical access restricted to farm owner
- Backup power (UPS) if critical

### Local Network

- WiFi: WPA3 encryption (or WPA2 minimum)
- Password: 20+ character, complex
- MAC filtering: optional additional layer
- Guest network: disabled

---

## Secrets Management

### Secure Storage

**Production Secrets** (stored in DigitalOcean Vault or AWS Secrets Manager):
- Database passwords
- JWT signing key
- API keys (3rd party services)
- Encryption keys (optional)
- MQTT credentials
- SSL certificates

**Rotation Policy**:
- Database password: every 90 days
- API keys: every 180 days
- JWT secret: every 1 year
- Encryption keys: yearly or after breach

**Access Control**:
- Only Kubernetes pods need permissions
- No secrets in git repository
- Development: Use .env.local (never commit)
- CI/CD: Use GitHub Secrets

---

## Security Checklist

Before deploying to production:

- [ ] All secrets stored in vault, not in code
- [ ] HTTPS enforced at load balancer
- [ ] Database encryption enabled
- [ ] Input validation implemented
- [ ] Output encoding implemented
- [ ] Rate limiting configured
- [ ] CORS policy defined
- [ ] Security headers set
- [ ] API authentication required
- [ ] Authorization checks implemented
- [ ] Logging and monitoring active
- [ ] Backup encryption verified
- [ ] SSL certificate valid and updated
- [ ] RBAC roles tested
- [ ] Penetration test completed
- [ ] Dependency scan passed
- [ ] Security audit completed

---

**Version History**
| Version | Date | Changes |
|---------|------|---------|
| 2.0 | Feb 2026 | Initial security specification |

**Related Documents**
- SPECIFICATIONS_ARCHITECTURE.md
- SPECIFICATIONS_DEPLOYMENT.md
- SPECIFICATIONS_DATABASE.md
