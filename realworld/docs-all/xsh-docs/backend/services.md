# Service 职责与外部依赖

---

## 1. 选品采集仓 (ProductService)

- **职责**：链接解析（ParseProductLink）、商品 CRUD、爆款榜单列表与同步。
- **依赖**：ProductRepo（xsh_products、xsh_product_links、xsh_hot_ranks）、可选 ProductParser（首期可仅做 URL 校验 + 人工录入，不调淘宝联盟）。
- **外部依赖**：**暂不对接淘宝联盟**。链接解析与爆款首期采用手动录入或简单解析/导入；文档预留后续对接联盟 API 的扩展点（如 ParseProductLink 内调用联盟转链/商品详情接口）。

---

## 2. AI 内容加工坊 (ContentService)

- **职责**：模板 CRUD、草稿 CRUD、GenerateCopy（AI 文案）、GenerateCover（图片合成）、ProcessVideo（视频去重）。
- **依赖**：TemplateRepo、DraftRepo、ProductRepo（读商品信息用于生成）、AI 客户端（文案/图像）、封面合成器、视频处理管道。
- **草稿 status 流转**：draft → ready → published。CreatePublishTask 仅接受 status=ready 的草稿（或按需允许 draft）；已 published 的草稿可编辑但不建议重复创建任务。
- **封面/视频 URL 与存储**：GenerateCover、ProcessVideo 的产出由后端生成后上传至对象存储（或本地静态服务），返回可访问 URL 并写回草稿对应字段；首期若无 OSS，可写本地路径或占位 URL，后续再接 OSS。草稿 **video_url** 为「原片」URL（用户上传或填写）；**ProcessVideo** 产出写回同一字段或单独字段「处理后视频 URL」，实现时二选一并在 API 中保持一致。
- **GenerateCopy 与模板关系**：输入为草稿（含商品信息）+ 所选 style（对应模板 style_type）；输出为整段文案，由 **Gemini 3 Flash** 按风格生成，可与模板占位符（如 {{title}}、{{price}}）合并后返回，或完全由 Gemini 自由生成，实现时二选一约定。
- **外部依赖**：**Gemini 3 Flash API** 用于文案生成与风格化洗稿（配合 xsh_content_templates 的 style_type）；封面合成为内部逻辑（商品图 + 价格爆炸贴 + 平台 Logo）；视频去重为 FFmpeg 或第三方，不依赖 Gemini。对接方式为 HTTP/gRPC 调用 Gemini 3 Flash，配置项见 [03-deployment.md](../03-deployment.md)（如 `xsh.ai.api_key`、`xsh.ai.model`）。

---

## 3. 多账号分发矩阵 (DispatchService)

- **职责**：平台账号 CRUD（含 proxy_url、cookie_storage_path）、定时规则 CRUD、发布任务创建/列表/取消。**不直接执行发布**，只写库；实际发布由 Playwright-go Worker 消费任务执行。
- **依赖**：PlatformAccountRepo、PublishTaskRepo、PublishScheduleRepo、DraftRepo（读草稿内容）。
- **外部依赖**：无（HTTP API 仅读写 DB）。发布执行见下节。

### 3.1 多账号发布执行（Playwright-go）

- **实现方式**：**Playwright-go Worker** 为独立进程（或与 Kratos 同进程内 goroutine），需独立配置文件或环境变量指向同一 MySQL，定时或常驻轮询 `xsh_publish_tasks` 中 status=pending 且 scheduled_at ≤ now 的任务。
- **定时与闭环**：定时由**系统 cron** 或 **Kratos 内定时任务** 二选一实现。可选方案一：cron 到点后根据 xsh_publish_schedules 中 is_enabled 且匹配时间点的规则，按 account_ids 与待选草稿（如指定草稿池或最新 ready 草稿）**自动创建**多条 xsh_publish_tasks，再由 Worker 消费；方案二：定时仅触发 Worker 拉取，任务全部由用户在后台手动创建。
- **流程**：按任务取 `account_id` → 查 `xsh_platform_accounts` 得到该账号的 **proxy_url**（及 cookie_storage_path）→ 使用该 **proxy** 启动/复用 Playwright Browser 或 Context（**一账号一 IP**）→ 恢复或执行登录（Cookie/Storage）→ 打开平台发布页 → 上传素材、填写文案并模拟提交 → 更新任务状态与 platform_content_id。
- **RetryPublishTask**：将 status=failed 的任务置回 pending（可选更新 scheduled_at 为当前或稍后），供 Worker 再次执行。
- **文档须说明**：
  - **一账号一 IP**：`xsh_platform_accounts` 表存每账号的 proxy 或独立 IP 配置，Playwright 启动时按账号注入，避免多账号同 IP 被封。
  - 与定时/队列的配合：可由 cron 或内部队列触发 Worker 拉取待执行任务。
  - Cookie/Storage 持久化：cookie_storage_path 指向本地文件或目录，用于保存/加载登录态。
  - 无头/有头与 Chromium：建议可配置，生产多用 headless；需安装 Chromium 或指定 executablePath。
  - 失败重试与日志：任务失败写 error_message、status=failed，可配置重试次数与告警。

---

## 4. 获客截流 (InboxService)

- **职责**：评论列表与筛选、标记已处理、同步评论、发送回复/私信、动作记录列表。
- **依赖**：CommentRepo、InboxActionRepo、PlatformAccountRepo。
- **SyncComments 主体**：由 **Kratos 内定时任务** 或单独「评论同步 Worker」（可与发布 Worker 同进程或分离）调用；按 account_id 使用该账号绑定的 proxy，通过 Playwright 模拟登录拉取评论后写入 xsh_comments。
- **want_link 写入时机**：评论同步入库后，由**异步任务**或**同步内**调用规则/关键词或 AI 写入 want_link 字段；若用 AI，可复用内容加工坊的 **Gemini 3 Flash** 做评论意图分类。同步时按 platform_comment_id 做「存在则更新、不存在则插入」，避免唯一键冲突。
- **外部依赖**：**不接官方开放平台**。评论同步与私信/回复均通过 **Playwright-go** 模拟登录后拉取评论、点击回复/私信并填写内容；需为每个账号使用该账号绑定的 **独立 proxy**（与发布一致，一账号一 IP），合规与风控注意同发布执行。「求链接」等意图识别采用 AI 或规则/关键词；若使用 AI，可复用内容加工坊的 **Gemini 3 Flash** 做评论意图分类（可选），便于实现与配置统一。

---

## 合规与风险汇总

- **当前不接淘宝联盟、不接抖音/小红书开放平台**；选品首期手动或简单解析，发布与登录、评论拉取与私信均为 **Playwright-go 模拟登录**。
- **一账号一 IP**：多账号发布与评论同步必须为每账号配置独立 IP/代理，在 schema（xsh_platform_accounts.proxy_url）、架构与部署文档中均需体现，避免同 IP 多号导致封号。
