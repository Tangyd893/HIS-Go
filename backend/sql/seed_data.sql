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

-- 科室数据（≥12 个科室，覆盖常见门诊）
INSERT INTO departments (id, name, parent_id, description, sort_order) VALUES
('dept_001', '内科', NULL, '内科门诊', 1),
('dept_002', '外科', NULL, '外科门诊', 2),
('dept_003', '儿科', NULL, '儿科门诊', 3),
('dept_004', '妇产科', NULL, '妇产科门诊', 4),
('dept_005', '急诊科', NULL, '急诊科室', 5),
('dept_006', '骨科', NULL, '骨科门诊', 6),
('dept_007', '眼科', NULL, '眼科门诊', 7),
('dept_008', '耳鼻喉科', NULL, '耳鼻喉科门诊', 8),
('dept_009', '口腔科', NULL, '口腔科门诊', 9),
('dept_010', '皮肤科', NULL, '皮肤科门诊', 10),
('dept_011', '中医科', NULL, '中医科门诊', 11),
('dept_012', '康复科', NULL, '康复理疗科', 12);

-- 演示患者
INSERT INTO patients (id, user_id, name, id_card, phone, gender, birth_date, address) VALUES
('patient_001', 'demo-patient', '王小明', '110101199001011234', '13900000001', 'M', '1990-01-01', '北京市朝阳区'),
('patient_002', NULL, '李小红', '110101199502022345', '13900000002', 'F', '1995-02-02', '北京市海淀区');

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

-- 清理过期排班（今日之前）
DELETE FROM schedules WHERE date < CURRENT_DATE;

-- 清理已取消/已完成的过期挂号
DELETE FROM registrations WHERE registration_date < to_char(CURRENT_DATE, 'YYYY-MM-DD');

-- ============================================================
-- 自动生成全科室 × 未来 7 天 × 上午/下午 排班号源
-- 科室-医生映射：
--   内科→张医生  外科→王医生  儿科→李医生  妇产科→赵医生
--   急诊科→陈医生  骨科→张医生  眼科→王医生  耳鼻喉科→李医生
--   口腔科→赵医生  皮肤科→陈医生  中医科→张医生  康复科→王医生
-- ============================================================
DO $$
DECLARE
    dept RECORD;
    d INT;
    slot INT;
    sched_id TEXT;
    doc_id TEXT;
    doc_name TEXT;
    fee_val NUMERIC;
    total_val INT;
BEGIN
    FOR dept IN
        SELECT id, name FROM (
            VALUES
                ('dept_001','内科'),('dept_002','外科'),('dept_003','儿科'),
                ('dept_004','妇产科'),('dept_005','急诊科'),('dept_006','骨科'),
                ('dept_007','眼科'),('dept_008','耳鼻喉科'),('dept_009','口腔科'),
                ('dept_010','皮肤科'),('dept_011','中医科'),('dept_012','康复科')
        ) AS t(id, name)
    LOOP
        -- 医生轮换分配
        CASE dept.id
            WHEN 'dept_001' THEN doc_id := 'demo-doctor'; doc_name := '张医生'; fee_val := 15.00; total_val := 30;
            WHEN 'dept_002' THEN doc_id := 'doctor-wang'; doc_name := '王医生'; fee_val := 20.00; total_val := 25;
            WHEN 'dept_003' THEN doc_id := 'doctor-li'; doc_name := '李医生'; fee_val := 15.00; total_val := 25;
            WHEN 'dept_004' THEN doc_id := 'doctor-zhao'; doc_name := '赵医生'; fee_val := 20.00; total_val := 20;
            WHEN 'dept_005' THEN doc_id := 'doctor-chen'; doc_name := '陈医生'; fee_val := 25.00; total_val := 50;
            WHEN 'dept_006' THEN doc_id := 'demo-doctor'; doc_name := '张医生'; fee_val := 20.00; total_val := 25;
            WHEN 'dept_007' THEN doc_id := 'doctor-wang'; doc_name := '王医生'; fee_val := 15.00; total_val := 30;
            WHEN 'dept_008' THEN doc_id := 'doctor-li'; doc_name := '李医生'; fee_val := 15.00; total_val := 25;
            WHEN 'dept_009' THEN doc_id := 'doctor-zhao'; doc_name := '赵医生'; fee_val := 15.00; total_val := 20;
            WHEN 'dept_010' THEN doc_id := 'doctor-chen'; doc_name := '陈医生'; fee_val := 15.00; total_val := 30;
            WHEN 'dept_011' THEN doc_id := 'demo-doctor'; doc_name := '张医生'; fee_val := 20.00; total_val := 25;
            WHEN 'dept_012' THEN doc_id := 'doctor-wang'; doc_name := '王医生'; fee_val := 15.00; total_val := 20;
        END CASE;

        FOR d IN 0..6 LOOP
            FOR slot IN 1..2 LOOP
                sched_id := 'sched_' || dept.id || '_' || d || '_' || slot;
                INSERT INTO schedules (id, dept_id, dept_name, doctor_id, doctor_name, date, time_slot, total_count, remain_count, fee, status)
                VALUES (sched_id, dept.id, dept.name, doc_id, doc_name,
                        to_char(CURRENT_DATE + d, 'YYYY-MM-DD'), slot,
                        total_val, total_val, fee_val, 1)
                ON CONFLICT (id) DO NOTHING;
            END LOOP;
        END LOOP;
    END LOOP;
