---
description: How to add a new frontend page to the Tokkatot PWA
---

1.  **Create the HTML Template**:
    *   Create `frontend/pages/new_page.html`.
    *   Link modular styles in order:
        1. `frontend/css/theme.css` (Design System)
        2. `frontend/css/components.css` (Standard UI Elements)
        3. `frontend/css/layout.css` (Common Layout)
        4. `frontend/css/new_page.css` (Page-specific styles)
    *   Initialize the Vue app container: `<div id="app" class="page-container">`.
2.  **Add Styles and Scripts**:
    *   Create `frontend/css/new_page.css` for unique styling Needs. Use Tokkatot design tokens (`--tok-*`).
    *   Scripts: Include `i18n.js`, `api.js`, and `components.js` from `frontend/js/utils/`.
    *   Create page logic using Vue.js 3.
3.  **Register Static Route**:
    *   In `middleware/main.go`, serve the new HTML file at the desired route.
4.  **Update Navigation**:
    *   Add the new link to the sidebar/navbar components if needed.
5.  **Service Worker**:
    *   Update `frontend/sw.js` to include the new page and its CSS in `ASSETS_TO_CACHE`.
6.  **Verify**:
    *   Open the app in a browser and navigate to the new route.
    *   Confirm responsive design and Khmer translation support.
