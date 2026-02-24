/**
 * Tokkatot API Helper
 * Global: window.API
 * Usage:  const data = await API.get('/v1/farms/1/coops')
 *         await API.post('/v1/devices/123/command', { command_type: 'turn_on' })
 *
 * Features:
 *  - Automatic JWT injection from localStorage ('access_token')
 *  - Automatic token refresh on 401
 *  - Redirects to /login on auth failure
 *  - Returns parsed JSON (always)
 */

(function (window) {
  'use strict';

  const BASE = '';   // same origin — Go middleware proxies /v1/...

  const API = {
    /* -----------------------------------------------------------------------
       Internal request method
    ----------------------------------------------------------------------- */
    async request(url, options) {
      options = options || {};
      const token = localStorage.getItem('access_token');

      const headers = Object.assign(
        { 'Content-Type': 'application/json' },
        token ? { 'Authorization': 'Bearer ' + token } : {},
        options.headers || {}
      );

      let res;
      try {
        res = await fetch(BASE + url, Object.assign({}, options, { headers }));
      } catch (err) {
        console.error('[API] Network error:', err);
        throw { success: false, message: 'Network error. Check your connection.' };
      }

      // Handle 401 — try refresh then retry once
      if (res.status === 401) {
        const refreshed = await this._refresh();
        if (refreshed) {
          return this.request(url, options); // retry with new token
        } else {
          localStorage.removeItem('access_token');
          localStorage.removeItem('refresh_token');
          window.location.href = '/login';
          return;
        }
      }

      let data;
      try {
        data = await res.json();
      } catch (_) {
        data = { success: res.ok, status: res.status };
      }

      return data;
    },

    /* -----------------------------------------------------------------------
       Token refresh
    ----------------------------------------------------------------------- */
    async _refresh() {
      const refresh = localStorage.getItem('refresh_token');
      if (!refresh) return false;
      try {
        const res = await fetch(BASE + '/v1/auth/refresh', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ refresh_token: refresh })
        });
        if (!res.ok) return false;
        const data = await res.json();
        if (data.success && data.data && data.data.access_token) {
          localStorage.setItem('access_token', data.data.access_token);
          return true;
        }
      } catch (_) {}
      return false;
    },

    /* -----------------------------------------------------------------------
       Public shorthand methods
    ----------------------------------------------------------------------- */
    get(url) {
      return this.request(url, { method: 'GET' });
    },

    post(url, body) {
      return this.request(url, {
        method: 'POST',
        body: JSON.stringify(body)
      });
    },

    put(url, body) {
      return this.request(url, {
        method: 'PUT',
        body: JSON.stringify(body)
      });
    },

    patch(url, body) {
      return this.request(url, {
        method: 'PATCH',
        body: JSON.stringify(body)
      });
    },

    delete(url) {
      return this.request(url, { method: 'DELETE' });
    },

    /* -----------------------------------------------------------------------
       File upload (multipart/form-data — disease detection)
    ----------------------------------------------------------------------- */
    async upload(url, formData) {
      const token = localStorage.getItem('access_token');
      const headers = token ? { 'Authorization': 'Bearer ' + token } : {};
      // Do NOT set Content-Type — browser sets it with boundary
      const res = await fetch(BASE + url, {
        method: 'POST',
        headers,
        body: formData
      });
      try { return await res.json(); }
      catch (_) { return { success: res.ok }; }
    },

    /* -----------------------------------------------------------------------
       Auth helpers
    ----------------------------------------------------------------------- */
    isAuthenticated() {
      return !!localStorage.getItem('access_token');
    },

    requireAuth() {
      if (!this.isAuthenticated()) {
        window.location.href = '/login';
        return false;
      }
      return true;
    },

    logout() {
      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
      localStorage.removeItem('selected_farm_id');
      localStorage.removeItem('selected_coop_id');
      window.location.href = '/login';
    },

    /* -----------------------------------------------------------------------
       Context helpers (farm / coop selection)
    ----------------------------------------------------------------------- */
    getSelectedFarmId() {
      return localStorage.getItem('selected_farm_id');
    },

    setSelectedFarmId(id) {
      localStorage.setItem('selected_farm_id', String(id));
    },

    getSelectedCoopId() {
      return localStorage.getItem('selected_coop_id');
    },

    setSelectedCoopId(id) {
      localStorage.setItem('selected_coop_id', String(id));
    }
  };

  window.API = API;

})(window);
