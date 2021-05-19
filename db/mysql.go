package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"log"
	"open-platform/utils"
)

var ORM *xorm.Engine

func init() {
	var err error
	ORM, err = xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true",
		utils.AppConfig.Mysql.User, utils.AppConfig.Mysql.Password, utils.AppConfig.Mysql.Host, utils.AppConfig.Mysql.Port, utils.AppConfig.Mysql.Database))

	// temp config for local debug
	//ORM, err = xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true",
	//	"unique", "P@ssw0rd", "localhost", "3306", "open"))

	if err != nil {
		log.Println(err)
		panic(err)
	}
	err = ORM.Sync(new(Reply), new(Status), new(Short_Url))
	if err != nil {
		log.Println(err)
		panic(err)
	}

}
