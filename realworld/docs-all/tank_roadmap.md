### 坦克大战小游戏实施里程碑与工作量预估

本文件将高层计划拆解为具体迭代里程碑，并给出相对工作量预估（以「人日」为粗略单位，供参考）。

---

### 迭代 0：环境与基础设施（约 0.5–1 人日）

- 确认后端 Kratos 工程结构、配置 MySQL 连接（已存在 `internal/conf/conf.proto` 与 `configs/config.yaml`）。\n- 通过 `make api` 生成 `auth/game/stats/todo` 等 proto 对应的 Go 代码。\n- 验证 HTTP 与 gRPC 服务均可启动，健康检查通过。\n- 前端 `front-vue` 可正常 `npm install`、`npm run serve` 启动空壳页面。

**交付物**：\n- 可启动的后端服务与前端骨架工程。

---

### 迭代 1：后端配置与数据模型（约 2–3 人日）

- 实现并验证数据库表结构迁移：\n  - `users`, `maps`, `map_tiles`, `game_sessions`, `match_records`, `list`。\n- 在 `internal/data` 实现对应的 GORM 实体（已完成）与基础仓储接口：\n  - 地图读取：加载地图列表与明细。\n  - 对局创建/更新：写入 `game_sessions`。\n  - 战绩写入：写入 `match_records`。\n- 在 `internal/biz` 实现对应 UseCase：\n  - `GameUsecase`：地图获取、对局创建/结束。\n  - `StatsUsecase`：战绩与排行榜查询。\n- 在 `internal/service` 实现 `GameService`、`StatsService`、`AuthService` 的骨架逻辑，并与 Kratos HTTP Server 完整集成。

**交付物**：\n- RESTful API：地图查询、对局创建与结束、战绩与排行榜接口均可通过 Postman 调用。

---

### 迭代 2：前端游戏核心原型（约 3–5 人日）

- 在 `front-vue/src/game` 下实现核心游戏循环与对象模型（参考 `docs/gameplay_core_loop.md`）：\n  - `GameLoop`、`InputManager`、`CollisionSystem`、`Tank`、`Bullet`、`Map` 等。\n  - 使用 Canvas 完成基础渲染（坦克移动、发射子弹、子弹命中销毁）。\n- 在 `HomeView` 与 `GameView` 之间完成参数传递（模式、地图 ID、本地玩家数量）。\n- 游戏结束后，将本地统计数据通过 `GameService.FinishSession` 上报后端。

**交付物**：\n- 可玩的基本单人模式原型（无复杂 UI 与动画，但完整闭环：开始 → 战斗 → 结束 → 记录结果）。

---

### 迭代 3：本地双人模式与 UI 强化（约 2–4 人日）

- 扩展 `InputManager` 以支持两套按键映射，实现本地双人对战模式。\n- 在 `GameView` 中完善 HUD：分别展示玩家 1、玩家 2 的血量与击杀统计。\n- 增加暂停、重新开始、返回主页等基础交互。\n- 对初版地图与坦克属性进行调优，保证体验流畅。

**交付物**：\n- 可玩的本地双人模式，基础 UI 完整。

---

### 迭代 4：战绩与排行榜页面（约 1.5–3 人日）

- 后端：完善 `StatsService` 的分页、排序与聚合逻辑（若尚未完成）。\n- 前端：实现 `HistoryView` 与 `LeaderboardView` 页面：\n  - 分页展示历史对局记录（模式、地图、时长、胜负、击毁数等）。\n  - 展示排行榜（总胜场、胜率等），支持简单过滤或排序。\n- 若启用账号系统：在登录态下区分用户的个人战绩与全局排行榜。

**交付物**：\n- 可用的战绩与排行榜页面，前后端打通。

---

### 迭代 5：WebSocket 与体验优化（可选，约 2–4 人日）

- 在 Kratos HTTP Server 中增加 WebSocket 端点 `/ws/game/session/{id}`，用于：\n  - 推送对局结束通知。\n  - 预留未来在线 PVP 的消息通道（状态同步协议可后续定义）。\n- 前端：在 `GameView` 中接入 WebSocket 客户端，监听并响应来自服务端的简单事件。\n- 性能优化：在较多 AI 敌人的情况下，优化碰撞检测与渲染逻辑。

**交付物**：\n- 初步连通的 WebSocket 通道与更平滑的游戏体验。

---

### 迭代 6：重构与文档（约 1–2 人日）

- 清理与重构关键模块代码，确保 `internal/biz` 与前端 `game` 模块的职责边界清晰。\n- 补充 README：\n  - 如何启动后端与前端。\n  - 如何初始化数据库（参考 `docs/tank_db_schema.md`）。\n  - 简要的架构说明与后续扩展方向。

**交付物**：\n- 更可维护的代码库与完整的项目文档。

