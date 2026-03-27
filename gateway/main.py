import os
import time
import requests
import sqlite3
import urllib3
import json
from datetime import datetime
from dotenv import load_dotenv
import provisioning

# Suppress insecure HTTPS warning for local ESP32 self-signed certs
urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)

# Load environment configuration
load_dotenv()

ESP32_IP = os.getenv("ESP32_IP", "192.168.1.100")
CLOUD_API_URL = os.getenv("CLOUD_API_URL", "http://localhost:3000")
GATEWAY_TOKEN = os.getenv("GATEWAY_TOKEN", "")
FARM_ID = os.getenv("FARM_ID", "")
COOP_ID = os.getenv("COOP_ID", "")
HARDWARE_ID = os.getenv("HARDWARE_ID", "PI_GATEWAY_001")
POLL_INTERVAL = int(os.getenv("POLL_INTERVAL", "5"))

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
        # verify=False is required for ESP32 self-signed certs
        response = requests.get(url, verify=False, timeout=3)
        if response.status_code == 200:
            return response.json()
        print(f"[{datetime.now()}] Warning: ESP32 returned status {response.status_code}")
    except requests.exceptions.RequestException as e:
        print(f"[{datetime.now()}] Error connecting to ESP32: {e}")
    return None

def push_to_cloud(payload):
    """Pushes formatted telemetry data to the cloud via persistent gateway token."""
    url = f"{CLOUD_API_URL}/v1/farms/{FARM_ID}/coops/{COOP_ID}/telemetry"
    headers = {
        "X-Gateway-Token": GATEWAY_TOKEN,
        "Content-Type": "application/json"
    }
    try:
        response = requests.post(url, headers=headers, json=payload, timeout=5)
        if response.status_code in [200, 201]:
            return True
        else:
            print(f"[{datetime.now()}] Cloud push failed. Status {response.status_code}: {response.text}")
    except requests.exceptions.RequestException as e:
        print(f"[{datetime.now()}] Network error pushing to cloud: {e}")
    return False

def send_heartbeat():
    """Reports gateway health status to the cloud."""
    url = f"{CLOUD_API_URL}/v1/gateway/heartbeat"
    headers = {"X-Gateway-Token": GATEWAY_TOKEN}
    try:
        requests.post(url, headers=headers, timeout=3)
    except:
        pass

def queue_locally(conn, payload):
    """Saves telemetry payload to the local database for later retry."""
    try:
        c = conn.cursor()
        c.execute('INSERT INTO telemetry_queue (payload) VALUES (?)', (json.dumps(payload),))
        conn.commit()
    except Exception as e:
        print(f"[{datetime.now()}] DB Sync error: {e}")

def process_queue(conn):
    """Attempts to push queued telemetry data if cloud is reachable."""
    try:
        c = conn.cursor()
        # Fetch the oldest 50 records
        c.execute('SELECT id, payload FROM telemetry_queue ORDER BY id ASC LIMIT 50')
        rows = c.fetchall()
        
        for row in rows:
            record_id, payload_str = row
            payload = json.loads(payload_str)
            
            if push_to_cloud(payload):
                # Successfully pushed, remove from queue
                c.execute('DELETE FROM telemetry_queue WHERE id = ?', (record_id,))
                conn.commit()
            else:
                # Still failing, stop processing queue until next cycle
                break
    except Exception as e:
        print(f"[{datetime.now()}] Queue process error: {e}")

def main():
    global GATEWAY_TOKEN, FARM_ID, COOP_ID
    
    # 1. Check if provisioning is needed
    if not GATEWAY_TOKEN or not FARM_ID or not COOP_ID:
        print("Starting Zero-Config Setup...")
        if not provisioning.run_setup_flow(CLOUD_API_URL):
            print("Setup failed. Please check network and try again.")
            return
        
        # Reload env after setup
        load_dotenv()
        GATEWAY_TOKEN = os.getenv("GATEWAY_TOKEN")
        FARM_ID = os.getenv("FARM_ID")
        COOP_ID = os.getenv("COOP_ID")

    print(f"[{datetime.now()}] Tokkatot Gateway started.")
    print(f"Farm: {FARM_ID} | Coop: {COOP_ID} | Interval: {POLL_INTERVAL}s")
    
    conn = init_db()
    last_heartbeat = 0
    
    try:
        while True:
            cycle_start = time.time()
            
            # 2. Send heartbeat every 60 seconds
            if time.time() - last_heartbeat > 60:
                send_heartbeat()
                last_heartbeat = time.time()

            # 3. Poll local ESP32
            raw_data = fetch_esp32_data()
            if raw_data:
                # Format payload for Tokkatot Cloud
                payload = {
                    "hardware_id": HARDWARE_ID,
                    "sensors": {
                        "temperature_c": float(raw_data.get("temperature", 0)),
                        "humidity_pct": float(raw_data.get("humidity", 0)),
                        "water_level_raw": float(raw_data.get("water_level", 0))
                    }
                }
                
                # Try to push directly, otherwise queue
                if not push_to_cloud(payload):
                    queue_locally(conn, payload)
                
                # Attempt to clear queue if cloud is back up
                process_queue(conn)
            
            # Maintain stable loop interval
            elapsed = time.time() - cycle_start
            sleep_time = max(0.0, POLL_INTERVAL - elapsed)
            time.sleep(sleep_time)

    except KeyboardInterrupt:
        print("\nGateway shutting down.")
    finally:
        conn.close()

if __name__ == "__main__":
    main()
