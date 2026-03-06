# Proto 与 Message/RPC 大纲

建议在 `api/xsh/v1/` 下按模块拆分为多个 proto 文件，便于维护。包名统一为 `xsh.v1`，`go_package = "realworld/api/xsh/v1;v1"`。

---

## 通用约定

- **时间字段**：统一使用 `google.protobuf.Timestamp`，与现有 realworld 一致；JSON 为 RFC3339 字符串。若项目统一用 int64（Unix 秒），则 xsh 全部时间字段保持一致。

---

## 1. product.proto（选品采集仓）

- **Product**：id, source_url, title, original_price, coupon_amount, commission_rate, image_url, is_hot, rank_source, created_at, updated_at
- **ProductLink**：id, url, platform, parse_status, product_id, summary, created_at, updated_at
- **HotRank**：id, rank_type, synced_at, snapshot_json, created_at
- **ParseProductLinkRequest**：url
- **ParseProductLinkReply**：product (Product), error_message
- **ListProductsRequest**：page, page_size, is_hot, rank_source
- **ListProductsReply**：items (repeated Product), total
- **GetProductRequest**：id
- **GetProductReply**：product
- **CreateProductRequest**：Product 字段（除 id、时间）
- **CreateProductReply**：product
- **UpdateProductRequest**：id + 可更新字段
- **UpdateProductReply**：product
- **ListHotRanksRequest**：rank_type, page, page_size
- **ListHotRanksReply**：items (repeated HotRank), total
- **DeleteProductRequest**：id
- **DeleteProductReply**：success
- **SyncHotRankRequest**：rank_type
- **SyncHotRankReply**：success, message
- **ProductService**：ParseProductLink, ListProducts, GetProduct, CreateProduct, UpdateProduct, DeleteProduct, ListHotRanks, SyncHotRank

---

## 2. content.proto（AI 内容加工坊）

- **Template**：id, name, style_type, content_template, sort_order, created_at, updated_at
- **Draft**：id, product_id, template_id, copy_text, cover_image_url, video_url, status, created_at, updated_at
- **ListTemplatesRequest**：无或 style_type
- **ListTemplatesReply**：items (repeated Template)
- **CreateTemplateRequest**：name, style_type, content_template, sort_order
- **CreateTemplateReply**：template
- **UpdateTemplateRequest**：id + 可更新字段
- **UpdateTemplateReply**：template
- **DeleteTemplateRequest**：id
- **DeleteTemplateReply**：success
- **CreateDraftRequest**：product_id, template_id, options
- **CreateDraftReply**：draft
- **ListDraftsRequest**：page, page_size, status, product_id
- **ListDraftsReply**：items (repeated Draft), total
- **GetDraftRequest**：id
- **GetDraftReply**：draft
- **UpdateDraftRequest**：id, copy_text, cover_image_url, video_url, status
- **UpdateDraftReply**：draft
- **DeleteDraftRequest**：id
- **DeleteDraftReply**：success
- **GenerateCopyRequest**：draft_id, style
- **GenerateCopyReply**：copy_text
- **GenerateCoverRequest**：draft_id
- **GenerateCoverReply**：cover_image_url
- **ProcessVideoRequest**：draft_id, options (抽帧/镜像/滤镜等)
- **ProcessVideoReply**：video_url
- **ContentService**：ListTemplates, CreateTemplate, UpdateTemplate, DeleteTemplate, CreateDraft, ListDrafts, GetDraft, UpdateDraft, DeleteDraft, GenerateCopy, GenerateCover, ProcessVideo

---

## 3. dispatch.proto（多账号分发矩阵）

- **PlatformAccount**：id, platform, nickname, proxy_url, cookie_storage_path, bind_at, created_at, updated_at
- **PublishTask**：id, draft_id, account_id, scheduled_at, published_at, status, platform_content_id, error_message, created_at, updated_at
- **PublishSchedule**：id, name, cron_expr, account_ids (repeated int64 或 JSON 字符串), is_enabled, created_at, updated_at
- **ListPlatformAccountsRequest**：platform (optional), page, page_size
- **ListPlatformAccountsReply**：items (repeated PlatformAccount), total
- **BindAccountRequest**：platform, nickname, proxy_url, cookie_storage_path
- **BindAccountReply**：account
- **UpdateAccountRequest**：id + 可更新字段
- **UpdateAccountReply**：account
- **UnbindAccountRequest**：id
- **UnbindAccountReply**：success
- **ListSchedulesRequest**：无
- **ListSchedulesReply**：items (repeated PublishSchedule)
- **CreateScheduleRequest**：name, cron_expr, account_ids, is_enabled
- **CreateScheduleReply**：schedule
- **UpdateScheduleRequest**：id + 可更新字段
- **UpdateScheduleReply**：schedule
- **DeleteScheduleRequest**：id
- **DeleteScheduleReply**：success
- **CreatePublishTaskRequest**：draft_id, account_ids (repeated), scheduled_at
- **CreatePublishTaskReply**：tasks (repeated PublishTask)，一条请求创建 N 条任务（每 account_id 一条），返回全部新建任务
- **ListPublishTasksRequest**：page, page_size, account_id, status
- **ListPublishTasksReply**：items (repeated PublishTask), total
- **RetryPublishTaskRequest**：id
- **RetryPublishTaskReply**：success
- **CancelPublishTaskRequest**：id
- **CancelPublishTaskReply**：success
- **DispatchService**：ListPlatformAccounts, BindAccount, UpdateAccount, UnbindAccount, ListSchedules, CreateSchedule, UpdateSchedule, DeleteSchedule, CreatePublishTask, ListPublishTasks, RetryPublishTask, CancelPublishTask

---

## 4. inbox.proto（获客截流）

- **Comment**：id, account_id, platform, platform_content_id, platform_comment_id, author_nickname, content, want_link, status, created_at, updated_at
- **InboxAction**：id, comment_id, action_type, content, sent_at, result_status, error_message, created_at
- **ListCommentsRequest**：account_id, want_link, status, page, page_size
- **ListCommentsReply**：items (repeated Comment), total
- **MarkCommentHandledRequest**：id, status (replied/ignored)
- **MarkCommentHandledReply**：success
- **SyncCommentsRequest**：account_id
- **SyncCommentsReply**：success, synced_count, message
- **SendReplyRequest**：comment_id, content
- **SendReplyReply**：success, action_id, error_message
- **SendPrivateMessageRequest**：comment_id, content
- **SendPrivateMessageReply**：success, action_id, error_message
- **ListInboxActionsRequest**：comment_id, page, page_size
- **ListInboxActionsReply**：items (repeated InboxAction), total
- **InboxService**：ListComments, MarkCommentHandled, SyncComments, SendReply, SendPrivateMessage, ListInboxActions

---

## 生成与路由

- 使用 Kratos 的 proto 生成：`kratos proto client/api` 或项目现有 Makefile 生成 `*_http.pb.go`、`*_grpc.pb.go`，并在 HTTP 路由中注册各 Service。
- 所有 HTTP 路径建议带 `/api/xsh/v1/` 前缀，与 [api-catalog.md](api-catalog.md) 一致。
