-- ============================================================
-- HIS-Go 扩展演示种子数据 — 覆盖尚未填充的业务表
-- 前置: seed_data_reset.sql + seed_data.sql
-- ============================================================

\c his_auth;

-- 权限与角色绑定
INSERT INTO permissions (id, perm_name, perm_code, parent_id, sort_order) VALUES
('perm_admin_all', '系统管理', 'system:admin', NULL, 1),
('perm_patient_read', '患者查看', 'patient:read', NULL, 2),
('perm_patient_write', '患者编辑', 'patient:write', NULL, 3),
('perm_registration', '挂号管理', 'registration:manage', NULL, 4),
('perm_clinic', '门诊诊疗', 'clinic:manage', NULL, 5),
('perm_prescription', '处方管理', 'prescription:manage', NULL, 6),
('perm_pharmacy', '药房发药', 'pharmacy:dispense', NULL, 7),
('perm_billing', '收费结算', 'billing:manage', NULL, 8),
('perm_examination', '检查报告', 'examination:manage', NULL, 9),
('perm_inpatient', '住院管理', 'inpatient:manage', NULL, 10),
('perm_emr', '电子病历', 'emr:manage', NULL, 11),
('perm_followup', '随访管理', 'followup:manage', NULL, 12),
('perm_statistics', '统计报表', 'statistics:view', NULL, 13)
ON CONFLICT (id) DO NOTHING;

INSERT INTO role_permissions (role_id, perm_id) VALUES
('role_admin', 'perm_admin_all'),
('role_admin', 'perm_statistics'),
('role_doctor', 'perm_patient_read'),
('role_doctor', 'perm_clinic'),
('role_doctor', 'perm_prescription'),
('role_doctor', 'perm_examination'),
('role_doctor', 'perm_emr'),
('role_doctor', 'perm_followup'),
('role_nurse', 'perm_patient_read'),
('role_nurse', 'perm_inpatient'),
('role_nurse', 'perm_followup'),
('role_pharmacist', 'perm_pharmacy'),
('role_pharmacist', 'perm_prescription'),
('role_patient', 'perm_patient_read')
ON CONFLICT DO NOTHING;

-- 补充演示药师账号（密码 demo123）
INSERT INTO users (id, username, password, real_name, phone, role, dept_id) VALUES
('demo-pharmacist', 'demo-pharmacist', '$2a$10$avEr2y6CrrENS8/NMWeeNOJcA2S76iJOzdkZzLLnrvbmor6fLiQVW', '王药师', '13800000004', 'pharmacist', 'dept_001')
ON CONFLICT (id) DO NOTHING;

\c his_system;

INSERT INTO dict_types (dict_name, dict_type, remark) VALUES
('性别', 'gender', '性别字典'),
('检查类型', 'exam_type', '检查类型字典'),
('医嘱类型', 'order_type', '住院医嘱类型')
ON CONFLICT (dict_type) DO NOTHING;

INSERT INTO dict_items (dict_type, label, value, sort_order) VALUES
('gender', '男', 'M', 1),
('gender', '女', 'F', 2),
('exam_type', 'CT', 'CT', 1),
('exam_type', 'DR', 'DR', 2),
('exam_type', '超声', 'US', 3),
('exam_type', '检验', 'LAB', 4),
('order_type', '长期医嘱', '1', 1),
('order_type', '临时医嘱', '2', 2),
('order_type', '护理医嘱', '3', 3)
ON CONFLICT DO NOTHING;

INSERT INTO system_params (id, param_name, param_key, param_value, remark) VALUES
('param_001', '医院名称', 'hospital.name', 'HIS-Go 演示医院', '演示环境医院名称'),
('param_002', '挂号费默认', 'registration.default_fee', '15.00', '默认挂号费（元）'),
('param_003', '就诊助手开关', 'triage.enabled', 'true', '就诊助手功能开关'),
('param_004', '随访提醒天数', 'followup.remind_days', '3', '随访任务提前提醒天数')
ON CONFLICT (id) DO NOTHING;

