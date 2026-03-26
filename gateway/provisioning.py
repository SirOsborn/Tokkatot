import os
import time
import requests
import json
import uuid
import platform
import socket

def get_hardware_info():
    """Returns a unique hardware ID for this device."""
    try:
        # Try to get Raspberry Pi serial number
        with open('/proc/cpuinfo','r') as f:
            for line in f:
                if line.startswith('Serial'):
                    return "PI_" + line.split(':')[1].strip()
    except:
        pass
    
    # Fallback to mac address + hostname
    mac = ':'.join(['{:02x}'.format((uuid.getnode() >> ele) & 0xff) for ele in range(0, 8*6, 8)][::-1])
    return f"GW_{socket.gethostname()}_{mac.replace(':','')}"

def run_setup_flow(api_url):
    """
    Zero-Config Setup Flow:
    1. Register with cloud to get a setup code.
    2. Display code to user.
    3. Poll cloud until user 'claims' this device.
    4. Save received config to .env.
    """
    hardware_id = get_hardware_info()
    print("\n" + "="*50)
    print("Welcome to Tokkatot Gateway Setup")
    print("="*50)
    print(f"Device ID: {hardware_id}")
    
    try:
        # 1. Request Setup Code
        response = requests.post(f"{api_url}/v1/gateway/provision/request", json={"hardware_id": hardware_id}, timeout=10)
        if response.status_code != 200:
            print(f"Error: Could not connect to Tokkatot Cloud ({response.status_code})")
            return False
        
        data = response.json().get("data", {})
        setup_code = data.get("setup_code")
        
        print("\nACTION REQUIRED:")
        print(f"Your Setup Code is: \033[1;32m{setup_code}\033[0m")
        print("\nPlease log into your Tokkatot Dashboard, go to 'Add Gateway',")
        print("and enter the code above to link this device to your farm.")
        print("\nWaiting for pairing...")
        
        # 2. Poll for status
        start_time = time.time()
        timeout = 15 * 60 # 15 minutes
        
        while time.time() - start_time < timeout:
            status_res = requests.get(f"{api_url}/v1/gateway/provision/status/{setup_code}", timeout=10)
            if status_res.status_code == 200:
                status_data = status_res.json().get("data", {})
                if status_data.get("is_claimed"):
                    print("\n\033[1;32mSUCCESS: Gateway Claimed!\033[0m")
                    
                    # 3. Save to .env
                    config = {
                        "CLOUD_API_URL": api_url,
                        "GATEWAY_TOKEN": status_data.get("token"),
                        "FARM_ID": status_data.get("farm_id"),
                        "COOP_ID": status_data.get("coop_id"),
                        "HARDWARE_ID": hardware_id
                    }
                    
                    with open(".env", "w") as f:
                        for key, value in config.items():
                            f.write(f"{key}={value}\n")
                    
                    print("Configuration saved to .env")
                    return True
            
            time.sleep(5)
            print(".", end="", flush=True)
            
        print("\nSetup timed out. Please restart the gateway to try again.")
        return False
        
    except Exception as e:
        print(f"\nConnection Error: {e}")
        return False
