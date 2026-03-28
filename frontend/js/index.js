console.log('[Index] Initializing Vue App...');
const { createApp } = Vue;

const app = createApp({
  data() {
    return {
      loading:        true,
      showAddFarm:    false,
      addFarmError:   '',
      addFarmLoading: false,
      showAddCoop:    false,
      addCoopError:   '',
      addCoopLoading: false,
      lang:           window.i18n ? window.i18n.getLang() : 'km',
      farms:          [],
      selectedFarmId: null,
      coops:          [],
      selectedCoopId: null,
      temp:           null,
      humidity:       null,
      coopChickens:   null,
      coopCapacity:   null,
      devices:        [],
      lastUpdated:    '--',
      wsConnected:    false,
      ws:             null,
      isAdmin:        false,
      userRole:       'farmer',
      workersCount:   0,
      newFarm:        { name: '', province: '' },
      newCoop:        { number: 1, name: '', capacity: null, current_count: null, chicken_type: '', description: '' },
      provinces:      window.i18n ? window.i18n.provinces() : [],
    };
  },

  computed: {
    currentFarm() {
      return this.farms.find(f => f.id === this.selectedFarmId) || null;
    }
  },

  async mounted() {
    if (!window.API.requireAuth()) return;
    this.userRole = localStorage.getItem('user_role') || 'farmer';
    this.isAdmin  = this.userRole === 'admin';
    if (this.isAdmin) {
      window.location.href = '/admin';
      return;
    }
    await loadComponents();
    this.lang = window.i18n ? window.i18n.getLang() : 'km';
    await this.loadFarms();
  },

  methods: {
    t(k) { return window.i18n ? window.i18n.t(k) : k; },

    formatFarmName(name) {
      if (!name) return '...';
      if (name.toLowerCase() === 'my farm') return this.t('my_farm');
      return name;
    },

    // ── Extract model slug from device_id ────────────────────────────────
    // device_id format: "SIM-ESP32-1234567890:conveyor_belt"
    // We split on ':' and take the last part → "conveyor_belt"
    getDeviceModel(device) {
      const deviceId = device.device_id || '';
      const parts    = deviceId.split(':');
      if (parts.length > 1) {
        return parts[parts.length - 1].toLowerCase().trim();
      }
      // Fallback: try model field, then name
      return (device.model || device.name || '').toLowerCase().trim();
    },

    // ── Translate device name using model slug ────────────────────────────
    getDeviceName(device) {
      const modelMap = {
        'conveyor_belt': 'device_conveyor_belt',
        'fan':           'device_cooling_fan',
        'cooling_fan':   'device_cooling_fan',
        'feeder_motor':  'device_feeder_motor',
        'feeder':        'device_feeder_motor',
        'heater':        'device_heater',
        'temp_humidity': 'device_temp_humidity',
        'water_level':   'device_water_level',
        'water_pump':    'device_water_level',
      };
      const model = this.getDeviceModel(device);
      const key   = modelMap[model];
      if (key) {
        const translated = this.t(key);
        if (translated && translated !== key) return translated;
      }
      // Fallback to raw name from DB
      return device.name || model || '—';
    },

    // ── Device icon using model slug ──────────────────────────────────────
    deviceIcon(device) {
      const iconMap = {
        'conveyor_belt': 'conveyor_belt',
        'fan':           'mode_fan',
        'cooling_fan':   'mode_fan',
        'feeder_motor':  'restaurant',
        'feeder':        'restaurant',
        'heater':        'mode_heat',
        'temp_humidity': 'thermostat',
        'water_level':   'water',
        'water_pump':    'water_pump',
        'sensor':        'sensors',
        'relay':         'electrical_services',
      };
      const model = this.getDeviceModel(device);
      return iconMap[model] || iconMap[device.type] || 'electrical_services';
    },

    // ── Farms ────────────────────────────────────────────────────────────
    async loadFarms() {
      this.loading = true;
      try {
        const data = await window.API.get('/v1/farms');
        this.farms = (data && data.data && data.data.data)          ? data.data.data
                   : (data && data.data && Array.isArray(data.data)) ? data.data : [];
        if (this.farms.length > 0) {
          const saved = window.API.getSelectedFarmId();
          const match = saved && this.farms.find(f => String(f.id) === saved);
          this.selectedFarmId = match ? match.id : this.farms[0].id;
          window.API.setSelectedFarmId(this.selectedFarmId);
          if (this.farms[0] && this.farms[0].name) localStorage.setItem('farm_name', this.farms[0].name);
          await this.loadCoops();
        }
      } catch(e) { console.error(e); }
      finally    { this.loading = false; }
    },

    async onFarmChange() {
      window.API.setSelectedFarmId(this.selectedFarmId);
      this.selectedCoopId = null;
      this.coops          = [];
      this.devices        = [];
      this.temp           = null;
      this.humidity       = null;
      this.coopChickens   = null;
      this.coopCapacity   = null;
      this.lastUpdated    = '--';
      if (this.ws) { try { this.ws.close(); } catch(_) {} this.ws = null; this.wsConnected = false; }
      await this.loadCoops();
    },

    // ── Coops ────────────────────────────────────────────────────────────
    async loadCoops() {
      const fid = this.selectedFarmId;
      if (!fid) return;
      try {
        const data = await window.API.get('/v1/farms/' + fid + '/coops');
        this.coops = (data && data.data && data.data.coops)          ? data.data.coops
                   : (data && data.data && Array.isArray(data.data)) ? data.data : [];
        if (this.coops.length > 0) {
          const saved = window.API.getSelectedCoopId();
          const match = saved && this.coops.find(c => String(c.id) === saved);
          await this.selectCoop(match || this.coops[0]);
        } else {
          this.selectedCoopId = null;
          this.devices        = [];
          this.temp           = null;
          this.humidity       = null;
          this.coopChickens   = null;
          this.coopCapacity   = null;
          this.lastUpdated    = '--';
          if (this.ws) { try { this.ws.close(); } catch(_) {} this.ws = null; this.wsConnected = false; }
        }
        await this.loadDevices();
        await this.loadWorkersCount();
      } catch(e) { console.error(e); }
    },

    async selectCoop(coop) {
      this.selectedCoopId = coop.id;
      window.API.setSelectedCoopId(coop.id);
      await this.loadCoopData(coop.id);
      this.connectWS();
    },

    async loadCoopData(coopId) {
      const fid = this.selectedFarmId;
      if (!fid || !coopId) return;
      try {
        const data = await window.API.get('/v1/farms/' + fid + '/coops/' + coopId);
        const c = data && data.data ? data.data : {};
        this.temp         = c.temperature   !== undefined ? c.temperature   : null;
        this.humidity     = c.humidity      !== undefined ? c.humidity      : null;
        this.coopChickens = c.current_count !== undefined ? c.current_count : null;
        this.coopCapacity = c.capacity      !== undefined ? c.capacity      : null;
        this.lastUpdated  = c.last_updated  ? new Date(c.last_updated).toLocaleTimeString() : '--';
      } catch(e) { console.error(e); }
    },

    // ── Devices ──────────────────────────────────────────────────────────
    async loadDevices() {
      const fid = this.selectedFarmId;
      if (!fid || !this.selectedCoopId) return;
      try {
        const data = await window.API.get('/v1/farms/' + fid + '/devices?coop_id=' + this.selectedCoopId);
        const raw  = (data && data.data && data.data.devices) ? data.data.devices : [];
        this.devices = raw.map(d => ({ ...d, loading: false }));
      } catch(e) { console.error(e); }
    },

    async toggleDevice(device) {
      device.loading = true;
      const cmd = device.last_state === 'on' ? 'turn_off' : 'turn_on';
      try {
        const fid = this.selectedFarmId;
        await window.API.post('/v1/farms/' + fid + '/devices/' + device.id + '/commands', { command_type: cmd });
        device.last_state = device.last_state === 'on' ? 'off' : 'on';
      } catch(e) { console.error(e); }
      finally    { device.loading = false; }
    },

    // ── Workers ──────────────────────────────────────────────────────────
    async loadWorkersCount() {
      const fid = this.selectedFarmId;
      if (!fid) { this.workersCount = 0; return; }
      try {
        const data    = await window.API.get('/v1/farms/' + fid + '/members');
        const members = (data && data.data && data.data.members) ? data.data.members : [];
        this.workersCount = members.filter(m => m.role && m.role !== 'farmer').length;
      } catch(e) { console.error(e); this.workersCount = 0; }
    },

    // ── Farm / Coop creation ─────────────────────────────────────────────
    async createFarm() {
      this.addFarmError = '';
      if (!this.newFarm.name.trim()) { this.addFarmError = this.t('farm_name_req'); return; }
      this.addFarmLoading = true;
      try {
        const res = await window.API.post('/v1/farms', { name: this.newFarm.name, province: this.newFarm.province || undefined });
        if (res && res.success) {
          this.showAddFarm = false;
          this.newFarm = { name: '', province: '' };
          await this.loadFarms();
        } else {
          this.addFarmError = (res && res.message) || this.t('error');
        }
      } catch(e) { this.addFarmError = this.t('error') + ': ' + (e.message || e); }
      finally    { this.addFarmLoading = false; }
    },

    async createCoop() {
      this.addCoopError = '';
      const fid = this.selectedFarmId;
      if (!fid) { this.addCoopError = this.t('error'); return; }
      if (!this.newCoop.name || !this.newCoop.number) {
        this.addCoopError = this.t('coop_name') + ' / ' + this.t('coop_number');
        return;
      }
      this.addCoopLoading = true;
      try {
        const payload = {
          number:        this.newCoop.number,
          name:          this.newCoop.name,
          capacity:      this.newCoop.capacity      || undefined,
          current_count: this.newCoop.current_count || undefined,
          chicken_type:  this.newCoop.chicken_type  || undefined,
          description:   this.newCoop.description   || undefined,
        };
        const res = await window.API.post('/v1/farms/' + fid + '/coops', payload);
        if (res && res.success) {
          this.showAddCoop = false;
          this.newCoop = { number: 1, name: '', capacity: null, current_count: null, chicken_type: '', description: '' };
          await this.loadCoops();
        } else {
          this.addCoopError = (res && res.message) || this.t('error');
        }
      } catch(e) { this.addCoopError = this.t('error') + ': ' + (e.message || e); }
      finally    { this.addCoopLoading = false; }
    },

    // ── WebSocket ────────────────────────────────────────────────────────
    connectWS() {
      if (this.ws) { try { this.ws.close(); } catch(_) {} }
      const token  = localStorage.getItem('access_token');
      const fid    = this.selectedFarmId;
      const coopId = this.selectedCoopId;
      if (!token || !fid || !coopId) return;
      try {
        const proto = location.protocol === 'https:' ? 'wss' : 'ws';
        this.ws = new WebSocket(proto + '://' + location.host + '/v1/ws?token=' + token + '&farm_id=' + fid + '&coop_id=' + coopId);
        this.ws.onopen    = () => { this.wsConnected = true; };
        this.ws.onclose   = () => { this.wsConnected = false; };
        this.ws.onerror   = () => { this.wsConnected = false; };
        this.ws.onmessage = (ev) => {
          try {
            const msg = JSON.parse(ev.data);
            if (msg.temperature !== undefined) this.temp     = msg.temperature;
            if (msg.humidity    !== undefined) this.humidity = msg.humidity;
            this.lastUpdated = new Date().toLocaleTimeString();
          } catch(_) {}
        };
      } catch(e) { console.warn('[WS]', e); }
    },
  },

  beforeUnmount() {
    if (this.ws) try { this.ws.close(); } catch(_) {}
  }

}).mount('#app');