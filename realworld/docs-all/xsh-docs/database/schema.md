# 表结构清单与字段说明

以下表均使用前缀 `xsh_`，字符集建议 `utf8mb4`。

---

## 选品采集仓

### xsh_products

解析后的商品信息（首期可来自手动录入或简单解析，不依赖淘宝联盟）。

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| id | BIGINT UNSIGNED AUTO_INCREMENT | 主键 | PK |
| source_url | VARCHAR(1024) | 来源链接 | 可选 INDEX |
| title | VARCHAR(512) NOT NULL | 商品标题 | |
| original_price | DECIMAL(10,2) | 原价 | |
| coupon_amount | DECIMAL(10,2) DEFAULT 0 | 优惠券面额 | |
| commission_rate | DECIMAL(5,2) | 佣金比例（百分比，可选） | |
| image_url | VARCHAR(1024) | 主图 URL | |
| is_hot | TINYINT(1) DEFAULT 0 | 是否神价/爆款 | INDEX |
| rank_source | VARCHAR(64) | 来源榜单类型（可选） | |
| created_at | DATETIME | 创建时间 | INDEX |
| updated_at | DATETIME | 更新时间 | |

---

### xsh_hot_ranks

爆款榜单快照（首期可不接联盟 API，可留空或手动导入）。

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| id | BIGINT UNSIGNED AUTO_INCREMENT | 主键 | PK |
| rank_type | VARCHAR(64) NOT NULL | 榜单类型 | INDEX |
| synced_at | DATETIME NOT NULL | 同步时间 | INDEX |
| snapshot_json | JSON/TEXT | 榜单原始数据或 product_id 列表 | |
| created_at | DATETIME | 创建时间 | |

---

### xsh_product_links

原始链接与解析记录。

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| id | BIGINT UNSIGNED AUTO_INCREMENT | 主键 | PK |
| url | VARCHAR(1024) NOT NULL | 原始链接 | INDEX |
| platform | VARCHAR(32) | 平台（如 taobao、tmall） | |
| parse_status | VARCHAR(32) | pending/success/failed | INDEX |
| product_id | BIGINT UNSIGNED | 解析成功后关联 xsh_products.id | FK, INDEX |
| summary | VARCHAR(512) | 解析结果摘要或错误信息 | |
| created_at | DATETIME | 创建时间 | |
| updated_at | DATETIME | 更新时间 | |

---

## AI 内容加工坊

### xsh_content_templates

文案模板（如暴躁省钱风、闺蜜安利风、超市比价风等）。

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| id | BIGINT UNSIGNED AUTO_INCREMENT | 主键 | PK |
| name | VARCHAR(128) NOT NULL | 模板名称 | |
| style_type | VARCHAR(64) | 风格类型 | INDEX |
| content_template | TEXT | 模板内容（占位符如 {{title}}、{{price}}） | |
| sort_order | INT DEFAULT 0 | 排序 | |
| created_at | DATETIME | 创建时间 | |
| updated_at | DATETIME | 更新时间 | |

---

### xsh_content_drafts

内容草稿（关联商品、模板，存文案/封面/视频）。

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| id | BIGINT UNSIGNED AUTO_INCREMENT | 主键 | PK |
| product_id | BIGINT UNSIGNED NOT NULL | 关联 xsh_products.id | FK, INDEX |
| template_id | BIGINT UNSIGNED | 关联 xsh_content_templates.id | FK, INDEX |
| copy_text | TEXT | 文案内容 | |
| cover_image_url | VARCHAR(1024) | 封面图 URL | |
| video_url | VARCHAR(1024) | 视频 URL（可选） | |
| status | VARCHAR(32) | draft/ready/published | INDEX |
| created_at | DATETIME | 创建时间 | |
| updated_at | DATETIME | 更新时间 | |

**status 流转**：draft → ready → published。CreatePublishTask 仅接受 status=ready 的草稿（或按需允许 draft）；已 published 可编辑，但不建议重复创建发布任务。

---

## 分发矩阵

### xsh_platform_accounts