END $$;

\c his_schedule;

-- 清理过期排班
DELETE FROM schedules WHERE work_date < CURRENT_DATE;

-- 自动生成 his_schedule 侧排班（与 registration 对应）
DO $$
DECLARE
    dept RECORD;
    d INT;
    slot INT;
    sched_id TEXT;
    doc_id TEXT;
    doc_name TEXT;
    max_val INT;
    room TEXT;
BEGIN
    FOR dept IN
        SELECT id, name FROM (
            VALUES
                ('dept_001','内科'),('dept_002','外科'),('dept_003','儿科'),
                ('dept_004','妇产科'),('dept_005','急诊科'),('dept_006','骨科'),
                ('dept_007','眼科'),('dept_008','耳鼻喉科'),('dept_009','口腔科'),
                ('dept_010','皮肤科'),('dept_011','中医科'),('dept_012','康复科')
        ) AS t(id, name)
    LOOP
        CASE dept.id
            WHEN 'dept_001' THEN doc_id := 'demo-doctor'; doc_name := '张医生'; max_val := 30; room := '201';
            WHEN 'dept_002' THEN doc_id := 'doctor-wang'; doc_name := '王医生'; max_val := 25; room := '302';
            WHEN 'dept_003' THEN doc_id := 'doctor-li'; doc_name := '李医生'; max_val := 25; room := '103';
            WHEN 'dept_004' THEN doc_id := 'doctor-zhao'; doc_name := '赵医生'; max_val := 20; room := '405';
            WHEN 'dept_005' THEN doc_id := 'doctor-chen'; doc_name := '陈医生'; max_val := 50; room := 'E01';
            WHEN 'dept_006' THEN doc_id := 'demo-doctor'; doc_name := '张医生'; max_val := 25; room := '306';
            WHEN 'dept_007' THEN doc_id := 'doctor-wang'; doc_name := '王医生'; max_val := 30; room := '207';
            WHEN 'dept_008' THEN doc_id := 'doctor-li'; doc_name := '李医生'; max_val := 25; room := '208';
            WHEN 'dept_009' THEN doc_id := 'doctor-zhao'; doc_name := '赵医生'; max_val := 20; room := '309';
            WHEN 'dept_010' THEN doc_id := 'doctor-chen'; doc_name := '陈医生'; max_val := 30; room := '410';
            WHEN 'dept_011' THEN doc_id := 'demo-doctor'; doc_name := '张医生'; max_val := 25; room := '511';
            WHEN 'dept_012' THEN doc_id := 'doctor-wang'; doc_name := '王医生'; max_val := 20; room := '612';
        END CASE;

        FOR d IN 0..6 LOOP
            FOR slot IN 1..2 LOOP
                sched_id := 'sched_' || dept.id || '_' || d || '_' || slot;
                INSERT INTO schedules (id, doctor_id, doctor_name, dept_id, dept_name, work_date, time_slot, max_patients, current_patients, room_no, status)
                VALUES (sched_id, doc_id, doc_name, dept.id, dept.name,
                        to_char(CURRENT_DATE + d, 'YYYY-MM-DD'), slot,
                        max_val, 0, room, 1)
                ON CONFLICT (id) DO NOTHING;
            END LOOP;
        END LOOP;
    END LOOP;
