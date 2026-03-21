# 数据库迁移与扩展说明

## 首期表（建议先上）

- users（扩展 nickname, avatar_url, signature, level, coins, heart_points）
- pets
- check_ins
- shop_items
- user_inventory
- breeds

在 realworld 的 `internal/data/data.go` 的 AutoMigrate 中依次加入对应 PO，如：

```go
&UserPO{},      // 已有，仅增加新字段
&PetPO{},
&CheckInPO{},
&ShopItemPO{},
&UserInventoryPO{},
&BreedPO{},
```

## 二期表（社交、成就、装扮等）

- pet_interactions
- user_friends
- friend_requests
- pet_appearance（或 pets 表增加 equipped_* 字段）
- achievements
- user_achievements
- user_purchases（可选，购买记录）

## 注意

- 新增字段时若表已存在，GORM AutoMigrate 会尝试添加列，需注意默认值与兼容性。
- 唯一索引（如 check_ins 的 user_id+date）在 GORM 中通过 tag 或迁移脚本保证。
