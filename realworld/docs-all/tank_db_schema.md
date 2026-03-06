### 坦克大战小游戏 - 数据库表结构设计

本文件描述坦克大战小游戏在 MySQL 中的核心表结构设计，对应后端 Kratos 项目中 `internal/data` 目录下的 GORM 实体：

- `TodoPO`       → `list`         （现有 Todo 示例表）
- `UserPO`       → `users`
- `MapPO`        → `maps`
- `MapTilePO`    → `map_tiles`
- `GameSessionPO`→ `game_sessions`
- `MatchRecordPO`→ `match_records`

所有表均通过 `internal/data/data.go` 中的 `NewData` 函数使用 `db.AutoMigrate(...)` 自动迁移。

---

### 1. 用户表 `users`

用于保存账号信息，与认证模块 `auth.v1.AuthService` 对接。

```sql
CREATE TABLE `users` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `username` VARCHAR(64) NOT NULL UNIQUE COMMENT '用户名',
  `password_hash` VARCHAR(255) NOT NULL COMMENT '密码哈希',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

对应 GORM 实体：`internal/data/user.go` 中的 `UserPO`。

---

### 2. 地图表 `maps`

存储地图的基础信息，与 `game.v1.MapBasicInfo`、`game.v1.MapDetail` 对应。

```sql
CREATE TABLE `maps` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` VARCHAR(100) NOT NULL COMMENT '地图名称',
  `description` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '地图描述',
  `width` INT NOT NULL COMMENT '地图宽度（格子数）',
  `height` INT NOT NULL COMMENT '地图高度（格子数）',
  `thumbnail_url` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '缩略图地址',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

对应 GORM 实体：`internal/data/map.go` 中的 `MapPO`。

---

### 3. 地图格子表 `map_tiles`

存储每张地图中每个格子的类型信息，对应 `game.v1.MapTile`。

```sql
CREATE TABLE `map_tiles` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
  `map_id` BIGINT NOT NULL COMMENT '关联 maps.id',
  `x` INT NOT NULL COMMENT '格子 X 坐标',
  `y` INT NOT NULL COMMENT '格子 Y 坐标',
  `tile` INT NOT NULL COMMENT '格子类型，对应 game.v1.TileType',
  PRIMARY KEY (`id`),
  KEY `idx_map_id` (`map_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

对应 GORM 实体：`internal/data/map.go` 中的 `MapTilePO`。

---

### 4. 对局表 `game_sessions`

记录每一局游戏对战的生命周期，与 `game.v1.CreateSessionReply.session_id` 等字段对应。

```sql
CREATE TABLE `game_sessions` (
  `session_id` VARCHAR(64) NOT NULL COMMENT '业务对局 ID，主键',
  `mode` INT NOT NULL COMMENT '游戏模式，对应 game.v1.GameMode',
  `map_id` BIGINT NOT NULL COMMENT '地图 ID',
  `local_player_count` INT NOT NULL COMMENT '本地玩家数量（1 或 2）',
  `status` INT NOT NULL COMMENT '对局状态，对应 game.v1.SessionStatus',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `started_at` DATETIME NULL DEFAULT NULL COMMENT '对局开始时间',
  `finished_at` DATETIME NULL DEFAULT NULL COMMENT '对局结束时间',
  PRIMARY KEY (`session_id`),
  KEY `idx_map_id` (`map_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

对应 GORM 实体：`internal/data/game_session.go` 中的 `GameSessionPO`。

---

### 5. 战绩记录表 `match_records`

存储每场对局的统计信息，供 `stats.v1.StatsService` 进行战绩查询与排行榜统计。

```sql
CREATE TABLE `match_records` (
  `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键',
  `session_id` VARCHAR(64) NOT NULL COMMENT '关联 game_sessions.session_id',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户 ID，游客可约定为 0',
  `mode` INT NOT NULL COMMENT '游戏模式，对应 stats.v1.GameMode',
  `map_id` BIGINT NOT NULL COMMENT '地图 ID',
  `is_win` TINYINT(1) NOT NULL COMMENT '是否获胜',
  `player1_kills` INT NOT NULL DEFAULT 0 COMMENT '玩家 1 击毁数',
  `player2_kills` INT NOT NULL DEFAULT 0 COMMENT '玩家 2 击毁数',
  `enemy_kills` INT NOT NULL DEFAULT 0 COMMENT '敌方击毁数',
  `duration_seconds` INT NOT NULL DEFAULT 0 COMMENT '对局时长（秒）',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_session_id` (`session_id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

对应 GORM 实体：`internal/data/match_record.go` 中的 `MatchRecordPO`。

---

### 6. 现有 Todo 表 `list`

为兼容现有 Todo 示例功能，继续保留 `list` 表：

```sql
CREATE TABLE `list` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `value` VARCHAR(255) NOT NULL COMMENT '待办内容',
  `is_completed` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否完成',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

对应 GORM 实体：`internal/data/todo.go` 中的 `TodoPO`。

