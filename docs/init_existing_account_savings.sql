-- 手工执行：为已有账本补齐主储蓄记录。
-- 注意：Codex 不直接执行 DDL/DML，本文件由你在数据库客户端里确认后执行。
-- 执行前请先确认已执行 docs/add_family_saving_table.sql。

USE family_ledger;

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
  a.id,
  CASE WHEN a.account_type = 'family' THEN 'family' ELSE 'personal' END AS scope,
  CASE WHEN a.account_type = 'family' THEN '' ELSE a.owner_openid END AS user_openid,
  0,
  0,
  '历史账本初始化',
  1,
  NOW(),
  NOW()
FROM t_family_account a
WHERE a.is_delete = 1
  AND NOT EXISTS (
    SELECT 1
    FROM t_family_saving s
    WHERE s.account_id = a.id
      AND s.scope = CASE WHEN a.account_type = 'family' THEN 'family' ELSE 'personal' END
      AND s.user_openid = CASE WHEN a.account_type = 'family' THEN '' ELSE a.owner_openid END
      AND s.is_delete = 1
  );
