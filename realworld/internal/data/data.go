package data

import (
	"context"
	"fmt"
	"realworld/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewGreeterRepo,
	NewDB,
	NewArticleRepo,
	NewTodoRepo,
	NewGameMapRepo,
	NewGameSessionRepo,
	NewGameStatsRepo,
	NewGameConfigProvider,
	NewStatsRepo,
)

// Data .
type Data struct {
	// TODO wrapped database client
	DB *gorm.DB
}

// NewData 创建 Data 结构体的实例，并在其中可以进行一次性的全局数据层初始化工作，
// 比如自动迁移数据表结构等。
func NewData(c *conf.Data, db *gorm.DB) (*Data, func(), error) {
	cleanup := func() {
		log.Info("closing the data resources")
	}

	d := &Data{
		DB: db,
	}

	// 使用 GORM 的 AutoMigrate 功能自动创建/更新所有与坦克大战相关的数据表结构。
	// 这一步在应用启动时执行一次，确保以下表存在并与对应的 PO 结构保持同步：
	//   - list           （待办事项示例表，来自现有功能）
	//   - users          （账号表）
	//   - maps           （地图基础信息表）
	//   - map_tiles      （地图格子明细表）
	//   - game_sessions  （游戏对局表）
	//   - match_records  （战绩记录表）
	if err := db.AutoMigrate(
		&TodoPO{},
		&UserPO{},
		&MapPO{},
		&MapTilePO{},
		&GameSessionPO{},
		&MatchRecordPO{},
	); err != nil {
		return nil, nil, err
	}
	if err := SeedMapsIfEmpty(context.Background(), db); err != nil {
		return nil, nil, err
	}

	return d, cleanup, nil
}

func NewDB(c *conf.Data) (*gorm.DB, error) {
	//
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Database.Username,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Database,
	)
	db, err := gorm.Open(mysql.Open(dsn))

	if err != nil {
		return nil, fmt.Errorf("failed opening connection to mysql: %v", err)
	}

	return db, nil
}
