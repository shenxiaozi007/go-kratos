# 其他未完善功能与 v2 处理方式

本文档列出 v1 中未完善的功能项，并标明在 v2 中采用「实现」「文档约定」或「延后」的处理方式，供排期与后续迭代参考。

---

## 一、封面合成（GenerateCover）

- **v1 状态**：Noop；[CopyGenerator] 已有接口，封面为 [CoverGenerator](internal/biz/xsh_content.go)。
- **v2 处理**：**文档约定**为主，实现标为 **v2 可选** 或 **v2 后续迭代**。
- **约定内容**：
  - 输入：draftID、商品图 URL（或本地路径）；输出：合成后的封面图 URL 或本地路径。
  - 存储：本地路径或上传 OSS 后返回 URL；若为本地路径，需约定相对目录与访问方式（如静态服务或 CDN）。
  - 逻辑：商品图 + 价格爆炸贴 + 平台 Logo 等元素合成（见 xsh-docs 业务描述）；具体实现可用图片库或调用第三方服务。

---

## 二、视频去重（ProcessVideo）

- **v1 状态**：Noop；[VideoProcessor](internal/biz/xsh_content.go) 接口已定义。
- **v2 处理**：**文档约定**为主，实现标为 **v2 可选** 或 **后续版本**。
- **约定内容**：
  - 输入：draftID、原片 URL/路径、可选参数（抽帧、镜像、滤镜等）；输出：处理后视频 URL 或路径。
  - 方案：FFmpeg 本地处理，或对接第三方视频处理服务；需约定输出格式、存储位置与访问方式。

---

## 三、评论同步（SyncComments）

- **v1 状态**：占位；无 Playwright 拉评、无 want_link 写入流水。
- **v2 处理**：**文档约定**为 **后续阶段**；v2 仅保证 Playwright 能**登录并抓取用户信息**（见 [03-playwright-xiaohongshu.md](03-playwright-xiaohongshu.md)）。
- **约定内容**：
  - 后续阶段：通过 Playwright 访问各账号评论列表/消息页，拉取评论并落库；AI 识别「求链接」等意图并标记或写入 Inbox 动作；与「登录 + 抓取用户信息」的 Worker/脚本边界分离。

---

## 四、发布任务执行（Worker 消费任务发帖）

- **v1 状态**：无 Worker 消费 `xsh_publish_tasks` 执行实际发帖。
- **v2 处理**：**文档约定**与「登录/抓取用户信息」边界；实现列在 **v2 后续** 或 **v3**。
- **约定内容**：
  - 发布 Worker：独立进程或任务队列消费者，读取 `xsh_publish_tasks` 中待执行任务，按账号使用 Playwright（含 proxy、Cookie）执行小红书/抖音发帖；与「模拟登录 + 抓取用户信息」的触发方式与执行主体区分，避免混在同一流程。
  - 任务状态：pending → running → success/failed；失败可重试（沿用 v1 重试接口与状态）。

---

## 五、定时规则与自动创建任务

- **v1 状态**：未实现 cron 与「按定时规则自动创建发布任务」的联动。
- **v2 处理**：**文档约定**方案；实现可标为 **v2 可选**。
- **约定内容**：
  - 方案 A：系统 cron 定时调用 Kratos 提供的「按规则创建任务」接口（如根据 `xsh_publish_schedules` 的 cron_expr、account_ids 在对应时间点创建一批 `xsh_publish_tasks`）。
  - 方案 B：独立 Worker 内嵌调度器，按规则表轮询并创建任务。
  - 需约定：规则启用/禁用、时区、与发布 Worker 的协作方式。

---

## 六、前端

- **登录与 Token、权限控制**：在 [02-auth-and-permissions.md](02-auth-and-permissions.md) 中已覆盖（登录页、Token 存储与请求头、按角色隐藏/禁用菜单或按钮）。
- **错误码与提示**：文案生成失败、限流、超时等在 [04-gemini-copywriting.md](04-gemini-copywriting.md) 已约定；其他通用错误码与前端提示在 02 与现有 [xsh-docs/backend/errors-and-conventions](../xsh-docs/backend/errors-and-conventions.md) 中统一即可，05 不再重复。

---

## 七、优先级建议（供排期）

| 项           | v2 建议       | 说明                         |
|--------------|----------------|------------------------------|
| 鉴权与权限   | 必做           | 见 01、02                    |
| Gemini 文案  | 必做           | 见 04                        |
| Playwright 登录与用户信息 | 必做 | 见 03                        |
| 封面合成     | 可选/后续      | 文档已约定，按资源排期       |
| 视频去重     | 可选/后续      | 文档已约定                   |
| 评论同步     | 后续           | 依赖登录与 Cookie 复用       |
| 发布 Worker  | 后续/v3        | 与登录流程分离               |
| 定时创建任务 | v2 可选        | 文档已约定方案               |
