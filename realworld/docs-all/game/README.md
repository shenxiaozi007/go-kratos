# 防守割草 · 微信小程序小游戏文档

**防守割草**是一款防守类割草微信小程序小游戏：玩家在屏幕下方防守，捣蛋小怪兽从上方涌来，玩家自动发射星星子弹净化怪兽，积累能量升级并选择随机技能强化自己。前后端分离：前端负责游戏逻辑与 UI，后端负责数据存档与高并发交互。

---

## 文档列表

| 文档 | 用途 |
|------|------|
| [01-project-overview.md](01-project-overview.md) | 项目开发文档：原型、界面流转、核心体验、技术架构 |
| [02-backend.md](02-backend.md) | 后端开发文档：API 设计、数据模型、服务分层、目录约定 |
| [03-frontend.md](03-frontend.md) | 前端开发文档：页面与路由、战斗 Canvas、技能系统、与后端对接 |
| [04-development-plan.md](04-development-plan.md) | 开发规划：阶段划分、里程碑、依赖关系 |
| [05-test-plan.md](05-test-plan.md) | 测试规划：策略、用例、环境与性能 |

---

## 技术栈速览

| 层级 | 技术 |
|------|------|
| 前端 | Vue 3、微信小程序（uni-app Vue3 或 Taro Vue3）、Canvas 2D、Pinia |
| 后端 | Go、Kratos v2、GORM、MySQL、Redis（排行榜） |
| 对接 | REST/JSON，错误与分页约定见 [xsh-docs/backend/errors-and-conventions.md](../xsh-docs/backend/errors-and-conventions.md) |

---

## 仓库内相关代码位置

- **API 定义**：`api/game/v1/` — 现有坦克游戏 proto；防守割草扩展见 [02-backend.md](02-backend.md)（建议新增 `defense.proto` 或独立 DefenseService）。
- **后端实现**：`internal/biz`、`internal/data`、`internal/service` 下 game 相关文件（如 `game_score.go`、`game_leaderboard.go`）。
- **前端项目**：建议独立目录（如 `front-game-defense` 或置于现有小程序工程内），页面与组件结构见 [03-frontend.md](03-frontend.md)。