END $$;

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

INSERT INTO followup_tasks (id, plan_id, assignee_id, execute_date, type, content, status, created_at) VALUES
('fut_001', 'fup_demo_001', 'demo-nurse', to_char(CURRENT_DATE - INTERVAL '21 days', 'YYYY-MM-DD'), 1, '自测血压并记录数值', 1, CURRENT_TIMESTAMP - INTERVAL '30 days'),
('fut_002', 'fup_demo_001', 'demo-nurse', to_char(CURRENT_DATE - INTERVAL '7 days', 'YYYY-MM-DD'), 1, '电话询问用药情况', 1, CURRENT_TIMESTAMP - INTERVAL '30 days'),
('fut_003', 'fup_demo_001', 'demo-doctor', to_char(CURRENT_DATE + INTERVAL '30 days', 'YYYY-MM-DD'), 2, '来院复查血压控制情况', 0, CURRENT_TIMESTAMP - INTERVAL '30 days');

-- ============================================================
-- 大量演示假数据 — 用于功能验证
-- ============================================================

\c his_user;

-- 补充患者
INSERT INTO patients (id, user_id, name, id_card, phone, gender, birth_date, address) VALUES
('patient_003', NULL, '张三', '110101198503031234', '13900000003', 'M', '1985-03-03', '北京市东城区'),
('patient_004', NULL, '赵四', '110101199207041234', '13900000004', 'M', '1992-07-04', '北京市西城区'),
('patient_005', NULL, '孙七', '110101198811051234', '13900000005', 'F', '1988-11-05', '北京市通州区'),
('patient_006', NULL, '周八', '110101199509061234', '13900000006', 'F', '1995-09-06', '北京市丰台区'),
('patient_007', NULL, '吴九', '110101197804071234', '13900000007', 'M', '1978-04-07', '北京市石景山区');

-- 补充员工
INSERT INTO employees (id, user_id, name, phone, dept_id, title, status) VALUES
('emp_001', 'demo-doctor', '张医生', '13800000001', 'dept_001', '主任医师', 1),
('emp_002', 'demo-nurse', '李护士', '13800000002', 'dept_001', '主管护师', 1),
('emp_007', 'demo-pharmacist', '王药师', '13800000004', 'dept_001', '主管药师', 1)
ON CONFLICT (id) DO NOTHING;

\c his_registration;

-- 挂号记录（基于自动生成的排班，引用新 ID 格式 sched_<dept>_<day>_<slot>）
-- 说明：day=0 为今天，slot=1 上午，slot=2 下午
INSERT INTO registrations (id, patient_id, patient_name, schedule_id, registration_date, queue_number, status, created_at, updated_at) VALUES
('reg_001', 'patient_001', '王小明', 'sched_dept_001_0_1', to_char(CURRENT_DATE, 'YYYY-MM-DD'), 1, 2, CURRENT_TIMESTAMP - INTERVAL '1 day', CURRENT_TIMESTAMP - INTERVAL '1 day'),
('reg_002', 'patient_003', '张三',   'sched_dept_001_0_1', to_char(CURRENT_DATE, 'YYYY-MM-DD'), 2, 0, CURRENT_TIMESTAMP - INTERVAL '2 hours', CURRENT_TIMESTAMP - INTERVAL '2 hours'),
('reg_003', 'patient_004', '赵四',   'sched_dept_001_0_1', to_char(CURRENT_DATE, 'YYYY-MM-DD'), 3, 0, CURRENT_TIMESTAMP - INTERVAL '1 hour', CURRENT_TIMESTAMP - INTERVAL '1 hour'),
('reg_004', 'patient_005', '孙七',   'sched_dept_003_0_1', to_char(CURRENT_DATE, 'YYYY-MM-DD'), 1, 1, CURRENT_TIMESTAMP - INTERVAL '3 hours', CURRENT_TIMESTAMP - INTERVAL '2 hours'),
('reg_005', 'patient_006', '周八',   'sched_dept_004_0_1', to_char(CURRENT_DATE, 'YYYY-MM-DD'), 1, 0, CURRENT_TIMESTAMP - INTERVAL '30 minutes', CURRENT_TIMESTAMP - INTERVAL '30 minutes'),
('reg_006', 'patient_007', '吴九',   'sched_dept_005_0_1', to_char(CURRENT_DATE, 'YYYY-MM-DD'), 5, 2, CURRENT_TIMESTAMP - INTERVAL '4 hours', CURRENT_TIMESTAMP - INTERVAL '3 hours');

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

