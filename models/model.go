package models

import (
	"github.com/jinzhu/gorm"
	"github.com/gin2/database"
)

type Model struct {
	orm *gorm.DB
}

func (model *Model) InitConnect() {
	mysql := database.NewMysql()
	model.orm = mysql.InitConnect()

}
