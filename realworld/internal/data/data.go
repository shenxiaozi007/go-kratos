package data

import (
	"fmt"
	"realworld/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewDB, NewArticleRepo, NewTodoRepo)

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

	// 使用 GORM 的 AutoMigrate 功能自动创建/更新 Todo 表结构。
	// 这一步在应用启动时执行一次，确保名为 "list" 的表存在并与 TodoPO 结构保持同步。
	if err := db.AutoMigrate(&TodoPO{}); err != nil {
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
