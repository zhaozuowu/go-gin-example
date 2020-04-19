package service

import (
	"github.com/gin2/models"
	"github.com/gin2/pkg/setting"
	"github.com/gin2/repository"
)

type UserService struct {
	userRepository repository.Repository
}

func (userService *UserService) GetUserList(requstParams map[string]interface{}, page int) ([]models.User, error) {

	offset := 0
	offset = (page - 1) * setting.PageSize
	if offset < 0 {
		offset = 0
	}

	return userService.userRepository.GetUserList(offset, setting.PageSize, requstParams)

}

func (userService *UserService) GetUserTotalNum(requstParams map[string]interface{}) (int,error) {

	return userService.userRepository.GetUserTotalNum(requstParams)

}
func NewUserService() *UserService {

	user := &UserService{}
	user.userRepository = repository.NewUserRepository()
	return user
}
