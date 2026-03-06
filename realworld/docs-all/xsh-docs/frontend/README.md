# 前端设计总览

推广引流后台为独立**管理后台 SPA**，可与现有 [front-vue](../../front-vue/) 同仓库不同路由前缀或入口，或新建目录 `front-xsh-admin`（二选一在项目中约定）。本文档以**新建 `front-xsh-admin`** 为例，便于与游戏前端隔离。

---

## 技术栈建议

| 类别 | 技术 |
|------|------|
| 框架 | Vue 3、Vue Router、Pinia |
| UI 库 | Element Plus 或 Naive UI（在项目中择一） |
| 请求 | axios，封装 baseURL 为 Kratos 后端（如 `/api/xsh/v1`） |
| 构建 | Vite（推荐）或 Vue CLI |

---

## 目录结构（front-xsh-admin）

```text
front-xsh-admin/
  src/
    api/              # 按模块封装后端接口
      product.js
      content.js
      dispatch.js
      inbox.js
    views/             # 页面
      product/         # 选品采集仓
      content/         # AI 内容加工坊
      dispatch/        # 多账号分发矩阵
      inbox/           # 获客截流
    components/        # 可复用组件
      common/
      product/
      content/
      dispatch/
      inbox/
    store/             # Pinia stores
    router/
      index.js
    assets/
    App.vue
    main.js
  package.json
  vite.config.js
```

---

## 与后端的对接

- 所有接口前缀：`/api/xsh/v1/`（与 [backend/api-catalog.md](../backend/api-catalog.md) 一致）。
- **响应约定**：列表接口响应与 proto 生成 JSON 一致，为 `items` + `total`；错误体与 [backend/errors-and-conventions.md](../backend/errors-and-conventions.md) 一致（如 `code`、`message`），前端可统一解析并提示。
- 文案生成、封面合成等由后端调用 **Gemini 3 Flash** 与内部逻辑完成，前端仅调用对应 API 展示结果。
- 认证：若后续加登录，可在 axios 拦截器中加 Token；首期可为内网或简单鉴权。
- 错误与 loading：统一在 axios 拦截器或各 api 模块内处理，页面只消费业务数据。

详见 [views-routes.md](views-routes.md) 与 [components.md](components.md)。
