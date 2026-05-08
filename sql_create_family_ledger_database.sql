CREATE DATABASE IF NOT EXISTS `family_ledger`
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_unicode_ci;

USE `family_ledger`;

CREATE TABLE IF NOT EXISTS `t_family_user` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `openid` varchar(64) NOT NULL DEFAULT '' COMMENT '微信小程序openid',
  `unionid` varchar(64) NOT NULL DEFAULT '' COMMENT '微信开放平台unionid',
  `nick_name` varchar(64) NOT NULL DEFAULT '' COMMENT '用户昵称，用户主动填写或授权后更新',
  `avatar_url` varchar(512) NOT NULL DEFAULT '' COMMENT '用户头像，用户主动选择后更新',
  `last_login` timestamp NULL DEFAULT NULL COMMENT '最近登录时间',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_openid` (`openid`),
  KEY `idx_unionid` (`unionid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='家庭记账小程序用户';


CREATE TABLE IF NOT EXISTS `t_family_transaction` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `account_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '账本ID',
  `owner_id` varchar(64) NOT NULL DEFAULT '' COMMENT '微信openid，后续可扩展为家庭账本ID',
  `type` varchar(16) NOT NULL DEFAULT 'expense' COMMENT 'expense支出 income收入',
  `amount` decimal(12,2) NOT NULL DEFAULT '0.00' COMMENT '金额',
  `category` varchar(32) NOT NULL DEFAULT '' COMMENT '分类',
  `note` varchar(80) NOT NULL DEFAULT '' COMMENT '备注',
  `date` char(10) NOT NULL DEFAULT '' COMMENT '记账日期 YYYY-MM-DD',
  `payment_method` varchar(32) NOT NULL DEFAULT '' COMMENT '支付方式',
  `member_name` varchar(32) NOT NULL DEFAULT '' COMMENT '成员名称',
  `is_delete` tinyint(1) NOT NULL DEFAULT '1' COMMENT '1正常 2删除',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_account_date` (`account_id`,`date`),
  KEY `idx_account_type` (`account_id`,`type`),
  KEY `idx_owner_date` (`owner_id`,`date`),
  KEY `idx_owner_type` (`owner_id`,`type`),
  KEY `idx_owner_category` (`owner_id`,`category`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='家庭记账流水';

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

CREATE TABLE IF NOT EXISTS `t_family_budget` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `owner_id` varchar(64) NOT NULL DEFAULT '' COMMENT '微信openid，后续可扩展为家庭账本ID',
  `amount` decimal(12,2) NOT NULL DEFAULT '0.00' COMMENT '月预算',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_owner_id` (`owner_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='家庭记账预算';
