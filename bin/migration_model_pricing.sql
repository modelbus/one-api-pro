-- ============================================================
-- 迁移脚本：从 Option 表的 ModelRatio/CompletionRatio/GroupRatio
-- 迁移到新的 model_prices / group_prices 表
--
-- 注意：MySQL 5.7 不支持 JSON_TABLE，此 SQL 脚本仅供参考。
-- 实际迁移请使用 Go 程序：cmd/migrate_pricing/main.go
--   编译：go build -o migrate_pricing ./cmd/migrate_pricing/
--   运行：SQL_DSN="user:pass@tcp(host:port)/dbname" ./migrate_pricing
--
-- 换算公式（旧体系 → 新体系）：
--   ModelRatio 是相对于基准价 $0.002/1K tokens 的倍率
--   InputPrice(¥/百万tokens) = ModelRatio × 0.002 × 1000 × 7 = ModelRatio × 14
--   CompletionRatio 是 output/input 的倍率（默认 1.0 即同价）
--   OutputPrice(¥/百万tokens) = InputPrice × CompletionRatio
--   GroupRatio 是折扣系数（直接迁移为 discount）
-- ============================================================

-- 1. 确保 model_prices 表存在
CREATE TABLE IF NOT EXISTS `model_prices` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `model_name` varchar(100) NOT NULL,
  `input_price` decimal(16,6) NOT NULL DEFAULT '0.000000',
  `output_price` decimal(16,6) NOT NULL DEFAULT '0.000000',
  `cached_price` decimal(16,6) NOT NULL DEFAULT '0.000000',
  `per_request_price` decimal(16,6) NOT NULL DEFAULT '0.000000',
  `billing_type` varchar(20) NOT NULL DEFAULT 'token',
  `enabled` tinyint(1) NOT NULL DEFAULT '1',
  `created_at` bigint NOT NULL DEFAULT '0',
  `updated_at` bigint NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_model_name` (`model_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='模型定价表（¥/百万tokens）';

-- 2. 确保 group_prices 表存在
CREATE TABLE IF NOT EXISTS `group_prices` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `group_name` varchar(32) NOT NULL,
  `model_name` varchar(100) NOT NULL DEFAULT '',
  `discount` decimal(10,4) NOT NULL DEFAULT '1.0000',
  `created_at` bigint NOT NULL DEFAULT '0',
  `updated_at` bigint NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_group_model` (`group_name`, `model_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='分组折扣表';

-- 3. 对于 MySQL 8.0+，可使用以下 JSON_TABLE 方式迁移（取消注释执行）：
-- CALL migrate_model_ratio();  -- 需要先创建存储过程（见下方）

-- ============================================================
-- MySQL 8.0+ 存储过程（仅供参考，MySQL 5.7 不支持 JSON_TABLE）
-- ============================================================
/*
DELIMITER //
DROP PROCEDURE IF EXISTS migrate_model_ratio//
CREATE PROCEDURE migrate_model_ratio()
BEGIN
  DECLARE model_ratio_json LONGTEXT;
  DECLARE completion_ratio_json LONGTEXT;
  DECLARE model_count INT;
  DECLARE now_ts BIGINT;
  SET now_ts = UNIX_TIMESTAMP();

  SELECT COUNT(*) INTO model_count FROM model_prices;
  IF model_count > 0 THEN
    SELECT 'model_prices 表已有数据，跳过 ModelRatio 迁移' AS message;
  ELSE
    SELECT value INTO model_ratio_json FROM options WHERE `key` = 'ModelRatio';
    IF model_ratio_json IS NOT NULL THEN
      SELECT IFNULL((SELECT value FROM options WHERE `key` = 'CompletionRatio'), '{}') INTO completion_ratio_json;
      INSERT INTO model_prices (model_name, input_price, output_price, cached_price, per_request_price, billing_type, enabled, created_at, updated_at)
      SELECT
        jt.model_name,
        ROUND(jt.ratio * 14, 6) AS input_price,
        ROUND(jt.ratio * 14 * IFNULL(c.ratio, 1.0), 6) AS output_price,
        0 AS cached_price,
        0 AS per_request_price,
        'token' AS billing_type,
        1 AS enabled,
        now_ts AS created_at,
        now_ts AS updated_at
      FROM JSON_TABLE(model_ratio_json, '$.*' COLUMNS(
        model_name VARCHAR(100) PATH '$.key',
        ratio DECIMAL(20,10) PATH '$.value'
      )) AS jt
      LEFT JOIN (
        SELECT jt2.model_name, jt2.ratio
        FROM JSON_TABLE(completion_ratio_json, '$.*' COLUMNS(
          model_name VARCHAR(100) PATH '$.key',
          ratio DECIMAL(20,10) PATH '$.value'
        )) AS jt2
      ) c ON jt.model_name = c.model_name
      WHERE jt.ratio >= 0
      ON DUPLICATE KEY UPDATE updated_at = now_ts;
      SELECT CONCAT('已从 ModelRatio 迁移 ', ROW_COUNT(), ' 条模型定价') AS message;
    END IF;
  END IF;
END//
DELIMITER ;

DELIMITER //
DROP PROCEDURE IF EXISTS migrate_group_ratio//
CREATE PROCEDURE migrate_group_ratio()
BEGIN
  DECLARE group_ratio_json LONGTEXT;
  DECLARE group_count INT;
  DECLARE now_ts BIGINT;
  SET now_ts = UNIX_TIMESTAMP();

  SELECT COUNT(*) INTO group_count FROM group_prices;
  IF group_count > 0 THEN
    SELECT 'group_prices 表已有数据，跳过 GroupRatio 迁移' AS message;
  ELSE
    SELECT value INTO group_ratio_json FROM options WHERE `key` = 'GroupRatio';
    IF group_ratio_json IS NOT NULL THEN
      INSERT INTO group_prices (group_name, model_name, discount, created_at, updated_at)
      SELECT jt.group_name, '', jt.discount, now_ts, now_ts
      FROM JSON_TABLE(group_ratio_json, '$.*' COLUMNS(
        group_name VARCHAR(32) PATH '$.key',
        discount DECIMAL(10,4) PATH '$.value'
      )) AS jt
      ON DUPLICATE KEY UPDATE discount = VALUES(discount), updated_at = now_ts;
      SELECT CONCAT('已从 GroupRatio 迁移 ', ROW_COUNT(), ' 条分组折扣') AS message;
    END IF;
  END IF;
END//
DELIMITER ;
*/

-- 4. 迁移后可删除旧 Option 记录（取消注释执行）
-- DELETE FROM options WHERE `key` IN ('ModelRatio', 'CompletionRatio', 'GroupRatio');