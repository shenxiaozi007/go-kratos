# Playwright-go 接入：小红书登录与抓取用户信息

本文档描述 v2 中通过 Playwright-go 模拟登录小红书网页，并在登录成功后抓取基本用户信息、落库与 Cookie 复用的设计。v2 不包含发帖、评论拉取、发布任务消费 Worker。

---

## 一、目标

- 使用 Playwright-go 打开小红书登录页，完成模拟登录。
- 登录成功后进入可获取用户信息的页面，解析 DOM 或接口拿到「基本用户信息」。
- 将抓取到的信息写入数据库，并持久化 Cookie 供后续复用（发帖、拉评等留待后续阶段）。

---

## 二、登录方式约定

- **主流程**：优先支持**扫码登录**（小红书创作者/网页端常见方式）。流程：打开登录页 → 等待用户扫码 → 检测登录成功（如 URL 跳转或特定元素出现）→ 继续抓取用户信息。
- **可选扩展**：账号密码登录（若平台开放且稳定）；若支持，在本文档或附录中单独描述流程与选择逻辑（如配置项 `login_mode: qrcode | password`）。

---

## 三、流程概述

1. **打开登录页**：使用约定域名，如 `https://creator.xiaohongshu.com` 或小红书网页版登录页 URL（以实际可用的创作者/登录入口为准）。
2. **注入 Proxy**：若账号已绑定 `xsh_platform_accounts.proxy_url`，启动 Playwright 时按该 proxy 配置出网，实现一账号一 IP。
3. **执行登录**：
   - 扫码：展示二维码或跳转至扫码页，等待用户扫码；轮询或监听页面变化以检测登录成功。
   - 密码（若实现）：填充账号密码并提交，处理验证码（若有）。
4. **登录成功判定**：通过 URL 跳转至创作者后台/个人中心，或页面出现用户昵称/头像等元素。
5. **抓取用户信息**：进入可获取基本信息的页面（如创作者中心首页、设置页或接口），解析 DOM 或拦截/请求接口，提取下表字段。
6. **落库**：将抓取结果写入 `xsh_platform_accounts` 扩展字段或关联表，见 [database/schema-v2-changes.md](database/schema-v2-changes.md)。
7. **Cookie/登录态持久化**：将当前浏览器上下文中的 Cookie（及可选 Storage）写入 `cookie_storage_path` 指定路径，供下次「免登录」复用。

---

## 四、抓取的基本用户信息字段

| 字段 | 说明 | 来源（示例） |
|------|------|--------------|
| platform_user_id | 小红书侧用户唯一 ID | 接口返回或 DOM 上的 data 属性 |
| platform_nickname | 昵称 | 页面展示或接口 |
| platform_avatar | 头像 URL | 头像图片 src 或接口 |
| platform_red_id | 小红书号（如有） | 个人主页或接口 |
| logged_at | 本次登录成功时间（Unix 秒或 datetime） | 服务端记录 |

以上字段在 v2 的存储方案见 [database/schema-v2-changes.md](database/schema-v2-changes.md)（方案 A：在 `xsh_platform_accounts` 增加字段）。

---

## 五、执行主体

- **方式**：可为「管理端触发的一次性任务」或「独立小脚本/Worker」；由后台提供接口（如「为某账号执行登录并抓取用户信息」），触发后同步或异步执行。
- **与发布 Worker 区分**：v2 仅实现「登录 + 抓取用户信息 + 回写」；发帖、评论同步、定时发布等由后续阶段实现，不在本阶段混入同一 Worker。

---

## 六、一账号一 IP 与 Cookie

- **Proxy**：沿用 v1 设计，使用 `xsh_platform_accounts.proxy_url`；Playwright 启动 Browser 或 BrowserContext 时按该账号配置 proxy，确保该账号所有请求经固定 IP。
- **Cookie 持久化**：登录成功后，将 context 的 Cookie 写入 `cookie_storage_path`（如 JSON 或 Playwright 支持的 storage state 格式）；下次「免登录」时先加载该路径的 Cookie/state 再访问页面，减少重复扫码。

---

## 七、错误与重试

- 登录超时（如长时间未扫码）、网络错误、页面结构变化导致解析失败：返回明确错误码与信息，便于管理端展示或重试。
- 可选：同一账号短时间内重复触发「登录并抓取」时，若已有有效 Cookie，可先尝试复用再决定是否重新扫码。

---

## 八、后续阶段（不在此版本实现）

- **评论同步（SyncComments）**：通过 Playwright 拉取各账号评论、识别「求链接」等并写入 Inbox，在 v2 后续或 v3 实现；本阶段仅保证能登录并抓取用户信息。
- **发布任务执行**：Worker 消费 `xsh_publish_tasks`、用 Playwright 执行发帖，与「登录/抓取用户信息」边界分离，见 [05-remaining-improvements.md](05-remaining-improvements.md)。
