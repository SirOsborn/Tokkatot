# Tokkatot 2.0: Deployment & Infrastructure Specification

**Document Version**: 2.0  
**Last Updated**: February 2026  
**Status**: Final Specification

**Deployment Model**: Cloud-Based Only  
**Cloud Provider Selection**: DigitalOcean (recommended) or AWS  
**Infrastructure Pattern**: Docker Containers + Kubernetes/Docker Compose  

---

## Overview

Tokkatot 2.0 is deployed exclusively on cloud infrastructure with edge computing fallback on-farm. This document specifies the complete deployment architecture, configurations, and operational procedures.

---

## ⚠️ DEPLOYMENT OPTION

### Only Option: Cloud-Based Deployment

Tokkatot 2.0 uses **cloud-only deployment** for all backend services:

- ✅ **Cloud Services**: All API services, databases, and monitoring run on cloud
- ✅ **Edge Fallback**: Local Raspberry Pi acts as backup (offline mode, local control)
- ✅ **No Self-Hosted Option**: Single-tenant or multi-tenant, always cloud-managed
- ✅ **Scalability**: Infinitely scalable for future expansion
- ✅ **Maintenance**: Platform handles all updates, backups, security patches
- ✅ **Cost-Effective**: Lower initial cost, pay-as-you-go, no capital expense

### Why Cloud-Only?

| Reason | Benefit |
|--------|---------|
| **Security** | Professional security team, automated patching |
| **Reliability** | 99.5% uptime SLA, automatic failover, backups |
| **Maintenance** | Zero downtime updates, automatic scaling |
| **Monitoring** | 24/7 system monitoring, automated alerts |
| **Compliance** | GDPR/CCPA ready, audit trails, encryption |
| **Cost** | No capital expense, predictable operating costs |
| **Updates** | Instant rollout of features and security patches |
| **Multi-Farm** | Easy to add new farms, centralized management |

---

## Recommended Provider: DigitalOcean

### Provider Comparison

| Feature | DigitalOcean | AWS | Linode |
|---------|-------------|-----|--------|
| **Pricing (Startup)** | $100-300/month | $200-500/month | $150-400/month |
| **Ease of Use** | Simple, minimal config | Complex, many options | Moderate |
| **Support** | Good community, support | 24/7 premium | Community-based |
| **Learning Curve** | Flat | Steep | Moderate |
| **Team Size** | Small teams | Enterprise | Small teams |
| **Scalability** | Excellent for SMB | Unlimited enterprise | Good for SMB |
| **Geographic Coverage** | 12 regions | 30+ regions | 11 regions |

**Recommendation**: **DigitalOcean** for v2.0 (simplicity, cost), migrate to AWS if needs exceed DigitalOcean capacity.

---

## Cloud Infrastructure Architecture

### DigitalOcean Deployment Stack

