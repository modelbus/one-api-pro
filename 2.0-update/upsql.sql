-- Create syntax for TABLE 'channels'
CREATE TABLE `channels` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `type` int(11) DEFAULT 0,
  `key` text,
  `status` int(11) DEFAULT 1,
  `name` varchar(255) DEFAULT '',
  `weight` int(11) DEFAULT 0,
  `created_time` bigint DEFAULT 0,
  `test_time` bigint DEFAULT 0,
  `response_time` int(11) DEFAULT 0,
  `base_url` varchar(255) DEFAULT '',
  `other` varchar(255) DEFAULT NULL,
  `balance` double DEFAULT 0,
  `balance_updated_time` bigint DEFAULT 0,
  `models` text,
  `group` varchar(32) DEFAULT 'default',
  `used_quota` bigint DEFAULT 0,
  `model_mapping` varchar(1024) DEFAULT '',
  `priority` bigint DEFAULT 0,
  `config` text,
  `system_prompt` text,
  `max_concurrency` int(11) DEFAULT 0,
  `cooldown_seconds` int(11) DEFAULT 60,
  `rpm` int(11) DEFAULT 0,
  `last_error` varchar(512) DEFAULT '',
  `last_error_time` bigint DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `idx_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Create syntax for TABLE 'logs' (with session_key for sticky sessions)
CREATE TABLE `logs` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) DEFAULT 0,
  `created_at` bigint DEFAULT 0,
  `type` int(11) DEFAULT 0,
  `content` text,
  `username` varchar(255) DEFAULT '',
  `token_name` varchar(255) DEFAULT '',
  `model_name` varchar(255) DEFAULT '',
  `quota` int(11) DEFAULT 0,
  `prompt_tokens` int(11) DEFAULT 0,
  `completion_tokens` int(11) DEFAULT 0,
  `cached_tokens` int(11) DEFAULT 0,
  `channel_id` int(11) DEFAULT 0,
  `request_id` varchar(255) DEFAULT '',
  `elapsed_time` bigint DEFAULT 0,
  `is_stream` tinyint(1) DEFAULT 0,
  `system_prompt_reset` tinyint(1) DEFAULT 0,
  `billing_source` int(11) DEFAULT 0,
  `plan_id` int(11) DEFAULT 0,
  `session_key` varchar(128) DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_created_at_type` (`created_at`, `type`),
  KEY `idx_username_model_name` (`username`, `model_name`),
  KEY `idx_token_name` (`token_name`),
  KEY `idx_model_name` (`model_name`),
  KEY `idx_channel_id` (`channel_id`),
  KEY `idx_logs_session_key` (`session_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Create syntax for TABLE 'model'
CREATE TABLE `model` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '模型ID',
  `name` varchar(100) NOT NULL COMMENT '模型名称',
  `model_key` varchar(200) NOT NULL DEFAULT '' COMMENT '模型标识键',
  `icon` varchar(500) DEFAULT NULL COMMENT '图标URL',
  `description` text COMMENT '描述',
  `context_window` int(11) DEFAULT '0' COMMENT '上下文窗口大小',
  `features` text COMMENT '特性(JSON)',
  `sort` int(11) DEFAULT '0' COMMENT '排序',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态: 0-禁用 1-正常',
  `create_time` int(11) DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_provider_id` (`provider_id`),
  KEY `idx_sort` (`sort`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='模型配置表';

-- Create syntax for TABLE 'plan'
CREATE TABLE `plan` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '套餐ID',
  `name` varchar(100) NOT NULL COMMENT '套餐名称',
  `price` decimal(10,2) DEFAULT '0.00' COMMENT '价格',
  `tokens` bigint(20) DEFAULT '0' COMMENT 'Token配额',
  `model_limit` text COMMENT '允许使用的模型ID列表(JSON)',
  `description` text COMMENT '套餐描述',
  `features` text COMMENT '功能特性(JSON)',
  `sort` int(11) DEFAULT '0' COMMENT '排序',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态: 0-下架 1-上架',
  `duration_days` int(11) NOT NULL DEFAULT '30' COMMENT '有效天数',
  `duration_text` varchar(50) DEFAULT NULL COMMENT '有效期显示文本',
  `recommended` tinyint(1) DEFAULT '0' COMMENT '是否推荐：0=否，1=是',
  `create_time` int(11) DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_sort` (`sort`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='套餐表';

-- Create syntax for TABLE 'usage'
CREATE TABLE `usage` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `model_id` int(11) NOT NULL,
  `tokens` int(11) DEFAULT '0',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Create syntax for TABLE 'user_plan'
CREATE TABLE `user_plan` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '记录ID',
  `user_id` int(11) NOT NULL COMMENT '用户ID（关联 user）',
  `plan_id` int(11) NOT NULL COMMENT '套餐ID（关联 plan）',
  `order_id` int(11) NOT NULL COMMENT '订单ID（关联 plan_order）',
  `start_time` datetime NOT NULL COMMENT '套餐开始时间',
  `end_time` datetime NOT NULL COMMENT '套餐结束时间',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态：1=生效中，0=已失效',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_user_status` (`user_id`,`status`),
  KEY `idx_end_time` (`end_time`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='用户-套餐关系表';

-- 模型定价表（替代旧的 ModelRatio/CompletionRatio Option）
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

-- 分组折扣表（替代旧的 GroupRatio Option）
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