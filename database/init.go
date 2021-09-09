package database

import (
	"github.com/UniqueStudio/open-platform/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	UserDB *gorm.DB
	OpenDB *gorm.DB
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