```
┌─────────────────────────────────────────────────────────┐
│                    DigitalOcean Account                  │
│                                                          │
│  ┌──────────────────────────────────────────────────┐  │
│  │        Kubernetes Cluster (DOKS)                 │  │
│  │  - 3+ nodes (high availability)                  │  │
│  │  - auto-scaling: 3-10 nodes                      │  │
│  │  - Network policy, RBAC enabled                  │  │
│  │                                                   │  │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐        │  │
│  │  │  Auth    │ │  Device  │ │ Schedule │        │  │
│  │  │ Service  │ │ Service  │ │ Service  │        │  │
│  │  └──────────┘ └──────────┘ └──────────┘        │  │
│  │                                                   │  │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐        │  │
│  │  │   Data   │ │Notification│ Config  │        │  │
│  │  │ Service  │ │ Service  │ │ Service │        │  │
│  │  └──────────┘ └──────────┘ └──────────┘        │  │
│  └──────────────────────────────────────────────────┘  │
│                                                          │
│  ┌──────────────────────────────────────────────────┐  │
│  │          Managed Databases                       │  │
│  │  - PostgreSQL Managed Database                   │  │
│  │    • Primary instance + standby replica          │  │
│  │    • Automated backups, point-in-time restore    │  │
│  │  - Redis Managed Cache                           │  │
│  │    • High availability mode                      │  │
│  │    • Manual and automatic failover               │  │
│  │  - InfluxDB Managed Database (Time-series)       │  │
│  │    • Clustered for reliability                   │  │
│  └──────────────────────────────────────────────────┘  │
│                                                          │
│  ┌──────────────────────────────────────────────────┐  │
│  │          Storage & CDN                           │  │
│  │  - Spaces (S3-compatible object storage)         │  │
│  │    • Firmware files                              │  │
│  │    • Backups and archives                        │  │
│  │  - CDN (for static assets)                       │  │
│  │    • Web app frontend delivery                   │  │
│  │    • Global edge locations                       │  │
│  └──────────────────────────────────────────────────┘  │
│                                                          │
│  ┌──────────────────────────────────────────────────┐  │
│  │       Networking & Load Balancing                │  │
│  │  - Load Balancer (Layer 4/7)                     │  │
│  │    • Routes traffic to K8s cluster               │  │
│  │    • SSL/TLS termination                         │  │
│  │    • Health checking                             │  │
│  │  - VPC (Virtual Private Cloud)                   │  │
│  │    • Isolated network for all resources          │  │
│  │    • Firewall rules (DDoS protection)            │  │
│  └──────────────────────────────────────────────────┘  │
│                                                          │
│  ┌──────────────────────────────────────────────────┐  │
│  │       Monitoring & Logging                       │  │
│  │  - Prometheus + Grafana (metrics)                │  │
│  │  - ELK Stack or Loki (logs)                      │  │
│  │  - AlertManager (alerting rules)                 │  │
│  └──────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
                        ↕ HTTPS/WSS
        ┌───────────────────────────────────────┐
        │    User Clients (Web + Mobile)        │
        │    ↓                                   │
        │    MQTT Brokers (Device Comms)        │
        │    ↓                                   │
        │    On-Farm Raspberry Pi Agents        │
        │    ↓                                   │
        │    ESP32 Devices (Local Hardware)     │
        └───────────────────────────────────────┘
```

---

## Service Deployment Architecture

### 1. Kubernetes Cluster Setup

**Cluster Name**: tokkatot-prod  
**Cloud**: DigitalOcean Kubernetes Service (DOKS)  
**Kubernetes Version**: 1.28+  

**Node Configuration**:
- **Node Type**: Standard (s-2vcpu-4gb) minimum
- **Node Count**: 3 nodes (high availability)
- **Auto-Scaling**: enabled, 3-10 nodes
- **Network Policy**: enabled (pod-to-pod security)
- **RBAC**: enabled (role-based access control)

**Kubernetes Namespaces**:
```
production/       # All production services
├── api            # API services
├── databases      # Database-related K8s jobs
├── monitoring     # Prometheus, Grafana
└── ingress        # Ingress controller

staging/          # Test environment (optional)
development/      # Development environment (optional)
```

---

### 2. Containerization (Docker)

**Base Images**:
```
Backend Services: golang:1.21-alpine (minimal, ~400MB)
Frontend:        node:18-alpine (minimal, ~400MB)
AI Service:      python:3.11-slim (minimal, ~800MB)
```

**Image Repository**: DigitalOcean Container Registry (DOCR)

**Docker Compose (for local development)**:
```yaml
Services:
  backend:
    Build: ./middleware
    Image: tokkatot-backend:latest
    Ports: 8080:8080
    
  frontend:
    Build: ./frontend
    Image: tokkatot-frontend:latest
    Ports: 3000:3000
    
  postgres:
    Image: postgres:15
    Ports: 5432:5432
    
  redis:
    Image: redis:7-alpine
    Ports: 6379:6379
    
  influxdb:
    Image: influxdb:2.7
    Ports: 8086:8086
    
  mosquitto:
    Image: eclipse-mosquitto:2.0
    Ports: 1883:1883
```

**Image Tagging Strategy**:
```
tokkatot-backend:latest              # Latest stable
tokkatot-backend:v2.0.0              # Version tag
tokkatot-backend:v2.0.0-rc1          # Release candidate
tokkatot-backend:dev-main-abc123     # Development build
```

---

### 3. Kubernetes Deployment Manifests

**Key Deployments**:

**A. Backend API Services**
```
Deployment: api-server
  Replicas: 3 (HA)
  Container Image: tokkatot-backend:latest
  Resources:
    CPU: 500m request, 1000m limit
    Memory: 512Mi request, 1Gi limit
  Liveness Probe: /health (http, 10s interval)
  Readiness Probe: /ready (http, 5s interval)
  Environment Variables:
    DATABASE_URL
    REDIS_URL
    JWT_SECRET
```

