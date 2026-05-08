USE `family_ledger`;

CREATE TABLE IF NOT EXISTS `t_family_account` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `account_type` varchar(16) NOT NULL DEFAULT 'personal' COMMENT '账本类型 personal个人 family家庭',
  `name` varchar(32) NOT NULL DEFAULT '' COMMENT '账本名称',
  `owner_openid` varchar(64) NOT NULL DEFAULT '' COMMENT '账本创建人openid',
  `budget_amount` decimal(12,2) NOT NULL DEFAULT '0.00' COMMENT '账本月额度',
  `is_delete` tinyint(1) NOT NULL DEFAULT '1' COMMENT '1正常 2删除',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_owner_type` (`owner_openid`,`account_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='家庭记账账本';

CREATE TABLE IF NOT EXISTS `t_family_account_member` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `account_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '账本ID',
  `user_openid` varchar(64) NOT NULL DEFAULT '' COMMENT '成员openid',
  `role` varchar(16) NOT NULL DEFAULT 'member' COMMENT 'owner创建人 member成员',
  `personal_budget_amount` decimal(12,2) NOT NULL DEFAULT '0.00' COMMENT '成员个人月额度',
  `display_name` varchar(64) NOT NULL DEFAULT '' COMMENT '成员展示名称',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '1正常 2退出',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_account_user` (`account_id`,`user_openid`),
  KEY `idx_user_status` (`user_openid`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='家庭记账账本成员';

ALTER TABLE `t_family_transaction`
  ADD COLUMN `account_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '账本ID' AFTER `id`,
  ADD KEY `idx_account_date` (`account_id`,`date`),
  ADD KEY `idx_account_type` (`account_id`,`type`);

-- 下面是历史数据迁移模板，请按实际情况人工执行：
-- 1. 为每个已有用户创建一个个人账本。
-- 2. 把旧 t_family_budget.amount 迁移到 t_family_account.budget_amount。
-- 3. 把旧 t_family_transaction.account_id 从 0 更新成对应个人账本 id。
--
-- INSERT INTO t_family_account (account_type, name, owner_openid, budget_amount)
-- SELECT 'personal', '个人账本', u.openid, COALESCE(b.amount, 0)
-- FROM t_family_user u
-- LEFT JOIN t_family_budget b ON b.owner_id = u.openid
-- WHERE NOT EXISTS (
--   SELECT 1 FROM t_family_account a WHERE a.owner_openid = u.openid AND a.account_type = 'personal' AND a.is_delete = 1
-- );
--
-- INSERT INTO t_family_account_member (account_id, user_openid, role, personal_budget_amount, display_name)
-- SELECT a.id, a.owner_openid, 'owner', a.budget_amount, COALESCE(NULLIF(u.nick_name, ''), '微信用户')
-- FROM t_family_account a
-- INNER JOIN t_family_user u ON u.openid = a.owner_openid
-- WHERE NOT EXISTS (
--   SELECT 1 FROM t_family_account_member m WHERE m.account_id = a.id AND m.user_openid = a.owner_openid
-- );
--
-- UPDATE t_family_transaction t
-- INNER JOIN t_family_account a ON a.owner_openid = t.owner_id AND a.account_type = 'personal' AND a.is_delete = 1
-- SET t.account_id = a.id
-- WHERE t.account_id = 0;