INSERT INTO operation_logs (id, user_id, username, module, action, method, url, ip, params, result, status, created_at) VALUES
('oplog_001', 'demo-admin', 'demo-admin', '系统管理', '登录', 'POST', '/api/auth/login', '127.0.0.1', '{"username":"demo-admin"}', '成功', 1, CURRENT_TIMESTAMP - INTERVAL '2 hours'),
('oplog_002', 'demo-doctor', 'demo-doctor', '门诊', '创建处方', 'POST', '/api/prescription', '192.168.1.10', '{"patientId":"patient_005"}', '成功', 1, CURRENT_TIMESTAMP - INTERVAL '1 hour'),
('oplog_003', 'demo-patient', 'demo-patient', '患者端', '预约挂号', 'POST', '/api/registration/register', '10.0.0.5', '{"scheduleId":"sched_001"}', '成功', 1, CURRENT_TIMESTAMP - INTERVAL '30 minutes')
ON CONFLICT (id) DO NOTHING;

\c his_clinic;

INSERT INTO examination_requests (id, clinic_record_id, patient_id, exam_type, exam_item, body_part, clinical_diagnosis, note, status, created_at) VALUES
('exreq_001', 'clinic_001', 'patient_001', 'CT', '胸部CT平扫', '胸部', '急性支气管炎', '排除肺炎', 2, CURRENT_TIMESTAMP - INTERVAL '1 day'),
('exreq_002', 'clinic_002', 'patient_005', 'DR', '胸部正位片', '胸部', '上呼吸道感染', '发热查因', 2, CURRENT_TIMESTAMP - INTERVAL '2 hours'),
('exreq_003', 'clinic_003', 'patient_007', 'DR', '右下肢正侧位片', '右下肢', '右下肢软组织挫伤', '排除骨折', 2, CURRENT_TIMESTAMP - INTERVAL '3 hours'),
('exreq_004', 'clinic_001', 'patient_001', '检验', '血常规', '血液', '急性支气管炎', '复查', 1, CURRENT_TIMESTAMP - INTERVAL '12 hours')
ON CONFLICT (id) DO NOTHING;

\c his_emr;

INSERT INTO record_templates (id, name, dept_id, content, type, created_at) VALUES
('tmpl_001', '内科门诊病历模板', 'dept_001', '主诉|现病史|既往史|体格检查|辅助检查|诊断|处理意见', 1, CURRENT_TIMESTAMP - INTERVAL '90 days'),
('tmpl_002', '外科门诊病历模板', 'dept_002', '主诉|现病史|外伤史|体格检查|影像检查|诊断|处理意见', 1, CURRENT_TIMESTAMP - INTERVAL '90 days'),
('tmpl_003', '急诊病历模板', 'dept_005', '主诉|现病史|生命体征|初步诊断|处置', 2, CURRENT_TIMESTAMP - INTERVAL '90 days')
ON CONFLICT (id) DO NOTHING;

INSERT INTO medical_records (id, patient_id, clinic_record_id, template_id, chief_complaint, present_illness, past_history, physical_exam, auxiliary_exam, diagnosis, treatment_plan, quality_level, status, created_at, updated_at) VALUES
('emr_001', 'patient_001', 'clinic_001', 'tmpl_001', '咳嗽三天，加重一天', '受凉后出现咳嗽，咳少量白痰，无发热', '体健', '咽充血，双肺呼吸音粗', '血常规正常', '急性支气管炎', '抗感染、止咳化痰，多饮水', 2, 2, CURRENT_TIMESTAMP - INTERVAL '1 day', CURRENT_TIMESTAMP - INTERVAL '1 day'),
('emr_002', 'patient_005', 'clinic_002', 'tmpl_001', '发热两天，体温38.5℃', '畏寒发热，伴咽痛', '无特殊', 'T 38.5℃，咽红', '胸片示右下肺炎症', '上呼吸道感染/肺炎', '抗感染、退热，3天后复查', 2, 2, CURRENT_TIMESTAMP - INTERVAL '2 hours', CURRENT_TIMESTAMP - INTERVAL '2 hours'),
('emr_003', 'patient_007', 'clinic_003', 'tmpl_002', '外伤后右下肢疼痛', '跌倒后右小腿肿痛，活动受限', '无', '右小腿肿胀压痛，活动受限', 'X线排除骨折', '右下肢软组织挫伤', '冷敷、止痛，休息', 1, 2, CURRENT_TIMESTAMP - INTERVAL '3 hours', CURRENT_TIMESTAMP - INTERVAL '3 hours')
ON CONFLICT (id) DO NOTHING;

