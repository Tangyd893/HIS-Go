-- Migration 001: 为 patients 表添加 user_id 列，关联 his_auth.users
-- 用途：患者端按当前登录用户查找对应患者档案（UX-6）
-- 执行时机：已部署环境需手动执行；新环境 init_all.sql 已包含

ALTER TABLE patients ADD COLUMN IF NOT EXISTS user_id VARCHAR(64);

-- 索引（可选，加速按 user_id 查询）
CREATE INDEX IF NOT EXISTS idx_patients_user_id ON patients(user_id);
