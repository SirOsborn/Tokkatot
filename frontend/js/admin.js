const { createApp } = Vue;

const app = createApp({
    data() {
        return {
            currentTab: 'keys',
            keys: [],
            farmers: [],
            gateways: [],
            showKeyModal: false,
            loading: true,
            generatedKey: null,
            newKey: {
                farm_name: '',
                customer_phone: '',
                national_id: '',
                full_name: '',
                sex: '',
                province: ''
            },
            stats: {
                total_farmers: 0,
                total_workers: 0,
                total_farms: 0,
                active_keys: 0,
                active_gateways: 0
            },
            unassigned: [],
            showAssignModal: false,
            selectedHardware: {},
            assignment: {
                farmer_id: '',
                name: ''
            },
            provinces: window.i18n ? window.i18n.provinces() : [],
            toast: { show: false, message: '', type: 'success' }
        };
    },
    async mounted() {
        this.checkAuth();
        if (window.i18n && window.i18n.applyAll) {
            window.i18n.applyAll();
        }
        await this.fetchData();
    },
    methods: {
        t(k) { return window.i18n ? window.i18n.t(k) : k; },
        showToast(message, type = 'success') {
            this.toast.message = message;
            this.toast.type = type;
            this.toast.show = true;
            setTimeout(() => { this.toast.show = false; }, 3000);
        },
        checkAuth() {
            const token = localStorage.getItem('access_token');
            const role = localStorage.getItem('user_role');
            
            if (!token || role !== 'admin') {
                console.warn('Access denied: Admin only');
                window.location.href = '/login';
            }
        },
        async fetchData() {
            this.loading = true;
            try {
                await this.fetchStats();
                await this.fetchUnassigned();
                if (this.currentTab === 'keys') {
                    await this.fetchKeys();
                } else if (this.currentTab === 'farmers') {
                    await this.fetchFarmers();
                } else if (this.currentTab === 'gateways') {
                    await this.fetchGateways();
                }
            } finally {
                this.loading = false;
            }
        },
        async fetchStats() {
            const response = await this.apiCall('/v1/admin/stats');
            if (response.success) {
                this.stats = response.data;
            }
        },
        async fetchKeys() {
            const response = await this.apiCall('/v1/admin/reg-keys');
            if (response.success) {
                this.keys = Array.isArray(response.data) ? response.data : [];
            } else {
                this.keys = [];
            }
        },
        async fetchFarmers() {
            const response = await this.apiCall('/v1/admin/farmers');
            if (response.success) {
                this.farmers = Array.isArray(response.data) ? response.data : [];
            } else {
                this.farmers = [];
            }
        },
        async fetchGateways() {
            const response = await this.apiCall('/v1/admin/gateways');
            if (response.success) {
                this.gateways = Array.isArray(response.data) ? response.data : [];
            } else {
                this.gateways = [];
            }
            await this.fetchUnassigned();
        },
        async fetchUnassigned() {
            const response = await this.apiCall('/v1/admin/unassigned-gateways');
            if (response.success && response.data) {
                this.unassigned = Array.isArray(response.data.unassigned) ? response.data.unassigned : [];
            } else {
                this.unassigned = [];
            }
        },
        openAssignModal(ug) {
            this.selectedHardware = ug;
            this.assignment = {
                farmer_id: '',
                name: 'Main Farm Gateway'
            };
            this.showAssignModal = true;
            // Ensure farmers are loaded for the dropdown
            this.fetchFarmers();
        },
        async assignGateway() {
            const farmer = this.farmers.find(f => f.id === this.assignment.farmer_id);
            if (!farmer || !farmer.farm_id) {
                this.showToast('Selected farmer has no farm associated', 'error');
                return;
            }

            this.loading = true;
            try {
                const response = await this.apiCall('/v1/admin/assign-gateway', 'POST', {
                    hardware_id: this.selectedHardware.hardware_id,
                    farm_id: farmer.farm_id,
                    name: this.assignment.name
                });
                if (response.success) {
                    this.showToast('Gateway assigned successfully!');
                    this.showAssignModal = false;
                    await this.fetchGateways();
                } else {
                    this.showToast(response.message || 'Assignment failed', 'error');
                }
            } finally {
                this.loading = false;
            }
        },
        onFarmerSelect() {
            const farmer = this.farmers.find(f => f.id === this.assignment.farmer_id);
            if (farmer && farmer.farm_name) {
                this.assignment.name = `${farmer.farm_name} Gateway`;
            }
        },
        async revokeGateway(gw) {
            if (!confirm(`Are you sure you want to revoke access for ${gw.name || 'this gateway'}?`)) return;
            
            this.loading = true;
            try {
                const response = await this.apiCall(`/v1/admin/gateways/${gw.id}`, 'DELETE');
                if (response.success) {
                    this.showToast('Gateway access revoked successfully');
                    await this.fetchGateways();
                } else {
                    this.showToast(response.message || 'Failed to revoke access', 'error');
                }
            } finally {
                this.loading = false;
            }
        },
        isOnline(lastUsedAt) {
            if (!lastUsedAt) return false;
            const lastSeen = new Date(lastUsedAt);
            const now = new Date();
            // Online if seen in the last 2 minutes (heartbeat is every 60s)
            return (now - lastSeen) < (2 * 60 * 1000);
        },
        async generateKey() {
            this.loading = true;
            try {
                const response = await this.apiCall('/v1/admin/reg-keys', 'POST', this.newKey);
                if (response.success) {
                    this.generatedKey = response.data;
                    await this.fetchKeys();
                    this.showToast(this.t('success'));
                    // Reset form
                    this.newKey = { farm_name: '', customer_phone: '', national_id: '', full_name: '', sex: '', province: '' };
                } else {
                    this.showToast(response.message || this.t('error'), 'error');
                }
            } finally {
                this.loading = false;
            }
        },
        async toggleUserStatus(user) {
            const newStatus = !user.is_active;
            const response = await this.apiCall(`/v1/admin/farmers/${user.id}`, 'DELETE', { active: newStatus });
            if (response.success) {
                user.is_active = newStatus;
                this.showToast(`Farmer ${newStatus ? 'activated' : 'deactivated'} successfully`);
                await this.fetchStats();
            } else {
                this.showToast(response.message || 'Action failed', 'error');
            }
        },
        async apiCall(endpoint, method = 'GET', body = null) {
            const token = localStorage.getItem('access_token');
            const options = {
                method,
                headers: {
                    'Authorization': `Bearer ${token}`,
                    'Content-Type': 'application/json'
                }
            };
            if (body) options.body = JSON.stringify(body);

            try {
                const response = await fetch(endpoint, options);
                const data = await response.json();
                
                if (response.status === 401 || response.status === 403) {
                    window.location.href = '/login';
                    return { success: false };
                }
                
                return { success: response.ok, data: data.data, message: data.message };
            } catch (error) {
                console.error('API Error:', error);
                this.showToast('Connection error', 'error');
                return { success: false };
            }
        },
        formatDate(dateStr) {
            if (!dateStr) return 'Never';
            const date = new Date(dateStr);
            return date.toLocaleDateString('en-GB', {
                day: '2-digit',
                month: 'short',
                year: 'numeric',
                hour: '2-digit',
                minute: '2-digit'
            });
        },
        logout() {
            localStorage.clear();
            window.location.href = '/login';
        }
    },
    watch: {
        currentTab() {
            this.fetchData();
        }
    }
});
app.mount('#app');
const appEl = document.getElementById('app');
if (appEl) {
    appEl.removeAttribute('v-cloak');
}
console.log('[Admin] App mounted successfully');