INSERT INTO quality_controls (id, record_id, reviewer_id, level, comment, created_at) VALUES
('qc_001', 'emr_001', 'demo-doctor', 2, '病历书写规范，诊断合理', CURRENT_TIMESTAMP - INTERVAL '20 hours'),
('qc_002', 'emr_002', 'demo-doctor', 2, '辅助检查完善，符合诊疗规范', CURRENT_TIMESTAMP - INTERVAL '1 hour'),
('qc_003', 'emr_003', 'demo-doctor', 1, '外伤史描述可更详细', CURRENT_TIMESTAMP - INTERVAL '2 hours')
ON CONFLICT (id) DO NOTHING;

\c his_inpatient;

INSERT INTO inpatient_records (id, patient_id, patient_name, admission_date, discharge_date, dept_id, room_no, bed_no, diagnosis, deposit, total_cost, status, created_at, updated_at) VALUES
('ip_001', 'patient_002', '李小红', CURRENT_TIMESTAMP - INTERVAL '5 days', CURRENT_TIMESTAMP - INTERVAL '2 days', 'dept_001', '301', '01', '社区获得性肺炎', 5000.00, 3860.50, 2, CURRENT_TIMESTAMP - INTERVAL '5 days', CURRENT_TIMESTAMP - INTERVAL '2 days'),
('ip_002', 'patient_003', '张三', CURRENT_TIMESTAMP - INTERVAL '1 day', NULL, 'dept_002', '402', '03', '急性阑尾炎术后', 8000.00, 2150.00, 1, CURRENT_TIMESTAMP - INTERVAL '1 day', CURRENT_TIMESTAMP)
ON CONFLICT (id) DO NOTHING;

INSERT INTO medical_orders (id, inpatient_id, doctor_id, order_type, content, start_time, end_time, status, created_at) VALUES
('mo_001', 'ip_001', 'demo-doctor', 1, '头孢曲松 2g ivgtt qd', CURRENT_TIMESTAMP - INTERVAL '5 days', CURRENT_TIMESTAMP - INTERVAL '2 days', 2, CURRENT_TIMESTAMP - INTERVAL '5 days'),
('mo_002', 'ip_001', 'demo-doctor', 2, '血常规 qd', CURRENT_TIMESTAMP - INTERVAL '4 days', CURRENT_TIMESTAMP - INTERVAL '3 days', 2, CURRENT_TIMESTAMP - INTERVAL '4 days'),
('mo_003', 'ip_002', 'demo-doctor', 1, '禁食水，补液', CURRENT_TIMESTAMP - INTERVAL '1 day', NULL, 1, CURRENT_TIMESTAMP - INTERVAL '1 day'),
('mo_004', 'ip_002', 'demo-doctor', 3, '心电监护', CURRENT_TIMESTAMP - INTERVAL '1 day', NULL, 1, CURRENT_TIMESTAMP - INTERVAL '1 day')
ON CONFLICT (id) DO NOTHING;

INSERT INTO nursing_records (id, inpatient_id, nurse_id, record_time, content, vital_signs, created_at) VALUES
('nr_001', 'ip_001', 'demo-nurse', CURRENT_TIMESTAMP - INTERVAL '4 days', '患者精神可，咳嗽减轻', 'T 36.8 P 78 R 18 BP 120/80', CURRENT_TIMESTAMP - INTERVAL '4 days'),
('nr_002', 'ip_001', 'demo-nurse', CURRENT_TIMESTAMP - INTERVAL '3 days', '夜间睡眠好，无发热', 'T 36.5 P 76 R 18 BP 118/78', CURRENT_TIMESTAMP - INTERVAL '3 days'),
('nr_003', 'ip_002', 'demo-nurse', CURRENT_TIMESTAMP - INTERVAL '12 hours', '术后6小时，切口敷料干燥', 'T 37.0 P 82 R 20 BP 125/82', CURRENT_TIMESTAMP - INTERVAL '12 hours')
ON CONFLICT (id) DO NOTHING;

\c his_outpatient;

INSERT INTO consultations (id, patient_id, doctor_id, type, description, status, created_at, updated_at) VALUES
('cons_001', 'patient_001', 'demo-doctor', 1, '咳嗽三天，是否需要来医院？', 2, CURRENT_TIMESTAMP - INTERVAL '2 days', CURRENT_TIMESTAMP - INTERVAL '2 days'),
('cons_002', 'patient_001', 'demo-doctor', 1, '头晕恶心，挂什么科？', 1, CURRENT_TIMESTAMP - INTERVAL '3 hours', CURRENT_TIMESTAMP - INTERVAL '1 hour')
ON CONFLICT (id) DO NOTHING;

