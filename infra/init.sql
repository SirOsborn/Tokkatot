-- ================================================================
-- TOKKATOT DATABASE INITIALIZATION SCRIPT
-- Auto-generated from middleware/database/schema.go
-- ================================================================

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email TEXT UNIQUE,
    phone TEXT UNIQUE,
    phone_country_code VARCHAR(5),
    password_hash TEXT NOT NULL,
    name TEXT NOT NULL,
    is_active BOOLEAN DEFAULT true,
    last_login TIMESTAMP,
    national_id_number VARCHAR(50),
    sex VARCHAR(10),
    province VARCHAR(100),
    full_name TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Farms table
CREATE TABLE IF NOT EXISTS farms (
    id UUID PRIMARY KEY,
    owner_id UUID NOT NULL REFERENCES users(id),
    name TEXT NOT NULL,
    location TEXT,
    province VARCHAR(100),
    timezone VARCHAR(40) DEFAULT 'Asia/Phnom_Penh',
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),
    description TEXT,
    image_url TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Admins table
CREATE TABLE IF NOT EXISTS admins (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE,
    phone VARCHAR(20) UNIQUE,
    password_hash TEXT NOT NULL,
    role VARCHAR(20) DEFAULT 'admin',
    language VARCHAR(10) DEFAULT 'km',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP
);

-- Coops table
CREATE TABLE IF NOT EXISTS coops (
    id UUID PRIMARY KEY,
    farm_id UUID NOT NULL REFERENCES farms(id),
    number INTEGER NOT NULL,
    name TEXT NOT NULL,
    capacity INTEGER,
    current_count INTEGER,
    chicken_type VARCHAR(20),
    main_device_id UUID,
    temp_min DECIMAL(5,2),
    temp_max DECIMAL(5,2),
    water_level_half_threshold DECIMAL(10,4),
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(farm_id, number)
);

-- Farm users (membership table)
CREATE TABLE IF NOT EXISTS farm_users (
    id UUID PRIMARY KEY,
    farm_id UUID NOT NULL REFERENCES farms(id),
    user_id UUID NOT NULL REFERENCES users(id),
    role VARCHAR(20) NOT NULL CHECK (role IN ('farmer', 'worker', 'viewer')),
    invited_by UUID NOT NULL REFERENCES users(id),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(farm_id, user_id)
);

-- Devices table
CREATE TABLE IF NOT EXISTS devices (
    id UUID PRIMARY KEY,
    farm_id UUID NOT NULL REFERENCES farms(id),
    coop_id UUID REFERENCES coops(id),
    device_id VARCHAR(50) NOT NULL UNIQUE,
    name TEXT NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('gpio', 'relay', 'pwm', 'adc', 'servo', 'sensor')),
    model TEXT,
    is_main_controller BOOLEAN DEFAULT false,
    firmware_version VARCHAR(20) NOT NULL,
    hardware_id TEXT NOT NULL,
    response TEXT,
    location TEXT,
    is_active BOOLEAN DEFAULT true,
    is_online BOOLEAN DEFAULT false,
    last_heartbeat TIMESTAMP,
    last_command_status VARCHAR(50),
    last_command_at TIMESTAMP,
    response TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Device commands
CREATE TABLE IF NOT EXISTS device_commands (
    id UUID PRIMARY KEY,
    device_id UUID NOT NULL REFERENCES devices(id),
    farm_id UUID NOT NULL REFERENCES farms(id),
    coop_id UUID REFERENCES coops(id),
    issued_by UUID NOT NULL REFERENCES users(id),
    command_type VARCHAR(50) NOT NULL,
    command_value TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'success', 'failed', 'timeout')),
    response TEXT,
    issued_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    executed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Schedules
