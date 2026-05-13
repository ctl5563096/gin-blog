CREATE TABLE IF NOT EXISTS `t_family_transaction_category_config` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `account_id` int(11) NOT NULL DEFAULT '0' COMMENT '归属账本',
  `creator_openid` varchar(64) NOT NULL DEFAULT '' COMMENT '创建人',
  `category_name` varchar(32) NOT NULL DEFAULT '' COMMENT '固定分类名称',
  `amount` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '默认金额',
  `type` varchar(16) NOT NULL DEFAULT '' COMMENT 'income/expense',
  `source` varchar(80) NOT NULL DEFAULT '' COMMENT '来源/说明',
  `is_delete` tinyint(1) NOT NULL DEFAULT '1',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_account_creator` (`account_id`, `creator_openid`),
  KEY `idx_account_category` (`account_id`, `category_name`, `type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='账单固定分类配置表';
