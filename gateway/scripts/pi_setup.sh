#!/bin/bash

# ==============================================================================
# Tokkatot Raspberry Pi Gateway - Cloud Migration Script (Ubuntu Server)
# ==============================================================================
# This script safely transitions a Pi from an Access Point (AP) to a 
# Cloud-Connected Gateway without getting locked out.
# ==============================================================================

set -e

# --- Configuration (EDIT THESE BEFORE RUNNING!) ---
WIFI_SSID="YOUR_ROUTER_SSID"
WIFI_PASSWORD="YOUR_ROUTER_PASSWORD"
TARGET_DIR="$HOME/tokkatot-gateway"
REPO_URL="https://github.com/SirOsborn/Tokkatot.git"

echo "🚀 Starting Tokkatot Gateway Migration..."

# 1. Install System Dependencies
echo "📦 Installing dependencies..."
sudo apt update && sudo apt install -y python3 python3-pip python3-venv sqlite3 git curl

# 2. Configure WiFi Client (Netplan) - SAFE MODE
echo "📡 Configuring WiFi Client..."
NETPLAN_FILE="/etc/netplan/99-tokkatot-wifi.yaml"

sudo bash -c "cat <<EOF > $NETPLAN_FILE
network:
  version: 2
  wifis:
    wlan0:
      optional: true
      access-points:
        \"$WIFI_SSID\":
          password: \"$WIFI_PASSWORD\"
      dhcp4: true
EOF"

sudo chmod 600 $NETPLAN_FILE
echo "Applying Netplan... (This may take a moment)"
sudo netplan apply

# 3. Verify Internet Connectivity
echo "🌐 Verifying internet connection..."
MAX_RETRIES=5
COUNT=0
until curl -s --head  --request GET http://www.google.com | grep "200 OK" > /dev/null || [ $COUNT -eq $MAX_RETRIES ]; do
  echo "Waiting for WiFi connection... ($((COUNT+1))/$MAX_RETRIES)"
  sleep 5
  COUNT=$((COUNT+1))
done

if [ $COUNT -eq $MAX_RETRIES ]; then
    echo "❌ ERROR: Could not connect to the internet. Aborting transition to keep AP active."
    exit 1
fi

echo "✅ Internet connection verified!"

# 4. Decommission Legacy AP Services
echo "🧹 Cleaning up legacy AP services (hostapd/dnsmasq)..."
sudo systemctl stop hostapd dnsmasq || true
sudo systemctl disable hostapd dnsmasq || true

# 5. Deploy New Gateway Code
echo "🚚 Deploying new gateway code..."
mkdir -p "$TARGET_DIR"

if [ -d "$TARGET_DIR/.git" ]; then
    cd "$TARGET_DIR"
    git fetch --all
    git reset --hard origin/main
else
    git clone $REPO_URL "$TARGET_DIR"
    cd "$TARGET_DIR"
fi

# 6. Setup Python Environment
echo "🐍 Setting up Python environment..."
# Move into the gateway directory
cd "$TARGET_DIR/gateway"
python3 -m venv venv
./venv/bin/pip install -r requirements.txt

# 7. Create Systemd Service
echo "⚙️ Creating Tokkatot Gateway Service..."
SERVICE_FILE="/etc/systemd/system/tokkatot-gateway.service"

sudo bash -c "cat <<EOF > $SERVICE_FILE
[Unit]
Description=Tokkatot Pi Gateway
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
echo "🎉 TRANSITION COMPLETE!"
echo "=============================================================================="
echo "1. Your Raspberry Pi is now a Cloud Gateway."
echo "2. The legacy Access Point has been disabled."
echo "3. TO LINK DEVICE: Run the command below and follow instructions on screen:"
echo ""
echo "   cd $(pwd) && ./venv/bin/python3 main.py"
echo ""
echo "Once the pairing is done, start the background service with:"
echo "   sudo systemctl start tokkatot-gateway"
echo "=============================================================================="