INSERT INTO followup_tasks (id, plan_id, assignee_id, execute_date, type, content, status, created_at) VALUES
('fut_004', 'fup_demo_003', 'demo-nurse', to_char(CURRENT_DATE + INTERVAL '3 days', 'YYYY-MM-DD'), 1, '询问体温恢复情况和用药依从性', 0, CURRENT_TIMESTAMP),
('fut_005', 'fup_demo_003', 'demo-doctor', to_char(CURRENT_DATE + INTERVAL '14 days', 'YYYY-MM-DD'), 2, '来院复查胸部正位片', 0, CURRENT_TIMESTAMP),
('fut_006', 'fup_demo_004', 'demo-nurse', to_char(CURRENT_DATE - INTERVAL '30 days', 'YYYY-MM-DD'), 2, '完成年度健康状况调查问卷', 1, CURRENT_TIMESTAMP - INTERVAL '60 days');

-- ============================================================
-- 第二轮演示数据扩充：多医生、慢病签约、健康档案时间轴
-- ============================================================

\c his_auth;

-- 补充医生账号（覆盖各科室演示）
INSERT INTO users (id, username, password, real_name, phone, role, dept_id) VALUES
('doctor-wang', 'doctor-wang', '$2a$10$avEr2y6CrrENS8/NMWeeNOJcA2S76iJOzdkZzLLnrvbmor6fLiQVW', '王医生', '13800000011', 'doctor', 'dept_002'),    -- NOSONAR
('doctor-li', 'doctor-li', '$2a$10$avEr2y6CrrENS8/NMWeeNOJcA2S76iJOzdkZzLLnrvbmor6fLiQVW', '李医生', '13800000012', 'doctor', 'dept_003'),      -- NOSONAR
('doctor-zhao', 'doctor-zhao', '$2a$10$avEr2y6CrrENS8/NMWeeNOJcA2S76iJOzdkZzLLnrvbmor6fLiQVW', '赵医生', '13800000013', 'doctor', 'dept_004'),  -- NOSONAR
('doctor-chen', 'doctor-chen', '$2a$10$avEr2y6CrrENS8/NMWeeNOJcA2S76iJOzdkZzLLnrvbmor6fLiQVW', '陈医生', '13800000014', 'doctor', 'dept_005');  -- NOSONAR

\c his_user;

-- 补充员工记录
INSERT INTO employees (id, user_id, name, phone, dept_id, title, status) VALUES
('emp_003', 'doctor-wang', '王医生', '13800000011', 'dept_002', '副主任医师', 1),
('emp_004', 'doctor-li', '李医生', '13800000012', 'dept_003', '主治医师', 1),
('emp_005', 'doctor-zhao', '赵医生', '13800000013', 'dept_004', '主任医师', 1),
('emp_006', 'doctor-chen', '陈医生', '13800000014', 'dept_005', '副主任医师', 1);

-- 补充患者（覆盖不同年龄段）
INSERT INTO patients (id, name, id_card, phone, gender, birth_date, address) VALUES
('patient_008', '郑十', '110101196505081234', '13900000008', 'M', '1965-05-08', '北京市大兴区'),
('patient_009', '陈一一', '110101200109091234', '13900000009', 'F', '2001-09-09', '北京市昌平区'),
('patient_010', '刘十二', '110101197210101234', '13900000010', 'M', '1972-10-10', '北京市顺义区');

