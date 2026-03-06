# v2 范围

本文档明确 v2 的三大块范围及 v1 未完善项在 v2 中的处理方式。

---

## 一、v2 三大块

### 1. 接入登录用户与权限模块

- **登录用户**：复用/扩展现有 auth 体系（api/auth/v1、users 表），实现 Register/Login（若未实现），并补齐 JWT 签发与校验；xsh 后台用户与现有 users 表统一或扩展，在 [02-auth-and-permissions.md](02-auth-and-permissions.md) 与 [database/schema-v2-changes.md](database/schema-v2-changes.md) 中约定。
- **权限模块**：角色为「运营」「管理员」；简单角色表或 user 表增加 role 字段；权限中间件按角色放行 xsh 接口；可选权限点与 API 映射（如仅管理员可绑定/解绑平台账号）。详见 [02-auth-and-permissions.md](02-auth-and-permissions.md) 与 [backend/permissions-design.md](backend/permissions-design.md)。
- **xsh API 鉴权**：所有 `/api/xsh/v1/*` 需经 JWT 中间件；公开接口白名单在 [backend/api-auth-xsh.md](backend/api-auth-xsh.md) 中写明。
- **前端**：登录页、Token 存储与 axios 请求头、登录态过期与跳转、按角色隐藏/禁用菜单或按钮。见 02 的「前端」小节。

### 2. 完善接入 Playwright-go：模拟登录小红书网页并抓取登录后基本用户信息

- **目标**：使用 Playwright-go 打开小红书登录页，完成模拟登录（扫码或账号密码，以 03 约定为准），在登录成功后抓取基本用户信息（如昵称、头像、平台用户 ID 等），并落库。
- **范围**：v2 仅包含「登录 + 抓取用户信息 + 回写存储」；不包含发帖、评论拉取、发布任务消费 Worker。详见 [03-playwright-xiaohongshu.md](03-playwright-xiaohongshu.md)。
- **存储**：在 xsh_platform_accounts 增加字段或新表存平台侧用户信息，见 [database/schema-v2-changes.md](database/schema-v2-changes.md)。

### 3. 接入 Gemini 3 Flash 文案生成功能

- **目标**：将内容加工坊「一键生成文案」从 Noop 换为真实调用 Gemini 3 Flash（或当前可用等价模型，如 gemini-2.0-flash）。
- **范围**：配置 xsh.ai、实现 CopyGenerator 的 Gemini 实现、Wire 替换 Noop、错误与限流约定。详见 [04-gemini-copywriting.md](04-gemini-copywriting.md)。

---

## 二、v1 未完善项与 v2 处理方式

| 项目 | v1 状态 | v2 处理方式 |
|------|--------|-------------|
| 登录/用户/权限 | 部分存在，未与 xsh 打通 | **v2 实现**：鉴权、角色、xsh 路由保护、前端登录与 Token |
| Playwright 小红书登录与用户信息 | 仅文档 | **v2 实现**：模拟登录、抓取用户信息、落库 |
| Gemini 文案 | Noop | **v2 实现**：真实接入 Gemini 3 Flash |
| 封面合成（GenerateCover） | Noop | **v2 文档约定**；实现标为 v2 可选或后续迭代，见 [05-remaining-improvements.md](05-remaining-improvements.md) |
| 视频去重（ProcessVideo） | Noop | **v2 文档约定**；实现标为 v2 可选或后续，见 05 |
| 评论同步（SyncComments） | 占位 | **v2 文档约定**为后续阶段；v2 仅保证 Playwright 能登录并抓取用户信息，见 03、05 |
| 发布任务执行（Worker 消费任务发帖） | 无 | **v2 文档约定**与「登录/抓取用户信息」边界；实现列 v2 后续或 v3，见 05 |
| 定时规则与自动创建任务 | 未实现 cron 联动 | **v2 文档约定**方案；实现可标为 v2 可选，见 05 |
| 前端鉴权 | 无 | **v2 实现**：登录页、Token、按角色控制菜单/按钮，见 02 |

---

## 三、实现顺序（供排期参考）

1. **阶段 1**：鉴权与权限（后端 + 前端登录与 Token）。
2. **阶段 2**：Gemini 文案接入。
3. **阶段 3**：Playwright 小红书登录与抓取用户信息。
4. **后续**：封面/视频/评论同步/发布 Worker/定时等，按 05 的优先级与资源排期。
