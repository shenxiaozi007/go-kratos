# 推广引流管理后台 (front-xsh-admin)

与 [xsh-docs](../xsh-docs/) 约定一致的管理后台 SPA，包含选品采集仓、AI 内容加工坊、多账号分发、获客截流四大模块。

## 技术栈

- Vue 3 + Vue Router
- Element Plus
- Axios（接口前缀 `/api/xsh/v1`）

## 开发

```bash
# 安装依赖
npm install

# 启动开发服务（默认端口 5174，/api 代理到后端 8000）
npm run serve
```

访问：http://localhost:5174 ，默认进入 `/xsh/product/list`。

**后端联通**：确保 Kratos 后端已启动（如 `go run ./cmd/realworld` 监听 8000），`vue.config.js` 中已配置 `devServer.proxy` 将 `/api` 转发到后端。

## 路由（/xsh 前缀）

| 模块 | 路由 | 说明 |
|------|------|------|
| 选品 | /xsh/product/parse | 链接解析 |
| 选品 | /xsh/product/list | 商品列表 |
| 选品 | /xsh/product/hot-ranks | 爆款榜单 |
| 内容 | /xsh/content/templates | 模板管理 |
| 内容 | /xsh/content/drafts | 草稿列表 |
| 内容 | /xsh/content/draft/edit/:id? | 草稿编辑 |
| 分发 | /xsh/dispatch/accounts | 账号绑定 |
| 分发 | /xsh/dispatch/schedules | 定时规则 |
| 分发 | /xsh/dispatch/tasks | 发布任务 |
| 获客 | /xsh/inbox/comments | 评论 Inbox |
| 获客 | /xsh/inbox/actions | 动作记录 |

## 构建

```bash
npm run build
```

产出在 `dist/`，部署时需将 `/api` 反向代理到后端服务。
