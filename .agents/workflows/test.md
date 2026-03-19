---
description: How to run the various test suites in Tokkatot
---

1.  **Backend API Tests**:
    *   Navigate to `middleware/`.
    *   Run specialized verification scripts: `go run scripts/verify_*.go`.
    *   Use the PowerShell test suite: `.\test_all_endpoints.ps1`.
2.  **Database Checks**:
    *   Run `go run scripts/check_db.go` to verify schema integrity and active users.
3.  **Frontend Testing**:
    *   Open the app in Chrome/Firefox.
    *   Use the "Toggle Device Toolbar" (F12, Ctrl+Shift+M) to test mobile responsiveness.
    *   Verify offline support by toggling "Offline" in the Network tab.
4.  **AI Service Tests**:
    *   `cd ai-service`.
    *   Test health: `curl http://localhost:8000/health`.
    *   Test prediction: Use a sample image with `curl -X POST -F "image=@sample.jpg" http://localhost:8000/predict`.
