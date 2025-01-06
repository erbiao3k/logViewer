package dao

import (
	"context"
	"github.com/go-redis/redis/v8"
	"logViewerServer/setting"
)

var Ctx = context.Background()

// InitRedis 初始化redis连接信息
func InitRedis(cfg *setting.RedisConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password, // 密码
		DB:       cfg.DB,       // 数据库
		PoolSize: cfg.PoolSize, // 连接池大小
	})
	return rdb
}

func Rdb() *redis.Client {
	return InitRedis(setting.Conf.RedisConfig)
}