\c his_registration;

-- 多医生排班（未来 3 天，每天上午）
INSERT INTO schedules (id, dept_id, dept_name, doctor_id, doctor_name, date, time_slot, total_count, remain_count, fee, status) VALUES
('sched_013', 'dept_002', '外科', 'doctor-wang', '王医生', to_char(CURRENT_DATE, 'YYYY-MM-DD'), 1, 20, 20, 20.00, 1),
('sched_014', 'dept_002', '外科', 'doctor-wang', '王医生', to_char(CURRENT_DATE, 'YYYY-MM-DD'), 2, 20, 20, 20.00, 1),
('sched_015', 'dept_003', '儿科', 'doctor-li', '李医生', to_char(CURRENT_DATE, 'YYYY-MM-DD'), 1, 25, 25, 15.00, 1),
('sched_016', 'dept_004', '妇产科', 'doctor-zhao', '赵医生', to_char(CURRENT_DATE, 'YYYY-MM-DD'), 1, 20, 20, 20.00, 1),
('sched_017', 'dept_005', '急诊科', 'doctor-chen', '陈医生', to_char(CURRENT_DATE, 'YYYY-MM-DD'), 1, 50, 50, 25.00, 1),
('sched_018', 'dept_001', '内科', 'demo-doctor', '张医生', to_char(CURRENT_DATE + INTERVAL '1 day', 'YYYY-MM-DD'), 1, 30, 30, 15.00, 1),
('sched_019', 'dept_002', '外科', 'doctor-wang', '王医生', to_char(CURRENT_DATE + INTERVAL '1 day', 'YYYY-MM-DD'), 1, 20, 20, 20.00, 1),
('sched_020', 'dept_003', '儿科', 'doctor-li', '李医生', to_char(CURRENT_DATE + INTERVAL '1 day', 'YYYY-MM-DD'), 1, 25, 25, 15.00, 1);

-- 更多挂号记录（覆盖各状态：0=待签到, 1=已签到, 2=已完成, 3=已取消）
INSERT INTO registrations (id, patient_id, patient_name, schedule_id, registration_date, queue_number, status, created_at, updated_at) VALUES
('reg_007', 'patient_008', '郑十',   'sched_013', to_char(CURRENT_DATE, 'YYYY-MM-DD'), 1, 1, CURRENT_TIMESTAMP - INTERVAL '4 hours', CURRENT_TIMESTAMP - INTERVAL '3 hours'),
('reg_008', 'patient_009', '陈一一', 'sched_015', to_char(CURRENT_DATE, 'YYYY-MM-DD'), 1, 0, CURRENT_TIMESTAMP - INTERVAL '1 hour', CURRENT_TIMESTAMP - INTERVAL '1 hour'),
('reg_009', 'patient_010', '刘十二', 'sched_016', to_char(CURRENT_DATE, 'YYYY-MM-DD'), 1, 3, CURRENT_TIMESTAMP - INTERVAL '2 hours', CURRENT_TIMESTAMP - INTERVAL '2 hours'),
('reg_010', 'patient_002', '李小红', 'sched_008', to_char(CURRENT_DATE, 'YYYY-MM-DD'), 2, 0, CURRENT_TIMESTAMP - INTERVAL '45 minutes', CURRENT_TIMESTAMP - INTERVAL '45 minutes'),
('reg_011', 'patient_001', '王小明', 'sched_001', to_char(CURRENT_DATE - INTERVAL '2 days', 'YYYY-MM-DD'), 2, 2, CURRENT_TIMESTAMP - INTERVAL '2 days', CURRENT_TIMESTAMP - INTERVAL '2 days');

\c his_clinic;

