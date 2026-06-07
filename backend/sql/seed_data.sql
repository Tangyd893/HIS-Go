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

-- 演示用户（密码: demo123，bcrypt 加密；仅用于本地/演示环境，生产环境请替换）
-- NOSONAR: 以下 bcrypt hash 为公开演示密码（demo123），非生产密钥
INSERT INTO users (id, username, password, real_name, phone, role, dept_id) VALUES
('demo-doctor', 'demo-doctor', '$2a$10$avEr2y6CrrENS8/NMWeeNOJcA2S76iJOzdkZzLLnrvbmor6fLiQVW', '张医生', '13800000001', 'doctor', 'dept_001'),        -- NOSONAR
('demo-nurse', 'demo-nurse', '$2a$10$avEr2y6CrrENS8/NMWeeNOJcA2S76iJOzdkZzLLnrvbmor6fLiQVW', '李护士', '13800000002', 'nurse', 'dept_001'),      -- NOSONAR
('demo-admin', 'demo-admin', '$2a$10$avEr2y6CrrENS8/NMWeeNOJcA2S76iJOzdkZzLLnrvbmor6fLiQVW', '管理员', '13800000003', 'admin', NULL),          -- NOSONAR
('demo-patient', 'demo-patient', '$2a$10$avEr2y6CrrENS8/NMWeeNOJcA2S76iJOzdkZzLLnrvbmor6fLiQVW', '王小明', '13900000001', 'patient', NULL);    -- NOSONAR

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

-- ============================================================
-- 演示业务数据：排班、药品
-- 用于支撑「挂号→接诊→处方→收费→发药」主路径演示
-- ============================================================

\c his_registration;

-- 排班号源（未来 7 天，每天上午/下午各一个号源）
INSERT INTO schedules (id, dept_id, dept_name, doctor_id, doctor_name, date, time_slot, total_count, remain_count, fee, status) VALUES
('sched_001', 'dept_001', '内科', 'demo-doctor', '张医生', to_char(CURRENT_DATE, 'YYYY-MM-DD'), 1, 30, 30, 15.00, 1),
('sched_002', 'dept_001', '内科', 'demo-doctor', '张医生', to_char(CURRENT_DATE, 'YYYY-MM-DD'), 2, 30, 30, 15.00, 1),
('sched_003', 'dept_001', '内科', 'demo-doctor', '张医生', to_char(CURRENT_DATE + INTERVAL '1 day', 'YYYY-MM-DD'), 1, 30, 30, 15.00, 1),
('sched_004', 'dept_001', '内科', 'demo-doctor', '张医生', to_char(CURRENT_DATE + INTERVAL '1 day', 'YYYY-MM-DD'), 2, 30, 30, 15.00, 1),
('sched_005', 'dept_002', '外科', 'demo-doctor', '张医生', to_char(CURRENT_DATE, 'YYYY-MM-DD'), 1, 20, 20, 20.00, 1);

\c his_schedule;

-- 排班管理（与 registration 侧对应）
INSERT INTO schedules (id, doctor_id, doctor_name, dept_id, dept_name, work_date, time_slot, max_patients, current_patients, room_no, status) VALUES
('sched_001', 'demo-doctor', '张医生', 'dept_001', '内科', to_char(CURRENT_DATE, 'YYYY-MM-DD'), 1, 30, 0, '201', 1),
('sched_002', 'demo-doctor', '张医生', 'dept_001', '内科', to_char(CURRENT_DATE, 'YYYY-MM-DD'), 2, 30, 0, '201', 1),
('sched_003', 'demo-doctor', '张医生', 'dept_001', '内科', to_char(CURRENT_DATE + INTERVAL '1 day', 'YYYY-MM-DD'), 1, 30, 0, '201', 1),
('sched_004', 'demo-doctor', '张医生', 'dept_001', '内科', to_char(CURRENT_DATE + INTERVAL '1 day', 'YYYY-MM-DD'), 2, 30, 0, '201', 1),
('sched_005', 'demo-doctor', '张医生', 'dept_002', '外科', to_char(CURRENT_DATE, 'YYYY-MM-DD'), 1, 20, 0, '302', 1);

