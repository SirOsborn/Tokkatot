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
        nav_home:         'ទំព័រដើម',
        nav_monitoring:   'ត្រួតពិនិត្យ',
        nav_disease:      'ជំងឺវិភាគ',
        nav_schedules:    'កាលវិភាគ',
        nav_settings:     'ការកំណត់',
        alerts:           'ការជូនដំណឹង',

      /* --- Auth --- */
      login:            'ចូលគណនី',
      logout:           'ចាកចេញ',
      register:         'ចុះឈ្មោះ',
      email:            'អ៊ីមែល',
      password:         'លេខសម្ងាត់',
      confirm_password: 'បញ្ជាក់លេខសម្ងាត់',
      full_name:        'ឈ្មោះពេញ',
      phone:            'លេខទូរស័ព្ទ',
      forgot_password:  'ភ្លេចលេខសម្ងាត់?',
      no_account:       'មិនទាន់មានគណនី?',
      have_account:     'មានគណនីហើយ?',
      login_success:    'ចូលគណនីបានជោគជ័យ',
      login_failed:     'ចូលគណនីមិនបាន។ សូមពិនិត្យអ៊ីមែល ឬ លេខសម្ងាត់',
      register_key:     'សោសម្រាប់ចុះឈ្មោះ',
      register_success: 'ការចុះឈ្មោះបានជោគជ័យ',
      register_failed:  'ការចុះឈ្មោះបានបរាជ័យ',
      national_id:      'អត្តសញ្ញាណប័ណ្ណ',
      register_as:      'ចុះឈ្មោះជា',
      farmer:           'ម្ចាស់កសិដ្ឋាន',
      worker:           'កម្មករ',
      req_reg_key:      'ត្រូវការសោសម្រាប់ចុះឈ្មោះពី Tokkatot',
      use_farmer_id:    'ប្រើ ID ម្ចាស់កសិដ្ឋានដើម្បីចូលប្រើ',
      farmer_id:        'ID ម្ចាស់កសិដ្ឋាន',
      paste_farmer_id:  'បញ្ចូល ID ម្ចាស់កសិដ្ឋាន',

      /* --- Dashboard / Home --- */
      dashboard:        'ទំព័រដើម',
      welcome:          'សូមស្វាគមន៍',
      your_farm:        'កសិដ្ឋានរបស់អ្នក',
      my_farm:          'កសិដ្ឋានរបស់ខ្ញុំ',
      no_farm:          'មិនទាន់មានកសិដ្ឋាន',
      no_farm_msg:      'បង្កើតកសិដ្ឋានដំបូងរបស់អ្នកដើម្បីចាប់ផ្តើម',
      add_farm:         'បង្កើតកសិដ្ឋាន',
      select_coop:      'ជ្រើសរើសទ្រុង',
      coop:             'ទ្រុង',
      coops:            'ទ្រុង',
      add_coop:         'បន្ថែមទ្រុង',
      coop_name:        'ឈ្មោះទ្រុង',
      coop_number:      'លេខទ្រុង',
      capacity:         'ចំនួនអតិបរមា',
      current_count:    'ចំនួនមានស្រាប់',
      chicken_type:     'ប្រភេទមាន់',
      description:      'សេចក្ដីពិពណ៌នា',
      chickens:         'ចំនួនមាន់',
      capacity:         'ចំណុះ',
      last_updated:     'កែប្រែចុងក្រោយ',

      /* --- Environment Metrics --- */
      temperature:      'សីតុណ្ហភាព',
      humidity:         'សំណើម',
      water_level:      'កម្រឹតទឹក',
      environment:      'បរិស្ថាន',
      in_coop:          'ក្នុងទ្រុង',

      /* --- Devices --- */
      devices:          'ឧបករណ៍',
      device_name:      'ឈ្មោះឧបករណ៍',
      device_type:      'ប្រភេទនៃឧបករណ៍',
      device_status:    'ស្ថានភាពនៃឧបករណ៍',
      turn_on:          'បើក',
      turn_off:         'បិទ',
      control:          'គ្រប់គ្រង',
      online:           'អនឡាញ',
      offline:          'ក្រៅបណ្តាញ',

      /* --- Schedules --- */
      schedules:        'កាលវិភាគ',
      schedule_name:    'ឈ្មោះនៃកាលវិភាគ',
      add_schedule:     'បន្ថែមកាលវិភាគ',
      edit_schedule:    'កែប្រែ',
      delete_schedule:  'លុប',
      enabled:          'បើក',
      disabled:         'បិទ',
      cron_expr:        'ពេលវេលាដែលត្រូវដំណើការ',
      action:           'សកម្មភាព',
      duration:         'រយៈពេល',
      priority:         'អាទិភាព',

      /* --- Disease Detection --- */
      disease_detection:'ជំងឺវិភាគ',
      take_photo:       'ថតរូប',
      upload_photo:     'បញ្ចូលរូប',
      analyze:          'វិភាគ',
      result:           'លទ្ធផល',
      confidence:       'ភាពជឿជាក់របស់ប្រព័ន្ធ AI',
      recommendations:  'ការណែនាំ',
      healthy:          'មានសុខភាពល្អ',
      disease_found:    'ជំងឺត្រូវបានរកឃើញ',
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
      profile:          'គណនី',
      my_profile:       'គណនីខ្ញុំ',
      role:             'តួនាទី',
      edit_profile:     'កែប្រែព័ត៌មាន',
      save:             'រក្សាទុក',
      cancel:           'បោះបង់',
      owner:            'ម្ចាស់',
      manager:          'អ្នកគ្រប់គ្រង',
      viewer:           'អ្នកមើល',

      /* --- Settings --- */
      settings:         'ការកំណត់',
      language:         'ភាសា',
      lang_km:          'ភាសាខ្មែរ',
      lang_en:          'ភាសាអង់គ្លេស',
      notifications:    'ការជូនដំណឹង',
      account:          'គណនី',
      system:           'ប្រព័ន្ធ',
      about:            'អំពីតុក្កតត​ - Tokkatot',
      version:          'កំណែទី',
      logout_confirm:   'ចង់ចេញពីគណនី?',
      yes:              'បាទ/ចាស',
      no:               'ទេ',

      // Profile & Settings
      settings_desc:    'គ្រប់គ្រងគណនី និងការកំណត់របស់អ្នក',
      save_changes:     'រក្សាទុកការផ្លាស់ប្តូរ',

      // Admin Portal
      admin_toolkit:    'ឧបករណ៍គ្រប់គ្រង',
      manage_keys_farmers: 'គ្រប់គ្រងកូដចុះឈ្មោះ និងកសិករ',
      registration_keys: 'កូដចុះឈ្មោះ',
      farmer_management: 'គ្រប់គ្រងកសិករ',
      total_farmers:    'កសិករសរុប',
      active_workers:   'បុគ្គលិកសកម្ម',
      total_farms:      'កសិដ្ឋានសរុប',
      available_keys:   'កូដនៅទំនេរ',
      active_reg_keys:  'កូដចុះឈ្មោះសកម្ម',
      generate_new_key:  'បង្កើតកូដថ្មី',
      used:             'ប្រើរួច',
      available:        'នៅទំនេរ',
      expires:          'ផុតកំណត់',
      created:          'បង្កើតនៅ',
      last_login:       'ចូលចុងក្រោយ',
      joined:           'ចូលរួមនៅ',
      deactivate:       'បិទដំណើរការ',
      activate:         'បើកដំណើរការ',
      generate_key_title: 'បង្កើតកូដចុះឈ្មោះ',
      customer_phone:   'លេខទូរស័ព្ទអតិថិជន',
      national_id:      'អត្តសញ្ញាណប័ណ្ណ',
      key_generated_success: 'បង្កើតកូដជោគជ័យ!',
      key_note:         'សូមកត់ត្រាកូដនេះ ឬផ្តល់ឱ្យបុគ្គលិកនៅការដ្ឋាន',

      /* --- Common / Errors --- */
      loading:          'កំពុងដំណើរការ...',
      error:            'មានបញ្ហា',
      retry:            'ព្យាយាមម្ដងទៀត',
      'Invalid credentials': 'អ៊ីមែល ឬ លេខសម្ងាត់ មិនត្រឹមត្រូវ',
      'Account is inactive': 'គណនីរបស់អ្នកត្រូវបានផ្អាកជាបណ្ដោះអាសន្ន',
      'User not found':      'រកមិនឃើញអ្នកប្រើប្រាស់នេះទេ',
      fill_all_fields:      'សូមបំពេញព័ត៌មានឱ្យបានគ្រប់ជ្រុងជ្រោយ',
      success:          'ជោគជ័យ',
      confirm:          'បញ្ជាក់',
      delete:           'លុប',
      edit:             'កែប្រែ',
      add:              'បន្ថែម',
      close:            'បិទ',
      back:             'ត្រឡប់ក្រោយ',
      next:             'ទៅមុខ',
      done:             'រួចរាល់',
      of:               'ក្នុង',
      chickens:         'មាន់',
      capacity:         'ចំណុះ',
      no_data:          'មិនមានទិន្នន័យ',
      no_coops:         'មិនទាន់មានទ្រុង',
      no_devices:       'មិនទាន់មានឧបករណ៍',
      no_devices_msg:   'មិនទាន់មានសារពីឧបករណ៍',
      fill_all_fields:  'សូមបំពេញព័ត៌មានទាំងអស់',
      reg_key_req:      'សូមបញ្ចូលសោចុះឈ្មោះ',
      user_not_found:   'រកមិនឃើញអ្នកប្រើប្រាស់',
      invalid_password: 'លេខសម្ងាត់មិនត្រឹមត្រូវ',
      email_already_exists: 'អីមែលនេះមានរួចហើយ',
      phone_already_exists: 'លេខទូរស័ព្ទនេះមានរួចហើយ',
      invalid_reg_key:  'សោចុះឈ្មោះមិនត្រឹមត្រូវ',
      reg_key_used:     'សោចុះឈ្មោះត្រូវបានប្រើរួចហើយ',
      reg_key_expired:  'សោចុះឈ្មោះបានហួសសម័យ',
      invalid_farmer_id: 'ID ម្ចាស់កសិដ្ឋានមិនត្រឹមត្រូវ',
      account_inactive: 'គណនីត្រូវបានផ្អាក',
      farmer_id_req:    'សូមបញ្ចូល ID ម្ចាស់កសិដ្ឋាន',
      pwd_mismatch:     'លេខសម្ងាត់មិនត្រូវគ្នា',
      pwd_min_len:      'លេខសម្ងាត់យ៉ាងតិច ៦ តួ',
        no_farm_access:   'មិនមានសិទ្ធិចូលមើលកសិដ្ឋាន',
        contact_owner:    'សូមទំនាក់ទំនងម្ចាស់កសិដ្ឋានដើម្បីទទួលបានការអនុញ្ញាត',
        add_farm:         'បង្កើតកសិដ្ឋាន',
        workers:          'បុគ្គលិក',
        manage_workers:   'គ្រប់គ្រងបុគ្គលិក',
        workers_count:    'ចំនួនបុគ្គលិក',
        invite_worker:    'អញ្ជើញបុគ្គលិក',
        member_role:      'តួនាទី',
        viewer_role:      'មើលតែប៉ុណ្ណោះ',
        worker_role:      'កម្មករ',
        remove_member:    'លុបចេញ',
        no_workers:       'មិនទាន់មានបុគ្គលិក',
        contact_or_phone: 'អ៊ីមែល ឬ លេខទូរស័ព្ទ',
        owner_only_manage:'តែម្ចាស់កសិដ្ឋានប៉ុណ្ណោះអាចគ្រប់គ្រងបុគ្គលិក',
        farm_name:        'ឈ្មោះកសិដ្ឋាន',
      province_city:    'ខេត្ត/ក្រុង',
      select_province:  '-- ជ្រើសរើសខេត្ត --',
      create:           'បង្កើត',
      create_farm_msg:  'បង្កើតកសិដ្ឋានដំបូងរបស់អ្នកដើម្បីចាប់ផ្តើម',
      farm_name_req:    'សូមបញ្ចូលឈ្មោះកសិដ្ឋាន',
      sex:              'ភេទ',
      male:             'ប្រុស',
      female:           'ស្រី',
      other:            'ផ្សេងៗ',
      farm_details:     'ព័ត៌មានកសិដ្ឋាន',
      hourly_today:     'ថ្ងៃនេះ - រៀងរាល់ម៉ោង',
      no_readings:      'មិនទាន់មានទិន្នន័យនៅឡើយទេ',
      no_sensor_found:  'រកមិនឃើញសេនស័រ',
      no_sensor_msg:    'ទ្រុងនេះមិនទាន់មានសេនស័រសីតុណ្ហភាពនៅឡើយទេ',
      ai_disease_title: 'AI ជំងឺវិភាគមាន់',
      ai_training_msg:  'យើងកំពុងបង្កើតប្រព័ន្ធ AI ដើម្បីរកជំងឺ​មាន់​ពីរូប​ភាព​។ មុខ​ងារ​នេះ​នឹង​អាច​ប្រើ​បាន​ក្នុង​កំណែ​បន្ទាប់​។',
      back_to_home:     'ទៅ​ទំព័រ​ដើម',
      detect_disease:   'រក​ជំងឺ​មាន់',
      upload_or_take:   'ថតរូប ឬ​ជ្រើស​រូបភាព',
      drag_or_click:    'ចុចដើម្បីជ្រើសរើសរូបភាព',
      select_file:      'ជ្រើសរើសរូបភាព',
      analyzing_msg:    'AI កំពុង​វិភាគ​រូបភាព…',
      retry_detection:  'សាកល្បង​ម្ដង​ទៀត',
      action:           'សកម្មភាព',
      turn_on_action:   'បើក',
      turn_off_action:  'បិទ',
      time_duration:    'ម៉ោងដំណើរការ',
      repeat:           'ម្ដងហើយម្ដងទៀត',
      every_day:        'រៀងរាល់ថ្ងៃ',
      auto_off_after:   'បិទដោយស្វ័យប្រវត្តិបន្ទាប់ពី',
      minutes:          'នាទី',
      advanced_settings: 'ការកំណត់ស៊ីជម្រៅ',
      multi_step_seq:   'លំដាប់ច្រើនជំហាន',
      add_step:         'បន្ថែមជំហាន',
      no_schedules:     'មិនមានកាលវិភាគ',
      tap_plus_schedule: 'ចុចប៊ូតុង + ដើម្បីបង្កើតកាលវិភាគ',
      run_now:          'ដំណើរការឥឡូវ',
      name_device_req:  'សូមបំពេញឈ្មោះឧបករណ៍',
      no_alerts:        'មិនមានការជូនដំណឹង',
      acknowledged:     'បានទទួលស្គាល់',
      acknowledge_btn:  'ទទួលស្គាល់',
      page_not_found:   'រកទំព័រមិនឃើញ',
      page_not_found_msg: 'ទំព័រដែលអ្នកកំពុងរកមិនមាន​ ឬត្រូវបានផ្លាស់ប្ដូរ',
      go_home:          'ទៅទំព័រដើម',
      days_short:       ['អា','ច','អ','ព','ព្រ','ស','សៅ'],
      duration_none:    'គ្មានរយៈពេល',
      duration_5m:      'រយៈពេល 5 នាទី',
      duration_15m:     'រយៈពេល 15 នាទី',
      duration_30m:     'រយៈពេល 30 នាទី',
      duration_1h:      'រយៈពេល 1 ម៉ោង',
      duration_custom:  'រយៈពេលផ្ទាល់ខ្លួន',
    },

      en: {
        /* --- Navigation --- */
        nav_home:         'Home',
        nav_monitoring:   'Dashboard',
        nav_disease:      'AI Disease',
        nav_schedules:    'Schedules',
        nav_settings:     'Settings',
        alerts:           'Alerts',
      
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
      register_failed:  'Registration failed',
      national_id:      'National ID',
      register_as:      'Register as',
      farmer:           'Farmer',
      worker:           'Worker',
      req_reg_key:      'Requires a registration key from Tokkatot',
      use_farmer_id:    "Use your farm owner's ID to join their farm",
      farmer_id:        "Farmer's ID",
      paste_farmer_id:  'Paste your farmer user ID',

      /* --- Dashboard / Home --- */
      dashboard:        'Home',
      welcome:          'Welcome',
      your_farm:        'Your Farm',
      my_farm:          'My Farm',
      no_farm:          'No Farm Yet',
      no_farm_msg:      'Create your first farm to get started',
      add_farm:         'Add Farm',
      select_coop:      'Select Coop',
      coop:             'Coop',
      coops:            'Coops',
      add_coop:         'Add Coop',
      coop_name:        'Coop Name',
      coop_number:      'Coop Number',
      capacity:         'Capacity',
      current_count:    'Current Count',
      chicken_type:     'Chicken Type',
      description:      'Description',
      chickens:         'Chickens',
      capacity:         'Capacity',
      last_updated:     'Last Updated',

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
      lang_km:          'Khmer',
      lang_en:          'English',
      notifications:    'Notifications',
      account:          'Account',
      system:           'System',
      about:            'About',
      version:          'Version',
      logout_confirm:   'Are you sure you want to logout?',
      yes:              'Yes',
      no:               'No',

      /* --- Common / Errors --- */
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
      no_devices_msg:   'No devices yet',
      fill_all_fields:  'Please fill all required fields',
      reg_key_req:      'Registration key is required',
      user_not_found:   'User not found',
      invalid_password: 'Invalid password',
      email_already_exists: 'Email already exists',
      phone_already_exists: 'Phone already exists',
      invalid_reg_key:  'Invalid registration key',
      reg_key_used:     'Registration key already used',
      reg_key_expired:  'Registration key expired',
      invalid_farmer_id: 'Invalid farmer ID',
      account_inactive: 'Account is deactivated',
      farmer_id_req:    "Farmer's ID is required",
      pwd_mismatch:     'Passwords do not match',
      pwd_min_len:      'Password must be at least 6 characters',
        no_farm_access:   'No farm access',
        contact_owner:    'Contact your farm owner to get access.',
        add_farm:         'Add Farm',
        workers:          'Workers',
        manage_workers:   'Manage Workers',
        workers_count:    'Workers',
        invite_worker:    'Invite Worker',
        member_role:      'Role',
        viewer_role:      'Viewer',
        worker_role:      'Worker',
        remove_member:    'Remove',
        no_workers:       'No workers yet',
        contact_or_phone: 'Email or phone',
        owner_only_manage:'Only the farm owner can manage workers',
        farm_name:        'Farm Name',
      province_city:    'Province / City',
      select_province:  '-- Select Province --',
      create:           'Create',
      create_farm_msg:  'Create your first farm to get started.',
      farm_name_req:    'Farm name is required.',
      sex:              'Sex',
      male:             'Male',
      female:           'Female',
      other:            'Other',
      farm_details:     'Farm Details',
      hourly_today:     'Today — Hourly',
      no_readings:      'No readings yet today',
      no_sensor_found:  'No Sensor Found',
      no_sensor_msg:    'This coop does not have an active temperature sensor attached.',
      ai_disease_title: 'AI Disease Detection',
      ai_training_msg:  'We are training our AI to detect poultry diseases from photos. This feature will be available in the next version.',
      back_to_home:     'Back to Home',
      detect_disease:   'Detect Disease',
      upload_or_take:   'Take Photo or Upload',
      drag_or_click:    'Drag and drop or click to choose',
      select_file:      'Select File',
      analyzing_msg:    'AI is analyzing the image...',
      retry_detection:  'Try Again',
      action:           'Action',
      turn_on_action:   'Turn ON',
      turn_off_action:  'Turn OFF',
      time_duration:    'Time',
      repeat:           'Repeat',
      every_day:        'Every day',
      auto_off_after:   'Auto-off after',
      minutes:          'minutes',
      advanced_settings: 'Advanced Settings',
      multi_step_seq:   'Multi-step Sequence',
      add_step:         'Add Step',
      no_schedules:     'No schedules yet',
      tap_plus_schedule: 'Tap + to create a schedule',
      run_now:          'Run Now',
      name_device_req:  'Name and device are required.',
      no_alerts:        'No alerts at the moment.',
      acknowledged:     'Acknowledged',
      acknowledge_btn:  'OK',
      page_not_found:   'Page Not Found',
      page_not_found_msg: 'The page you are looking for does not exist <br> or has been moved.',
      go_home:          'Go to Home',
      days_short:       ['Sun','Mon','Tue','Wed','Thu','Fri','Sat'],
      duration_none:    'None',
      duration_5m:      '5m',
      duration_15m:     '15m',
      duration_30m:     '30m',
      duration_1h:      '1h',
      duration_custom:  'Custom',
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
      const val = dict[key] || translations.en[key];
      if (!val) {
        console.warn(`[i18n] Missing translation for key: "${key}" in lang: ${this._lang}`);
        return key;
      }
      return val;
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
    },

    /**
     * Returns all 25 Cambodian provinces as an array of { km, en } objects.
     * Use i18n.provinces() to populate dropdowns — do NOT add these to the
     * flat translation map because they are data, not UI strings.
     *
     * Usage:
     *   i18n.provinces()             // full list
     *   i18n.provinces()[0].km       // 'ភ្នំពេញ'
     *   i18n.provinces()[0].en       // 'Phnom Penh'
     *   // For a <select> in the current language:
     *   i18n.provinces().map(p => p[i18n.getLang()])
     */
    provinces() {
      return [
        { km: 'ភ្នំពេញ',        en: 'Phnom Penh' },
        { km: 'កណ្ដាល',          en: 'Kandal' },
        { km: 'សៀមរាប',          en: 'Siem Reap' },
        { km: 'បាត់ដំបង',        en: 'Battambang' },
        { km: 'កំពង់ចាម',        en: 'Kampong Cham' },
        { km: 'កំពង់ឆ្នាំង',      en: 'Kampong Chhnang' },
        { km: 'កំពង់ស្ពឺ',        en: 'Kampong Speu' },
        { km: 'កំពង់ធំ',          en: 'Kampong Thom' },
        { km: 'កំពត',             en: 'Kampot' },
        { km: 'ក្រចេះ',           en: 'Kratie' },
        { km: 'មណ្ឌលគិរី',        en: 'Mondulkiri' },
        { km: 'ព្រះវិហារ',        en: 'Preah Vihear' },
        { km: 'ព្រៃវែង',          en: 'Prey Veng' },
        { km: 'ពោធិ៍សាត់',        en: 'Pursat' },
        { km: 'រតនគិរី',          en: 'Ratanakiri' },
        { km: 'ស្វាយរៀង',         en: 'Svay Rieng' },
        { km: 'តាកែវ',            en: 'Takeo' },
        { km: 'ឧត្តរមានជ័យ',      en: 'Oddar Meanchey' },
        { km: 'កែប',              en: 'Kep' },
        { km: 'កោះកុង',           en: 'Koh Kong' },
        { km: 'ប៉ៃលិន',           en: 'Pailin' },
        { km: 'សំបូរ',            en: 'Sambour' },
        { km: 'ត្បូងឃ្មុំ',        en: 'Tboung Khmum' },
        { km: 'ព្រះសីហនុ',        en: 'Sihanoukville' },
        { km: 'ស្ទឹងត្រែង',        en: 'Stung Treng' },
      ];
    }
  };

  /* Set lang attribute on init */
  document.documentElement.lang = i18n.getLang();

  console.log('[i18n] Initialized with', Object.keys(translations.km).length, 'keys');
  window.i18n = i18n;
  /* Shorthand */
  window.t = function (key) { return i18n.t(key); };

})(window);