-- 补充门诊记录
INSERT INTO clinic_records (id, registration_id, patient_id, patient_name, doctor_id, chief_complaint, diagnosis, status, created_at, updated_at) VALUES
('clinic_004', 'reg_007', 'patient_008', '郑十', 'doctor-wang', '右上腹疼痛一周', '胆囊炎', 2, CURRENT_TIMESTAMP - INTERVAL '3 hours', CURRENT_TIMESTAMP - INTERVAL '3 hours'),
('clinic_005', 'reg_011', 'patient_001', '王小明', 'demo-doctor', '头晕乏力', '高血压', 2, CURRENT_TIMESTAMP - INTERVAL '2 days', CURRENT_TIMESTAMP - INTERVAL '2 days'),
('clinic_006', 'reg_009', 'patient_010', '刘十二', 'doctor-zhao', '下腹痛', '盆腔炎', 2, CURRENT_TIMESTAMP - INTERVAL '2 hours', CURRENT_TIMESTAMP - INTERVAL '2 hours');

\c his_prescription;

-- 补充处方
INSERT INTO prescriptions (id, patient_id, patient_name, doctor_id, prescription_type, status, note, created_at, updated_at) VALUES
('pres_demo_006', 'patient_008', '郑十', 'doctor-wang', 1, 2, '胆囊炎用药', CURRENT_TIMESTAMP - INTERVAL '3 hours', CURRENT_TIMESTAMP - INTERVAL '3 hours'),
('pres_demo_007', 'patient_001', '王小明', 'demo-doctor', 1, 2, '降压药', CURRENT_TIMESTAMP - INTERVAL '2 days', CURRENT_TIMESTAMP - INTERVAL '2 days'),
('pres_demo_008', 'patient_010', '刘十二', 'doctor-zhao', 1, 1, '', CURRENT_TIMESTAMP - INTERVAL '2 hours', CURRENT_TIMESTAMP - INTERVAL '2 hours');

INSERT INTO prescription_details (id, prescription_id, drug_id, drug_name, specification, dosage, usage, frequency, days, quantity, unit_price, note) VALUES
('presd_008', 'pres_demo_006', 'drug_004', '奥美拉唑肠溶胶囊', '20mg×14粒', 20, '口服', '每日一次', 14, 1, 38.00, '晨起空腹'),
('presd_009', 'pres_demo_007', 'drug_003', '阿司匹林肠溶片', '100mg×30片', 100, '口服', '每日一次', 30, 1, 15.00, '早餐后'),
('presd_010', 'pres_demo_008', 'drug_001', '阿莫西林胶囊', '0.5g×24粒', 0.5, '口服', '每日三次', 7, 2, 12.00, '');

\c his_billing;

-- 补充账单
INSERT INTO bills (id, patient_id, registration_id, bill_no, total_amount, paid_amount, pay_method, status, created_at, updated_at) VALUES
('bill_006', 'patient_008', 'reg_007', 'BL20260605001', 38.00, 38.00, 4, 1, CURRENT_TIMESTAMP - INTERVAL '3 hours', CURRENT_TIMESTAMP - INTERVAL '3 hours'),
('bill_007', 'patient_001', 'reg_011', 'BL20260603001', 15.00, 15.00, 2, 1, CURRENT_TIMESTAMP - INTERVAL '2 days', CURRENT_TIMESTAMP - INTERVAL '2 days'),
('bill_008', 'patient_010', 'reg_009', 'BL20260605002', 24.00, 0, 0, 0, CURRENT_TIMESTAMP - INTERVAL '2 hours', CURRENT_TIMESTAMP - INTERVAL '2 hours'),
('bill_009', 'patient_002', 'reg_010', 'BL20260605003', 20.00, 20.00, 1, 1, CURRENT_TIMESTAMP - INTERVAL '1 hour', CURRENT_TIMESTAMP - INTERVAL '1 hour');

INSERT INTO bill_details (id, bill_id, item_type, item_name, unit_price, quantity, amount) VALUES
('bld_008', 'bill_006', 1, '奥美拉唑肠溶胶囊', 38.00, 1, 38.00),
('bld_009', 'bill_007', 1, '阿司匹林肠溶片', 15.00, 1, 15.00),
('bld_010', 'bill_008', 1, '阿莫西林胶囊', 12.00, 2, 24.00),
('bld_011', 'bill_009', 2, '挂号费', 20.00, 1, 20.00);

\c his_pharmacy;

