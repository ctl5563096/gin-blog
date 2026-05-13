-- 手工执行：创建导入用家庭账本，并授权指定微信用户。
-- 注意：Codex 不直接执行 DDL/DML，本文件由你在数据库客户端里确认后执行。

USE family_ledger;

SET @open_id = 'oJFRa5GFViu8utWLNPQNO92opDMo';
SET @account_name = '家庭账本';

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
  'family',
  @account_name,
  @open_id,
  8000,
  1,
  NOW(),
  NOW()
WHERE NOT EXISTS (
  SELECT 1
  FROM t_family_account
  WHERE account_type = 'family'
    AND owner_openid = @open_id
    AND name = @account_name
    AND is_delete = 1
);

SELECT @account_id := id
FROM t_family_account
WHERE account_type = 'family'
  AND owner_openid = @open_id
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
  1000,
  '微信用户',
  1,
  NOW(),
  NOW()
) ON DUPLICATE KEY UPDATE
  role = VALUES(role),
  personal_budget_amount = VALUES(personal_budget_amount),
  status = VALUES(status),
  updated_at = NOW();
