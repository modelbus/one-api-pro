-- migration_plan_virtual_amount_billing_type.sql
-- 套餐虚拟余额（virtual_amount）+ 计费维度（billing_type）迁移脚本
-- Plan virtual amount (virtual_amount) + billing type (billing_type) migration
--
-- 版本: v0.0.25
-- 日期: 2026-10-05
-- 作者: opencode
--
-- 说明：
--   1. 程序启动时 GORM AutoMigrate 会自动补齐这些列，本脚本仅用于
--      「不方便重启 / 需要先手动改库」的场景，幂等可重复执行。
--   2. plans.virtual_amount 与 user_plans.* 均为「微元」整数，
--      1 元 = 1_000_000 微元（config.QuotaPerUnit）。
--   3. plans.model_limits 必须是对象 JSON，例如：
--      {"gpt-4o": {"period_h": 5, "token_period": 50000}}
--      历史错误写法 {"gpt-4o": 1000} 会导致套餐不可用（不再静默免费）。
--   4. virtual_amount = 0 表示「不限额度」，兼容存量套餐。
--   5. 存量 user_plans 的快照列留空 / 0：运行时以 plans 行实时配置为准，
--      快照仅在 plans 行被删除时兜底。

-- ---------------------------------------------------------------------------
-- plans：套餐默认计费维度 + 虚拟余额总池
-- ---------------------------------------------------------------------------
SET @db := DATABASE();

SET @sql := IF(
  (SELECT COUNT(*) FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = @db AND TABLE_NAME = 'plans' AND COLUMN_NAME = 'billing_type') = 0,
  'ALTER TABLE `plans` ADD COLUMN `billing_type` VARCHAR(20) NOT NULL DEFAULT ''token''',
  'SELECT ''plans.billing_type exists'''
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := IF(
  (SELECT COUNT(*) FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = @db AND TABLE_NAME = 'plans' AND COLUMN_NAME = 'virtual_amount') = 0,
  'ALTER TABLE `plans` ADD COLUMN `virtual_amount` BIGINT NOT NULL DEFAULT 0',
  'SELECT ''plans.virtual_amount exists'''
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- ---------------------------------------------------------------------------
-- user_plans：购买快照（model_limits / virtual_amount）+ 累计已用额度
-- ---------------------------------------------------------------------------
SET @sql := IF(
  (SELECT COUNT(*) FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = @db AND TABLE_NAME = 'user_plans' AND COLUMN_NAME = 'model_limits') = 0,
  'ALTER TABLE `user_plans` ADD COLUMN `model_limits` TEXT NULL',
  'SELECT ''user_plans.model_limits exists'''
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := IF(
  (SELECT COUNT(*) FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = @db AND TABLE_NAME = 'user_plans' AND COLUMN_NAME = 'virtual_amount') = 0,
  'ALTER TABLE `user_plans` ADD COLUMN `virtual_amount` BIGINT NOT NULL DEFAULT 0',
  'SELECT ''user_plans.virtual_amount exists'''
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := IF(
  (SELECT COUNT(*) FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = @db AND TABLE_NAME = 'user_plans' AND COLUMN_NAME = 'used_amount') = 0,
  'ALTER TABLE `user_plans` ADD COLUMN `used_amount` BIGINT NOT NULL DEFAULT 0',
  'SELECT ''user_plans.used_amount exists'''
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- ---------------------------------------------------------------------------
-- 存量数据回填：把当前 plans 配置写入活跃订阅快照，保证套餐被删后仍可计费
-- ---------------------------------------------------------------------------
UPDATE `user_plans` up
JOIN `plans` p ON p.`id` = up.`plan_id`
SET up.`model_limits`   = p.`model_limits`,
    up.`virtual_amount` = p.`virtual_amount`
WHERE up.`status` = 1
  AND (up.`model_limits` IS NULL OR up.`model_limits` = '');

-- 存量套餐若缺少 billing_type，回填为 token（与历史硬编码行为保持一致）
UPDATE `plans` SET `billing_type` = 'token'
WHERE `billing_type` IS NULL OR `billing_type` = '';

-- 存量套餐若 model_limits 仍是数字值写法（历史错误格式），提醒人工修正
SELECT `id`, `name`, `model_limits`
FROM `plans`
WHERE `model_limits` IS NOT NULL
  AND `model_limits` <> ''
  AND JSON_VALID(`model_limits`) = 1
  AND JSON_TYPE(JSON_EXTRACT(`model_limits`, '$.*')) <> 'OBJECT';