平台账号；**每账号单独 proxy 配置，用于 Playwright 出网，防封（一账号一 IP）**。

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| id | BIGINT UNSIGNED AUTO_INCREMENT | 主键 | PK |
| platform | VARCHAR(32) NOT NULL | 平台：douyin / xiaohongshu | INDEX |
| nickname | VARCHAR(128) | 昵称/备注 | |
| proxy_url | VARCHAR(512) | 代理 URL（如 http://user:pass@host:port），该账号发布与拉评均用此 IP | |
| cookie_storage_path | VARCHAR(512) | 登录态存储路径（Cookie/Storage 持久化） | |
| bind_at | DATETIME | 绑定时间 | |
| created_at | DATETIME | 创建时间 | |
| updated_at | DATETIME | 更新时间 | |

---

### xsh_publish_tasks

发布任务（关联草稿与账号，计划时间与执行状态）。

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| id | BIGINT UNSIGNED AUTO_INCREMENT | 主键 | PK |
| draft_id | BIGINT UNSIGNED NOT NULL | 关联 xsh_content_drafts.id | FK, INDEX |
| account_id | BIGINT UNSIGNED NOT NULL | 关联 xsh_platform_accounts.id | FK, INDEX |
| scheduled_at | DATETIME NOT NULL | 计划发布时间 | INDEX |
| published_at | DATETIME | 实际发布时间 | |
| status | VARCHAR(32) | pending/running/success/failed/cancelled | INDEX |
| platform_content_id | VARCHAR(256) | 平台侧内容 ID（发布成功后回写） | |
| error_message | VARCHAR(512) | 失败原因 | |
| created_at | DATETIME | 创建时间 | |
| updated_at | DATETIME | 更新时间 | |

---

### xsh_publish_schedules

定时规则（如每天 10:00、12:00、20:00）。

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| id | BIGINT UNSIGNED AUTO_INCREMENT | 主键 | PK |
| name | VARCHAR(128) | 规则名称 | |
| cron_expr | VARCHAR(128) | cron 表达式或时间点配置 | |
| account_ids | JSON/VARCHAR(512) | 关联账号 ID 列表，空则表示全局 | |
| is_enabled | TINYINT(1) DEFAULT 1 | 是否启用 | |
| created_at | DATETIME | 创建时间 | |
| updated_at | DATETIME | 更新时间 | |

---

## 获客截流

### xsh_comments

评论汇总（从各账号拉取后入库）。

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| id | BIGINT UNSIGNED AUTO_INCREMENT | 主键 | PK |
| account_id | BIGINT UNSIGNED NOT NULL | 关联 xsh_platform_accounts.id | FK, INDEX |
| platform | VARCHAR(32) NOT NULL | 平台 | INDEX |
| platform_content_id | VARCHAR(256) | 平台侧内容 ID | INDEX |
| platform_comment_id | VARCHAR(256) | 平台侧评论 ID | UNIQUE |
| author_nickname | VARCHAR(128) | 评论者昵称 | |
| content | TEXT | 评论内容 | |
| want_link | TINYINT(1) DEFAULT 0 | 是否「求链接」等意图（AI 或规则标记） | INDEX |
| status | VARCHAR(32) | pending/replied/ignored | INDEX |
| created_at | DATETIME | 拉取时间 | INDEX |
| updated_at | DATETIME | 更新时间 | |

**同步去重策略**：同步评论时按 platform_comment_id 做「存在则更新、不存在则插入」，避免 UNIQUE 冲突；或仅插入不更新，由业务保证同一评论不重复拉取。

---

### xsh_inbox_actions

私信/回复动作记录。

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| id | BIGINT UNSIGNED AUTO_INCREMENT | 主键 | PK |
| comment_id | BIGINT UNSIGNED NOT NULL | 关联 xsh_comments.id | FK, INDEX |
| action_type | VARCHAR(32) NOT NULL | reply / private_message | |
| content | TEXT | 发送内容 | |
| sent_at | DATETIME | 发送时间 | |
| result_status | VARCHAR(32) | success/failed | |
| error_message | VARCHAR(512) | 失败原因 | |
| created_at | DATETIME | 创建时间 | |