CREATE TABLE IF NOT EXISTS schedules (
    id UUID PRIMARY KEY,
    farm_id UUID NOT NULL REFERENCES farms(id),
    coop_id UUID REFERENCES coops(id),
    device_id UUID NOT NULL REFERENCES devices(id),
    name TEXT NOT NULL,
    schedule_type VARCHAR(20) NOT NULL CHECK (schedule_type IN ('time_based', 'duration_based', 'condition_based')),
    cron_expression TEXT,
    on_duration INTEGER,
    off_duration INTEGER,
    condition_json JSONB,
    action VARCHAR(20) NOT NULL CHECK (action IN ('on', 'off', 'set_value')),
    action_value TEXT,
    action_duration INTEGER,
    action_sequence JSONB,
    priority INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    next_execution TIMESTAMP,
    last_execution TIMESTAMP,
    execution_count INTEGER DEFAULT 0,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Schedule executions
CREATE TABLE IF NOT EXISTS schedule_executions (
    id UUID PRIMARY KEY,
    schedule_id UUID NOT NULL REFERENCES schedules(id),
    device_id UUID NOT NULL REFERENCES devices(id),
    scheduled_time TIMESTAMP NOT NULL,
    actual_execution_time TIMESTAMP,
    status VARCHAR(20) NOT NULL CHECK (status IN ('executed', 'failed', 'skipped')),
    execution_duration_ms INTEGER,
    device_response JSONB,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Event logs
CREATE TABLE IF NOT EXISTS event_logs (
    id UUID PRIMARY KEY,
    farm_id UUID NOT NULL REFERENCES farms(id),
    user_id UUID NOT NULL REFERENCES users(id),
    event_type VARCHAR(50) NOT NULL,
    resource_id UUID,
    old_value JSONB,
    new_value JSONB,
    ip_address VARCHAR(45),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Registration keys
CREATE TABLE IF NOT EXISTS registration_keys (
    id UUID PRIMARY KEY,
    key_code VARCHAR(50) UNIQUE NOT NULL,
    farm_id UUID REFERENCES farms(id) ON DELETE SET NULL,
    farm_name VARCHAR(255),
    customer_phone VARCHAR(20),
    national_id_number VARCHAR(50),
    full_name TEXT,
    sex VARCHAR(10),
    province VARCHAR(100),
    is_used BOOLEAN DEFAULT false,
    used_by_user_id UUID REFERENCES users(id),
    used_at TIMESTAMP,
    expires_at TIMESTAMP,
    created_by VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- User sessions
CREATE TABLE IF NOT EXISTS user_sessions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    device_name VARCHAR(255),
    ip_address VARCHAR(45),
    user_agent TEXT,
    refresh_token TEXT NOT NULL UNIQUE,
    last_activity TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Alerts
CREATE TABLE IF NOT EXISTS alerts (
    id UUID PRIMARY KEY,
    farm_id UUID NOT NULL REFERENCES farms(id),
    coop_id UUID REFERENCES coops(id),
    device_id UUID REFERENCES devices(id),
    alert_type VARCHAR(50) NOT NULL,
    severity VARCHAR(20) NOT NULL CHECK (severity IN ('info', 'warning', 'critical')),
    message TEXT NOT NULL,
    threshold_value DECIMAL(10,4),
    actual_value DECIMAL(10,4),
    is_active BOOLEAN DEFAULT true,
    is_acknowledged BOOLEAN DEFAULT false,
    triggered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    acknowledged_by UUID REFERENCES users(id),
    acknowledged_at TIMESTAMP,
    resolved_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Alert subscriptions
CREATE TABLE IF NOT EXISTS alert_subscriptions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    alert_type VARCHAR(50) NOT NULL,
    channel VARCHAR(20) NOT NULL DEFAULT 'push',
    is_enabled BOOLEAN DEFAULT true,
    quiet_hours_start VARCHAR(5),
    quiet_hours_end VARCHAR(5),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, alert_type, channel)
);

-- Web Push Subscriptions
CREATE TABLE IF NOT EXISTS web_push_subscriptions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    endpoint TEXT NOT NULL UNIQUE,
    p256dh TEXT NOT NULL,
    auth TEXT NOT NULL,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_used TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Device configurations
CREATE TABLE IF NOT EXISTS device_configurations (
    id UUID PRIMARY KEY,
    device_id UUID NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    parameter_name VARCHAR(100) NOT NULL,
    parameter_value TEXT NOT NULL,
    unit VARCHAR(20),
    min_value DECIMAL(10,4),
    max_value DECIMAL(10,4),
    is_calibrated BOOLEAN DEFAULT false,
    calibrated_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(device_id, parameter_name)
);

-- Device readings
CREATE TABLE IF NOT EXISTS device_readings (
    id UUID PRIMARY KEY,
    device_id UUID NOT NULL REFERENCES devices(id),
    sensor_type VARCHAR(50) NOT NULL,
    value DECIMAL(10,4) NOT NULL,
    unit VARCHAR(20) NOT NULL DEFAULT '',
    quality VARCHAR(20) NOT NULL DEFAULT 'good',
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Gateway persistent tokens (API Keys)
CREATE TABLE IF NOT EXISTS gateway_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    farm_id UUID NOT NULL REFERENCES farms(id) ON DELETE CASCADE,
    device_id UUID REFERENCES devices(id) ON DELETE SET NULL,
    user_id UUID NOT NULL REFERENCES users(id),
    token_hash TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    is_active BOOLEAN DEFAULT true,
    last_used_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Gateway temporary provisioning codes
CREATE TABLE IF NOT EXISTS gateway_provisions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    setup_code VARCHAR(12) NOT NULL UNIQUE,
    hardware_id TEXT NOT NULL,
    farm_id UUID REFERENCES farms(id) ON DELETE CASCADE,
    coop_id UUID REFERENCES coops(id) ON DELETE CASCADE,
    is_claimed BOOLEAN DEFAULT false,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ===== PERFORMANCE INDEXES =====

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);
CREATE INDEX IF NOT EXISTS idx_users_is_active_created ON users(is_active, created_at);
CREATE INDEX IF NOT EXISTS idx_farms_owner_id ON farms(owner_id);
CREATE INDEX IF NOT EXISTS idx_farms_is_active ON farms(is_active);
CREATE INDEX IF NOT EXISTS idx_farms_created_at ON farms(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_coops_farm_id ON coops(farm_id);
CREATE INDEX IF NOT EXISTS idx_coops_farm_number ON coops(farm_id, number);
CREATE INDEX IF NOT EXISTS idx_coops_is_active ON coops(is_active);
CREATE INDEX IF NOT EXISTS idx_coops_main_device ON coops(main_device_id);
CREATE INDEX IF NOT EXISTS idx_farm_users_farm_id ON farm_users(farm_id);
CREATE INDEX IF NOT EXISTS idx_farm_users_user_id ON farm_users(user_id);
CREATE INDEX IF NOT EXISTS idx_farm_users_active ON farm_users(is_active);
CREATE INDEX IF NOT EXISTS idx_devices_farm_id_active ON devices(farm_id, is_active);
CREATE INDEX IF NOT EXISTS idx_devices_coop_id ON devices(coop_id);
CREATE INDEX IF NOT EXISTS idx_devices_device_id ON devices(device_id);
CREATE INDEX IF NOT EXISTS idx_devices_hardware_id ON devices(hardware_id);
CREATE INDEX IF NOT EXISTS idx_devices_is_online ON devices(is_online);
CREATE INDEX IF NOT EXISTS idx_devices_main_controller ON devices(coop_id, is_main_controller);
CREATE INDEX IF NOT EXISTS idx_device_commands_device_id ON device_commands(device_id);
CREATE INDEX IF NOT EXISTS idx_device_commands_farm_id ON device_commands(farm_id);
CREATE INDEX IF NOT EXISTS idx_device_commands_coop_id ON device_commands(coop_id);
CREATE INDEX IF NOT EXISTS idx_device_commands_status ON device_commands(status);
CREATE INDEX IF NOT EXISTS idx_device_commands_created ON device_commands(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_schedules_farm_device ON schedules(farm_id, device_id);
CREATE INDEX IF NOT EXISTS idx_schedules_coop_id ON schedules(coop_id);
CREATE INDEX IF NOT EXISTS idx_schedules_is_active ON schedules(is_active);
CREATE INDEX IF NOT EXISTS idx_schedules_next_execution ON schedules(next_execution);
CREATE INDEX IF NOT EXISTS idx_schedule_executions_schedule_id ON schedule_executions(schedule_id);
CREATE INDEX IF NOT EXISTS idx_schedule_executions_time ON schedule_executions(scheduled_time DESC);
CREATE INDEX IF NOT EXISTS idx_schedule_executions_status ON schedule_executions(status);
CREATE INDEX IF NOT EXISTS idx_event_logs_farm_user ON event_logs(farm_id, user_id);
CREATE INDEX IF NOT EXISTS idx_event_logs_created ON event_logs(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_registration_keys_code ON registration_keys(key_code);
CREATE INDEX IF NOT EXISTS idx_registration_keys_is_used ON registration_keys(is_used);
CREATE INDEX IF NOT EXISTS idx_registration_keys_expires ON registration_keys(expires_at);
CREATE INDEX IF NOT EXISTS idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_sessions_expires ON user_sessions(expires_at);
CREATE INDEX IF NOT EXISTS idx_user_sessions_refresh ON user_sessions(refresh_token);
CREATE INDEX IF NOT EXISTS idx_alerts_farm_id ON alerts(farm_id);
CREATE INDEX IF NOT EXISTS idx_alerts_coop_id ON alerts(coop_id);
CREATE INDEX IF NOT EXISTS idx_alerts_device_id ON alerts(device_id);
CREATE INDEX IF NOT EXISTS idx_alerts_is_active ON alerts(is_active);
CREATE INDEX IF NOT EXISTS idx_alerts_severity ON alerts(severity);
CREATE INDEX IF NOT EXISTS idx_alerts_triggered_at ON alerts(triggered_at DESC);
CREATE INDEX IF NOT EXISTS idx_alert_subscriptions_user_id ON alert_subscriptions(user_id);
CREATE INDEX IF NOT EXISTS idx_web_push_user_id ON web_push_subscriptions(user_id);
CREATE INDEX IF NOT EXISTS idx_device_configs_device_id ON device_configurations(device_id);
CREATE INDEX IF NOT EXISTS idx_device_readings_device_id ON device_readings(device_id);
CREATE INDEX IF NOT EXISTS idx_device_readings_timestamp ON device_readings(timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_device_readings_sensor_type ON device_readings(device_id, sensor_type);
CREATE INDEX IF NOT EXISTS idx_gateway_tokens_farm_id ON gateway_tokens(farm_id);
CREATE INDEX IF NOT EXISTS idx_gateway_tokens_token_hash ON gateway_tokens(token_hash);
CREATE INDEX IF NOT EXISTS idx_gateway_provisions_setup_code ON gateway_provisions(setup_code);
CREATE INDEX IF NOT EXISTS idx_gateway_provisions_expires_at ON gateway_provisions(expires_at);
