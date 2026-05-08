-- ============================================================
-- 迁移脚本 007：检查检验服务 (his_examination)
-- ============================================================
\c his_examination;

CREATE TABLE IF NOT EXISTS examination_reports (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    patient_id VARCHAR(64) NOT NULL,
    patient_name VARCHAR(100),
    exam_request_id VARCHAR(64),
    exam_type VARCHAR(50),
    exam_item VARCHAR(100),
    body_part VARCHAR(50),
    findings TEXT,
    impression TEXT,
    conclusion TEXT,
    technician_id VARCHAR(64),
    reviewer_id VARCHAR(64),
    status SMALLINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_exam_reports_patient ON examination_reports (patient_id);
CREATE INDEX idx_exam_reports_request ON examination_reports (exam_request_id);

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
