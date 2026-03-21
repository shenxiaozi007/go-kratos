# 萌宠之家 - 前端

云养宠伴侣 Web 端，React + Vite + Tailwind CSS。接口与设计见 `docs-all/cat-admin/`。

## 开发

```bash
cd realworld/front-pet-app
npm install
npm run dev
```

默认端口 5174，API 通过 Vite 代理到 `http://localhost:8000`（需先启动 realworld 后端）。

## 构建

```bash
npm run build
npm run preview
```

## 配置

- 接口地址：开发时使用 Vite 代理 `/api` → `http://localhost:8000`；生产可在 `index.html` 前注入 `window.APP_CONFIG = { apiBaseUrl: 'https://your-api.com' }`。
- 登录态：Token 存于 `localStorage.pet_app_token`，请求头自动带 `Authorization: Bearer <token>`。

## 路由

| 路径 | 页面 |
|------|------|
| / | 首页（宠物卡片、签到/商店/领养入口） |
| /interaction | 互动页 |
| /social | 社交（好友/排行榜/申请） |
| /profile | 我的 |
| /login | 登录/注册 |
| /shop | 商店 |
| /inventory | 背包 |
| /adoption | 领养（猫/狗品种列表） |
| /breed-detail/:id | 品种详情与领养 |
| /edit-profile | 编辑资料 |
| /purchase-success | 购买成功 |
| /achievements | 成就 |
| /dressup | 换装 |
| /ranking | 排行榜 |
| /requests | 好友申请列表 |

## 设计

- 字体：Plus Jakarta Sans
- 主色：`#ef3957`（primary）
- 背景：`#f8f6f6`；玻璃效果 `backdrop-blur`
- 最大宽度 448px 居中，底部 4 Tab 导航
