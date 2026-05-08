-- ============================================================
-- 迁移脚本 006：药房管理服务 (his_pharmacy)
-- ============================================================
\c his_pharmacy;

CREATE TABLE IF NOT EXISTS drugs (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    name VARCHAR(100) NOT NULL,
    generic_name VARCHAR(100),
    specification VARCHAR(100),
    manufacturer VARCHAR(200),
    batch_no VARCHAR(64),
    purchase_price NUMERIC(10,2) DEFAULT 0,
    retail_price NUMERIC(10,2) DEFAULT 0,
    stock INT DEFAULT 0,
    min_stock INT DEFAULT 0,
    expiry_date VARCHAR(10),
    status SMALLINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_drugs_name ON drugs (name);
CREATE INDEX idx_drugs_expiry ON drugs (expiry_date);

CREATE TABLE IF NOT EXISTS dispense_records (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    prescription_id VARCHAR(64) NOT NULL,
    patient_id VARCHAR(64) NOT NULL,
    drug_id VARCHAR(64) NOT NULL,
    quantity INT DEFAULT 0,
    dispenser_id VARCHAR(64),
    checker_id VARCHAR(64),
    status SMALLINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_dispense_prescription ON dispense_records (prescription_id);

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