**B. Cronjobs (Scheduled Tasks)**
```
CronJob: schedule-executor
  Schedule: */5 * * * * (every 5 minutes)
  Job: Execute pending schedules
  
CronJob: data-archival
  Schedule: 0 2 * * * (daily at 2 AM)
  Job: Archive old logs and sensor data
```

**C. WebSocket Server (Optional)**
```
Deployment: websocket-server
  Replicas: 2
  Type: StatefulSet (maintain connections)
  Service: ClusterIP (internal routing)
```

---

### 4. Managed Databases

**PostgreSQL (Primary)**
- **Instance Class**: db-s-2vcpu-4gb (2 CPU, 4GB RAM)
- **Storage**: 100GB (auto-expands if needed)
- **Replica**: 1 standby replica (automatic failover)
- **Backups**: Daily, retained 7 days
- **Connection Pooling**: PgBouncer (100 connections)
- **SSL**: Required for all connections (TLS 1.3)
- **Monitoring**: CPU, memory, connections, query performance

**Redis (Cache)**
- **Instance Class**: db-s-1vcpu-1gb (1 CPU, 1GB RAM)
- **Mode**: High availability (primary + replica)
- **Eviction Policy**: allkeys-lru (evict least recently used)
- **Backups**: Automated snapshots
- **SSL**: Required (TLS 1.3)

**InfluxDB (Time-Series)**
- **Instance Class**: db-s-2vcpu-8gb (2 CPU, 8GB RAM)
- **Storage**: 500GB
- **Retention**: Default 30 days (configurable per measurement)
- **Backups**: Daily automated snapshots
- **Clustering**: 3-node cluster for reliability

---

### 5. Object Storage (Spaces)

**DigitalOcean Spaces Configuration**:
```
Bucket: tokkatot-production
├── firmware/              # ESP32 firmware binaries
│   ├── v2.0.0/
│   ├── v2.0.1/
│   └── versions.json
├── backups/              # Database backups
│   ├── postgres/daily/
│   ├── influxdb/weekly/
├── configs/              # Device configurations
│   ├── farm-001/
│   └── farm-002/
└── media/                # User uploads (optional)
    ├── farm-photos/
    └── reports/

CDN: Enabled (edge caching)
Versioning: Enabled (object versioning)
CORS: Configured for frontend access
SSL: Required
```

**Access Control**:
- API access via S3-compatible keys
- Signed URLs for firmware downloads
- Public read for static assets
- Private for sensitive files

---

### 6. Load Balancing & Networking

**Load Balancer**:
- **Type**: DigitalOcean Load Balancer
- **Protocol**: HTTPS only (redirect HTTP to HTTPS)
- **Certificates**: Let's Encrypt via cert-manager (auto-renewal)
- **Health Check**: TCP 8080 every 10 seconds
- **Sticky Sessions**: Not enabled (stateless service)
- **Rate Limiting**: 1000 req/min per IP at LB level

**Firewall Rules**:
```
Ingress:
  - HTTPS (443) from anywhere (users)
  - HTTP (80) from anywhere (redirect only)
  - MQTT (1883) from on-farm devices
  - Custom MQTT over TLS (8883) for secure device comms

Egress:
  - All outbound allowed (external APIs, email, etc)

Kubernetes Internal:
  - Pod-to-pod communication allowed
  - Service-to-service communication via DNS
```

**DNS**:
- **Domain**: tokkatot.farm (example)
- **Provider**: DigitalOcean DNS
- **Records**:
  ```
  tokkatot.farm       A → Load Balancer IP
  api.tokkatot.farm   CNAME → tokkatot.farm
  mqtt.tokkatot.farm  A → MQTT Broker IP
  cdn.tokkatot.farm   CNAME → Spaces CDN
  ```

---

## CI/CD Pipeline (GitHub Actions)

### Build Process

**Trigger**: Push to main branch or manual trigger

**Steps**:
1. **Checkout code** from GitHub
2. **Build Docker images**
   - Backend: `docker build -t tokkatot-backend:v${VERSION} .`
   - Frontend: `docker build -t tokkatot-frontend:v${VERSION} .`
3. **Run tests**
   - Unit tests: `go test ./...`
   - Integration tests: `npm run test:integration`
   - Linting: `golangci-lint run`
4. **Push to registry**
   - Tag: `v2.0.0` → DOCR as `tokkatot-backend:v2.0.0`
5. **Deploy to staging** (manual approval)
6. **Deploy to production** (manual approval)

