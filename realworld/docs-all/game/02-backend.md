# 后端开发文档：防守割草

本文档约定防守割草在 Kratos 后端的服务边界、API 设计、数据模型与目录约定。与现有 [api/game/v1](../../api/game/v1) 坦克游戏并存，建议通过**独立 Proto 与 Service**（如 `defense.proto` / `DefenseService`）扩展，避免与 `GameService` 混在一起。

---

## 服务边界

- 在 realworld 仓库内扩展，不新建独立服务。
- **api/game/v1**：现有 `game.proto` 保留坦克地图/对局相关；防守割草新增 `defense.proto`，包名仍为 `game.v1` 或单独 `defense.v1`（按项目统一），定义 `DefenseService`。
- 防守割草仅需：**提交分数**、**全服排行榜**、**我的最高分**；可选：**微信登录/鉴权**（code 换 openid/session）。

---

## API 设计（建议 RPC 与 Message）

以下为建议的 RPC 与 Message 列表，便于落地为 proto 与 HTTP 路由。

### 1. 提交分数

- **SubmitScore**
  - 请求：`user_id` 或 `openid`（若未登录可传匿名 id）、`score`（本局得分）、`wave`（生存波数）、可选 `session_id` 防重放。
  - 响应：`success`、可选 `new_high_score`（是否刷新个人最高分）。
  - 行为：校验参数（如 score/wave 合理范围），更新用户最高分（若本局更高），可写分数记录表供审计；排行榜数据可同步写 Redis 或异步刷新。
  - HTTP：`POST /api/v1/game/defense/score` 或 `POST /api/game/v1/defense/submit-score`（以实际 proto 路径为准）。

### 2. 全服排行榜

- **GetLeaderboard**
  - 请求：`page`（从 1）、`page_size`（默认 20，最大 100）。
  - 响应：`items[]`（每条：排名、user_id/openid、昵称、头像 URL、分数）、`total`。
  - 行为：从 Redis 有序集合按分数倒序分页；若 Redis 未命中则从 MySQL/分数表聚合后回填 Redis。
  - HTTP：`GET /api/v1/game/defense/leaderboard?page=1&page_size=20`。

### 3. 我的最高分

- **GetMyBest**
  - 请求：当前用户标识（从 session/token 或 openid 取）。
  - 响应：`best_score`、`best_wave`、`rank`（当前在全服排名，可选）。
  - 行为：查用户表或分数表得最高分与波数；排名可从 Redis 查。
  - HTTP：`GET /api/v1/game/defense/my-best`（需鉴权或传 openid）。

### 4. 登录/鉴权（可选，首期可不做）

- **WechatLogin** 或由统一网关处理：`code` 换 `openid`、`session_key`，后端发自定义 session 或 JWT 给前端，后续请求带 token，后端解析出 user_id/openid。
- 首期若无登录：提交分数与排行榜可使用**匿名 id**（如 deviceId 或临时 id），我的最高分可省略或用同一匿名 id 查询。

---

## 数据模型

### 用户/玩家表（如 `game_defense_user`）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint PK | 主键 |
| openid | varchar(64) | 微信 openid，匿名可为空或占位 |
| nickname | varchar(64) | 昵称，可选 |
| avatar_url | varchar(256) | 头像 URL，可选 |
| best_score | int | 历史最高分 |
| best_wave | int | 历史最高波数 |
| updated_at | timestamp | 最近一次更新（如提交分数时） |

- 唯一索引：`openid`（或匿名 id 占位字段）。

### 分数记录表（如 `game_defense_score_log`）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint PK | 主键 |
| user_id | bigint | 关联用户表 |
| score | int | 本局得分 |
| wave | int | 本局波数 |
| created_at | timestamp | 提交时间 |

- 用途：排行榜聚合、防作弊、审计；可配合 Redis 缓存减少对 MySQL 的排行榜查询压力。

---

## 存储

- **MySQL**：用户表、分数记录表；GORM 管理，在 `internal/data` 的迁移中注册新表。
- **Redis**：
  - 有序集合（ZADD key=leaderboard, score=分数, member=user_id），按分数倒序；ZREVRANGE 分页查榜。
  - 提交分数时更新用户最高分后，ZADD 更新该用户在排行榜中的分数。
- 排行榜可设置过期时间或定时从 MySQL 回填，视 QPS 与一致性要求而定。

---

## 错误与分页约定

- 与 [xsh-docs/backend/errors-and-conventions.md](../xsh-docs/backend/errors-and-conventions.md) 一致：
  - 业务错误：HTTP 200，body 内 `code`、`message`；如 `40002` 参数校验失败、`40001` 资源不存在。
  - 分页：`page` 从 1，`page_size` 默认 20、最大 100；列表响应含 `items` 与 `total`。
  - 时间字段：Proto 使用 `google.protobuf.Timestamp`，JSON 为 RFC3339 字符串。

---

## 防作弊与校验（建议）

- **分数合理性**：单局 score/wave 上限（如 999999、999），超出拒绝或截断。
- **频率限制**：同一用户/openid 短时间（如 1 分钟）内提交次数上限，避免刷榜。
- **去重**：同一局可带 `session_id` 或客户端生成的局唯一 id，后端幂等处理，避免重复提交同一局。

---

## 目录约定

| 层级 | 文件/目录 | 说明 |
|------|-----------|------|
| api | api/game/v1/defense.proto | 防守割草 Service 与 Message 定义（或 api/defense/v1/defense.proto 独立包） |
| biz | internal/biz/game_defense.go 或 game_score.go、game_leaderboard.go | 分数提交、排行榜、我的最高分 Usecase |
| data | internal/data/game_defense.go | PO、Repo，用户表与分数表 CRUD、Redis 排行榜读写 |
| service | internal/service/game_defense.go | 实现 DefenseService，调用 biz |
| server | internal/server/http.go | 注册 RegisterDefenseHTTPServer（以生成代码为准） |
| wire | cmd/realworld/wire.go、wire_gen.go | 增加 defense 相关 Provider 与注入 |

Proto 生成后，在 `internal/server/http.go` 中注册路由，在 `wire.go` 中注入 biz/data/service 依赖。
