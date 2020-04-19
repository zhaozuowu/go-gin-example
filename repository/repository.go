package repository

import "github.com/gin2/models"

type Repository interface {
	GetUserList(offset int, pageSize int, whereParams map[string]interface{}) ([]models.User, error)
	GetUserTotalNum(map[string]interface{}) (int, error)
}
