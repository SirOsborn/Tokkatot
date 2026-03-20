console.log('[Workers] Initializing Vue App...');
const { createApp } = Vue;

createApp({
  data() {
    return {
      loading: true,
      farms: [],
      selectedFarmId: null,
      members: [],
      showInvite: false,
      inviteContact: '',
      inviteRole: 'viewer',
      inviteError: '',
      inviteLoading: false,
      userRole: localStorage.getItem('user_role') || 'farmer'
    };
  },
  computed: {
    canManage() {
      return this.userRole === 'farmer';
    },
    workers() {
      return this.members.filter(m => m.role && m.role !== 'farmer');
    },
    currentFarmName() {
      const f = this.farms.find(x => x.id === this.selectedFarmId);
      return f ? f.name : '';
    }
  },
  async mounted() {
    if (!window.API.requireAuth()) return;
    await loadComponents();
    await this.loadFarms();
    if (window.i18n && window.i18n.applyAll) window.i18n.applyAll();
  },
  methods: {
    t(k) { return window.i18n ? window.i18n.t(k) : k; },
    formatDate(dateStr) {
      if (!dateStr) return '--';
      const d = new Date(dateStr);
      return d.toLocaleDateString('en-GB', { day: '2-digit', month: 'short', year: 'numeric' });
    },
    async loadFarms() {
      this.loading = true;
      try {
        const data = await window.API.get('/v1/farms');
        this.farms = (data && data.data && data.data.data) ? data.data.data
                   : (data && data.data && Array.isArray(data.data)) ? data.data : [];
        if (this.farms.length > 0) {
          const saved = window.API.getSelectedFarmId();
          const match = saved && this.farms.find(f => String(f.id) === saved);
          this.selectedFarmId = match ? match.id : this.farms[0].id;
          window.API.setSelectedFarmId(this.selectedFarmId);
          await this.loadMembers();
        }
      } catch (e) { console.error(e); }
      finally { this.loading = false; }
    },
    async onFarmChange() {
      window.API.setSelectedFarmId(this.selectedFarmId);
      await this.loadMembers();
    },
    async loadMembers() {
      const fid = this.selectedFarmId;
      if (!fid) return;
      try {
        const data = await window.API.get('/v1/farms/' + fid + '/members');
        this.members = (data && data.data && data.data.members) ? data.data.members : [];
      } catch (e) { console.error(e); this.members = []; }
    },
    async inviteMember() {
      this.inviteError = '';
      const fid = this.selectedFarmId;
      if (!fid) return;
      const contact = (this.inviteContact || '').trim();
      if (!contact) { this.inviteError = this.t('contact_or_phone'); return; }

      const payload = { role: this.inviteRole };
      if (contact.includes('@')) {
        payload.email = contact;
      } else {
        payload.phone = contact;
      }

      this.inviteLoading = true;
      try {
        const res = await window.API.post('/v1/farms/' + fid + '/members', payload);
        if (res && res.success) {
          this.showInvite = false;
          this.inviteContact = '';
          await this.loadMembers();
        } else {
          this.inviteError = (res && res.message) || this.t('error');
        }
      } catch (e) {
        this.inviteError = this.t('error');
      } finally {
        this.inviteLoading = false;
      }
    },
    async updateRole(member) {
      const fid = this.selectedFarmId;
      if (!fid || !member || !member.user_id) return;
      try {
        await window.API.put('/v1/farms/' + fid + '/members/' + member.user_id, { role: member.role });
      } catch (e) {
        console.error(e);
      }
    },
    async removeMember(member) {
      const fid = this.selectedFarmId;
      if (!fid || !member || !member.user_id) return;
      try {
        await window.API.delete('/v1/farms/' + fid + '/members/' + member.user_id);
        await this.loadMembers();
      } catch (e) {
        console.error(e);
      }
    }
  }
}).mount('#app');
