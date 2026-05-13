-- 手工执行：新增储蓄流水表。
-- 注意：Codex 不直接执行 DDL/DML，本文件由你在数据库客户端里确认后执行。
-- 执行前请先确认已执行 docs/add_family_saving_table.sql。

USE family_ledger;

CREATE TABLE IF NOT EXISTS `t_family_saving_flow` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `account_id` int unsigned NOT NULL DEFAULT 0 COMMENT '账本ID',
  `transaction_id` int unsigned NOT NULL DEFAULT 0 COMMENT '关联账单ID',
  `scope` varchar(16) NOT NULL DEFAULT 'personal' COMMENT '储蓄范围：family家庭 personal个人',
  `user_openid` varchar(64) NOT NULL DEFAULT '' COMMENT '个人储蓄所属用户openid，家庭储蓄为空',
  `action` varchar(16) NOT NULL DEFAULT '' COMMENT 'create/import/update/delete',
  `delta_amount` decimal(12,2) NOT NULL DEFAULT 0.00 COMMENT '本次储蓄变动金额',
  `balance_after` decimal(12,2) NOT NULL DEFAULT 0.00 COMMENT '变动后储蓄余额',
  `record_type` varchar(16) NOT NULL DEFAULT '' COMMENT 'income/expense',
  `category_name` varchar(32) NOT NULL DEFAULT '' COMMENT '账单分类',
  `description` varchar(80) NOT NULL DEFAULT '' COMMENT '账单说明',
  `is_delete` tinyint NOT NULL DEFAULT 1 COMMENT '1正常 2删除',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_account_created` (`account_id`, `created_at`),
  KEY `idx_transaction` (`transaction_id`),
  KEY `idx_scope_user` (`scope`, `user_openid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='家庭账本储蓄流水表';
