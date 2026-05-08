-- ============================================================
-- 迁移脚本 002：挂号预约服务 (his_registration)
-- ============================================================
\c his_registration;

CREATE TABLE IF NOT EXISTS schedules (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    dept_id VARCHAR(64) NOT NULL,
    dept_name VARCHAR(100),
    doctor_id VARCHAR(64) NOT NULL,
    doctor_name VARCHAR(50),
    date VARCHAR(10) NOT NULL,
    time_slot INT NOT NULL,
    total_count INT NOT NULL DEFAULT 0,
    remain_count INT NOT NULL DEFAULT 0,
    fee NUMERIC(10,2) DEFAULT 0,
    status SMALLINT DEFAULT 1,
    version INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_schedules_unique ON schedules (dept_id, doctor_id, date, time_slot);
CREATE INDEX idx_schedules_dept_date ON schedules (dept_id, date);

CREATE TABLE IF NOT EXISTS registrations (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    patient_id VARCHAR(64) NOT NULL,
    patient_name VARCHAR(100),
    schedule_id VARCHAR(64) NOT NULL,
    registration_date VARCHAR(10),
    queue_number INT DEFAULT 0,
    status SMALLINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_registrations_patient ON registrations (patient_id);
CREATE INDEX idx_registrations_schedule ON registrations (schedule_id);

CREATE TABLE IF NOT EXISTS msg_record (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    msg_id VARCHAR(128) NOT NULL UNIQUE,
    business_key VARCHAR(128),
    exchange VARCHAR(100),
    routing_key VARCHAR(100),
    body TEXT NOT NULL,
    status SMALLINT DEFAULT 0,
    retry_count INT DEFAULT 0,
    max_retry INT DEFAULT 3,
    next_retry_time TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
