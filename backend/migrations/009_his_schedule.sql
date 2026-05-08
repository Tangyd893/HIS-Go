-- ============================================================
-- 迁移脚本 009：排班管理服务 (his_schedule)
-- ============================================================
\c his_schedule;

CREATE TABLE IF NOT EXISTS schedules (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    doctor_id VARCHAR(64) NOT NULL,
    doctor_name VARCHAR(50),
    dept_id VARCHAR(64) NOT NULL,
    dept_name VARCHAR(100),
    work_date VARCHAR(10) NOT NULL,
    time_slot INT NOT NULL,
    max_patients INT DEFAULT 0,
    current_patients INT DEFAULT 0,
    room_no VARCHAR(20),
    status SMALLINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE UNIQUE INDEX idx_schedules_unique ON schedules (doctor_id, work_date, time_slot);
CREATE INDEX idx_schedules_dept ON schedules (dept_id, work_date);

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
