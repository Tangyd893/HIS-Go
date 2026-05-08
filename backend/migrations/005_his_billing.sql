-- ============================================================
-- 迁移脚本 005：收费结算服务 (his_billing)
-- ============================================================
\c his_billing;

CREATE TABLE IF NOT EXISTS bills (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    patient_id VARCHAR(64) NOT NULL,
    registration_id VARCHAR(64),
    bill_no VARCHAR(64) NOT NULL UNIQUE,
    total_amount NUMERIC(12,2) DEFAULT 0,
    paid_amount NUMERIC(12,2) DEFAULT 0,
    pay_method SMALLINT DEFAULT 0,
    status SMALLINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_bills_patient ON bills (patient_id);
CREATE INDEX idx_bills_bill_no ON bills (bill_no);

CREATE TABLE IF NOT EXISTS bill_details (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    bill_id VARCHAR(64) NOT NULL,
    item_type SMALLINT DEFAULT 0,
    item_name VARCHAR(100),
    unit_price NUMERIC(10,2) DEFAULT 0,
    quantity INT DEFAULT 0,
    amount NUMERIC(10,2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_bill_details_bill ON bill_details (bill_id);

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
