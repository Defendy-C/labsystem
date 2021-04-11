package db

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"labsystem/util/logger"
	"strconv"

	"labsystem/configs"
)

type RedisDB struct {
	Config *configs.RedisConfig
	Cli    *redis.Client
	ctx    *context.Context
}

func NewRedis() *RedisDB {
	config := configs.NewRedisConfig()
	logger.Log.Info("connecting redis...", zap.Any("config", config))
	addr := config.Host + ":" + strconv.Itoa(int(config.Port))
	ctx := context.Background()
	return &RedisDB{config, redis.NewClient(&redis.Options{Addr: addr, Password: config.Password, DB: int(config.DB)}), &ctx}
}

func (db *RedisDB) Ping() error {
	return db.Cli.Ping(*db.ctx).Err()
}
