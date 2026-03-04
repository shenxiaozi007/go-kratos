### 坦克大战前端架构设计（Vue）

本项目前端基于 Vue 3（`front-vue` 目录），采用模块化与分层设计，整体结构建议如下：

```text
front-vue/
  src/
    api/           # 与后端 HTTP 接口交互封装
      auth.js
      game.js
      stats.js
    components/    # 可复用 UI 组件
      ui/
      game/
    views/         # 路由页面
      HomeView.vue
      GameView.vue
      HistoryView.vue
      LeaderboardView.vue
    game/          # 游戏核心逻辑与对象模型
      core/
        GameLoop.ts       # 游戏主循环（初始化、更新、渲染）
        InputManager.ts   # 键盘输入管理
        CollisionSystem.ts# 碰撞检测
      entities/
        Tank.ts
        Bullet.ts
        Obstacle.ts
        Map.ts
      renderer/
        CanvasRenderer.ts # 使用 HTML5 Canvas 渲染
    store/         # 状态管理（可使用简单的组合式 API 或引入 Pinia）
      gameStore.ts
      userStore.ts
    router/
      index.ts     # 定义路由：主页、游戏页、战绩页、排行榜页
    assets/
    main.ts
```

---

### 1. 路由与页面职责划分

- `HomeView.vue`：模式选择（单人 / 本地双人）、地图选择、开始游戏按钮。
- `GameView.vue`：承载 Canvas 画布与游戏 HUD（血量、击杀数、计时器等）。
- `HistoryView.vue`：调用 `/api/v1/stats/matches` 展示历史对局记录。
- `LeaderboardView.vue`：调用 `/api/v1/stats/leaderboard` 展示排行榜。

路由示例：

```ts
// src/router/index.ts
import { createRouter, createWebHistory } from 'vue-router';
import HomeView from '@/views/HomeView.vue';
import GameView from '@/views/GameView.vue';
import HistoryView from '@/views/HistoryView.vue';
import LeaderboardView from '@/views/LeaderboardView.vue';

const routes = [
  { path: '/', name: 'home', component: HomeView },
  { path: '/game', name: 'game', component: GameView },
  { path: '/history', name: 'history', component: HistoryView },
  { path: '/leaderboard', name: 'leaderboard', component: LeaderboardView },
];

export default createRouter({
  history: createWebHistory(),
  routes,
});
```

---

### 2. API 层设计

在 `src/api` 中为每个后端模块提供一个文件，封装 HTTP 请求与返回类型：

- `auth.js`：`login`, `register`, `getProfile`。
- `game.js`：`listMaps`, `getMap`, `getTankConfig`, `createSession`, `finishSession`。
- `stats.js`：`listMatches`, `getLeaderboard`。

示例（伪代码）：

```js
// src/api/game.js
import request from './request'; // 对 axios 的二次封装

export function listMaps() {
  return request.get('/api/v1/game/maps');
}

export function getMap(id) {
  return request.get(`/api/v1/game/maps/${id}`);
}
```

---

### 3. 游戏核心模块分层

#### 3.1 `game/core/GameLoop.ts`

- 负责创建并管理主循环：`init()`, `start()`, `stop()`, `update(dt)`, `render(ctx)`。
- 通过 requestAnimationFrame 调度刷新，维持 60 FPS 目标。

#### 3.2 `game/core/InputManager.ts`

- 统一处理键盘事件，将按键映射为玩家动作（前进、后退、左转、右转、开火）。
- 支持两套按键映射，分别对应玩家 1 与玩家 2。

#### 3.3 `game/core/CollisionSystem.ts`

- 负责坦克与地图格子、坦克之间、子弹与坦克/墙体之间的碰撞检测。
- 对外暴露简单接口：`checkCollisions(state)`，并返回需要处理的碰撞事件列表。

#### 3.4 `game/entities/*`

- `Tank.ts`：封装坦克的状态（位置、朝向、速度、血量等）与更新逻辑。
- `Bullet.ts`：子弹的移动与生命周期。
- `Obstacle.ts`：墙体、河流、草丛等地形元素。
- `Map.ts`：从后端获取的地图数据转换为前端可用于碰撞与渲染的网格结构。

#### 3.5 `game/renderer/CanvasRenderer.ts`

- 负责将当前游戏状态绘制到 Canvas 上。
- 将逻辑层与渲染层解耦，逻辑层只关心坐标与状态，渲染层负责实际的像素绘制。

---

### 4. 状态管理

- `gameStore.ts`：保存当前选中的模式、地图、对局 session_id，以及游戏运行中的统计信息（击毁数、剩余敌人等）。\n- `userStore.ts`：保存当前用户信息与登录状态（如启用账号系统）。\n\n可以使用 Vue 3 组合式 API + `reactive` 简单实现，或引入 Pinia 进行更规范的状态管理。

---

### 5. 与后端的交互流程（示例）

1. 在 `HomeView` 中选择模式和地图，调用 `GameService.CreateSession` 创建对局，得到 `session_id`。\n2. 跳转到 `GameView`，初始化 `GameLoop`、加载地图与配置（调用 `listMaps` / `getMap` / `getTankConfig`）。\n3. 游戏结束后，在前端根据本地计算结果构造 `FinishSessionRequest`，调用 `GameService.FinishSession`。\n4. 统计数据通过 `StatsService` 在 `HistoryView` 和 `LeaderboardView` 中展示。

