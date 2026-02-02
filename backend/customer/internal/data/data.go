package data

import (
	"context"
	"customer/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewCustomerData)

// Data .
type Data struct {
	// TODO wrapped database client
	Rdb *redis.Client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	data := &Data{}

	ctxRedis := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: c.Redis.Addr,
	})
	// 测试链接
	errRedis := rdb.Ping(ctxRedis).Err()

	if errRedis != nil {
		data.Rdb = nil
	} else {
		data.Rdb = rdb
	}

	cleanup := func() {
		_ = rdb.Close()
		log.NewHelper(logger).Info("closing the data resources")
	}

	return data, cleanup, nil
}
