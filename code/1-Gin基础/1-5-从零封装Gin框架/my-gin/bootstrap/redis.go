package bootstrap

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"my-gin/global"
)

func InitializeRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.App.Config.Redis.Host, global.App.Config.Redis.Port),
		Password: global.App.Config.Redis.Password,
		DB:       global.App.Config.Redis.DB,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		global.App.Log.Error("Redis connect ping failed, err: ", zap.Any("err", err))
		return nil
	}
	return client
}
