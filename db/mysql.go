package db

import (
	"fmt"
	"log"
	"open-platform/utils"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var ORM *xorm.Engine
var ORM_ShortUrl *xorm.Engine

func init() {
	var err error
	ORM, err = xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true",
		utils.AppConfig.Mysql.User, utils.AppConfig.Mysql.Password, utils.AppConfig.Mysql.Host, utils.AppConfig.Mysql.Port, utils.AppConfig.Mysql.Database))

	if err != nil {
		log.Println(err)
		panic(err)
	}
	err = ORM.Sync(new(Reply), new(Status))
	if err != nil {
		log.Println(err)
		panic(err)
	}

	ORM_ShortUrl, err = xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true",
		utils.AppConfig.Mysql.User, utils.AppConfig.Mysql.Password, utils.AppConfig.Mysql.Host, utils.AppConfig.Mysql.Port, utils.AppConfig.Mysql.Database))

	if err != nil {
		log.Println(err)
		panic(err)
	}
	err = ORM_ShortUrl.Sync2(new(shorter))
	if err != nil {
		log.Println(err)
		panic(err)
	}
}
