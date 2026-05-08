-- ============================================================
-- 迁移脚本 015：电子病历服务 (his_emr)
-- ============================================================
\c his_emr;

CREATE TABLE IF NOT EXISTS medical_records (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    patient_id VARCHAR(64) NOT NULL,
    clinic_record_id VARCHAR(64),
    template_id VARCHAR(64),
    chief_complaint TEXT,
    present_illness TEXT,
    past_history TEXT,
    physical_exam TEXT,
    auxiliary_exam TEXT,
    diagnosis TEXT,
    treatment_plan TEXT,
    quality_level INT DEFAULT 0,
    status SMALLINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_medical_records_patient ON medical_records (patient_id);
CREATE INDEX idx_medical_records_clinic ON medical_records (clinic_record_id);

CREATE TABLE IF NOT EXISTS record_templates (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    name VARCHAR(100) NOT NULL,
    dept_id VARCHAR(64),
    content TEXT,
    type INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS quality_controls (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    record_id VARCHAR(64) NOT NULL,
    reviewer_id VARCHAR(64),
    level INT DEFAULT 0,
    comment TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_qc_record ON quality_controls (record_id);

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
