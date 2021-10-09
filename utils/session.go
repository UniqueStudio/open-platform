package utils

import (
	"github.com/UniqueStudio/open-platform/config"
	"github.com/gin-contrib/sessions/redis"
)

var (
	RedisSessionStore redis.Store
)

func SetupSessionStore() error {
	store, err := redis.NewStore(
		32, "tcp",
		config.Config.Database.Redis.Addr,
		config.Config.Database.Redis.Password,
		[]byte(config.Config.Application.SessionKey),
	)
	if err != nil {
		return err
	}
	RedisSessionStore = store
	return nil
}
