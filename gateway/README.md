# Tokkatot Raspberry Pi Gateway

This Python service runs on the Raspberry Pi located at the farm. It acts as the bridge between the local ESP32 controller (which doesn't have direct internet access) and the Tokkatot Cloud Backend.

## Features
- Periodically polls the ESP32 for telemetry (Temperature, Humidity, Water Level).
- Pushes formatted telemetry to the Tokkatot Cloud Backend.
- Local SQLite queue cache: If the farm's internet 3G/4G connection drops, readings are stored locally and re-attempted when the connection is restored.

## Installation

1. Install Python 3 on the Raspberry Pi.
2. Clone/Copy this folder.
3. Install dependencies:
   ```bash
   pip install -r requirements.txt
   ```
4. Copy `.env.example` to `.env` and configure your keys.
   ```bash
   cp .env.example .env
   # Edit .env with nano or vim
   ```
5. Run the gateway:
   ```bash
   python main.py
   ```

## Production (Systemd Service)
For production, run this as a systemd service:

```ini
[Unit]
Description=Tokkatot Pi Gateway
After=network.target

[Service]
ExecStart=/usr/bin/python3 /path/to/gateway/main.py
WorkingDirectory=/path/to/gateway
StandardOutput=inherit
StandardError=inherit
Restart=always
User=pi

[Install]
WantedBy=multi-user.target
```
