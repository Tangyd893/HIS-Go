-- ============================================================
-- 迁移脚本 008：住院管理服务 (his_inpatient)
-- ============================================================
\c his_inpatient;

CREATE TABLE IF NOT EXISTS inpatient_records (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    patient_id VARCHAR(64) NOT NULL,
    patient_name VARCHAR(100),
    admission_date TIMESTAMP,
    discharge_date TIMESTAMP,
    dept_id VARCHAR(64),
    room_no VARCHAR(20),
    bed_no VARCHAR(20),
    diagnosis TEXT,
    deposit NUMERIC(12,2) DEFAULT 0,
    total_cost NUMERIC(12,2) DEFAULT 0,
    status SMALLINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_inpatient_patient ON inpatient_records (patient_id);
CREATE INDEX idx_inpatient_status ON inpatient_records (status);

CREATE TABLE IF NOT EXISTS medical_orders (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    inpatient_id VARCHAR(64) NOT NULL,
    doctor_id VARCHAR(64),
    order_type SMALLINT DEFAULT 0,
    content TEXT,
    start_time TIMESTAMP,
    end_time TIMESTAMP,
    status SMALLINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_orders_inpatient ON medical_orders (inpatient_id);

CREATE TABLE IF NOT EXISTS nursing_records (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    inpatient_id VARCHAR(64) NOT NULL,
    nurse_id VARCHAR(64),
    record_time TIMESTAMP,
    content TEXT,
    vital_signs TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_nursing_inpatient ON nursing_records (inpatient_id);

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
