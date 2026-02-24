/**
 * Tokkatot Component Loader
 * Injects shared header + navbar, highlights active nav, applies i18n.
 * Global: window.loadComponents(options)
 *
 * Usage in each page (before </body>):
 *   await loadComponents();
 *
 * HTML placeholders:
 *   <div id="header-placeholder"></div>
 *   <div id="navbar-placeholder"></div>
 */

(function (window) {
  'use strict';

  /* ---------- helpers ---------- */

  function injectHTML(el, html) {
    /* createContextualFragment executes <script> tags — innerHTML does not */
    var range = document.createRange();
    range.selectNode(el);
    el.appendChild(range.createContextualFragment(html));
  }

  function highlightActiveNav() {
    var path = window.location.pathname;
    var map = {
      '/':                  'home',
      '/monitoring':        'monitoring',
      '/disease-detection': 'disease',
      '/schedules':         'schedules',
      '/settings':          'settings',
      '/profile':           'settings'   /* profile lives under settings tab */
    };
    var active = map[path] || 'home';
    document.querySelectorAll('#app-navbar .nav-item').forEach(function (item) {
      if (item.getAttribute('data-nav') === active) {
        item.classList.add('active');
      } else {
        item.classList.remove('active');
      }
    });
  }

  function updateLangButton() {
    var btn = document.getElementById('lang-toggle-btn');
    if (btn) btn.textContent = (window.i18n ? window.i18n.getLang() : 'km').toUpperCase();
  }

  function loadFarmName() {
    var name = localStorage.getItem('farm_name');
    var el = document.getElementById('header-farm-name');
    if (el && name) el.textContent = name;
  }

  /* ---------- public API ---------- */

  window.loadComponents = async function (options) {
    options = options || {};
    var doHeader = options.header !== false;
    var doNavbar = options.navbar !== false;

    var tasks = [];

    if (doHeader) {
      tasks.push(
        fetch('/components/header.html')
          .then(function (r) { return r.text(); })
          .then(function (html) {
            var el = document.getElementById('header-placeholder');
            if (el) injectHTML(el, html);
          })
          .catch(function (e) { console.warn('[Components] header:', e); })
      );
    }

    if (doNavbar) {
      tasks.push(
        fetch('/components/navbar.html')
          .then(function (r) { return r.text(); })
          .then(function (html) {
            var el = document.getElementById('navbar-placeholder');
            if (el) injectHTML(el, html);
          })
          .catch(function (e) { console.warn('[Components] navbar:', e); })
      );
    }

    await Promise.all(tasks);

    /* Post-inject setup */
    highlightActiveNav();
    updateLangButton();
    loadFarmName();

    if (window.i18n && window.i18n.applyAll) {
      window.i18n.applyAll();
    }
  };

  /* Expose toggle so header button can call it */
  window.headerToggleLang = function () {
    if (window.i18n) {
      window.i18n.toggleLang();
      updateLangButton();
      if (window.i18n.applyAll) window.i18n.applyAll();
    }
  };

})(window);
