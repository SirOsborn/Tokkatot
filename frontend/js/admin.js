const { createApp } = Vue;

const app = createApp({
    data() {
        return {
            currentTab: 'keys',
            keys: [],
            farmers: [],
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
                active_keys: 0
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
                if (this.currentTab === 'keys') {
                    await this.fetchKeys();
                } else {
                    await this.fetchFarmers();
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
                this.keys = response.data;
            }
        },
        async fetchFarmers() {
            const response = await this.apiCall('/v1/admin/farmers');
            if (response.success) {
                this.farmers = response.data;
            }
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
            if (!newStatus) {
                const response = await this.apiCall(`/v1/admin/farmers/${user.id}`, 'DELETE');
                if (response.success) {
                    user.is_active = false;
                    this.showToast('Farmer deactivated');
                } else {
                    this.showToast(response.message || 'Action failed', 'error');
                }
            } else {
                this.showToast('Activation not implemented yet', 'error');
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
console.log('[Admin] App mounted successfully');
