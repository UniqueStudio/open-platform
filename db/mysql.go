package db

import (
	"fmt"
	"log"
	"open-platform/utils"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var ORM *xorm.Engine

func init() {
	fmt.Println(utils.AppConfig.Mysql.Host)
	var err error
	ORM, err = xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true",
		utils.AppConfig.Mysql.User, utils.AppConfig.Mysql.Password, utils.AppConfig.Mysql.Host, utils.AppConfig.Mysql.Port, utils.AppConfig.Mysql.Database))

	if err != nil {
		log.Println(err)
		panic(err)
	}
	err = ORM.Sync(new(Reply), new(Status),new(Short_Url))
	if err != nil {
		log.Println(err)
		panic(err)
	}

}
