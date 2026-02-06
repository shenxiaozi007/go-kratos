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
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewDB)

// Data .
type Data struct {
	// TODO wrapped database client
	DB *gorm.DB
}

// NewData .
func NewData(c *conf.Data, db *gorm.DB) (*Data, func(), error) {
	cleanup := func() {
		log.Info("closing the data resources")
	}

	return &Data{
		DB: db,
	}, cleanup, nil
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