INSERT INTO consultation_messages (id, consultation_id, sender_id, sender_name, content, msg_type, created_at) VALUES
('cmsg_001', 'cons_001', 'patient_001', '王小明', '医生您好，我咳嗽三天了，需要来医院吗？', 'text', CURRENT_TIMESTAMP - INTERVAL '2 days'),
('cmsg_002', 'cons_001', 'demo-doctor', '张医生', '建议来院呼吸内科就诊，如发热加重请及时就医。', 'text', CURRENT_TIMESTAMP - INTERVAL '2 days' + INTERVAL '30 minutes'),
('cmsg_003', 'cons_002', 'patient_001', '王小明', '最近头晕恶心，应该挂什么科？', 'text', CURRENT_TIMESTAMP - INTERVAL '3 hours'),
('cmsg_004', 'cons_002', 'demo-doctor', '张医生', '建议先挂神经内科或内科，如症状加重请急诊。', 'text', CURRENT_TIMESTAMP - INTERVAL '1 hour')
ON CONFLICT (id) DO NOTHING;

INSERT INTO chronic_contracts (id, patient_id, doctor_id, disease_type, contract_date, end_date, status, created_at) VALUES
('cc_001', 'patient_001', 'demo-doctor', '高血压', to_char(CURRENT_DATE - INTERVAL '180 days', 'YYYY-MM-DD'), to_char(CURRENT_DATE + INTERVAL '185 days', 'YYYY-MM-DD'), 1, CURRENT_TIMESTAMP - INTERVAL '180 days'),
('cc_002', 'patient_001', 'demo-doctor', '2型糖尿病', to_char(CURRENT_DATE - INTERVAL '90 days', 'YYYY-MM-DD'), to_char(CURRENT_DATE + INTERVAL '275 days', 'YYYY-MM-DD'), 1, CURRENT_TIMESTAMP - INTERVAL '90 days')
ON CONFLICT (id) DO NOTHING;

INSERT INTO health_data (id, patient_id, data_type, value, unit, measure_time, abnormal, created_at) VALUES
('hd_001', 'patient_001', 'blood_pressure', '135/85', 'mmHg', to_char(CURRENT_DATE - INTERVAL '1 day', 'YYYY-MM-DD HH24:MI'), false, CURRENT_TIMESTAMP - INTERVAL '1 day'),
('hd_002', 'patient_001', 'blood_glucose', '7.2', 'mmol/L', to_char(CURRENT_DATE, 'YYYY-MM-DD 08:00'), true, CURRENT_TIMESTAMP - INTERVAL '6 hours'),
('hd_003', 'patient_001', 'heart_rate', '78', 'bpm', to_char(CURRENT_DATE, 'YYYY-MM-DD 08:00'), false, CURRENT_TIMESTAMP - INTERVAL '6 hours'),
('hd_004', 'patient_001', 'weight', '72.5', 'kg', to_char(CURRENT_DATE - INTERVAL '7 days', 'YYYY-MM-DD'), false, CURRENT_TIMESTAMP - INTERVAL '7 days')
ON CONFLICT (id) DO NOTHING;

\c his_health_record;

INSERT INTO health_record_summaries (id, patient_id, patient_name, total_visits, total_prescriptions, total_examinations, updated_at) VALUES
('hrs_001', 'patient_001', '王小明', 5, 2, 2, CURRENT_TIMESTAMP),
('hrs_002', 'patient_002', '李小红', 3, 0, 1, CURRENT_TIMESTAMP - INTERVAL '2 days'),
('hrs_003', 'patient_005', '孙七', 2, 1, 1, CURRENT_TIMESTAMP - INTERVAL '2 hours')
ON CONFLICT (id) DO NOTHING;

INSERT INTO timeline_events (id, patient_id, date, event_type, description, related_id, created_at) VALUES
('tl_001', 'patient_001', to_char(CURRENT_DATE - INTERVAL '7 days', 'YYYY-MM-DD'), 'examination', '完成胸部CT平扫', 'exam_demo_001', CURRENT_TIMESTAMP - INTERVAL '7 days'),
('tl_002', 'patient_001', to_char(CURRENT_DATE - INTERVAL '3 days', 'YYYY-MM-DD'), 'prescription', '开具感冒用药处方', 'pres_demo_001', CURRENT_TIMESTAMP - INTERVAL '3 days'),
('tl_003', 'patient_001', to_char(CURRENT_DATE - INTERVAL '1 day', 'YYYY-MM-DD'), 'visit', '内科就诊：急性支气管炎', 'clinic_001', CURRENT_TIMESTAMP - INTERVAL '1 day'),
('tl_004', 'patient_001', to_char(CURRENT_DATE - INTERVAL '2 days', 'YYYY-MM-DD'), 'examination', '血常规检验', 'exam_demo_002', CURRENT_TIMESTAMP - INTERVAL '2 days'),
('tl_005', 'patient_005', to_char(CURRENT_DATE, 'YYYY-MM-DD'), 'visit', '儿科就诊：上呼吸道感染', 'clinic_002', CURRENT_TIMESTAMP - INTERVAL '2 hours')
ON CONFLICT (id) DO NOTHING;

