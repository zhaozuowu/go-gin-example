package database

import (
	"fmt"
	"github.com/gin2/pkg/setting"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

type Mysql struct {
}

func (mysql *Mysql) InitConnect() *gorm.DB {

	sec, err := setting.Cfg.GetSection("database")
	if err != nil {
		log.Fatalf("Fail to get section 'database':%v\n", err)
	}

	dbType := sec.Key("TYPE").String()
	dbUser := sec.Key("USER").String()
	dbPassword := sec.Key("PASSWORD").String()
	dbName := sec.Key("DATABASE").String()
	host := sec.Key("HOST").String()

	db, err := gorm.Open(dbType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		dbUser,
		dbPassword,
		host,
		dbName,
	))

	if err != nil {
		log.Fatalf("Fail to connect mysql:%v", err)
	}
	return db

}

func NewMysql() Db {

	db := &Mysql{}
	return db
}
