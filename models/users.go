package models

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

type User struct {
	Model
	ID        uint   `gorm:"primary key;not null;auto_incrment;comment:'主键'" json:"id"`
	Name      string `gorm:"type:varchar(32);not null;default:'';comment:'姓名'" json:"name"`
	Email     string `gorm:"type:varchar(32);not null;default:'';comment:'邮箱'" json:"email"`
	Password  string `gorm:"type:varchar(60);not null;default:'';comment:'密码'" json:"password"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (user *User) GetUserList(offset int, pageSize int, whereParams map[string]interface{}) (userList []User, err error) {

	if err = user.orm.Where(whereParams).Offset(offset).Limit(pageSize).Find(&userList).Error; err != nil {
		return
	}
	return

}

func (user *User) TableName() string {

	return "users"
}

func (user *User) GetUserTotalNum(whereParams map[string]interface{}) (total int, err error) {

	if err = user.orm.Model(&User{}).Where(whereParams).Count(&total).Error; err != nil {
		return
	}
	return
}
func (user *User) CloseDb() {

	defer user.orm.Close()
}

func (user *User) CheckEmailHaveExists(email string) bool {

	var userInfo User
	user.orm.Select("id").Where("email = ?", email).First(&userInfo)
	if userInfo.ID > 0 {
		return true
	}
	return false

}

func (user *User) CreateUser(name, email string) error {

	if err := user.orm.Create(&User{Name: name, Email: email}).Error; err != nil {
		return err
	}
	return nil
}

func (user *User) GetUserInfoById(id int) (userInfo User, err error) {

	if err = user.orm.Where("id = ?", id).First(&userInfo).Error; err != nil {
		return
	}
	return
}

func (user *User) CheckLogin(name string, password string) bool {

	user.orm.Where(User{Name: name, Password: password}).First(&user)
	if user.ID > 0 {
		return true
	}
	return false
}
func NewUserModel() *User {
	user := &User{}
	user.InitConnect()
	return user
}
