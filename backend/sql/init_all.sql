-- ============================================================
-- HIS-Go 数据库初始化脚本
-- 说明：为每个微服务创建独立的 PostgreSQL 数据库
-- 对应的详细表结构将在后续版本补充
-- ============================================================

-- 创建各微服务数据库
CREATE DATABASE his_auth;
CREATE DATABASE his_user;
CREATE DATABASE his_registration;
CREATE DATABASE his_clinic;
CREATE DATABASE his_emr;
CREATE DATABASE his_prescription;
CREATE DATABASE his_billing;
CREATE DATABASE his_pharmacy;
CREATE DATABASE his_examination;
CREATE DATABASE his_inpatient;
CREATE DATABASE his_schedule;
CREATE DATABASE his_outpatient;
CREATE DATABASE his_followup;
CREATE DATABASE his_health_record;
CREATE DATABASE his_notification;
CREATE DATABASE his_statistics;
CREATE DATABASE his_system;

-- ============================================================
-- his_auth 数据库表结构
-- ============================================================
\c his_auth;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    real_name VARCHAR(50),
    phone VARCHAR(20),
    email VARCHAR(100),
    avatar VARCHAR(255),
    role VARCHAR(20) DEFAULT 'doctor',
    dept_id VARCHAR(64),
    status SMALLINT DEFAULT 1,
    last_login_time TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 角色表
CREATE TABLE IF NOT EXISTS roles (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    role_name VARCHAR(50) NOT NULL,
    role_code VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    status SMALLINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 权限表
CREATE TABLE IF NOT EXISTS permissions (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    perm_name VARCHAR(100) NOT NULL,
    perm_code VARCHAR(100) NOT NULL UNIQUE,
    parent_id VARCHAR(64),
    sort_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 角色权限关联
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id VARCHAR(64) NOT NULL,
    perm_id VARCHAR(64) NOT NULL,
    PRIMARY KEY (role_id, perm_id)
);

-- 本地消息表（Transactional Outbox）
CREATE TABLE IF NOT EXISTS msg_record (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    msg_id VARCHAR(128) NOT NULL UNIQUE,
    business_key VARCHAR(128),
    exchange VARCHAR(100),
    routing_key VARCHAR(100),
    body TEXT NOT NULL,
    status SMALLINT DEFAULT 0, -- 0-待发送 1-已发送 2-消费成功 3-消费失败
    retry_count INT DEFAULT 0,
    max_retry INT DEFAULT 3,
    next_retry_time TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================
-- his_user 数据库表结构
-- ============================================================
\c his_user;

CREATE TABLE IF NOT EXISTS patients (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    user_id VARCHAR(64),
    name VARCHAR(100) NOT NULL,
    id_card VARCHAR(18) UNIQUE,
    phone VARCHAR(20),
    gender CHAR(1),
    birth_date DATE,
    address TEXT,
    allergy_history TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS employees (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    user_id VARCHAR(64),
    name VARCHAR(50) NOT NULL,
    phone VARCHAR(20),
    dept_id VARCHAR(64),
    title VARCHAR(50),
    specialty TEXT,
    introduction TEXT,
    status SMALLINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS departments (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    name VARCHAR(100) NOT NULL,
    parent_id VARCHAR(64),
    description TEXT,
    sort_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
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

-- ============================================================
-- 其他数据库占位
-- 各微服务的详细表结构将在后续迭代中补充
-- ============================================================
\c his_registration;
-- 挂号记录表、号源表等

\c his_clinic;
-- 接诊记录表、检查申请表等

\c his_emr;
-- 病历表、病历段落表、质控记录表、模板表等

\c his_prescription;
-- 处方表、处方明细表等

\c his_billing;
-- 收费单表、收费明细表等

\c his_pharmacy;
-- 药品表、入库记录表、发药记录表等

\c his_examination;
-- 检查申请表、检查报告表等

\c his_inpatient;
-- 住院记录表、医嘱表、护理记录表等

\c his_schedule;
-- 排班表、号源表等

\c his_outpatient;
-- 问诊记录表、慢病签约表、健康数据表等

\c his_followup;
-- 随访计划表、随访任务表、满意度调查表等

\c his_health_record;
-- 健康档案汇总表、时间轴事件表等

\c his_notification;
-- 通知记录表、通知模板表等

\c his_statistics;
-- 统计结果缓存表等

\c his_system;
-- 字典类型表、字典项表、系统参数表、操作日志表等
