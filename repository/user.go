package repository

import "github.com/gin2/models"

type UserRepository struct {
	userModel *models.User
}

func (userRepository *UserRepository) GetUserList(offset int, pageSize int, whereParams map[string]interface{}) ([]models.User, error) {

	return userRepository.userModel.GetUserList(offset, pageSize, whereParams)
}

func (userRepository *UserRepository) GetUserTotalNum(whereParams map[string]interface{}) (int, error) {
	return userRepository.userModel.GetUserTotalNum(whereParams)
}

func NewUserRepository() Repository {

	obj := &UserRepository{}
	obj.userModel = models.NewUserModel()
	return obj
}
