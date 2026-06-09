-- ============================================================
-- 管理端演示数据同步脚本（幂等，可重复执行）
-- 用法: docker exec his-postgres-demo psql -U his_admin -f /sql-init/seed_admin_sync.sql
-- ============================================================

\c his_user;

INSERT INTO employees (id, user_id, name, phone, dept_id, title, status) VALUES
('emp_001', 'demo-doctor', '张医生', '13800000001', 'dept_001', '主任医师', 1),
('emp_002', 'demo-nurse', '李护士', '13800000002', 'dept_001', '主管护师', 1),
('emp_003', 'doctor-wang', '王医生', '13800000011', 'dept_002', '副主任医师', 1),
('emp_004', 'doctor-li', '李医生', '13800000012', 'dept_003', '主治医师', 1),
('emp_005', 'doctor-zhao', '赵医生', '13800000013', 'dept_004', '主任医师', 1),
('emp_006', 'doctor-chen', '陈医生', '13800000014', 'dept_005', '副主任医师', 1),
('emp_007', 'demo-pharmacist', '王药师', '13800000004', 'dept_001', '主管药师', 1)
ON CONFLICT (id) DO UPDATE SET
  name = EXCLUDED.name, phone = EXCLUDED.phone, dept_id = EXCLUDED.dept_id,
  title = EXCLUDED.title, status = EXCLUDED.status;

\c his_schedule;

-- 唯一约束: (doctor_id, work_date, time_slot) — 每医生每天每时段仅一条
INSERT INTO schedules (id, doctor_id, doctor_name, dept_id, dept_name, work_date, time_slot, max_patients, current_patients, room_no, status) VALUES
('sched_001', 'demo-doctor', '张医生', 'dept_001', '内科', to_char(CURRENT_DATE, 'YYYY-MM-DD'), 1, 30, 2, '201', 1),
('sched_002', 'demo-doctor', '张医生', 'dept_001', '内科', to_char(CURRENT_DATE, 'YYYY-MM-DD'), 2, 30, 0, '201', 1),
('sched_003', 'demo-doctor', '张医生', 'dept_001', '内科', to_char(CURRENT_DATE + INTERVAL '1 day', 'YYYY-MM-DD'), 1, 30, 0, '201', 1),
('sched_004', 'demo-doctor', '张医生', 'dept_001', '内科', to_char(CURRENT_DATE + INTERVAL '1 day', 'YYYY-MM-DD'), 2, 30, 0, '201', 1),
('sched_005', 'doctor-wang', '王医生', 'dept_002', '外科', to_char(CURRENT_DATE, 'YYYY-MM-DD'), 1, 20, 1, '302', 1),
('sched_006', 'doctor-li', '李医生', 'dept_003', '儿科', to_char(CURRENT_DATE, 'YYYY-MM-DD'), 1, 25, 2, '103', 1),
('sched_007', 'doctor-li', '李医生', 'dept_003', '儿科', to_char(CURRENT_DATE, 'YYYY-MM-DD'), 2, 25, 0, '103', 1),
('sched_008', 'doctor-zhao', '赵医生', 'dept_004', '妇产科', to_char(CURRENT_DATE, 'YYYY-MM-DD'), 1, 20, 2, '401', 1),
('sched_009', 'doctor-chen', '陈医生', 'dept_005', '急诊科', to_char(CURRENT_DATE, 'YYYY-MM-DD'), 1, 50, 5, '急诊1', 1),
('sched_010', 'demo-doctor', '张医生', 'dept_001', '内科', to_char(CURRENT_DATE + INTERVAL '2 days', 'YYYY-MM-DD'), 1, 30, 0, '201', 1),
('sched_011', 'doctor-wang', '王医生', 'dept_002', '外科', to_char(CURRENT_DATE, 'YYYY-MM-DD'), 2, 20, 0, '302', 1),
('sched_012', 'doctor-zhao', '赵医生', 'dept_004', '妇产科', to_char(CURRENT_DATE, 'YYYY-MM-DD'), 2, 20, 0, '401', 1),
('sched_013', 'doctor-chen', '陈医生', 'dept_005', '急诊科', to_char(CURRENT_DATE, 'YYYY-MM-DD'), 2, 50, 0, '急诊2', 1),
('sched_014', 'demo-doctor', '张医生', 'dept_001', '内科', to_char(CURRENT_DATE + INTERVAL '3 days', 'YYYY-MM-DD'), 1, 30, 0, '201', 1),
('sched_015', 'doctor-wang', '王医生', 'dept_002', '外科', to_char(CURRENT_DATE + INTERVAL '1 day', 'YYYY-MM-DD'), 1, 20, 0, '302', 1),
('sched_016', 'doctor-li', '李医生', 'dept_003', '儿科', to_char(CURRENT_DATE + INTERVAL '1 day', 'YYYY-MM-DD'), 1, 25, 0, '103', 1),
('sched_017', 'demo-doctor', '张医生', 'dept_001', '内科', to_char(CURRENT_DATE + INTERVAL '4 days', 'YYYY-MM-DD'), 1, 30, 0, '201', 1),
('sched_018', 'demo-doctor', '张医生', 'dept_001', '内科', to_char(CURRENT_DATE + INTERVAL '5 days', 'YYYY-MM-DD'), 2, 30, 0, '201', 1),
('sched_019', 'doctor-wang', '王医生', 'dept_002', '外科', to_char(CURRENT_DATE + INTERVAL '2 days', 'YYYY-MM-DD'), 1, 20, 0, '302', 1),
('sched_020', 'doctor-chen', '陈医生', 'dept_005', '急诊科', to_char(CURRENT_DATE + INTERVAL '1 day', 'YYYY-MM-DD'), 1, 50, 0, '急诊1', 1)
ON CONFLICT (doctor_id, work_date, time_slot) DO UPDATE SET
  id = EXCLUDED.id, dept_id = EXCLUDED.dept_id, dept_name = EXCLUDED.dept_name,
  doctor_name = EXCLUDED.doctor_name, max_patients = EXCLUDED.max_patients,
  current_patients = EXCLUDED.current_patients, room_no = EXCLUDED.room_no, status = EXCLUDED.status;
