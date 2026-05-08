-- ============================================================
-- 迁移脚本 004：处方管理服务 (his_prescription)
-- ============================================================
\c his_prescription;

CREATE TABLE IF NOT EXISTS prescriptions (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    patient_id VARCHAR(64) NOT NULL,
    patient_name VARCHAR(100),
    doctor_id VARCHAR(64),
    diagnosis_id VARCHAR(64),
    prescription_type SMALLINT DEFAULT 0,
    status SMALLINT DEFAULT 0,
    note TEXT,
    version INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_prescriptions_patient ON prescriptions (patient_id);
CREATE INDEX idx_prescriptions_status ON prescriptions (status);

CREATE TABLE IF NOT EXISTS prescription_details (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    prescription_id VARCHAR(64) NOT NULL,
    drug_id VARCHAR(64),
    drug_name VARCHAR(100),
    specification VARCHAR(100),
    dosage NUMERIC(10,2),
    usage VARCHAR(100),
    frequency VARCHAR(50),
    days INT DEFAULT 0,
    quantity INT DEFAULT 0,
    unit_price NUMERIC(10,2) DEFAULT 0,
    note TEXT
);

CREATE INDEX idx_prescription_details_pid ON prescription_details (prescription_id);

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
