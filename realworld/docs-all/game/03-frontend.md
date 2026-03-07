# 前端开发文档：防守割草小游戏

本文档约定微信小程序前端的技术选型、页面与路由、战斗场景 Canvas、技能系统及与后端的对接方式。

---

## 技术选型

| 类别 | 技术 |
|------|------|
| 框架 | Vue 3、Vue Router、Pinia |
| 小程序形态 | 微信小程序，通过 **uni-app（Vue3）** 或 **Taro（Vue3）** 编译输出 |
| 游戏渲染 | Canvas 2D（小程序内 `<canvas type="2d">`） |
| 请求 | 封装 HTTP 客户端，baseURL 指向 Kratos 后端（如 `/api/v1/game/defense`） |

**选型说明**：uni-app 与 Taro 均支持 Vue3 与微信小程序；可根据团队熟悉度、Canvas 兼容性及生态选一。Canvas 2D 需使用各框架推荐的小程序 Canvas API（如 `wx.createCanvasContext` 或新版 2D 画布），保证帧率与触摸响应。

---

## 页面与路由

与 [01-project-overview.md](01-project-overview.md) 中的流转一致：

| 页面 | 路径示例 | 说明 |
|------|----------|------|
| 登录/闪屏 | `/pages/login/index` 或首页 | 微信授权或展示 Logo 后跳主菜单 |
| 主菜单 | `/pages/index/index` 或 `/pages/menu/index` | Logo、开始游戏、排行榜、当前最高分 |
| 战斗场景 | `/pages/battle/index` | HUD + Canvas 游戏区 + 升级弹窗 |
| 结算页 | `/pages/result/index` 或 `/pages/game-over/index` | 得分、波数、再玩一次、返回主页 |

路由命名与框架约定保持一致（如 uni-app 的 `pages.json`、Taro 的 `app.config`）。

---

## 战斗场景

### Canvas 尺寸与适配

- 使用 rpx 或设计稿比例换算为逻辑宽高（如 750rpx 宽、固定宽高比），在 `onReady` 中获取 Canvas 节点并设置 `width`/`height`，避免拉伸模糊。
- 游戏坐标系统：以 Canvas 左上角为原点，玩家固定在底部中央（如 `canvasHeight - 80`），怪兽生成在顶部随机 x（如 `0 ~ canvasWidth`）。

### 游戏循环

- 使用 `requestAnimationFrame`（或小程序等效 API，如 `canvas.requestAnimationFrame`）驱动主循环；每帧更新：
  - 玩家（位置固定，仅表现可做轻微动效）
  - 子弹（星星）位置、与怪兽碰撞检测
  - 怪兽位置（向下移动）、生成新怪兽（按波次/时间间隔）
- 碰撞检测：子弹与怪兽 AABB 或半径检测，命中后移除子弹与怪兽、增加分数与能量。

### 实体与逻辑

- **玩家**：底部中央，自动发射星星子弹（定时或按固定间隔生成子弹实体）。
- **子弹**：从玩家位置向上（或按技能多方向）运动，速度与伤害由配置或技能决定。
- **怪兽**：从顶部随机 x 生成，向下移动；血量或碰撞次数由配置决定，净化后加分与能量。
- **能量与等级**：击败怪兽增加能量值；能量条满则触发升级，暂停循环并弹出升级面板。

### 升级面板

- 游戏循环暂停（停止 requestAnimationFrame 或设 `paused = true`）。
- 展示 3 张技能卡牌（从技能池随机 3 个，不重复），点击其一后：
  - 应用技能效果（如增加弹道数、星星弹射、伤害提升等），写入本地状态或 Pinia。
  - 关闭弹窗，清空或部分清空能量条，恢复游戏循环。

---

## 技能系统

- **技能数据**：技能 ID、名称、描述、效果类型（如 `multi_shot`、`star_bounce`、`damage_up`）。
- **效果实现**：在前端游戏逻辑中根据技能 ID 修改发射逻辑（子弹数量、角度）、碰撞后行为（弹射）、伤害系数等；可维护一份「技能配置表」在前端，或由后端下发可选技能列表（见 [02-backend.md](02-backend.md) 可选「技能列表接口」）。
- **随机与去重**：每次升级从全技能池中随机 3 个（或排除已选），保证单局内可选技能丰富且不重复。

---

## 与后端对接

| 场景 | 接口 | 时机 |
|------|------|------|
| 主菜单展示最高分 | GetMyBest | 进入主菜单时调用，展示 best_score、best_wave、rank |
| 全服排行榜 | GetLeaderboard | 点击「全服排行榜」时请求，分页展示 items |
| 提交本局分数 | SubmitScore | 进入结算页时自动调用，传 score、wave、user_id/openid |

- **baseURL**：与项目环境一致（如开发 `/api/v1/game/defense`，生产同域名）。
- **错误处理**：统一解析后端返回的 `code`、`message`（见 [xsh-docs/backend/errors-and-conventions.md](../xsh-docs/backend/errors-and-conventions.md)），提交失败时可在结算页提示「上传失败，可重试」并允许重试。
- **登录态**：若后端要求鉴权，请求头带 token 或 session；首期匿名时可不带或带设备 id。

---

## 目录建议

```
front-game-defense/  或  miniprogram/
  pages/
    index/           # 主菜单
    battle/          # 战斗
    result/          # 结算
    leaderboard/     # 排行榜（可选独立页）
  components/
    game-canvas/     # 战斗 Canvas 封装
    upgrade-modal/   # 升级技能弹窗
    hud/             # 分数、等级、暂停
  game/              # 游戏逻辑（与框架解耦）
    loop.js          # 主循环、帧更新
    entities/        # 玩家、子弹、怪兽类或工厂
    collision.js     # 碰撞检测
    skills.js        # 技能配置与效果应用
  api/
    defense.js       # SubmitScore、GetLeaderboard、GetMyBest 封装
  store/
    user.js          # 用户信息、最高分（Pinia）
    game.js          # 本局分数、波数、技能列表等
  app.vue / main.js
```

游戏核心逻辑尽量放在 `game/` 下，便于单测与复用；页面只负责挂载 Canvas、调用循环与展示 HUD/弹窗。
