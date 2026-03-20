# Automation Use Cases

This document captures real automation scenarios supported today.

---

## 1) Feeder Motor (Scheduled)
**Goal:** deliver feed multiple times per day.
- Start time: 06:00
- Run for: 8 minutes
- Off for: 4 hours
- Repeat: 3 times

Expected behavior:
- Gateway runs sequence locally on the coop.
- If internet drops, the sequence still runs.

---

## 2) Conveyor Belt (Scheduled)
**Goal:** remove manure periodically.
- Start time: 08:00
- Run for: 10 minutes
- Off for: 2 hours
- Repeat: 4 times

---

## 3) Temperature Threshold Control
**Goal:** keep coop within farmer’s desired range.
- If temp < min → heater ON
- If temp > max → fan ON
- Otherwise → both OFF

Thresholds are stored per coop and executed by the gateway.

---

## 4) Water Level Alert (Monitor-only)
**Goal:** detect broken floating valve.
- If water level below half threshold for 1 minute → alert
- No actuator control (monitor-only)

---

## 5) Mixed Hardware (Optional Devices)
**Goal:** support partial hardware deployments.
- If a device is not installed, it is marked **inactive** by the gateway.
- Schedules only list **active actuators**.