INSERT INTO record_authorizations (id, patient_id, doctor_id, auth_time, expire_time, status, created_at) VALUES
('rauth_001', 'patient_001', 'demo-doctor', to_char(CURRENT_DATE - INTERVAL '30 days', 'YYYY-MM-DD HH24:MI'), to_char(CURRENT_DATE + INTERVAL '335 days', 'YYYY-MM-DD HH24:MI'), 1, CURRENT_TIMESTAMP - INTERVAL '30 days')
ON CONFLICT (id) DO NOTHING;

\c his_notification;

INSERT INTO notification_templates (id, name, title_template, content_template, channel, params, created_at) VALUES
('ntpl_001', '挂号成功通知', '挂号成功提醒', '您已成功预约{{deptName}} {{doctorName}} {{date}} {{timeSlot}}，请按时到院。', 1, 'deptName,doctorName,date,timeSlot', CURRENT_TIMESTAMP - INTERVAL '30 days'),
('ntpl_002', '处方审核通知', '处方状态更新', '您的处方（{{prescriptionId}}）已{{status}}。', 1, 'prescriptionId,status', CURRENT_TIMESTAMP - INTERVAL '30 days'),
('ntpl_003', '随访提醒', '随访任务提醒', '您有一项随访任务「{{taskName}}」将于{{executeDate}}执行。', 2, 'taskName,executeDate', CURRENT_TIMESTAMP - INTERVAL '30 days')
ON CONFLICT (id) DO NOTHING;

INSERT INTO notifications (id, template_id, receiver_id, title, content, channel, status, send_time, created_at) VALUES
('ntf_001', 'ntpl_001', 'demo-patient', '挂号成功提醒', '您已成功预约内科 张医生 ' || to_char(CURRENT_DATE - INTERVAL '1 day', 'YYYY-MM-DD') || ' 上午，请按时到院。', 1, 1, CURRENT_TIMESTAMP - INTERVAL '1 day', CURRENT_TIMESTAMP - INTERVAL '1 day'),
('ntf_002', 'ntpl_002', 'demo-patient', '处方状态更新', '您的处方（pres_demo_001）已发药。', 1, 1, CURRENT_TIMESTAMP - INTERVAL '1 day', CURRENT_TIMESTAMP - INTERVAL '1 day'),
('ntf_003', 'ntpl_003', 'demo-patient', '随访任务提醒', '您有一项随访任务「门诊复查」将于' || to_char(CURRENT_DATE + INTERVAL '30 days', 'YYYY-MM-DD') || '执行。', 2, 1, CURRENT_TIMESTAMP - INTERVAL '1 hour', CURRENT_TIMESTAMP - INTERVAL '1 hour'),
('ntf_004', 'ntpl_001', 'patient_003', '挂号成功提醒', '您已成功预约内科，请按时到院。', 1, 1, CURRENT_TIMESTAMP - INTERVAL '2 hours', CURRENT_TIMESTAMP - INTERVAL '2 hours')
ON CONFLICT (id) DO NOTHING;

\c his_followup;

INSERT INTO satisfaction_surveys (id, followup_task_id, patient_id, score, feedback, created_at) VALUES
('sat_001', 'fut_002', 'patient_001', 5, '护士随访很耐心，用药指导清楚', CURRENT_TIMESTAMP - INTERVAL '7 days'),
('sat_002', 'fut_006', 'patient_002', 4, '问卷填写方便，希望增加短信提醒', CURRENT_TIMESTAMP - INTERVAL '30 days')
ON CONFLICT (id) DO NOTHING;

\c his_user;

INSERT INTO employees (id, user_id, name, phone, dept_id, title, specialty, status) VALUES
('emp_003', 'demo-pharmacist', '王药师', '13800000004', 'dept_001', '主管药师', '临床药学', 1)
ON CONFLICT (id) DO NOTHING;
