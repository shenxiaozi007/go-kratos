# 小红书推广引流后台 v2 — 计划文档

本目录为推广引流后台**第二版（v2）**的完整计划文档，与 [xsh-docs](../xsh-docs/)（v1 规格）并列。v1 已完成首版开发；v2 在 v1 基础上增加登录与权限、Playwright 小红书登录与用户信息抓取、Gemini 文案真实接入，并明确其余未完善项的处理方式。

---

## v2 相对 v1 的增量

| 维度 | v1 | v2 |
|------|----|----|
| **登录与权限** | 无；xsh 接口未鉴权 | 接入登录用户；JWT 校验；角色（运营/管理员）；xsh 路由保护；前端登录与 Token |
| **Playwright 小红书** | 仅文档描述 | 模拟登录小红书网页；抓取登录后基本用户信息并落库；Cookie 持久化与复用 |
| **Gemini 文案** | Noop 占位 | 真实接入 Gemini 3 Flash（或等价模型），替换 NoopCopyGenerator |
| **其他** | 封面/视频/评论同步/发布 Worker 为占位或未实现 | 在 05 中列清单并标明 v2 实现 / 文档约定 / 延后 |

---

## 文档索引与阅读顺序

1. **[01-v2-scope.md](01-v2-scope.md)** — v2 范围：三大块 + 其他完善项清单及 v2 处理方式。
2. **[02-auth-and-permissions.md](02-auth-and-permissions.md)** — 登录用户与权限模块设计（复用 auth、JWT、角色、xsh 鉴权、前端）。
3. **[03-playwright-xiaohongshu.md](03-playwright-xiaohongshu.md)** — Playwright-go 接入：模拟登录小红书网页、抓取登录后基本用户信息、存储与执行主体。
4. **[04-gemini-copywriting.md](04-gemini-copywriting.md)** — Gemini 3 Flash 文案生成接入（配置、实现、错误与限流）。
5. **[05-remaining-improvements.md](05-remaining-improvements.md)** — 其他未完善功能与 v2 处理方式（封面/视频/评论同步/发布 Worker/定时/前端）。
6. **[backend/api-auth-xsh.md](backend/api-auth-xsh.md)** — xsh 鉴权方式与路由保护约定。
7. **[backend/permissions-design.md](backend/permissions-design.md)** — 角色/权限表与中间件设计。
8. **[database/schema-v2-changes.md](database/schema-v2-changes.md)** — v2 新增/变更表与字段（用户与角色、小红书账号扩展信息等）。

---

## v2 实现顺序建议

1. **阶段 1**：鉴权与权限（后端 JWT + 角色/权限 + xsh 路由保护）、前端登录与 Token。
2. **阶段 2**：Gemini 文案接入（配置、CopyGenerator 实现、Wire 替换 Noop）。
3. **阶段 3**：Playwright 小红书登录与抓取用户信息（登录流程、用户信息落库、Cookie 持久化）。
4. **后续**：封面合成、视频去重、评论同步、发布 Worker、定时规则与自动创建任务等，按 [05-remaining-improvements.md](05-remaining-improvements.md) 的优先级排期。

---

## 与 v1 文档的关系

- **v1 规格**：[docs-all/xsh-docs](../xsh-docs/) — 业务概览、架构、库表、API、前端路由、部署等仍为基线。
- **v2 计划**：本目录仅描述 v2 新增与变更；未提及的部分沿用 v1 设计与实现。