### Example GitHub Actions Workflow

```yaml
name: Build & Deploy

on:
  push:
    branches: [main]
  workflow_dispatch:  # Manual trigger

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Build Backend
        run: docker build -t tokkatot-backend:${{ github.sha }} ./middleware
        
      - name: Push to DOCR
        run: |
          docker tag tokkatot-backend:${{ github.sha }} \
            registry.digitalocean.com/tokkatot/backend:${{ github.sha }}
          docker push registry.digitalocean.com/tokkatot/backend:${{ github.sha }}
        env:
          REGISTRY_TOKEN: ${{ secrets.DOCR_TOKEN }}
          
      - name: Deploy to K8s
        run: |
          kubectl set image deployment/api-server \
            api-server=registry.digitalocean.com/tokkatot/backend:${{ github.sha }} \
            -n production
```

---

## Monitoring & Observability

### Prometheus Metrics

**Scrape Targets**:
- Kubernetes nodes (CPU, memory, disk)
- Pod metrics (requests, latency, errors)
- Service endpoints (/metrics at port 9090)
- PostgreSQL exporter
- Redis exporter

**Key Metrics**:
```
api_request_duration_seconds (histogram)
api_request_total (counter)
api_errors_total (counter)
device_commands_queued (gauge)
database_connections_active (gauge)
cache_hit_ratio (gauge)
```

### Grafana Dashboards

**Dashboards Available**:
1. **System Health** - CPU, memory, disk, network
2. **API Performance** - Latency, error rates, throughput
3. **Database** - Connections, query performance, replication lag
4. **Devices** - Online/offline, command success, data flow
5. **Business Metrics** - Active farms, active users, commands/day

### AlertManager Alerts

**Alert Rules**:
```
alert: HighErrorRate
  if: rate(api_errors_total[5m]) > 0.05
  for: 5m
  action: send to Slack/email
  
alert: DatabaseDown
  if: pg_up != 1
  for: 1m
  action: page on-call engineer
  
alert: PodRestartingTooOften
  if: rate(kube_pod_container_status_restarts_total[15m]) > 0.1
  for: 5m
  action: send to Slack
```

---

## Backup & Disaster Recovery

### Backup Strategy

**PostgreSQL**:
- **Frequency**: Daily at 2 AM UTC
- **Type**: Full backup + WAL archiving
- **Retention**: 7 days for daily, 4 weeks for weekly
- **Location**: DigitalOcean Spaces (geo-redundant)
- **Verification**: Monthly restore drill

**Redis**:
- **Frequency**: Every 6 hours
- **Type**: RDB snapshots
- **Retention**: 14 days
- **Location**: Spaces object storage

**InfluxDB**:
- **Frequency**: Daily
- **Type**: Full database backup
- **Retention**: 30 days
- **Location**: Spaces object storage

**Application Code**:
- Git repository with all history
- Release tags for each version

### Disaster Recovery RTO/RPO

| Component | RTO | RPO | Method |
|-----------|-----|-----|--------|
| **API Service** | 5 min | 0 | Kubernetes auto-restart |
| **PostgreSQL** | 15 min | 1 hour | Automated failover |
| **Redis** | 10 min | 0 | Automated failover |
| **InfluxDB** | 20 min | 1 hour | Backup restore |
| **Full Cluster** | 4 hours | 1 hour | Infrastructure rebuild |

### Recovery Procedures

**Database Corruption**:
1. Detect via monitoring alert
2. Stop affected service
3. Restore from last clean backup
4. Validate data integrity
5. Resume service
6. Estimated time: 30-60 minutes

**Data Center Outage**:
1. Kubernetes saves state in etcd
2. All managed databases have failover ready
3. Spin up new cluster in different region
4. Restore volumes from snapshots
5. Update DNS to new load balancer
6. Estimated time: 2-4 hours

---

## Infrastructure as Code (Terraform)

**IaC Tool**: Terraform or Pulumi (recommended)

**Terraform Module Structure**:
```
terraform/
├── main.tf            # Main infrastructure
├── kubernetes.tf      # K8s cluster setup
├── databases.tf       # Managed DB provisioning
├── networking.tf      # Load balancer, firewall
├── storage.tf         # Spaces configuration
├── monitoring.tf      # Prometheus, Grafana
├── variables.tf       # Input variables
├── terraform.tfvars   # Variable values
└── outputs.tf         # Output values
```

