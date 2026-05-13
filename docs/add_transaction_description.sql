-- 手工执行：给账单新增收入/支出说明字段。
-- 注意：Codex 不直接执行 DDL/DML，本文件由你在数据库客户端里确认后执行。

USE family_ledger;

ALTER TABLE `t_family_transaction`
  ADD COLUMN `description` varchar(80) NOT NULL DEFAULT '' COMMENT '收入/支出说明' AFTER `category`;
