#!/bin/bash
# ==============================================================================
# Tokkatot Raspberry Pi Gateway - Production Installation Script
# ==============================================================================
# This script installs the Tokkatot Gateway software on a fresh Ubuntu Server.
# It assumes you are already connected to your home WiFi.
# ==============================================================================

set -e

# --- Configuration ---
TARGET_DIR="$HOME/tokkatot-gateway"
REPO_URL="https://github.com/SirOsborn/Tokkatot.git"

echo "🚀 Starting Tokkatot Gateway Production Installation..."

# 1. Install System Dependencies
echo "📦 Installing system dependencies..."
sudo apt update && sudo apt install -y python3 python3-pip python3-venv sqlite3 git curl

# 2. Deploy Gateway Code
echo "🚚 Deploying gateway code..."
mkdir -p "$TARGET_DIR"

if [ -d "$TARGET_DIR/.git" ]; then
    echo "Updating existing repository..."
    cd "$TARGET_DIR"
    git fetch --all
    git reset --hard origin/main
else
    echo "Cloning repository..."
    git clone $REPO_URL "$TARGET_DIR"
    cd "$TARGET_DIR"
fi

# 3. Setup Python Virtual Environment
echo "🐍 Setting up Python environment..."
cd "$TARGET_DIR/gateway"
python3 -m venv venv
./venv/bin/pip install --upgrade pip
./venv/bin/pip install -r requirements.txt

# 4. Create Systemd Service for Persistence
echo "⚙️ Creating Tokkatot Gateway Service..."
SERVICE_FILE="/etc/systemd/system/tokkatot-gateway.service"

sudo bash -c "cat <<EOF > $SERVICE_FILE
[Unit]
Description=Tokkatot Smart Gateway
After=network-online.target
Wants=network-online.target

[Service]
ExecStart=$(pwd)/venv/bin/python3 $(pwd)/main.py
WorkingDirectory=$(pwd)
StandardOutput=inherit
StandardError=inherit
Restart=always
User=$USER

[Install]
WantedBy=multi-user.target
EOF"

sudo systemctl daemon-reload
sudo systemctl enable tokkatot-gateway

echo ""
echo "=============================================================================="
echo "🎉 INSTALLATION COMPLETE!"
echo "=============================================================================="
echo "1. Your Raspberry Pi is now ready to act as a Smart Gateway."
echo "2. The 'tokkatot-gateway' service is enabled (will start on boot)."
echo ""
echo "🔒 FINAL STEP: Link your Pi to your Cloud Account now:"
echo "------------------------------------------------------------------------------"
echo "   cd $(pwd) && ./venv/bin/python3 main.py"
echo "------------------------------------------------------------------------------"
echo ""
echo "Once the 10-digit pairing is successful, start the background service with:"
echo "   sudo systemctl start tokkatot-gateway"
echo "=============================================================================="
