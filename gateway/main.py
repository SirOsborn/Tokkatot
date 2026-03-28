import os
import time
import requests
import sqlite3
import urllib3
import json
import uuid
import socket
from datetime import datetime
from dotenv import load_dotenv

# Suppress insecure HTTPS warning for local ESP32 self-signed certs
urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)

# Load environment configuration
load_dotenv()

ESP32_IP = os.getenv("ESP32_IP", "tokkatot-sensor.local")
CLOUD_API_URL = os.getenv("CLOUD_API_URL", "http://localhost:3000")
GATEWAY_TOKEN = os.getenv("GATEWAY_TOKEN", "")
FARM_ID = os.getenv("FARM_ID", "")
COOP_ID = os.getenv("COOP_ID", "")
POLL_INTERVAL = int(os.getenv("POLL_INTERVAL", "5"))

def get_unique_hardware_id():
    """Generates a unique hardware ID from CPU serial or MAC address."""
    serial = "UNKNOWN"
    try:
        # Try to get Raspberry Pi serial number
        if os.path.exists('/proc/cpuinfo'):
            with open('/proc/cpuinfo', 'r') as f:
                for line in f:
                    if line.startswith('Serial'):
                        serial = line.split(':')[1].strip()
                        return f"PI_{serial}"
    except:
        pass
    
    # Fallback to MAC address
    mac = ':'.join(['{:02x}'.format((uuid.getnode() >> ele) & 0xff) for ele in range(0,8*6,8)][::-1])
    return f"HW_{mac.replace(':', '').upper()}"

HARDWARE_ID = os.getenv("HARDWARE_ID", get_unique_hardware_id())

# Warning for localhost on Pi
if "localhost" in CLOUD_API_URL and os.path.exists('/proc/cpuinfo'):
    print(f"[!] WARNING: Using localhost for CLOUD_API_URL on a Raspberry Pi.")
    print(f"    Your Pi cannot reach the Middleware at localhost. Update your .env!")

DB_FILE = "telemetry_queue.db"

def init_db():
    """Initializes the local SQLite database for offline queuing."""
    conn = sqlite3.connect(DB_FILE)
    c = conn.cursor()
    c.execute('''
        CREATE TABLE IF NOT EXISTS telemetry_queue (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
            payload TEXT
        )
    ''')
    conn.commit()
    return conn

def fetch_esp32_data():
    """Polls the local ESP32 for current sensor readings."""
    url = f"https://{ESP32_IP}/get-current-data"
    try:
        response = requests.get(url, verify=False, timeout=3)
        if response.status_code == 200:
            return response.json()
    except Exception as e:
        print(f"[{datetime.now()}] Connection failed to ESP32: {e}")
    return None

def push_to_cloud(payload):
    """Pushes formatted telemetry data to the cloud."""
    if not GATEWAY_TOKEN or not FARM_ID or not COOP_ID:
        return False

    url = f"{CLOUD_API_URL}/v1/farms/{FARM_ID}/coops/{COOP_ID}/telemetry"
    headers = {
        "X-Gateway-Token": GATEWAY_TOKEN,
        "Content-Type": "application/json"
    }
    try:
        response = requests.post(url, headers=headers, json=payload, timeout=5)
        return response.status_code in [200, 201]
    except Exception as e:
        print(f"[{datetime.now()}] Cloud push error: {e}")
    return False

def send_heartbeat():
    """Reports gateway health status for Discovery or Active tracking."""
    if not GATEWAY_TOKEN:
        # DISCOVERY MODE: Check-in with hardware ID so Admin can see us
        url = f"{CLOUD_API_URL}/v1/devices/{HARDWARE_ID}/heartbeat"
        payload = {"status": "online", "response": "Discovery Mode (Unassigned)"}
        try:
            res = requests.post(url, json=payload, timeout=5)
            if res.status_code == 200:
                print(f"[{datetime.now()}] Discovery Check-in: {HARDWARE_ID} (Waiting for Admin)")
        except:
            pass
        return

    # ACTIVE MODE: Normal heartbeat with token
    url = f"{CLOUD_API_URL}/v1/gateway/heartbeat"
    headers = {"X-Gateway-Token": GATEWAY_TOKEN}
    payload = {"status": "online"}
    try:
        requests.post(url, headers=headers, json=payload, timeout=5)
        print(f"[{datetime.now()}] Heartbeat sent.")
    except Exception as e:
        print(f"[{datetime.now()}] Heartbeat error: {e}")

