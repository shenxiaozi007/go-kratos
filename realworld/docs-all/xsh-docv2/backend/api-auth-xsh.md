# xsh 鉴权方式与路由保护

本文档约定 v2 中 xsh 相关 API 的鉴权方式与路由保护规则。

---

## 一、鉴权原则

- 所有 **`/api/xsh/v1/*`** 接口均需经过 **JWT 鉴权中间件**；请求须携带有效 `Authorization: Bearer <access_token>`，否则返回 **401**。
- 鉴权通过后，中间件将 `user_id`（及可选 `username`、`role`）写入 context，供业务层与权限中间件使用。

---

## 二、公开接口白名单

- **当前约定**：xsh 模块**无公开接口**；所有 xsh 接口均需登录后访问。
- 若后续需增加无需登录的接口（如健康检查、回调、公开展示），须在本文档中**明确列出**路径与用途，并在路由注册时将该路径排除在 JWT 中间件之外。

---

## 三、路由注册方式

- 在 HTTP Server 中，将 xsh 的 router 组（或所有 xsh 的 HTTP 映射）统一挂在 **JWT 中间件** 之下；可选在 JWT 之后再挂 **权限中间件**（按角色放行部分敏感接口）。
- 与 xsh 无关的路由（如 `/api/v1/auth/login`、`/api/v1/auth/register`）不经过 JWT，保持可匿名访问；`/api/v1/auth/profile` 需经 JWT。

---

## 四、错误响应

- **401 Unauthorized**：未提供 Token、Token 无效或已过期；Body 可统一为项目约定的错误格式（如 `code: 401, reason: "invalid or expired token"`）。
- **403 Forbidden**：Token 有效但角色/权限不足（由权限中间件返回）；Body 注明权限不足原因（如 "admin only"）。

---

## 五、与 02 的衔接

- JWT 签发、校验逻辑及 context 写入见 [02-auth-and-permissions.md](../02-auth-and-permissions.md)。
- 权限点与 API 的映射（如「仅管理员可绑定/解绑账号」）见 [permissions-design.md](permissions-design.md)。
