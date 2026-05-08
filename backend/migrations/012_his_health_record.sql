-- ============================================================
-- 迁移脚本 012：健康档案服务 (his_health_record)
-- ============================================================
\c his_health_record;

CREATE TABLE IF NOT EXISTS health_record_summaries (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    patient_id VARCHAR(64) NOT NULL UNIQUE,
    patient_name VARCHAR(100),
    total_visits INT DEFAULT 0,
    total_prescriptions INT DEFAULT 0,
    total_examinations INT DEFAULT 0,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS timeline_events (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    patient_id VARCHAR(64) NOT NULL,
    date VARCHAR(10),
    event_type VARCHAR(50),
    description TEXT,
    related_id VARCHAR(64),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_timeline_patient ON timeline_events (patient_id);

CREATE TABLE IF NOT EXISTS record_authorizations (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    patient_id VARCHAR(64) NOT NULL,
    doctor_id VARCHAR(64),
    auth_time VARCHAR(20),
    expire_time VARCHAR(20),
    status SMALLINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

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
