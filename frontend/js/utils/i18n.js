/**
 * Tokkatot i18n — Khmer / English translations
 * Global: window.i18n
 *
 * Usage:
 *   i18n.t('home')           // returns Khmer or English string
 *   i18n.setLang('en')       // switch to English
 *   i18n.getLang()           // 'km' or 'en'
 *
 * Default language: Khmer ('km') — stored in localStorage('language')
 */

(function (window) {
  'use strict';

  const translations = {
    km: {
      /* --- Navigation --- */
      nav_home:         'ផ្ទះ',
      nav_monitoring:   'ការត្រួតពិនិត្យ',
      nav_disease:      'ជំងឺ',
      nav_schedules:    'កាលវិភាគ',
      nav_settings:     'ការកំណត់',

      /* --- Auth --- */
      login:            'ចូលគណនី',
      logout:           'ចាកចេញ',
      register:         'បង្កើតគណនី',
      email:            'អ៊ីមែល',
      password:         'លេខសម្ងាត់',
      confirm_password: 'បញ្ជាក់លេខសម្ងាត់',
      full_name:        'ឈ្មោះពេញ',
      phone:            'លេខទូរស័ព្ទ',
      forgot_password:  'ភ្លេចលេខសម្ងាត់?',
      no_account:       'មិនទាន់មានគណនី?',
      have_account:     'មានគណនីហើយ?',
      login_success:    'ចូលបានជោគជ័យ',
      login_failed:     'ចូលមិនបាន។ ពិនិត្យអ៊ីមែល ឬ លេខសម្ងាត់',
      register_key:     'លេខកូដចុះឈ្មោះ',
      register_success: 'ចុះឈ្មោះបានជោគជ័យ',

      /* --- Dashboard / Home --- */
      dashboard:        'ផ្ទះ',
      welcome:          'សូមស្វាគមន៍',
      your_farm:        'ចំការរបស់អ្នក',
      no_farm:          'មិនទាន់មានចំការ',
      select_coop:      'ជ្រើសរើសទ្រុង',
      coop:             'ទ្រុង',
      coops:            'ទ្រុង',

      /* --- Environment Metrics --- */
      temperature:      'សីតុណ្ហភាព',
      humidity:         'សំណើម',
      water_level:      'ទឹក',
      environment:      'បរិស្ថាន',
      in_coop:          'ក្នុងទ្រុង',

      /* --- Devices --- */
      devices:          'ឧបករណ៍',
      device_name:      'ឈ្មោះឧបករណ៍',
      device_type:      'ប្រភេទ',
      device_status:    'ស្ថានភាព',
      turn_on:          'បើក',
      turn_off:         'បិទ',
      control:          'គ្រប់គ្រង',
      online:           'អនឡាញ',
      offline:          'ក្រៅបណ្តាញ',

      /* --- Schedules --- */
      schedules:        'កាលវិភាគ',
      schedule_name:    'ឈ្មោះ',
      add_schedule:     'បន្ថែមកាលវិភាគ',
      edit_schedule:    'កែប្រែ',
      delete_schedule:  'លុប',
      enabled:          'បើក',
      disabled:         'បិទ',
      cron_expr:        'ពេលវេលា (Cron)',
      action:           'សកម្មភាព',
      duration:         'រយៈពេល',
      priority:         'អាទិភាព',

      /* --- Disease Detection --- */
      disease_detection:'ការរកឃើញជំងឺ',
      take_photo:       'ថតរូប',
      upload_photo:     'ផ្ទុករូប',
      analyze:          'វិភាគ',
      result:           'លទ្ធផល',
      confidence:       'ភាពជឿជាក់',
      recommendations:  'អនុសាសន៍',
      healthy:          'មានសុខភាព',
      disease_found:    'រកឃើញជំងឺ',
      no_disease:       'គ្មានជំងឺ',
      coming_soon:      'នឹងមានឆាប់ៗ',

      /* --- Monitoring --- */
      monitoring:       'ការត្រួតពិនិត្យ',
      temperature_history: 'ប្រវត្តិសីតុណ្ហភាព',
      today:            'ថ្ងៃនេះ',
      high:             'ខ្ពស់',
      low:              'ទាប',
      now:              'ឥឡូវ',

      /* --- Profile --- */
      profile:          'ប្រវត្តិ',
      my_profile:       'គណនីខ្ញុំ',
      role:             'តួនាទី',
      edit_profile:     'កែប្រែ',
      save:             'រក្សាទុក',
      cancel:           'បោះបង់',
      owner:            'ម្ចាស់',
      manager:          'អ្នកគ្រប់គ្រង',
      viewer:           'អ្នកមើល',

      /* --- Settings --- */
      settings:         'ការកំណត់',
      language:         'ភាសា',
      notifications:    'ការជូនដំណឹង',
      account:          'គណនី',
      system:           'ប្រព័ន្ធ',
      about:            'អំពី',
      version:          'កំណែ',
      logout_confirm:   'តើអ្នកចង់ចాកចេញ?',
      yes:              'បាទ/ចាស',
      no:               'ទេ',

      /* --- Common --- */
      loading:          'កំពុងផ្ទុក...',
      error:            'កំហុស',
      retry:            'ព្យាយាមម្ដងទៀត',
      success:          'ជោគជ័យ',
      confirm:          'បញ្ជាក់',
      delete:           'លុប',
      edit:             'កែប្រែ',
      add:              'បន្ថែម',
      close:            'បិទ',
      back:             'ត្រឡប់',
      next:             'បន្ទាប់',
      done:             'រួចរាល់',
      of:               'ក្នុង',
      chickens:         'មាន់',
      capacity:         'ចំណុះ',
      no_data:          'មិនមានទិន្នន័យ',
      no_coops:         'មិនទាន់មានទ្រុង',
      no_devices:       'មិនទាន់មានឧបករណ៍',
    },

    en: {
      /* --- Navigation --- */
      nav_home:         'Home',
      nav_monitoring:   'Monitoring',
      nav_disease:      'Disease',
      nav_schedules:    'Schedules',
      nav_settings:     'Settings',

      /* --- Auth --- */
      login:            'Login',
      logout:           'Logout',
      register:         'Register',
      email:            'Email',
      password:         'Password',
      confirm_password: 'Confirm Password',
      full_name:        'Full Name',
      phone:            'Phone Number',
      forgot_password:  'Forgot Password?',
      no_account:       "Don't have an account?",
      have_account:     'Already have an account?',
      login_success:    'Logged in successfully',
      login_failed:     'Login failed. Check email or password.',
      register_key:     'Registration Key',
      register_success: 'Registration successful',

      /* --- Dashboard / Home --- */
      dashboard:        'Home',
      welcome:          'Welcome',
      your_farm:        'Your Farm',
      no_farm:          'No farm yet',
      select_coop:      'Select Coop',
      coop:             'Coop',
      coops:            'Coops',

      /* --- Environment Metrics --- */
      temperature:      'Temperature',
      humidity:         'Humidity',
      water_level:      'Water',
      environment:      'Environment',
      in_coop:          'In Coop',

      /* --- Devices --- */
      devices:          'Devices',
      device_name:      'Device Name',
      device_type:      'Type',
      device_status:    'Status',
      turn_on:          'Turn On',
      turn_off:         'Turn Off',
      control:          'Control',
      online:           'Online',
      offline:          'Offline',

      /* --- Schedules --- */
      schedules:        'Schedules',
      schedule_name:    'Name',
      add_schedule:     'Add Schedule',
      edit_schedule:    'Edit',
      delete_schedule:  'Delete',
      enabled:          'Enabled',
      disabled:         'Disabled',
      cron_expr:        'Time (Cron)',
      action:           'Action',
      duration:         'Duration',
      priority:         'Priority',

      /* --- Disease Detection --- */
      disease_detection: 'Disease Detection',
      take_photo:       'Take Photo',
      upload_photo:     'Upload Photo',
      analyze:          'Analyze',
      result:           'Result',
      confidence:       'Confidence',
      recommendations:  'Recommendations',
      healthy:          'Healthy',
      disease_found:    'Disease Found',
      no_disease:       'No Disease',
      coming_soon:      'Coming Soon',

      /* --- Monitoring --- */
      monitoring:       'Monitoring',
      temperature_history: 'Temperature History',
      today:            'Today',
      high:             'High',
      low:              'Low',
      now:              'Now',

      /* --- Profile --- */
      profile:          'Profile',
      my_profile:       'My Profile',
      role:             'Role',
      edit_profile:     'Edit Profile',
      save:             'Save',
      cancel:           'Cancel',
      owner:            'Owner',
      manager:          'Manager',
      viewer:           'Viewer',

      /* --- Settings --- */
      settings:         'Settings',
      language:         'Language',
      notifications:    'Notifications',
      account:          'Account',
      system:           'System',
      about:            'About',
      version:          'Version',
      logout_confirm:   'Are you sure you want to logout?',
      yes:              'Yes',
      no:               'No',

      /* --- Common --- */
      loading:          'Loading...',
      error:            'Error',
      retry:            'Try Again',
      success:          'Success',
      confirm:          'Confirm',
      delete:           'Delete',
      edit:             'Edit',
      add:              'Add',
      close:            'Close',
      back:             'Back',
      next:             'Next',
      done:             'Done',
      of:               'of',
      chickens:         'Chickens',
      capacity:         'Capacity',
      no_data:          'No data',
      no_coops:         'No coops yet',
      no_devices:       'No devices yet',
    }
  };

  const i18n = {
    _lang: localStorage.getItem('language') || 'km',

    getLang() {
      return this._lang;
    },

    setLang(lang) {
      if (lang !== 'km' && lang !== 'en') return;
      this._lang = lang;
      localStorage.setItem('language', lang);
      document.documentElement.lang = lang;
    },

    toggleLang() {
      this.setLang(this._lang === 'km' ? 'en' : 'km');
    },

    t(key) {
      const dict = translations[this._lang] || translations.km;
      return dict[key] || translations.en[key] || key;
    },

    /* Apply all data-i18n attributes on the page */
    applyAll() {
      document.querySelectorAll('[data-i18n]').forEach(el => {
        const key = el.getAttribute('data-i18n');
        el.textContent = this.t(key);
      });
      document.querySelectorAll('[data-i18n-placeholder]').forEach(el => {
        const key = el.getAttribute('data-i18n-placeholder');
        el.placeholder = this.t(key);
      });
    }
  };

  /* Set lang attribute on init */
  document.documentElement.lang = i18n.getLang();

  window.i18n = i18n;
  /* Shorthand */
  window.t = function (key) { return i18n.t(key); };

})(window);