\c his_pharmacy;

-- 常用药品
INSERT INTO drugs (id, name, generic_name, specification, manufacturer, batch_no, purchase_price, retail_price, stock, min_stock, expiry_date, status) VALUES
('drug_001', '阿莫西林胶囊', '阿莫西林', '0.5g×24粒', '华北制药', 'B20260601', 8.50, 12.00, 500, 50, '2028-06-01', 1),
('drug_002', '布洛芬缓释胶囊', '布洛芬', '0.3g×20粒', '中美史克', 'B20260515', 12.00, 18.50, 300, 30, '2028-05-15', 1),
('drug_003', '阿司匹林肠溶片', '阿司匹林', '100mg×30片', '拜耳医药', 'B20260420', 10.00, 15.00, 400, 40, '2028-04-20', 1),
('drug_004', '奥美拉唑肠溶胶囊', '奥美拉唑', '20mg×14粒', '阿斯利康', 'B20260310', 25.00, 38.00, 200, 20, '2028-03-10', 1),
('drug_005', '盐酸氨溴索片', '氨溴索', '30mg×20片', '上海医药', 'B20260501', 6.00, 9.50, 600, 60, '2028-05-01', 1);

-- ============================================================
-- 患者端演示数据：处方、检查报告、随访计划
-- ============================================================

\c his_prescription;

INSERT INTO prescriptions (id, patient_id, patient_name, doctor_id, prescription_type, status, note, created_at, updated_at) VALUES
('pres_demo_001', 'patient_001', '王小明', 'demo-doctor', 1, 2, '感冒用药', CURRENT_TIMESTAMP - INTERVAL '3 days', CURRENT_TIMESTAMP - INTERVAL '3 days'),
('pres_demo_002', 'patient_001', '王小明', 'demo-doctor', 1, 1, '胃药', CURRENT_TIMESTAMP - INTERVAL '1 day', CURRENT_TIMESTAMP - INTERVAL '1 day');

INSERT INTO prescription_details (id, prescription_id, drug_id, drug_name, specification, dosage, usage, frequency, days, quantity, unit_price, note) VALUES
('presd_001', 'pres_demo_001', 'drug_001', '阿莫西林胶囊', '0.5g×24粒', 0.5, '口服', '每日三次', 3, 2, 12.00, '饭后服用'),
('presd_002', 'pres_demo_001', 'drug_005', '盐酸氨溴索片', '30mg×20片', 30, '口服', '每日两次', 5, 1, 9.50, ''),
('presd_003', 'pres_demo_002', 'drug_004', '奥美拉唑肠溶胶囊', '20mg×14粒', 20, '口服', '每日一次', 14, 1, 38.00, '晨起空腹');

\c his_examination;

INSERT INTO examination_reports (id, patient_id, patient_name, exam_type, exam_item, body_part, findings, impression, conclusion, technician_id, reviewer_id, status, created_at, updated_at) VALUES
('exam_demo_001', 'patient_001', '王小明', 'CT', '胸部CT平扫', '胸部', '双肺纹理清晰，未见明显异常密度影', '胸部CT未见异常', '胸部CT未见异常', 'technician_001', 'demo-doctor', 3, CURRENT_TIMESTAMP - INTERVAL '7 days', CURRENT_TIMESTAMP - INTERVAL '7 days'),
('exam_demo_002', 'patient_001', '王小明', '检验', '血常规', '血液', '白细胞 6.5×10⁹/L，中性粒细胞 62%', '血常规正常', '血常规各项指标在正常范围内', 'technician_001', 'demo-doctor', 2, CURRENT_TIMESTAMP - INTERVAL '2 days', CURRENT_TIMESTAMP - INTERVAL '2 days');

\c his_followup;

