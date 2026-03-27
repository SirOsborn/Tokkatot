---
description: How to add a new frontend page to the Tokkatot PWA
---

# Tokkatot Frontend Development

Tokkatot uses a zero-build, mobile-first, Vue.js-based PWA architecture. Follow these strict rules when adding new pages.

## 🧱 1. HTML Template Setup
1.  **File Location**: Always Create `frontend/pages/new_page.html`.
2.  **Auth & Security (CRITICAL)**:
    -   You MUST include the following **Blocking Auth Script** inside the `<head>` of any protected page. This script redirects unauthorized users immediately:
        ```html
        <script>
          // Pre-render auth check (no UI flash)
          if (!localStorage.getItem('access_token')) {
            window.location.href = '/login';
          }
        </script>
        ```
3.  **Styles & Assets**:
    -   Include the Tokkatot design system files in order: `theme.css`, `components.css`, `layout.css`.
    -   Initialize the Vue app within `<div id="app" class="page-container">`.

## 🎨 2. Design System Tokens
-   Never use ad-hoc colors. Use the **Tokkatot Design System** variables:
    -   `--tok-teal`: Core brand primary.
    -   `--tok-surface-light`: Soft card backgrounds.
    -   `--tok-alert-red`: Warning and threshold violations.

## 🔗 3. Registration & Routing
1.  **Register Route**: Add the endpoint to `middleware/main.go` under `setupRoutes`.
    ```go
    app.Get("/new-page", func(c *fiber.Ctx) error {
        return c.SendFile(filepath.Join(frontendPath, "pages", "new_page.html"))
    })
    ```
2.  **Update Navigation**: Include the new page link in the Sidebar or Navbar components.

---
**Proprietary Software - Tokkatot Startup**
*For internal use only. Unauthorized distribution is prohibited.*
