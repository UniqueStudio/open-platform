package database

import (
	"context"
	"time"

	"github.com/UniqueStudio/open-platform/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/go-redis/redis/v8"
)

var (
	UserDB      *gorm.DB
	OpenDB      *gorm.DB
	RedisClient *redis.Client
)

func SetupDatabase() error {
	opendb, err := gorm.Open(postgres.Open(config.Config.Database.Postgres.OpenDSN), &gorm.Config{})
	if err != nil {
		return err
	}
	OpenDB = opendb

	userdb, err := gorm.Open(postgres.Open(config.Config.Database.Postgres.UserDSN), &gorm.Config{})
	if err != nil {
		return err
	}
	UserDB = userdb

	rc := redis.NewClient(&redis.Options{
		Addr:     config.Config.Database.Redis.Addr,
		Password: config.Config.Database.Redis.Password,
		DB:       config.Config.Database.Redis.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := rc.Ping(ctx).Err(); err != nil {
		return err
	}

	RedisClient = rc

	return initTables()
}

func initTables() error {
	if !OpenDB.Migrator().HasTable(&SMSSignature{}) {
		if err := OpenDB.AutoMigrate(&SMSSignature{}); err != nil {
			return err
		}
	}

	if !OpenDB.Migrator().HasTable(&SMSTemplate{}) {
		if err := OpenDB.AutoMigrate(&SMSTemplate{}); err != nil {
			return err
		}
	}

	return nil
}
