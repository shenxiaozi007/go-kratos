# 登录用户与权限模块设计

本文档描述 v2 中登录用户与权限模块的设计，包括复用现有 auth 体系、JWT 签发与校验、角色与权限、xsh API 鉴权及前端登录与 Token。

---

## 一、登录用户

### 1.1 复用/扩展现有 auth 体系

- **API**：沿用 [api/auth/v1/auth.proto](../../api/auth/v1/auth.proto) 的 Register、Login、GetProfile。
- **数据**：沿用 [internal/data/user.go](../../internal/data/user.go) 的 UserPO 与 `users` 表；v2 若需扩展（如 role），见 [database/schema-v2-changes.md](database/schema-v2-changes.md)。
- **路由**：`POST /api/v1/auth/register`、`POST /api/v1/auth/login`、`GET /api/v1/auth/profile`（GetProfile 需经 JWT 中间件）。

### 1.2 JWT 签发与校验

- **签发**：Login 成功后在 AuthService 中签发 JWT（如 HS256），payload 至少包含：`user_id`、`username`、`exp`（过期时间）；`expires_in` 与 LoginReply 中一致（秒）。
- **校验**：在 [pkg/middeware/auth.go](../../pkg/middeware/auth.go) 中：
  - 从 `Authorization` Header 解析 `Bearer <token>`；
  - 验签并解析 payload，将 `user_id`（及可选 `username`、`role`）写入 context，供后续业务与权限中间件使用；
  - 无 Token 或无效/过期返回 401，错误信息统一（如 "invalid or expired token"）。
- **密钥与配置**：JWT 密钥（及可选 issuer/audience）从 conf 读取，注入到 auth 服务与中间件。

### 1.3 Register/Login 实现（若尚未实现）

- **Register**：校验 username 唯一、密码强度（可选）；密码经 bcrypt 等哈希后写入 `users` 表；返回 user_id、username（不返回 token，需再调 Login）。
- **Login**：根据 username 查 UserPO，校验密码哈希；通过则签发 JWT，返回 access_token、expires_in。

---

## 二、权限模块

### 2.1 角色

- **运营**：可管理选品、内容草稿与模板、查看/操作发布任务与评论等；不可绑定/解绑平台账号、不可管理其他后台用户与角色。
- **管理员**：拥有运营全部权限，并可绑定/解绑平台账号、管理用户与角色（若实现用户管理接口）。

与 [xsh-docs/01-overview](../xsh-docs/01-overview.md) 中角色约定一致。

### 2.2 角色与存储

- **方案**：在 `users` 表增加 `role` 字段（如 `operator` / `admin`），或使用简单角色表（如 `user_roles`：user_id, role）。在 [database/schema-v2-changes.md](database/schema-v2-changes.md) 与 [backend/permissions-design.md](backend/permissions-design.md) 中二选一并统一。
- **默认**：首个用户或指定初始化用户可为 admin；新注册用户默认 operator（或按业务配置）。

### 2.3 权限中间件

- 在 JWT 校验之后增加「角色/权限」中间件（或与 JWT 合并为同一中间件）：
  - 从 context 读取 `user_id`、`role`；
  - 按路由或接口标识判断所需角色（如「绑定/解绑账号」仅 admin）；
  - 不满足则返回 403。
- 权限点与 API 映射见 [backend/permissions-design.md](backend/permissions-design.md)。

---

## 三、xsh API 鉴权

- **原则**：所有 `/api/xsh/v1/*` 接口均需经过 JWT 中间件；未带有效 Token 的请求返回 401。
- **公开接口白名单**：若存在无需登录的 xsh 接口（如健康检查、回调），在 [backend/api-auth-xsh.md](backend/api-auth-xsh.md) 中明确列出；当前建议无白名单，即全部需鉴权。
- **路由注册**：在 HTTP 路由中，将 xsh 的 router 组统一挂在 JWT（及可选权限）中间件之下。

---

## 四、前端

### 4.1 登录页

- 提供登录表单：用户名、密码；提交调用 `POST /api/v1/auth/login`，成功后保存返回的 `access_token`、`expires_in`。
- **注册入口**：提供注册页（如 `/register`）或登录页内入口；提交调用 `POST /api/v1/auth/register`，注册成功后可自动调用 `POST /api/v1/auth/login` 直接登录。

### 4.2 Token 存储与请求头

- Token 存储：localStorage 或 sessionStorage（key 如 `xsh_access_token`）；可选同时存 `expires_at`（当前时间 + expires_in）用于前端判断是否即将过期。
- Axios（或请求封装）：在请求拦截器中为所有发往后端的请求添加 Header：`Authorization: Bearer <access_token>`；若未登录则不发 Token，或跳转登录页。

### 4.3 登录态过期与跳转

- 收到 401 时：清除本地 Token，跳转至登录页；可选提示「登录已过期，请重新登录」。
- 可选：在 `expires_at` 前 N 分钟主动刷新 Token（若后端提供 refresh 接口）或提示重新登录。

### 4.4 按角色隐藏/禁用菜单或按钮

- 登录后调用 `GET /api/v1/auth/profile`，返回中需包含当前用户 `role`（需在 GetProfileReply 或扩展接口中增加 role 字段）。
- 前端根据 role 控制：
  - 仅管理员可见「账号绑定/解绑」等菜单或按钮；
  - 运营不可见或禁用上述操作。
- 后端仍须做权限校验，前端仅做展示与交互上的隐藏/禁用。

---

## 五、与 v1 的衔接

- v1 的 xsh 接口当前无鉴权；v2 在 HTTP 层为 xsh 路由统一加上 JWT（及权限）中间件后，所有现有 xsh 接口自动受保护。
- 若 v1 已有 AuthService 实现但未签发真实 JWT，v2 需补全 JWT 签发与 middeware 中的验签与 context 写入。
