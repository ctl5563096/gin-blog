# 家庭记账接口

基础地址：`http://127.0.0.1:8081`

统一返回：

```json
{
  "code": 200,
  "msg": "OK",
  "data": {}
}
```

当前版本使用微信小程序登录鉴权：

1. 小程序调用 `wx.login()` 获取临时 `code`。
2. 后端调用微信 `jscode2session` 换取 `openid`。
3. 后端生成业务 token，写入 Redis。
4. 小程序后续请求通过 `Authorization: Bearer <token>` 访问记账接口。

后端需要配置环境变量：

- `WECHAT_MINIAPP_APPID`
- `WECHAT_MINIAPP_SECRET`

也兼容旧命名：

- `WX_APPID`
- `WX_SECRET`

## 数据表

先执行项目根目录下的 `sql_create_family_ledger_database.sql`，创建专用数据库和数据表。

数据库名默认是：

- `family_ledger`

后端配置项：

- `LEDGER_DATABASE=family_ledger`

会创建：

- `t_family_user`
- `t_family_account`
- `t_family_account_member`
- `t_family_transaction`
- `t_family_budget`

本次账本类型改造需要人工执行：

- `sql_family_ledger_account_migration.sql`

## 接口列表

### 微信登录

`POST /v1/ledger/auth/login`

请求：

```json
{
  "code": "wx.login 返回的 code"
}
```

返回：

```json
{
  "token": "业务登录态 token",
  "openid": "微信 openid",
  "unionid": "",
  "expiresIn": 86400
}
```

下面所有 `/v1/ledger/*` 业务接口都需要请求头：

```text
Authorization: Bearer <token>
```

除账本列表和创建账本外，账单、汇总、预算接口还需要当前账本：

```text
X-Ledger-Account-Id: <账本ID>
```

如果用户没有账本，小程序必须引导用户创建个人账本或家庭账本，或者退出小程序。

### 账本列表

`GET /v1/ledger/accounts`

返回当前用户加入的个人/家庭账本。

### 创建账本

`POST /v1/ledger/accounts`

```json
{
  "accountType": "family",
  "name": "家庭账本",
  "budgetAmount": 5000,
  "personalBudgetAmount": 1000
}
```

`accountType` 可选值：

- `personal` 个人账本
- `family` 家庭账本

### 账单列表

`GET /v1/ledger/transactions`

参数：

- `month` 可选，格式 `YYYY-MM`
- `type` 可选，`all`、`expense`、`income`

### 新增账单

`POST /v1/ledger/transactions`

```json
{
  "type": "expense",
  "amount": 36.5,
  "category": "餐饮",
  "note": "晚饭",
  "date": "2026-05-07",
  "paymentMethod": "微信",
  "memberName": "我"
}
```

### 更新账单

`PUT /v1/ledger/transactions`

参数同新增账单，额外传 `id`。

### 删除账单

`DELETE /v1/ledger/transactions?id=1`

删除为软删除，`is_delete` 从 `1` 更新为 `2`。

### 月度汇总

`GET /v1/ledger/summary?month=2026-05`

返回：

```json
{
  "income": 8000,
  "expense": 128.5,
  "balance": 7871.5,
  "todayExpense": 42,
  "count": 3
}
```

### 获取预算

`GET /v1/ledger/budget`

返回：

```json
{
  "accountId": 1,
  "accountType": "family",
  "amount": 5000,
  "personalAmount": 1000,
  "personalBudgetAmount": 1000
}
```

### 保存预算

`PUT /v1/ledger/budget`

```json
{
  "scope": "account",
  "amount": 5000
}
```

`scope` 可选值：

- `account` 账本额度，默认值
- `member` 当前成员个人额度

### 获取用户资料

`GET /v1/ledger/user/profile`

需要登录 token。

### 更新用户资料

`PUT /v1/ledger/user/profile`

需要登录 token。

```json
{
  "nickName": "小陈",
  "avatarUrl": "wxfile://tmp/avatar.png"
}
```

说明：

- 昵称来自小程序 `type="nickname"` 输入框。
- 头像来自小程序 `open-type="chooseAvatar"`。
- 头像应先上传 OSS，再把返回的 URL 保存到资料。

### 上传用户头像

`POST /v1/ledger/user/avatar`

需要登录 token，`multipart/form-data` 上传。

字段：

- `avatar`：图片文件，支持 `.png`、`.jpg`、`.jpeg`、`.gif`

返回：

```json
{
  "path": "https://ctl-blog-1.oss-cn-beijing.aliyuncs.com/ledger/avatar/202605/xxx.jpg"
}
```

OSS 路径规则：

```text
ledger/avatar/YYYYMM/<随机文件名>.<扩展名>
```
