# Tokkatot Deployment Guide (AWS Cloud)

This guide walks you through deploying the Tokkatot Agri-Tech Platform to an AWS EC2 instance with a fully automated CI/CD pipeline.

## 1. Server Prerequisites (Ubuntu 22.04+ Recommended)

Connect to your EC2 instance and install Docker/Docker Compose:

```bash
# Update and install Docker
sudo apt update && sudo apt install -y docker.io docker-compose
sudo usermod -aG docker $USER
newgrp docker
```

## 2. Generate Production Secrets

### 2a. VAPID Keys (for Push Notifications)
Run the generator tool to create your unique keys:

```bash
# From the root of the project
go run middleware/scripts/generate_vapid/main.go
```

### 2b. Database & JWT
Prepare your random passwords for your production `.env`.

## 3. GitHub Secrets Configuration

Go to your GitHub Repository -> **Settings** -> **Secrets and variables** -> **Actions** and add:

| Secret Name | Description |
| :--- | :--- |
| `SERVER_HOST` | Your EC2 Public IP or Domain (`app.tokkatot.com`) |
| `SERVER_USER` | The SSH user (usually `ubuntu`) |
| `SSH_PRIVATE_KEY` | Your `.pem` file content (RSA/ED25519) |
| `DB_PASSWORD` | Secure password for PostgreSQL |
| `JWT_SECRET` | Long random string for authentication |
| `VAPID_PUBLIC_KEY` | From Step 2a |
| `VAPID_PRIVATE_KEY` | From Step 2a |

## 4. Initial SSL Setup (Certbot)

The Nginx configuration expects Let's Encrypt certificates. On your first deployment:

1.  Clone the repository to the server: `git clone https://github.com/USER/tokkatot.git ~/tokkatot`
2.  Populate `~/tokkatot/.env` with your secrets.
3.  Run the stack temporarily without SSL to get certs:
    ```bash
    docker-compose up -d nginx certbot
    ```
4.  Run Certbot to generate certs:
    ```bash
    docker exec tokkatot-certbot certbot certonly --webroot -w /var/www/certbot -d app.tokkatot.com --email info@tokkatot.com --agree-tos --no-eff-email
    ```
5.  Restart the full stack: `docker-compose up -d`

## 5. CI/CD Workflow (GitHub Actions)

- **Push to `dev`**: Automatically deploys the latest fixes to your development environment.
- **Push to `stage`**: Automatically deploys to your staging server/environment.
- **Push to `main`**: Automatically deploys to the production-grade `app.tokkatot.com`.

The pipeline builds Docker images, pushes them to **GitHub Container Registry (GHCR)**, and triggers a pull-and-restart on your EC2 instance via SSH.
