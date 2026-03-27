# Tokkatot Deployment Guide (AWS Cloud)

This guide walks you through the 3-tier production-grade deployment of the Tokkatot Agri-Tech Platform to AWS EC2.

## 1. Multi-Stage Pipeline Architecture

The platform uses a strict promotion flow via GitHub Actions:
- **`dev` branch**: Automated linting and Docker build checks.
- **`stage` branch**: Deploys to `staging.tokkatot.com`. Automatically seeds **test farmers and demo coops**.
- **`main` branch**: Deploys to `app.tokkatot.com`. **Production environment** — all test seeding is disabled.

## 2. Server Infrastructure (Ubuntu 22.04+)

Connect to your EC2 instance and install Docker:
```bash
sudo apt update && sudo apt install -y docker.io docker-compose
sudo usermod -aG docker $USER && newgrp docker
```

## 3. Initial "Chicken & Egg" SSL Setup
Since Nginx requires SSL certificates to start, but Certbot often needs Nginx to verify them, followed this standalone sequence on a fresh server:

1. **Stop any existing web server**: `sudo systemctl stop nginx || true`
2. **Generate Certs via Standalone Mode**:
   ```bash
   docker run --rm -v tokkatot_letsencrypt-certs:/etc/letsencrypt -p 80:80 certbot/certbot certonly --standalone -d app.tokkatot.com -d staging.tokkatot.com -m admin@tokkatot.com --agree-tos
   ```
3. **Deploy the Stack**: `docker-compose up -d`

## 4. Environment-Aware Seeding (TEST_FARMER)
The middleware is programmed to skip data seeding in production.
- **Staging only:** Add `TEST_FARMER_EMAIL` and `TEST_FARMER_PASSWORD` to your server's `~/tokkatot-staging/.env`. 
- On every container start, the middleware will idempotently ensure the demo farm/coop exists.

## 5. IoT Gateway Provisioning
The Raspberry Pi gateway uses a **Zero-Config Setup**. 
1. Run `./gateway/scripts/pi_setup.sh` on your Pi.
2. Enter the **10-digit pairing code** on the Tokkatot Web Dashboard to link the hardware to your account.
