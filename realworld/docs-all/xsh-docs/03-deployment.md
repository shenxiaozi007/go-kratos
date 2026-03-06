# 部署与配置

## 鉴权

- **首期**：无登录或内网/IP 白名单访问；若管理后台仅内网可访问，可不做鉴权。
- **后续**：可加 JWT 或 Session，与现有 realworld middleware 一致。

---

## 后端

- 沿用现有 realworld 的启动方式：`go run ./cmd/realworld` 或编译后运行，配置文件为 `configs/config.yaml`。
- HTTP/GRPC 端口与现有一致（如 HTTP 8000、GRPC 9000）；xsh 相关 API 挂载在同一 Server 下，前缀 `/api/xsh/v1/`。
- 数据库：与现有项目共用 MySQL（如 `kratos-blog`），或单独库 `xsh_admin`（需在配置中增加数据源并在 data 层使用）。Redis 沿用现有配置，若 xsh 需缓存可复用同一实例。

## 前端

- 管理后台（front-xsh-admin）构建：`npm run build`，产物部署至 Nginx 或与现有 front-vue 同站不同 path（如 `/xsh-admin/`）。
- 开发环境：`npm run dev`，通过代理将 `/api` 转发到 Kratos 后端，避免跨域。

## 环境变量与配置项（占位）

可在 `configs/config.yaml` 或环境变量中预留/扩展以下项（具体 key 以项目为准）：

| 配置项 | 说明 |
|--------|------|
| server.http.addr | HTTP 监听地址 |
| data.database.* | MySQL 连接（driver、source、host、port、database） |
| data.redis.addr | Redis 地址 |
| xsh.playwright.headless | 是否无头模式（可选） |
| xsh.playwright.browser_path | Chromium 可执行路径（可选） |
| xsh.ai.api_key | Gemini 3 Flash API Key（用于内容加工坊文案生成） |
| xsh.ai.model / xsh.ai.endpoint | Gemini 模型 ID 或 API 端点（可选，默认由 SDK/官方约定） |

## 配置示例（YAML 片段）

Worker 独立部署时读取同一 config 或环境变量，示例：

```yaml
xsh:
  playwright:
    headless: true
    browser_path: ""  # 可选，Chromium 路径
  ai:
    api_key: "your-gemini-api-key"
    model: "gemini-2.0-flash"  # 或实际模型 ID
    # endpoint: ""  # 可选
```

## 运行依赖

- **封面合成**：若用本地制图（如 ImageMagick、Canvas 等），需在部署环境安装对应依赖并配置路径。
- **视频去重**：若用 FFmpeg，需安装 FFmpeg 并配置可执行路径；若用 OSS 或第三方处理服务，需配置 bucket/endpoint 等（占位，后续接 OSS 时补充）。

## 部署与防封：一账号一 IP

- **Playwright 发布 Worker** 需为每个平台账号配置**独立代理或独立出口 IP**，与 [02-architecture.md](02-architecture.md) 一致。
- **实现方式**：
  - **单机多代理**：在 `xsh_platform_accounts` 中为每个账号配置不同 `proxy_url`（如 HTTP/SOCKS5），Worker 执行任务时按 `account_id` 读取该账号的 proxy，在 Playwright 的 `BrowserType.Launch` 或 `NewContext` 中传入 `Proxy: account.ProxyURL`。
  - **多机多 IP**：每台机器使用不同出口 IP，将账号与机器绑定（或在账号表增加 machine_id/agent_id），由该机器上的 Worker 只处理绑定到本机的账号任务。
- **原则**：同一 IP 不同时用于多个平台账号的发布与评论拉取，避免平台限流或封号。在运维文档中说明 proxy 配置方式与账号-IP 绑定原则，并定期检查代理可用性。
