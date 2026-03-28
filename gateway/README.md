# 🥧 Tokkatot Raspberry Pi Gateway

This Python service runs on the Raspberry Pi gateway located at the farm. It acts as the "Reliable Bridge" between small-footprint ESP32 sensor nodes and the high-performance Tokkatot Cloud Backend.

---

## 🚀 Key Features

*   **Zero-Config Discovery:** Automatically discovers local sensors using mDNS (**`tokkatot-sensor.local`**). No manual IP entries required.
*   **Full-Loop Relay:** Simultaneously polls sensors for telemetry and fetches pending control commands (Fan/Heater) from the cloud.
*   **Offline Resilience:** Integrated **SQLite Queue**. If your 3G/4G connection drops, telemetry is cached locally and synced automatically when the link is restored.
*   **Production Secure:** Pairs with your cloud dashboard via a **Secure 10-Digit Code** for zero-config provisioning.

---

## 📦 Installation (Production)

For a fast setup on a fresh Raspberry Pi (Ubuntu Server 24.04), follow these steps:

### 1. Flash the OS
Use **Raspberry Pi Imager** to flash Ubuntu Server.
- **Settings Cog:** Enable SSH, set username `tokkatot`, and pre-configure your WiFi.

### 2. Run the Tokkatot Setup
SSH into your Pi (`ssh tokkatot@tokkatot.local`) and paste the following:
```bash
# Clone the repository
git clone https://github.com/SirOsborn/Tokkatot.git
cd Tokkatot/gateway/scripts

# Run the automated production setup
chmod +x pi_setup.sh
./pi_setup.sh
```

### 3. Pair with Cloud
Wait for the setup to finish, then run:
```bash
cd ~/tokkatot-gateway/gateway
./venv/bin/python3 main.py
```
- Enter your **10-digit Pairing Code** from the Tokkatot web dashboard.
- Once paired, press **Ctrl+C**.

### 4. Enable Persistence
To ensure the gateway starts on every reboot:
```bash
sudo systemctl start tokkatot-gateway
```

---

## 🛠️ Development & Debugging

- **Logs:** `journalctl -u tokkatot-gateway -f`
- **Manual Poll:** `./venv/bin/python3 main.py --poll-once`
- **DB Check:** `sqlite3 telemetry_queue.db "SELECT * FROM telemetry_queue;"`

© 2026 Tokkatot Agri-Tech. All rights reserved.
