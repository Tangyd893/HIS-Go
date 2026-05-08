-- ============================================================
-- 迁移脚本 011：随访管理服务 (his_followup)
-- ============================================================
\c his_followup;

CREATE TABLE IF NOT EXISTS followup_plans (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    patient_id VARCHAR(64) NOT NULL,
    template_id VARCHAR(64),
    plan_name VARCHAR(100),
    start_date VARCHAR(10),
    end_date VARCHAR(10),
    frequency INT DEFAULT 0,
    status SMALLINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_followup_plans_patient ON followup_plans (patient_id);

CREATE TABLE IF NOT EXISTS followup_tasks (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    plan_id VARCHAR(64) NOT NULL,
    assignee_id VARCHAR(64),
    execute_date VARCHAR(10),
    type SMALLINT DEFAULT 0,
    content TEXT,
    status SMALLINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_followup_tasks_plan ON followup_tasks (plan_id);

CREATE TABLE IF NOT EXISTS satisfaction_surveys (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    followup_task_id VARCHAR(64),
    patient_id VARCHAR(64) NOT NULL,
    score INT DEFAULT 0,
    feedback TEXT,
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