**Key Resources**:
```terraform
# Kubernetes Cluster
resource "digitalocean_kubernetes_cluster" "main" {
  name    = "tokkatot-prod"
  version = "1.28"
  nodes   = 3
  ...
}

# PostgreSQL Database
resource "digitalocean_database_cluster" "postgres" {
  name    = "tokkatot-postgres"
  size    = "db-s-2vcpu-4gb"
  ...
}

# Spaces Bucket
resource "digitalocean_spaces_bucket" "firmware" {
  name   = "tokkatot-production"
  region = "sgp1"
  ...
}
```

---

## Environment Configuration

### Environment Variables

**Production Environment (.env.production)**:
```
# Server
API_HOST=0.0.0.0
API_PORT=8080
LOG_LEVEL=info

# Database
DATABASE_URL=postgresql://user:pass@db.cloud.digitalocean.com:25060/tokkatot_prod
DATABASE_POOL_SIZE=50
DATABASE_SSL_MODE=require

# Redis
REDIS_URL=rediss://default:pass@redis.db.cloud.digitalocean.com:25061

# InfluxDB
INFLUXDB_URL=https://db.cloud.digitalocean.com:443
INFLUXDB_ORG=tokkatot
INFLUXDB_BUCKET=sensor_readings

# JWT
JWT_SECRET=<64-char-random-secret>
JWT_EXPIRATION=24h

# MQTT
MQTT_BROKER=mqtt://localhost:1883
MQTT_USERNAME=mqtt_user
MQTT_PASSWORD=<mqtt-password>
MQTT_TLS_ENABLED=false  # false if internal, true if external

# Storage
STORAGE_ENDPOINT=tokkatot-production.nyc3.digitaloceanspaces.com
STORAGE_BUCKET=tokkatot-production
STORAGE_KEY=<DO-access-key>
STORAGE_SECRET=<DO-secret>

# Monitoring
PROMETHEUS_URL=http://prometheus:9090
GRAFANA_URL=https://grafana.tokkatot.farm

# In-App Notifications Only (No external providers needed)
NOTIFICATION_RETENTION_DAYS=90
NOTIFICATION_AUDIT_RETENTION_YEARS=2
```

---

## Security Considerations

### Network Security
- All traffic encrypted in transit (TLS 1.3+)
- DDoS protection at load balancer level
- Rate limiting on API endpoints
- VPC isolation for database access

### Access Control
- Kubernetes RBAC for pod permissions
- Database user with limited privileges
- Service accounts with minimal permissions
- API key rotation quarterly

### Data Security
- Encryption at rest for databases
- Encryption at rest for object storage
- Password hashing with bcrypt
- Sensitive data (tokens, keys) in secrets manager

---

## Cost Estimation (DigitalOcean)

### Monthly Infrastructure Cost

| Service | Size | Cost |
|---------|------|------|
| Kubernetes Cluster | 3x s-2vcpu-4gb nodes | $90 |
| PostgreSQL | db-s-2vcpu-4gb | $60 |
| Redis | db-s-1vcpu-1gb | $15 |
| InfluxDB | db-s-2vcpu-8gb | $120 |
| Spaces | 500 GB storage | $25 |
| Load Balancer | 1x LB | $12 |
| Bandwidth | (~500 GB/month) | $50 |
| **Total** | | **~$372/month** |

**Scaling**:
- Additional K8s nodes: +$30 each
- Database storage: +$0.05/GB
- Bandwidth: +$0.10/GB overage

---

## Deployment Checklist

- [ ] Infrastructure provisioned and tested
- [ ] DNS configured and pointing to load balancer
- [ ] SSL certificates installed
- [ ] Database backups tested
- [ ] Monitoring dashboards created
- [ ] Alert thresholds configured
- [ ] Deployment procedures documented
- [ ] Rollback procedures tested
- [ ] Security audit completed
- [ ] Load testing passed (target: 1000 concurrent users)
- [ ] Disaster recovery drill completed
- [ ] Team trained on operational procedures

---

**Version History**
| Version | Date | Changes |
|---------|------|---------|
| 2.0 | Feb 2026 | Cloud-only deployment specification |
| | | Removed self-hosted option (Option 3) |
| | | DigitalOcean as primary recommendation |

**Related Documents**
- SPECIFICATIONS_ARCHITECTURE.md
- SPECIFICATIONS_SECURITY.md
- SPECIFICATIONS_TECHNOLOGY_STACK.md
