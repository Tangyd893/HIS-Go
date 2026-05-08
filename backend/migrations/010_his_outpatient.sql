-- ============================================================
-- 迁移脚本 010：院外患者服务 (his_outpatient)
-- ============================================================
\c his_outpatient;

CREATE TABLE IF NOT EXISTS consultations (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    patient_id VARCHAR(64) NOT NULL,
    doctor_id VARCHAR(64),
    type SMALLINT DEFAULT 0,
    description TEXT,
    images TEXT,
    status SMALLINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_consultations_patient ON consultations (patient_id);

CREATE TABLE IF NOT EXISTS consultation_messages (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    consultation_id VARCHAR(64) NOT NULL,
    sender_id VARCHAR(64),
    sender_name VARCHAR(50),
    content TEXT,
    msg_type VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_cons_msg_consultation ON consultation_messages (consultation_id);

CREATE TABLE IF NOT EXISTS chronic_contracts (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    patient_id VARCHAR(64) NOT NULL,
    doctor_id VARCHAR(64),
    disease_type VARCHAR(100),
    contract_date VARCHAR(10),
    end_date VARCHAR(10),
    status SMALLINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS health_data (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    patient_id VARCHAR(64) NOT NULL,
    data_type VARCHAR(50),
    value VARCHAR(50),
    unit VARCHAR(20),
    measure_time VARCHAR(20),
    abnormal BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_health_data_patient ON health_data (patient_id);

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
