-- ============================================================
-- 演示数据重置脚本（按依赖顺序删除，便于重复导入 seed）
-- 用法: psql -U his_admin -f seed_data_reset.sql
-- ============================================================

\c his_notification;
DELETE FROM notifications WHERE id LIKE 'ntf_%';
DELETE FROM notification_templates WHERE id LIKE 'ntpl_%';

\c his_health_record;
DELETE FROM timeline_events WHERE id LIKE 'tl_%';
DELETE FROM record_authorizations WHERE id LIKE 'rauth_%';
DELETE FROM health_record_summaries WHERE id LIKE 'hrs_%';

\c his_outpatient;
DELETE FROM consultation_messages WHERE id LIKE 'cmsg_%';
DELETE FROM consultations WHERE id LIKE 'cons_%';
DELETE FROM health_data WHERE id LIKE 'hd_%';
DELETE FROM chronic_contracts WHERE id LIKE 'cc_%';

\c his_emr;
DELETE FROM quality_controls WHERE id LIKE 'qc_%';
DELETE FROM medical_records WHERE id LIKE 'emr_%';
DELETE FROM record_templates WHERE id LIKE 'tmpl_%';

\c his_inpatient;
DELETE FROM nursing_records WHERE id LIKE 'nr_%';
DELETE FROM medical_orders WHERE id LIKE 'mo_%';
DELETE FROM inpatient_records WHERE id LIKE 'ip_%';

\c his_followup;
DELETE FROM satisfaction_surveys WHERE id LIKE 'sat_%';
DELETE FROM followup_tasks WHERE id LIKE 'fut_%';
DELETE FROM followup_plans WHERE id LIKE 'fup_%';

\c his_examination;
DELETE FROM examination_reports WHERE id LIKE 'exam_%';

\c his_pharmacy;
DELETE FROM dispense_records WHERE id LIKE 'disp_%';
DELETE FROM drugs WHERE id LIKE 'drug_%';

\c his_billing;
DELETE FROM bill_details WHERE id LIKE 'bld_%';
DELETE FROM bills WHERE id LIKE 'bill_%';

\c his_prescription;
DELETE FROM prescription_details WHERE id LIKE 'presd_%';
DELETE FROM prescriptions WHERE id LIKE 'pres_%';

\c his_clinic;
DELETE FROM examination_requests WHERE id LIKE 'exreq_%';
DELETE FROM clinic_records WHERE id LIKE 'clinic_%';

\c his_registration;
DELETE FROM registrations WHERE id LIKE 'reg_%';
DELETE FROM schedules WHERE id LIKE 'sched_%';

\c his_schedule;
DELETE FROM schedules WHERE id LIKE 'sched_%';

\c his_user;
DELETE FROM employees WHERE id LIKE 'emp_%';
DELETE FROM patients WHERE id LIKE 'patient_%';
DELETE FROM departments WHERE id LIKE 'dept_%';

\c his_system;
DELETE FROM operation_logs WHERE id LIKE 'oplog_%';
DELETE FROM system_params WHERE id LIKE 'param_%';
DELETE FROM dict_items WHERE dict_type IN ('prescription_type', 'pay_method', 'followup_type', 'gender', 'exam_type', 'order_type');
DELETE FROM dict_types WHERE dict_type IN ('prescription_type', 'pay_method', 'followup_type', 'gender', 'exam_type', 'order_type');

\c his_auth;
DELETE FROM role_permissions WHERE role_id LIKE 'role_%';
DELETE FROM permissions WHERE id LIKE 'perm_%';
DELETE FROM users WHERE id LIKE 'demo-%';
DELETE FROM roles WHERE id LIKE 'role_%';
