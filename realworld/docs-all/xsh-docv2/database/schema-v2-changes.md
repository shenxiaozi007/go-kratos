# v2 数据库新增/变更

本文档描述 v2 相对 v1 的数据库变更，与 [xsh-docs/database/schema](../xsh-docs/database/schema.md) 配合使用；未列出的表与字段沿用 v1。

---

## 一、用户与角色（鉴权与权限）

### 1.1 users 表扩展

在现有 `users` 表上增加字段（若尚未存在）：

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| role | VARCHAR(32) DEFAULT 'operator' | 角色：operator / admin | INDEX |

- 默认值：新注册用户为 `operator`；首个或指定初始化用户可为 `admin`。
- 若 v1 的 `users` 表为 `internal/data/user.go` 中的 UserPO（无 role），则在本版本迁移中增加 `role` 字段。

### 1.2 可选：user_roles 表（方案 B）

若采用「一用户多角色」或独立角色表，可新增：

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| id | BIGINT UNSIGNED AUTO_INCREMENT | 主键 | PK |
| user_id | BIGINT UNSIGNED NOT NULL | 关联 users.id | FK, INDEX |
| role | VARCHAR(32) NOT NULL | operator / admin | INDEX |
| created_at | DATETIME | 创建时间 | |

v2 建议优先使用 **users.role** 单字段（方案 A），无需本表。

---

## 二、小红书账号扩展信息（Playwright 登录与抓取用户信息）

采用 **方案 A**：在 `xsh_platform_accounts` 上增加字段，存储登录后抓取到的平台侧用户信息。

### 2.1 xsh_platform_accounts 新增字段

在 v1 表结构基础上增加：

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| platform_user_id | VARCHAR(128) | 平台侧用户唯一 ID（如小红书 uid） | INDEX |
| platform_nickname | VARCHAR(128) | 平台昵称 | |
| platform_avatar | VARCHAR(512) | 头像 URL | |
| platform_red_id | VARCHAR(64) | 小红书号（如有） | |
| logged_at | DATETIME 或 BIGINT | 最近一次登录成功时间（用于判断登录态是否过期等） | INDEX |

- 上述字段均为可选（NULL）；在 Playwright 完成「登录 + 抓取用户信息」后由后端回写。
- 与 [03-playwright-xiaohongshu.md](../03-playwright-xiaohongshu.md) 中「抓取的基本用户信息字段」一致。

### 2.2 方案 B 备选（本版本不采用）

若后续希望将「平台档案」与「账号绑定」分离，可新增表 `xsh_platform_account_profiles`（account_id, platform_user_id, platform_nickname, platform_avatar, platform_red_id, logged_at, updated_at）。v2 统一采用方案 A，减少表与 JOIN。

---

## 三、迁移与兼容

- 新增字段均允许 NULL 或设默认值，保证对已有数据兼容。
- 迁移方式沿用项目现有策略（如 GORM AutoMigrate 或独立 SQL 迁移脚本）；具体执行在实现阶段完成，本文档仅约定结构。

---

## 四、其他

- v2 不新增与「评论」「发布任务」「内容草稿」等业务表结构变更；仅用户/角色与 `xsh_platform_accounts` 扩展。
- 若 auth 相关需增加「刷新 Token」「登出」等，可后续在 users 或独立 token 表中扩展，不在本版文档强制约定。
