-- ============================================================
-- 迁移脚本 003：门诊诊疗服务 (his_clinic)
-- ============================================================
\c his_clinic;

CREATE TABLE IF NOT EXISTS clinic_records (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    registration_id VARCHAR(64) NOT NULL,
    patient_id VARCHAR(64) NOT NULL,
    patient_name VARCHAR(100),
    doctor_id VARCHAR(64),
    chief_complaint TEXT,
    present_illness TEXT,
    diagnosis TEXT,
    icd_code VARCHAR(20),
    advice TEXT,
    status SMALLINT DEFAULT 0,
    visit_time TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_clinic_patient ON clinic_records (patient_id);
CREATE INDEX idx_clinic_registration ON clinic_records (registration_id);

CREATE TABLE IF NOT EXISTS examination_requests (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    clinic_record_id VARCHAR(64) NOT NULL,
    patient_id VARCHAR(64) NOT NULL,
    exam_type VARCHAR(50),
    exam_item VARCHAR(100),
    body_part VARCHAR(50),
    clinical_diagnosis TEXT,
    note TEXT,
    status SMALLINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_exam_req_clinic ON examination_requests (clinic_record_id);
CREATE INDEX idx_exam_req_patient ON examination_requests (patient_id);

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