-- 补充发药
INSERT INTO dispense_records (id, prescription_id, patient_id, drug_id, quantity, dispenser_id, status, created_at) VALUES
('disp_005', 'pres_demo_006', 'patient_008', 'drug_004', 1, 'demo-nurse', 1, CURRENT_TIMESTAMP - INTERVAL '3 hours'),
('disp_006', 'pres_demo_007', 'patient_001', 'drug_003', 1, 'demo-nurse', 1, CURRENT_TIMESTAMP - INTERVAL '2 days');

-- 补充药品
INSERT INTO drugs (id, name, generic_name, specification, manufacturer, batch_no, purchase_price, retail_price, stock, min_stock, expiry_date, status) VALUES
('drug_006', '头孢克洛胶囊', '头孢克洛', '0.25g×12粒', '广州白云山', 'B20260401', 15.00, 22.00, 300, 30, '2028-04-01', 1),
('drug_007', '硝苯地平控释片', '硝苯地平', '30mg×7片', '拜耳医药', 'B20260320', 28.00, 42.00, 150, 20, '2028-03-20', 1),
('drug_008', '盐酸二甲双胍片', '二甲双胍', '0.5g×20片', '中美上海施贵宝', 'B20260215', 8.00, 12.50, 400, 40, '2028-02-15', 1);

\c his_outpatient;

-- 慢病签约（表名 chronic_contracts）
INSERT INTO chronic_contracts (id, patient_id, doctor_id, disease_type, contract_date, end_date, status) VALUES
('ctr_001', 'patient_001', 'demo-doctor', 'hypertension', to_char(CURRENT_DATE - INTERVAL '30 days', 'YYYY-MM-DD'), to_char(CURRENT_DATE + INTERVAL '335 days', 'YYYY-MM-DD'), 1),
('ctr_002', 'patient_010', 'demo-doctor', 'diabetes', to_char(CURRENT_DATE - INTERVAL '60 days', 'YYYY-MM-DD'), to_char(CURRENT_DATE + INTERVAL '305 days', 'YYYY-MM-DD'), 1),
('ctr_003', 'patient_008', 'demo-doctor', 'hypertension', to_char(CURRENT_DATE - INTERVAL '90 days', 'YYYY-MM-DD'), to_char(CURRENT_DATE + INTERVAL '275 days', 'YYYY-MM-DD'), 1);

-- 健康自测数据（患者端健康档案时间轴用，列名 measure_time）
INSERT INTO health_data (id, patient_id, data_type, value, unit, measure_time) VALUES
('hd_001', 'patient_001', 'blood_pressure', '135/85', 'mmHg', to_char(CURRENT_TIMESTAMP - INTERVAL '30 days', 'YYYY-MM-DD HH24:MI')),
('hd_002', 'patient_001', 'blood_sugar', '5.6', 'mmol/L', to_char(CURRENT_TIMESTAMP - INTERVAL '28 days', 'YYYY-MM-DD HH24:MI')),
('hd_003', 'patient_001', 'blood_pressure', '128/82', 'mmHg', to_char(CURRENT_TIMESTAMP - INTERVAL '14 days', 'YYYY-MM-DD HH24:MI')),
('hd_004', 'patient_001', 'blood_pressure', '142/90', 'mmHg', to_char(CURRENT_TIMESTAMP - INTERVAL '2 days', 'YYYY-MM-DD HH24:MI')),
('hd_005', 'patient_001', 'weight', '72.5', 'kg', to_char(CURRENT_TIMESTAMP - INTERVAL '7 days', 'YYYY-MM-DD HH24:MI')),
('hd_006', 'patient_010', 'blood_sugar', '7.2', 'mmol/L', to_char(CURRENT_TIMESTAMP - INTERVAL '14 days', 'YYYY-MM-DD HH24:MI')),
('hd_007', 'patient_010', 'blood_sugar', '6.8', 'mmol/L', to_char(CURRENT_TIMESTAMP - INTERVAL '7 days', 'YYYY-MM-DD HH24:MI')),
('hd_008', 'patient_008', 'blood_pressure', '150/95', 'mmHg', to_char(CURRENT_TIMESTAMP - INTERVAL '7 days', 'YYYY-MM-DD HH24:MI'));
