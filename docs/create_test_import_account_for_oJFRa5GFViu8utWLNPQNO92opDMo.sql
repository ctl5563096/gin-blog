-- 手工执行：为 open_id=oJFRa5GFViu8utWLNPQNO92opDMo 创建一个导入测试账本。
-- 注意：Codex 不直接执行 DDL/DML，本文件由你在数据库客户端里确认后执行。
-- 执行前请先确认已执行 docs/add_family_saving_table.sql。

USE family_ledger;

SET @open_id = 'oJFRa5GFViu8utWLNPQNO92opDMo';
SET @account_name = '导入测试账本';
SET @account_type = 'family';
SET @budget_amount = 8000;
SET @personal_budget_amount = 1000;
SET @initial_saving_amount = 0;

INSERT INTO t_family_user (
  openid,
  unionid,
  nick_name,
  avatar_url,
  last_login,
  created_at,
  updated_at
) VALUES (
  @open_id,
  '',
  '微信用户',
  '',
  NOW(),
  NOW(),
  NOW()
) ON DUPLICATE KEY UPDATE
  last_login = VALUES(last_login),
  updated_at = NOW();

INSERT INTO t_family_account (
  account_type,
  name,
  owner_openid,
  budget_amount,
  is_delete,
  created_at,
  updated_at
)
SELECT
  @account_type,
  @account_name,
  @open_id,
  @budget_amount,
  1,
  NOW(),
  NOW()
WHERE NOT EXISTS (
  SELECT 1
  FROM t_family_account
  WHERE owner_openid = @open_id
    AND account_type = @account_type
    AND name = @account_name
    AND is_delete = 1
);

SELECT @account_id := id
FROM t_family_account
WHERE owner_openid = @open_id
  AND account_type = @account_type
  AND name = @account_name
  AND is_delete = 1
ORDER BY id DESC
LIMIT 1;

INSERT INTO t_family_account_member (
  account_id,
  user_openid,
  role,
  personal_budget_amount,
  display_name,
  status,
  created_at,
  updated_at
) VALUES (
  @account_id,
  @open_id,
  'owner',
  @personal_budget_amount,
  '微信用户',
  1,
  NOW(),
  NOW()
) ON DUPLICATE KEY UPDATE
  role = VALUES(role),
  personal_budget_amount = VALUES(personal_budget_amount),
  status = VALUES(status),
  updated_at = NOW();

INSERT INTO t_family_saving (
  account_id,
  scope,
  user_openid,
  target_amount,
  current_amount,
  note,
  is_delete,
  created_at,
  updated_at
)
SELECT
  @account_id,
  'family',
  '',
  0,
  @initial_saving_amount,
  '导入测试账本初始化',
  1,
  NOW(),
  NOW()
WHERE NOT EXISTS (
  SELECT 1
  FROM t_family_saving
  WHERE account_id = @account_id
    AND scope = 'family'
    AND user_openid = ''
    AND is_delete = 1
);

SELECT
  @account_id AS account_id,
  @open_id AS owner_openid,
  @account_name AS account_name;
