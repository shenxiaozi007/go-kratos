# API 列表（按模块）

所有接口建议统一前缀，例如 `/api/xsh/v1/`，具体以 proto 中 `google.api.http` 为准。以下为 RPC 级列表，对应到 [proto-outline.md](proto-outline.md)。

---

## 通用约定

- **分页**：`page` 从 1 起，`page_size` 默认 20、最大 100；列表响应为 `items` + `total`。
- **错误**：格式与业务错误码见 [errors-and-conventions.md](errors-and-conventions.md)。
- **时间**：Proto 时间字段约定见 [errors-and-conventions.md](errors-and-conventions.md)。

---

## 1. 选品采集仓 (Input)

| 方法 | HTTP | 说明 |
|------|------|------|
| ParseProductLink | POST /api/xsh/v1/products/parse | 解析链接，返回商品名、原价、券额等；首期可简单解析或人工录入，不依赖淘宝联盟 |
| ListProducts | GET /api/xsh/v1/products | 分页列表，支持按 is_hot、rank_source 等筛选 |
| GetProduct | GET /api/xsh/v1/products/{id} | 单商品详情 |
| CreateProduct | POST /api/xsh/v1/products | 手动创建商品（录入） |
| UpdateProduct | PUT /api/xsh/v1/products/{id} | 更新商品 |
| DeleteProduct | DELETE /api/xsh/v1/products/{id} | 删除商品 |
| ListHotRanks | GET /api/xsh/v1/hot-ranks | 爆款榜单列表（首期可不接联盟，可空或手动数据） |
| SyncHotRank | POST /api/xsh/v1/hot-ranks/sync | 触发同步（首期可空实现） |

---

## 2. AI 内容加工坊 (Core)

| 方法 | HTTP | 说明 |
|------|------|------|
| ListTemplates | GET /api/xsh/v1/templates | 模板列表 |
| CreateTemplate | POST /api/xsh/v1/templates | 创建模板 |
| UpdateTemplate | PUT /api/xsh/v1/templates/{id} | 更新模板 |
| DeleteTemplate | DELETE /api/xsh/v1/templates/{id} | 删除模板 |
| CreateDraft | POST /api/xsh/v1/drafts | 创建草稿（关联 product_id、template_id） |
| ListDrafts | GET /api/xsh/v1/drafts | 分页列表 |
| GetDraft | GET /api/xsh/v1/drafts/{id} | 草稿详情 |
| UpdateDraft | PUT /api/xsh/v1/drafts/{id} | 更新文案/封面/视频 |
| DeleteDraft | DELETE /api/xsh/v1/drafts/{id} | 删除草稿 |
| GenerateCopy | POST /api/xsh/v1/drafts/{id}/generate-copy | 调用 **Gemini 3 Flash** 按风格生成文案（style 参数对应模板风格） |
| GenerateCover | POST /api/xsh/v1/drafts/{id}/generate-cover | 合成封面图，返回 URL |
| ProcessVideo | POST /api/xsh/v1/drafts/{id}/process-video | 视频去重（抽帧/镜像/滤镜），返回新视频 URL（可选） |

---

## 3. 多账号分发矩阵 (Dispatch)

| 方法 | HTTP | 说明 |
|------|------|------|
| ListPlatformAccounts | GET /api/xsh/v1/accounts | 按 platform 筛选列表 |
| BindAccount | POST /api/xsh/v1/accounts | 绑定账号（含 proxy_url、cookie_storage_path 等，一账号一 IP） |
| UpdateAccount | PUT /api/xsh/v1/accounts/{id} | 更新账号（含代理） |
| UnbindAccount | DELETE /api/xsh/v1/accounts/{id} | 解绑账号 |
| ListSchedules | GET /api/xsh/v1/schedules | 定时规则列表 |
| CreateSchedule | POST /api/xsh/v1/schedules | 创建定时规则 |
| UpdateSchedule | PUT /api/xsh/v1/schedules/{id} | 更新 |
| DeleteSchedule | DELETE /api/xsh/v1/schedules/{id} | 删除 |
| CreatePublishTask | POST /api/xsh/v1/publish-tasks | 创建发布任务：请求 draft_id、account_ids[]、scheduled_at；一条请求创建 N 条任务（每账号一条），响应返回 tasks[] 或 task_ids[] |
| ListPublishTasks | GET /api/xsh/v1/publish-tasks | 分页列表，按状态/账号筛选 |
| RetryPublishTask | POST /api/xsh/v1/publish-tasks/{id}/retry | 重试：将 status=failed 的任务置回 pending，供 Worker 再次执行（可选更新 scheduled_at） |
| CancelPublishTask | POST /api/xsh/v1/publish-tasks/{id}/cancel | 取消任务 |

---

## 4. 获客截流 (Inbox)

| 方法 | HTTP | 说明 |
|------|------|------|
| ListComments | GET /api/xsh/v1/comments | 按 account_id、want_link、status 筛选分页 |
| MarkCommentHandled | POST /api/xsh/v1/comments/{id}/handled | 标记已处理/忽略 |
| SyncComments | POST /api/xsh/v1/comments/sync | 按账号同步评论（Playwright 或定时拉取） |
| SendReply | POST /api/xsh/v1/comments/{id}/reply | 回复评论 |
| SendPrivateMessage | POST /api/xsh/v1/comments/{id}/private-message | 私信引导 |
| ListInboxActions | GET /api/xsh/v1/inbox-actions | 动作记录列表，可按 comment_id 筛选 |
