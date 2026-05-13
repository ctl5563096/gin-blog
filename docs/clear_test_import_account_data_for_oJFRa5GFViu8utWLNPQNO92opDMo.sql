-- 手工执行：清空 open_id=oJFRa5GFViu8utWLNPQNO92opDMo 的“导入测试账本”业务数据，但保留账本本身与成员关系。
-- 注意：Codex 不直接执行 DDL/DML，本文件由你在数据库客户端里确认后执行。

USE family_ledger;

SET @open_id = 'oJFRa5GFViu8utWLNPQNO92opDMo';
SET @account_name = '导入测试账本';

SELECT @account_id := id
FROM t_family_account
WHERE owner_openid = @open_id
  AND name = @account_name
  AND is_delete = 1
ORDER BY id DESC
LIMIT 1;

-- 清空账单
DELETE FROM t_family_transaction
WHERE account_id = @account_id;

-- 清空储蓄流水
DELETE FROM t_family_saving_flow
WHERE account_id = @account_id;

-- 清空固定分类配置
DELETE FROM t_family_transaction_category_config
WHERE account_id = @account_id;

-- 清空每月固定储蓄配置
DELETE FROM t_family_monthly_saving_config
WHERE account_id = @account_id;

-- 重置账本储蓄余额，保留储蓄记录本身
UPDATE t_family_saving
SET target_amount = 0,
    current_amount = 0,
    note = '测试账本数据已清空',
    updated_at = NOW()
WHERE account_id = @account_id
  AND is_delete = 1;

SELECT
  @account_id AS account_id,
  @open_id AS owner_openid,
  @account_name AS account_name;
