-- ============================================================
-- 迁移脚本 014：系统管理服务 (his_system)
-- ============================================================
\c his_system;

CREATE TABLE IF NOT EXISTS dict_types (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    dict_name VARCHAR(100) NOT NULL,
    dict_type VARCHAR(100) NOT NULL UNIQUE,
    status SMALLINT DEFAULT 1,
    remark VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS dict_items (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    dict_type VARCHAR(100) NOT NULL,
    label VARCHAR(100) NOT NULL,
    value VARCHAR(100) NOT NULL,
    sort_order INT DEFAULT 0,
    status SMALLINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_dict_items_type ON dict_items (dict_type);

CREATE TABLE IF NOT EXISTS system_params (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    param_name VARCHAR(100),
    param_key VARCHAR(100) NOT NULL UNIQUE,
    param_value TEXT,
    remark VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS operation_logs (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    user_id VARCHAR(64),
    username VARCHAR(50),
    module VARCHAR(50),
    action VARCHAR(50),
    method VARCHAR(10),
    url VARCHAR(200),
    ip VARCHAR(50),
    params TEXT,
    result TEXT,
    status SMALLINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_operation_logs_user ON operation_logs (user_id);
CREATE INDEX idx_operation_logs_created ON operation_logs (created_at);

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
