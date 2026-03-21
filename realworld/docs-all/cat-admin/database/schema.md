# 数据库表结构：萌宠之家

以下表在 realworld 的 MySQL 中新增或扩展，建议使用 GORM AutoMigrate 与现有 data 层一致。

---

## 1. 用户表扩展（users）

沿用现有 `users` 表，可扩展字段（若尚未存在）：

| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint PK | 主键 |
| username | varchar(64) unique | 登录名 |
| password_hash | varchar(255) | 密码哈希 |
| role | varchar(32) | 角色 |
| nickname | varchar(64) | 昵称（展示用） |
| avatar_url | varchar(512) | 头像 URL |
| signature | varchar(256) | 个性签名 |
| level | int | 等级 |
| coins | int | 金币 |
| heart_points | int | 爱心积分（排行榜等） |
| created_at | timestamp | 创建时间 |
| updated_at | timestamp | 更新时间 |

---

## 2. 宠物表（pets）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint PK | 主键 |
| user_id | bigint | 所属用户 ID |
| name | varchar(64) | 宠物名字 |
| species | varchar(16) | CAT / DOG |
| breed_id | bigint | 关联品种（可选） |
| avatar_url | varchar(512) | 宠物头像/大图 URL |
| background_url | varchar(512) | 互动背景图 URL |
| mood | varchar(64) | 心情文案（如「非常开心」） |
| affection | int | 亲密度 0–1000 |
| fullness | int | 饱食度 0–100 |
| happiness | int | 心情值 0–100 |
| cleanliness | int | 清洁度 0–100 |
| age_years | int | 年龄（岁，可选） |
| health_status | varchar(64) | 健康状态（可选） |
| created_at | timestamp | 创建时间 |
| updated_at | timestamp | 更新时间 |

- 索引：user_id（查询当前用户宠物）。

---

## 3. 互动记录（pet_interactions，可选）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint PK | 主键 |
| pet_id | bigint | 宠物 ID |
| user_id | bigint | 用户 ID |
| action | varchar(32) | STROKE / TEASE / TREAT / BATH |
| created_at | timestamp | 发生时间 |

- 用途：统计、审计；首期可不建，仅更新 pets 表即可。

---

## 4. 每日签到（check_ins）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint PK | 主键 |
| user_id | bigint | 用户 ID |
| date | date | 签到日期（唯一约束 user_id+date） |
| continuous_days | int | 当日连续签到天数 |
| created_at | timestamp | 创建时间 |

- 唯一索引：(user_id, date)。

---

## 5. 好友关系（user_friends，二期）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint PK | 主键 |
| user_id | bigint | 用户 ID |
| friend_id | bigint | 好友用户 ID |
| status | varchar(16) | ACCEPTED / PENDING |
| created_at | timestamp | 创建时间 |

---

## 6. 好友申请（friend_requests，二期）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint PK | 主键 |
| from_user_id | bigint | 申请人 ID |
| to_user_id | bigint | 被申请人 ID |
| message | varchar(256) | 申请留言 |
| status | varchar(16) | PENDING / ACCEPTED / REJECTED |
| created_at | timestamp | 创建时间 |

---

## 7. 商店商品（shop_items）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint PK | 主键 |
| name | varchar(128) | 商品名称 |
| category | varchar(32) | 食物 / 玩具 / 装扮 等 |
| price_coins | int | 售价（金币） |
| image_url | varchar(512) | 图片 URL |
| description | varchar(512) | 描述（可选） |
| created_at | timestamp | 创建时间 |

---

## 8. 用户背包（user_inventory）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint PK | 主键 |
| user_id | bigint | 用户 ID |
| item_id | bigint | 商品/道具 ID |
| quantity | int | 数量 |
| created_at | timestamp | 创建时间 |
| updated_at | timestamp | 更新时间 |

- 唯一索引：(user_id, item_id)，便于合并数量。

---

## 9. 品种（breeds）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint PK | 主键 |
| species | varchar(16) | CAT / DOG |
| name_cn | varchar(64) | 中文名（如「金毛寻回犬」） |
| name_en | varchar(64) | 英文名（可选） |
| size_tag | varchar(32) | 体型（如「大型犬」） |
| life_span | varchar(32) | 寿命（如「10-12年」） |
| image_url | varchar(512) | 图片 URL |
| created_at | timestamp | 创建时间 |

---

## 10. 宠物装扮（pet_appearance，或简化为 pets 表字段）

**方案 A**：独立表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint PK | 主键 |
| pet_id | bigint | 宠物 ID |
| slot | varchar(16) | HEAD / BODY / NECK |
| item_id | bigint | 装备的物品 ID |
| created_at | timestamp | 创建时间 |

**方案 B**：在 pets 表增加 equipped_head_id、equipped_body_id、equipped_neck_id（bigint，可为 NULL）。

---

## 11. 成就定义（achievements，二期）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint PK | 主键 |
| code | varchar(64) | 成就编码 |
| name | varchar(128) | 成就名称 |
| description | varchar(512) | 描述 |
| icon_url | varchar(512) | 图标 URL |
| condition_type | varchar(32) | 解锁条件类型 |
| condition_value | int | 条件值 |
| created_at | timestamp | 创建时间 |

---

## 12. 用户成就（user_achievements，二期）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint PK | 主键 |
| user_id | bigint | 用户 ID |
| achievement_id | bigint | 成就 ID |
| unlocked_at | timestamp | 解锁时间 |

---

## 首期建议

首期仅建：**users 扩展**、**pets**、**check_ins**、**shop_items**、**user_inventory**、**breeds**。  
社交（user_friends、friend_requests）、成就（achievements、user_achievements）、pet_interactions、pet_appearance 可二期再加。
