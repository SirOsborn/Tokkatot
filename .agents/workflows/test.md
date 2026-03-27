---
description: How to run the various test suites in Tokkatot
---

# Tokkatot Test Suite

Tokkatot employs a comprehensive testing strategy that covers all parts of the application: frontend logic, API endpoints, and IoT gateway simulations.

## 🔗 1. API Health Check
Confirm that the middleware service is healthy by calling the `/v1/health` endpoint:
```bash
# Locally
curl http://localhost:3000/v1/health

# Production
curl https://app.tokkatot.com/v1/health
```

## 📡 2. IoT Gateway Simulation
To test the full telemetry flow from sensor to dashboard, use the `gateway-sim` Docker service:

1.  **Start the Simulator**:
    ```bash
    docker-compose up -d gateway-sim
    ```
2.  **Verify Telemetry**:
    -   Check the `gateway-sim` logs: `docker-compose logs -f gateway-sim`.
    -   Verify that it correctly posts telemetry to `/v1/farms/claimed-gateway/telemetry`.
    -   Confirm that data appears in the **Monitoring Dashboard** on the PWA.

## 🧪 3. Backend Unit Tests (Go)
Run the built-in Go tests:
```bash
cd middleware
go test ./...
```

## 🎨 4. Frontend Verification
1.  **Auth & Security**:
    -   Open a private browser window.
    -   Attempt to visit `http://localhost:3000/admin`.
    -   Confirm that you are **instantly redirected** to `/login` with no UI flash.
2.  **Responsive Design**:
    -   Use Chrome DevTools (F12) to test the app in "Mobile" mode (iPhone 12/Pixel 5).
    -   Confirm all buttons and dashboards fit the screen.

---
**Proprietary Software - Tokkatot Startup**
*For internal use only. Unauthorized distribution is prohibited.*
