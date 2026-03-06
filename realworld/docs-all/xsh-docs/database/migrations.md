# 迁移策略

## 方案一：GORM AutoMigrate（推荐）

在现有 [internal/data/data.go](../../internal/data/data.go) 的 `NewData` 中，在已有 `AutoMigrate` 调用里追加 xsh 相关 PO，例如：

```go
if err := db.AutoMigrate(
    // ... 现有 TodoPO, UserPO, MapPO, ...
    &XshProductPO{},
    &XshHotRankPO{},
    &XshProductLinkPO{},
    &XshContentTemplatePO{},
    &XshContentDraftPO{},
    &XshPlatformAccountPO{},
    &XshPublishTaskPO{},
    &XshPublishSchedulePO{},
    &XshCommentPO{},
    &XshInboxActionPO{},
); err != nil {
    return nil, nil, err
}
```

- 与现有 realworld 共用同一 DB 实例与库（如 `kratos-blog`），仅新增表，表名由各 PO 的 `TableName()` 返回 `xsh_products` 等。
- 无需单独执行 SQL 脚本，启动时自动建表/补列（注意 GORM 的 AutoMigrate 会加列但不会删列或改类型，大变更需手写迁移）。

## 方案二：独立库与 SQL 脚本（可选）

若希望 xsh 使用独立库 `xsh_admin`：

1. 在 `configs/config.yaml` 的 `data.database` 中增加一套配置（或新 key，如 `xsh_database`），指向 `xsh_admin`。
2. 在 data 层初始化第二个 `*gorm.DB` 仅用于 xsh 表，并在对应 Repo 中使用该 DB。
3. 首次建表可导出 GORM 生成的 DDL，或按 [schema.md](schema.md) 手写 `CREATE TABLE`，保存为 `scripts/migrations/001_xsh_tables.sql`，在部署或 CI 中执行一次。

## 建议

- 首期采用**方案一**，与现有项目同库，在 `internal/data` 中新增 xsh 的 PO 与 Repo，并在 `data.go` 的 AutoMigrate 中注册上述 10 张表，保证表结构与 [schema.md](schema.md) 一致（字段名、类型、索引在 PO 的 gorm tag 中体现）。
