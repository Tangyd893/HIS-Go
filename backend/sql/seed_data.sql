-- ============================================================
-- HIS-Go 种子数据脚本
-- 说明：基础字典数据和演示账号
-- ============================================================

\c his_auth;

-- 角色数据
INSERT INTO roles (id, role_name, role_code, description) VALUES
('role_admin', '管理员', 'admin', '系统管理员'),
('role_doctor', '医生', 'doctor', '医生角色'),
('role_nurse', '护士', 'nurse', '护士角色'),
('role_patient', '患者', 'patient', '患者角色'),
('role_pharmacist', '药师', 'pharmacist', '药师角色');

-- 演示用户（密码: demo123，bcrypt加密）
INSERT INTO users (id, username, password, real_name, phone, role, dept_id) VALUES
('demo-doctor', 'demo-doctor', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '张医生', '13800000001', 'doctor', 'dept_001'),
('demo-nurse', 'demo-nurse', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '李护士', '13800000002', 'nurse', 'dept_001'),
('demo-admin', 'demo-admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '管理员', '13800000003', 'admin', NULL);

\c his_user;

-- 科室数据
INSERT INTO departments (id, name, parent_id, description, sort_order) VALUES
('dept_001', '内科', NULL, '内科门诊', 1),
('dept_002', '外科', NULL, '外科门诊', 2),
('dept_003', '儿科', NULL, '儿科门诊', 3),
('dept_004', '妇产科', NULL, '妇产科门诊', 4),
('dept_005', '急诊科', NULL, '急诊科室', 5);

-- 演示患者
INSERT INTO patients (id, name, id_card, phone, gender, birth_date, address) VALUES
('patient_001', '王小明', '110101199001011234', '13900000001', 'M', '1990-01-01', '北京市朝阳区'),
('patient_002', '李小红', '110101199502022345', '13900000002', 'F', '1995-02-02', '北京市海淀区');

\c his_system;

-- 字典类型
CREATE TABLE IF NOT EXISTS dict_types (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    dict_name VARCHAR(100) NOT NULL,
    dict_type VARCHAR(100) NOT NULL UNIQUE,
    status SMALLINT DEFAULT 1,
    remark TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 字典项
CREATE TABLE IF NOT EXISTS dict_items (
    id VARCHAR(64) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    dict_type VARCHAR(100) NOT NULL,
    label VARCHAR(100) NOT NULL,
    value VARCHAR(100) NOT NULL,
    sort_order INT DEFAULT 0,
    status SMALLINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 处方类型
INSERT INTO dict_types (dict_name, dict_type, remark) VALUES ('处方类型', 'prescription_type', '处方类型字典');
INSERT INTO dict_items (dict_type, label, value, sort_order) VALUES
('prescription_type', '西药', '1', 1),
('prescription_type', '中成药', '2', 2),
('prescription_type', '中草药', '3', 3);

-- 支付方式
INSERT INTO dict_types (dict_name, dict_type, remark) VALUES ('支付方式', 'pay_method', '支付方式字典');
INSERT INTO dict_items (dict_type, label, value, sort_order) VALUES
('pay_method', '现金', '1', 1),
('pay_method', '微信', '2', 2),
('pay_method', '支付宝', '3', 3),
('pay_method', '医保', '4', 4);

-- 随访类型
INSERT INTO dict_types (dict_name, dict_type, remark) VALUES ('随访类型', 'followup_type', '随访类型字典');
INSERT INTO dict_items (dict_type, label, value, sort_order) VALUES
('followup_type', '电话随访', '1', 1),
('followup_type', '问卷随访', '2', 2),
('followup_type', '上门随访', '3', 3);