INSERT INTO followup_plans (id, patient_id, plan_name, start_date, end_date, frequency, status, created_at, updated_at) VALUES
('fup_demo_001', 'patient_001', '高血压随访计划', to_char(CURRENT_DATE - INTERVAL '30 days', 'YYYY-MM-DD'), to_char(CURRENT_DATE + INTERVAL '60 days', 'YYYY-MM-DD'), 2, 1, CURRENT_TIMESTAMP - INTERVAL '30 days', CURRENT_TIMESTAMP - INTERVAL '30 days'),
('fup_demo_002', 'patient_001', '术后康复随访', to_char(CURRENT_DATE, 'YYYY-MM-DD'), to_char(CURRENT_DATE + INTERVAL '90 days', 'YYYY-MM-DD'), 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT INTO followup_tasks (id, plan_id, task_name, description, execute_date, status, created_at) VALUES
('fut_001', 'fup_demo_001', '血压测量', '自测血压并记录数值', to_char(CURRENT_DATE - INTERVAL '21 days', 'YYYY-MM-DD'), 1, CURRENT_TIMESTAMP - INTERVAL '30 days'),
('fut_002', 'fup_demo_001', '电话随访', '电话询问用药情况', to_char(CURRENT_DATE - INTERVAL '7 days', 'YYYY-MM-DD'), 1, CURRENT_TIMESTAMP - INTERVAL '30 days'),
('fut_003', 'fup_demo_001', '门诊复查', '来院复查血压控制情况', to_char(CURRENT_DATE + INTERVAL '30 days', 'YYYY-MM-DD'), 0, CURRENT_TIMESTAMP - INTERVAL '30 days');

-- ============================================================
-- 大量演示假数据 — 用于功能验证
-- ============================================================

\c his_user;

-- 补充患者
INSERT INTO patients (id, name, id_card, phone, gender, birth_date, address) VALUES
('patient_003', '张三', '110101198503031234', '13900000003', 'M', '1985-03-03', '北京市东城区'),
('patient_004', '赵四', '110101199207041234', '13900000004', 'M', '1992-07-04', '北京市西城区'),
('patient_005', '孙七', '110101198811051234', '13900000005', 'F', '1988-11-05', '北京市通州区'),
('patient_006', '周八', '110101199509061234', '13900000006', 'F', '1995-09-06', '北京市丰台区'),
('patient_007', '吴九', '110101197804071234', '13900000007', 'M', '1978-04-07', '北京市石景山区');

-- 补充员工
INSERT INTO employees (id, user_id, name, phone, dept_id, title, status) VALUES
('emp_001', 'demo-doctor', '张医生', '13800000001', 'dept_001', '主任医师', 1),
('emp_002', 'demo-nurse', '李护士', '13800000002', 'dept_001', '主管护师', 1);

\c his_registration;

-- 扩展排班 (未来 5 天，3 个科室)
INSERT INTO schedules (id, dept_id, dept_name, doctor_id, doctor_name, date, time_slot, total_count, remain_count, fee, status) VALUES
('sched_006', 'dept_003', '儿科', 'demo-doctor', '张医生', to_char(CURRENT_DATE, 'YYYY-MM-DD'), 1, 25, 23, 15.00, 1),
('sched_007', 'dept_003', '儿科', 'demo-doctor', '张医生', to_char(CURRENT_DATE, 'YYYY-MM-DD'), 2, 25, 25, 15.00, 1),
('sched_008', 'dept_004', '妇产科', 'demo-doctor', '张医生', to_char(CURRENT_DATE, 'YYYY-MM-DD'), 1, 20, 18, 20.00, 1),
('sched_009', 'dept_005', '急诊科', 'demo-doctor', '张医生', to_char(CURRENT_DATE, 'YYYY-MM-DD'), 1, 50, 45, 25.00, 1),
('sched_010', 'dept_001', '内科', 'demo-doctor', '张医生', to_char(CURRENT_DATE + INTERVAL '2 days', 'YYYY-MM-DD'), 1, 30, 28, 15.00, 1),
('sched_011', 'dept_002', '外科', 'demo-doctor', '张医生', to_char(CURRENT_DATE + INTERVAL '2 days', 'YYYY-MM-DD'), 1, 20, 19, 20.00, 1),
('sched_012', 'dept_003', '儿科', 'demo-doctor', '张医生', to_char(CURRENT_DATE + INTERVAL '2 days', 'YYYY-MM-DD'), 1, 25, 25, 15.00, 1);

-- 挂号记录 (过去几天的模拟数据)
INSERT INTO registrations (id, patient_id, patient_name, schedule_id, registration_date, queue_number, status, created_at, updated_at) VALUES
('reg_001', 'patient_001', '王小明', 'sched_001', to_char(CURRENT_DATE - INTERVAL '1 day', 'YYYY-MM-DD'), 1, 2, CURRENT_TIMESTAMP - INTERVAL '1 day', CURRENT_TIMESTAMP - INTERVAL '1 day'),
('reg_002', 'patient_003', '张三',   'sched_001', to_char(CURRENT_DATE, 'YYYY-MM-DD'), 2, 0, CURRENT_TIMESTAMP - INTERVAL '2 hours', CURRENT_TIMESTAMP - INTERVAL '2 hours'),
('reg_003', 'patient_004', '赵四',   'sched_001', to_char(CURRENT_DATE, 'YYYY-MM-DD'), 3, 0, CURRENT_TIMESTAMP - INTERVAL '1 hour', CURRENT_TIMESTAMP - INTERVAL '1 hour'),
('reg_004', 'patient_005', '孙七',   'sched_006', to_char(CURRENT_DATE, 'YYYY-MM-DD'), 1, 1, CURRENT_TIMESTAMP - INTERVAL '3 hours', CURRENT_TIMESTAMP - INTERVAL '2 hours'),
('reg_005', 'patient_006', '周八',   'sched_008', to_char(CURRENT_DATE, 'YYYY-MM-DD'), 1, 0, CURRENT_TIMESTAMP - INTERVAL '30 minutes', CURRENT_TIMESTAMP - INTERVAL '30 minutes'),
('reg_006', 'patient_007', '吴九',   'sched_009', to_char(CURRENT_DATE, 'YYYY-MM-DD'), 5, 2, CURRENT_TIMESTAMP - INTERVAL '4 hours', CURRENT_TIMESTAMP - INTERVAL '3 hours');

\c his_clinic;

-- 门诊记录
INSERT INTO clinic_records (id, registration_id, patient_id, patient_name, doctor_id, chief_complaint, diagnosis, status, created_at, updated_at) VALUES
('clinic_001', 'reg_001', 'patient_001', '王小明', 'demo-doctor', '咳嗽三天，加重一天', '急性支气管炎', 2, CURRENT_TIMESTAMP - INTERVAL '1 day', CURRENT_TIMESTAMP - INTERVAL '1 day'),
('clinic_002', 'reg_004', 'patient_005', '孙七', 'demo-doctor', '发热两天，体温38.5℃', '上呼吸道感染', 2, CURRENT_TIMESTAMP - INTERVAL '2 hours', CURRENT_TIMESTAMP - INTERVAL '2 hours'),
('clinic_003', 'reg_006', 'patient_007', '吴九', 'demo-doctor', '外伤后右下肢疼痛', '右下肢软组织挫伤', 2, CURRENT_TIMESTAMP - INTERVAL '3 hours', CURRENT_TIMESTAMP - INTERVAL '3 hours');

\c his_prescription;

-- 更多处方
INSERT INTO prescriptions (id, patient_id, patient_name, doctor_id, prescription_type, status, note, created_at, updated_at) VALUES
('pres_demo_003', 'patient_005', '孙七', 'demo-doctor', 1, 4, '上感用药', CURRENT_TIMESTAMP - INTERVAL '2 hours', CURRENT_TIMESTAMP - INTERVAL '1 hour'),
('pres_demo_004', 'patient_007', '吴九', 'demo-doctor', 1, 2, '外伤', CURRENT_TIMESTAMP - INTERVAL '3 hours', CURRENT_TIMESTAMP - INTERVAL '2 hours'),
('pres_demo_005', 'patient_003', '张三', 'demo-doctor', 2, 0, '', CURRENT_TIMESTAMP - INTERVAL '1 hour', CURRENT_TIMESTAMP - INTERVAL '1 hour');

INSERT INTO prescription_details (id, prescription_id, drug_id, drug_name, specification, dosage, usage, frequency, days, quantity, unit_price, note) VALUES
('presd_004', 'pres_demo_003', 'drug_001', '阿莫西林胶囊', '0.5g×24粒', 0.5, '口服', '每日三次', 5, 2, 12.00, ''),
('presd_005', 'pres_demo_003', 'drug_002', '布洛芬缓释胶囊', '0.3g×20粒', 0.3, '口服', '必要时', 3, 1, 18.50, '体温>38.5℃时服用'),
('presd_006', 'pres_demo_004', 'drug_002', '布洛芬缓释胶囊', '0.3g×20粒', 0.3, '口服', '每日两次', 3, 1, 18.50, '止痛'),
('presd_007', 'pres_demo_005', 'drug_003', '阿司匹林肠溶片', '100mg×30片', 100, '口服', '每日一次', 30, 1, 15.00, '早餐后服用');

\c his_billing;

-- 账单
INSERT INTO bills (id, patient_id, registration_id, bill_no, total_amount, paid_amount, pay_method, status, created_at, updated_at) VALUES
('bill_001', 'patient_001', 'reg_001', 'BL20260601001', 33.50, 33.50, 1, 1, CURRENT_TIMESTAMP - INTERVAL '1 day', CURRENT_TIMESTAMP - INTERVAL '1 day'),
('bill_002', 'patient_005', 'reg_004', 'BL20260604001', 42.50, 42.50, 2, 1, CURRENT_TIMESTAMP - INTERVAL '1 hour', CURRENT_TIMESTAMP - INTERVAL '1 hour'),
('bill_003', 'patient_007', 'reg_006', 'BL20260604002', 18.50, 18.50, 3, 1, CURRENT_TIMESTAMP - INTERVAL '2 hours', CURRENT_TIMESTAMP - INTERVAL '2 hours'),
('bill_004', 'patient_003', 'reg_002', 'BL20260604003', 15.00, 0, 0, 0, CURRENT_TIMESTAMP - INTERVAL '1 hour', CURRENT_TIMESTAMP - INTERVAL '1 hour'),
('bill_005', 'patient_004', 'reg_003', 'BL20260604004', 12.00, 0, 0, 0, CURRENT_TIMESTAMP - INTERVAL '30 minutes', CURRENT_TIMESTAMP - INTERVAL '30 minutes');

-- 账单明细
INSERT INTO bill_details (id, bill_id, item_type, item_name, unit_price, quantity, amount) VALUES
('bld_001', 'bill_001', 1, '阿莫西林胶囊', 12.00, 2, 24.00),
('bld_002', 'bill_001', 1, '盐酸氨溴索片', 9.50, 1, 9.50),
('bld_003', 'bill_002', 1, '阿莫西林胶囊', 12.00, 2, 24.00),
('bld_004', 'bill_002', 1, '布洛芬缓释胶囊', 18.50, 1, 18.50),
('bld_005', 'bill_003', 1, '布洛芬缓释胶囊', 18.50, 1, 18.50),
('bld_006', 'bill_004', 1, '阿司匹林肠溶片', 15.00, 1, 15.00),
('bld_007', 'bill_005', 2, '挂号费', 12.00, 1, 12.00);

\c his_pharmacy;

-- 发药记录
INSERT INTO dispense_records (id, prescription_id, patient_id, drug_id, quantity, dispenser_id, status, created_at) VALUES
('disp_001', 'pres_demo_001', 'patient_001', 'drug_001', 2, 'demo-nurse', 1, CURRENT_TIMESTAMP - INTERVAL '1 day'),
('disp_002', 'pres_demo_003', 'patient_005', 'drug_001', 2, 'demo-nurse', 1, CURRENT_TIMESTAMP - INTERVAL '1 hour'),
('disp_003', 'pres_demo_003', 'patient_005', 'drug_002', 1, 'demo-nurse', 1, CURRENT_TIMESTAMP - INTERVAL '1 hour'),
('disp_004', 'pres_demo_004', 'patient_007', 'drug_002', 1, 'demo-nurse', 1, CURRENT_TIMESTAMP - INTERVAL '2 hours');

\c his_examination;

-- 更多检查报告
INSERT INTO examination_reports (id, patient_id, patient_name, exam_type, exam_item, body_part, findings, impression, conclusion, technician_id, reviewer_id, status, created_at, updated_at) VALUES
('exam_demo_003', 'patient_005', '孙七', 'DR', '胸部正位片', '胸部', '双肺纹理增粗，右下肺可见片状阴影', '右下肺炎症', '右下叶肺炎，建议抗感染治疗后复查', 'technician_001', 'demo-doctor', 3, CURRENT_TIMESTAMP - INTERVAL '3 hours', CURRENT_TIMESTAMP - INTERVAL '2 hours'),
('exam_demo_004', 'patient_007', '吴九', 'DR', '右下肢正侧位片', '右下肢', '右胫骨未见明确骨折线，软组织肿胀', '右下肢软组织损伤', '排除骨折，软组织挫伤', 'technician_001', 'demo-doctor', 3, CURRENT_TIMESTAMP - INTERVAL '3 hours', CURRENT_TIMESTAMP - INTERVAL '3 hours'),
('exam_demo_005', 'patient_002', '李小红', '超声', '腹部超声', '腹部', '肝脏大小形态正常，胆囊壁光滑，未见异常回声', '腹部超声未见异常', '腹部超声未见异常', 'technician_001', 'demo-doctor', 2, CURRENT_TIMESTAMP - INTERVAL '5 days', CURRENT_TIMESTAMP - INTERVAL '5 days');

\c his_followup;

-- 更多随访
INSERT INTO followup_plans (id, patient_id, plan_name, start_date, end_date, frequency, status, created_at, updated_at) VALUES
('fup_demo_003', 'patient_005', '肺炎康复随访', to_char(CURRENT_DATE, 'YYYY-MM-DD'), to_char(CURRENT_DATE + INTERVAL '30 days', 'YYYY-MM-DD'), 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('fup_demo_004', 'patient_002', '年度健康随访', to_char(CURRENT_DATE - INTERVAL '60 days', 'YYYY-MM-DD'), to_char(CURRENT_DATE + INTERVAL '120 days', 'YYYY-MM-DD'), 3, 1, CURRENT_TIMESTAMP - INTERVAL '60 days', CURRENT_TIMESTAMP - INTERVAL '60 days');

INSERT INTO followup_tasks (id, plan_id, task_name, description, execute_date, status, created_at) VALUES
('fut_004', 'fup_demo_003', '首次电话随访', '询问体温恢复情况和用药依从性', to_char(CURRENT_DATE + INTERVAL '3 days', 'YYYY-MM-DD'), 0, CURRENT_TIMESTAMP),
('fut_005', 'fup_demo_003', '复查胸片', '来院复查胸部正位片', to_char(CURRENT_DATE + INTERVAL '14 days', 'YYYY-MM-DD'), 0, CURRENT_TIMESTAMP),
('fut_006', 'fup_demo_004', '健康问卷', '完成年度健康状况调查问卷', to_char(CURRENT_DATE - INTERVAL '30 days', 'YYYY-MM-DD'), 1, CURRENT_TIMESTAMP - INTERVAL '60 days');
