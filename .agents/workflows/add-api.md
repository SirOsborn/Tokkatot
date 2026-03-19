---
description: How to add a new API endpoint in the Go middleware
---

1.  **Define the Schema**:
    *   Add request/response structs in `middleware/schemas/`.
2.  **Implement Service Logic**:
    *   Add the necessary business logic in the corresponding file in `middleware/services/`.
3.  **Create the Handler**:
    *   Add the HTTP handler in `middleware/api/`.
    *   **CRITICAL**: Use `GetUserIDFromContext(c)` or `GetFarmIDFromContext(c)` for ID extraction.
    *   **CRITICAL**: Add full Swagger annotations (Summary, Description, Tags, Produce, Param, Success, Router).
4.  **Register the Route**:
    *   Add the route in `middleware/main.go`.
5.  **Verify**:
    *   Rebuild the backend: `go build -o backend.exe main.go`.
    *   Run a verification script or use `curl`/`Postman` to test the endpoint.