def fetch_cloud_commands():
    if not GATEWAY_TOKEN: return []
    url = f"{CLOUD_API_URL}/v1/gateway/commands/{HARDWARE_ID}"
    headers = {"X-Gateway-Token": GATEWAY_TOKEN}
    try:
        res = requests.get(url, headers=headers, timeout=5)
        if res.status_code == 200:
            return res.json().get("data", [])
    except:
        pass
    return []

def update_command_status(command_id, status, response_text):
    if not GATEWAY_TOKEN: return
    url = f"{CLOUD_API_URL}/v1/gateway/commands/{command_id}/status"
    headers = {"X-Gateway-Token": GATEWAY_TOKEN, "Content-Type": "application/json"}
    payload = {"status": status, "response": response_text}
    try:
        requests.post(url, headers=headers, json=payload, timeout=5)
    except:
        pass

def relay_to_esp32(command):
    cmd_type = command.get("command_type")
    mapping = {
        "fan_on": ("fan", True), "fan_off": ("fan", False),
        "heater_on": ("heater", True), "heater_off": ("heater", False),
        "feeder_on": ("feeder_motor", True), "feeder_off": ("feeder_motor", False),
        "conveyor_on": ("conveyor_belt", True), "conveyor_off": ("conveyor_belt", False)
    }
    if cmd_type not in mapping: return False, "Unknown command"
    
    endpoint, state = mapping[cmd_type]
    url = f"https://{ESP32_IP}/actuators/{endpoint}"
    payload = {"state": state, "duration": command.get("action_duration")}
    
    try:
        res = requests.post(url, json=payload, verify=False, timeout=5)
        return res.status_code == 200, res.text if res.status_code == 200 else "ESP32 Error"
    except Exception as e:
        return False, str(e)

def process_commands():
    cmds = fetch_cloud_commands()
    for cmd in cmds:
        print(f"[{datetime.now()}] Executing: {cmd.get('command_type')}")
        success, msg = relay_to_esp32(cmd)
        update_command_status(cmd.get("id"), "executed" if success else "failed", msg)

def queue_locally(conn, payload):
    try:
        conn.execute('INSERT INTO telemetry_queue (payload) VALUES (?)', (json.dumps(payload),))
        conn.commit()
    except:
        pass

def process_queue(conn):
    if not GATEWAY_TOKEN: return
    c = conn.cursor()
    c.execute('SELECT id, payload FROM telemetry_queue LIMIT 10')
    for row in c.fetchall():
        if push_to_cloud(json.loads(row[1])):
            conn.execute('DELETE FROM telemetry_queue WHERE id = ?', (row[0],))
            conn.commit()
        else:
            break

def main():
    print(f"Tokkatot Gateway Started | ID: {HARDWARE_ID}")
    conn = init_db()
    last_heartbeat = 0
    
    try:
        while True:
            # Heartbeat every 30s
            if time.time() - last_heartbeat > 30:
                send_heartbeat()
                last_heartbeat = time.time()

            if GATEWAY_TOKEN:
                data = fetch_esp32_data()
                if data:
                    payload = {"hardware_id": HARDWARE_ID, "sensors": data}
                    if not push_to_cloud(payload):
                        queue_locally(conn, payload)
                process_queue(conn)
                process_commands()
            
            time.sleep(POLL_INTERVAL)
    except KeyboardInterrupt:
        print("Shutting down...")
    finally:
        conn.close()

if __name__ == "__main__":
    main()
