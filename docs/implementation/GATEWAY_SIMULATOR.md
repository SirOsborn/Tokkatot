# Gateway Simulator (No Raspberry Required)

Use this script to simulate the Raspberry Pi gateway and test cloud communication.

---

## Location
`middleware/scripts/gateway_sim.go`

---

## Basic Run
```bash
go run middleware/scripts/gateway_sim.go --phone "+85512345678" --password "your_password" --create-farm --create-coop
```

---

## Common Options
- `--base-url` (default `http://127.0.0.1:3000`)
- `--email` or `--phone`
- `--password` (required)
- `--create-farm` / `--create-coop` if none exist
- `--farm-id` / `--coop-id` to target an existing farm/coop
- `--devices` comma list, e.g. `feeder_motor,fan,temp_humidity`
- `--inactive` comma list to simulate missing devices
- `--interval` seconds between telemetry (default 10)
- `--once` send a single telemetry event and exit
- `--poll-schedules` to fetch schedules each cycle
- `--e2e` run the full end-to-end test sequence once

---

## Example: Partial Hardware
```bash
go run middleware/scripts/gateway_sim.go \
  --phone "+85512345678" \
  --password "your_password" \
  --farm-id "YOUR_FARM_ID" \
  --coop-id "YOUR_COOP_ID" \
  --devices "feeder_motor,conveyor_belt,temp_humidity,water_level" \
  --inactive "water_level"
```

---

## What It Tests
- Device auto‑report
- Telemetry ingestion
- Monitoring timeline
- Water alert rules
- Schedule list
