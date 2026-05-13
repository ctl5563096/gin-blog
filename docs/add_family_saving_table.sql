-- 手工执行：新增储蓄状态表。
-- 注意：Codex 不直接执行 DDL/DML，本文件由你在数据库客户端里确认后执行。

USE family_ledger;

CREATE TABLE IF NOT EXISTS `t_family_saving` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `account_id` int unsigned NOT NULL DEFAULT 0 COMMENT '账本ID',
  `scope` varchar(16) NOT NULL DEFAULT 'personal' COMMENT '储蓄范围：family家庭 personal个人',
  `user_openid` varchar(64) NOT NULL DEFAULT '' COMMENT '个人储蓄所属用户openid，家庭储蓄为空',
  `target_amount` decimal(12,2) NOT NULL DEFAULT 0.00 COMMENT '储蓄目标金额',
  `current_amount` decimal(12,2) NOT NULL DEFAULT 0.00 COMMENT '当前已储蓄金额',
  `note` varchar(80) NOT NULL DEFAULT '' COMMENT '备注',
  `is_delete` tinyint NOT NULL DEFAULT 1 COMMENT '1正常 2删除',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_account_scope_user` (`account_id`, `scope`, `user_openid`),
  KEY `idx_account_scope` (`account_id`, `scope`),
  KEY `idx_user_openid` (`user_openid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='家庭账本储蓄状态表';
