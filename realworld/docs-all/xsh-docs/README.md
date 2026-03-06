# 小红书/抖音推广引流后台 — 规格文档

本目录为「推广引流后台」的完整规格文档，供后续开发使用。后台面向小红书、抖音的选品、内容加工、多账号分发与获客截流，**暂不对接淘宝联盟**，发布与登录**全部采用模拟登录**（Playwright-go），并采用**一账号一 IP** 防封策略。

---

## 四大模块

| 模块 | 说明 |
|------|------|
| **选品采集仓 (Input)** | 链接解析（粘贴淘宝/猫超链接）、爆款监控；首期手动录入或简单解析，不接淘宝联盟 API。 |
| **AI 内容加工坊 (Core)** | 文案模板与洗稿（由 **Gemini 3 Flash** 生成）、封面为商品图 + 价格爆炸贴 + 平台 Logo 合成、视频去重（可选）。 |
| **多账号分发矩阵 (Dispatch)** | 账号绑定（含代理/IP 配置）、定时发布、Playwright-go 模拟登录发布，一账号一 IP。 |
| **获客截流管理 (Inbox)** | 评论汇总、AI 识别「求链接」、私信/回复引导。 |

---

## 文档索引与阅读顺序

1. **[01-overview.md](01-overview.md)** — 业务背景、用户角色、核心流程（选品→加工→发布→截流）。
2. **[02-architecture.md](02-architecture.md)** — 系统架构图、技术栈、与现有 realworld 的集成、发布执行器与一账号一 IP。
3. **[database/](database/)** — 数据库设计  
   - [README.md](database/README.md)：库表总览、ER 说明  
   - [schema.md](database/schema.md)：表结构清单与字段  
   - [migrations.md](database/migrations.md)：迁移策略  
4. **[backend/](backend/)** — 后端设计  
   - [README.md](backend/README.md)：模块划分、目录约定  
   - [errors-and-conventions.md](backend/errors-and-conventions.md)：错误格式、分页与时间约定  
   - [api-catalog.md](backend/api-catalog.md)：API 列表  
   - [proto-outline.md](backend/proto-outline.md)：Proto 与 Message/RPC 大纲  
   - [services.md](backend/services.md)：Service 职责与外部依赖（含 Playwright-go、一账号一 IP）  
5. **[frontend/](frontend/)** — 前端设计  
   - [README.md](frontend/README.md)：技术栈、目录结构  
   - [views-routes.md](frontend/views-routes.md)：页面与路由  
   - [components.md](frontend/components.md)：关键组件  
6. **[03-deployment.md](03-deployment.md)** — 部署与配置、Playwright 与代理/IP 绑定。

---

## 约定摘要

- **淘宝联盟**：暂不对接；选品首期手动录入或简单解析。
- **抖音/小红书**：不接官方开放平台；发布、评论、私信均为**模拟登录**（Playwright-go）。
- **防封**：多账号发布**一账号一 IP**，每账号在 `xsh_platform_accounts` 绑定独立代理或独立出口 IP。
- **AI**：内容加工坊文案生成等采用 **Gemini 3 Flash**。

**建议实现顺序**见 [01-overview.md](01-overview.md) 末节，便于分阶段交付（选品+内容 → 分发 → 定时 → 评论与 Inbox）。
