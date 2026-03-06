# 后端模块划分与目录约定

推广引流后台在现有 **realworld** Kratos 项目中扩展，不新建独立服务。所有 xsh 相关代码使用 `xsh` 或 `xsh_` 前缀便于区分。

---

## 目录结构

```
realworld/
  api/
    xsh/
      v1/
        product.proto      # 选品采集仓
        content.proto      # AI 内容加工坊
        dispatch.proto     # 多账号分发矩阵
        inbox.proto        # 获客截流
        *.pb.go / *_http.pb.go / *_grpc.pb.go
  internal/
    biz/
      xsh_product.go       # 选品 Usecase
      xsh_content.go      # 草稿/模板 Usecase
      xsh_dispatch.go      # 账号/任务/定时 Usecase
      xsh_inbox.go         # 评论/私信 Usecase
    data/
      xsh_product.go      # xsh_products, xsh_product_links, xsh_hot_ranks PO & Repo
      xsh_content.go      # xsh_content_templates, xsh_content_drafts PO & Repo
      xsh_dispatch.go     # xsh_platform_accounts, xsh_publish_tasks, xsh_publish_schedules PO & Repo
      xsh_inbox.go        # xsh_comments, xsh_inbox_actions PO & Repo
    service/
      xsh_product.go      # 选品 HTTP Service
      xsh_content.go      # 内容加工 HTTP Service
      xsh_dispatch.go      # 分发矩阵 HTTP Service
      xsh_inbox.go        # 获客截流 HTTP Service
  cmd/realworld/
    wire.go               # 增加 xsh 的 Provider
    wire_gen.go           # 生成注入
```

---

## 模块与职责

| 模块 | Proto | Biz | Data 表 | Service | 说明 |
|------|-------|-----|---------|---------|------|
| 选品 | product.proto | ProductUsecase | xsh_products, xsh_product_links, xsh_hot_ranks | ProductService | 链接解析、商品列表、爆款（首期可不接联盟） |
| 内容 | content.proto | ContentUsecase | xsh_content_templates, xsh_content_drafts | ContentService | 模板 CRUD、草稿 CRUD、**Gemini 3 Flash** 文案生成、封面合成、视频去重 |
| 分发 | dispatch.proto | DispatchUsecase | xsh_platform_accounts, xsh_publish_tasks, xsh_publish_schedules | DispatchService | 账号 CRUD（含 proxy）、定时规则、发布任务；执行由 Playwright-go Worker 消费 |
| 截流 | inbox.proto | InboxUsecase | xsh_comments, xsh_inbox_actions | InboxService | 评论列表、同步、标记、私信/回复 |

---

## 发布执行器（Playwright-go）

发布与登录**不接官方 API**，由独立 **Playwright-go Worker**（可与主进程同机或独立部署）消费 `xsh_publish_tasks`，按 `account_id` 读取 `xsh_platform_accounts` 中的 **proxy 配置**，按**一账号一 IP** 启动浏览器完成模拟登录与发布。详见 [services.md](services.md)。

---

## 注册与依赖注入

- 在 [internal/server/http.go](../../internal/server/http.go) 中注册 `RegisterXshProductHTTPServer`、`RegisterXshContentHTTPServer`、`RegisterXshDispatchHTTPServer`、`RegisterXshInboxHTTPServer`（具体命名以生成的 Register 为准）。
- 在 [cmd/realworld/wire.go](../../cmd/realworld/wire.go) 的 `wire.Build` 中加入 xsh 的 data/biz/service Provider，并在 `wire_gen.go` 中为 `NewHTTPServer` 注入对应 xsh Service 参数。
