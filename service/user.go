package service

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gin2/models"
	error2 "github.com/gin2/pkg/error"
	"github.com/gin2/pkg/logging"
	"github.com/gin2/pkg/redis"
	"github.com/gin2/pkg/setting"
	"github.com/gin2/repository"
	"github.com/json-iterator/go"   // 引入
	"github.com/unknwon/com"
	"strconv"
	"strings"
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

	keys := []string{
		error2.CACHE_USER_LIST,
		"list",
	}

	if requstParams["name"] != "" {
		keys = append(keys, com.ToStr(requstParams["name"]))
	}

	if requstParams["email"] != "" {
		keys = append(keys, com.ToStr(requstParams["email"]))
	}

	if page > 0 {
		keys = append(keys, strconv.Itoa(page))
	}

	str := strings.Join(keys, "_")
	h := md5.New()
	h.Write([]byte(str))
	cacheKey :=  hex.EncodeToString(h.Sum(nil))

	cacheResult, err := redis.Get(cacheKey)

	if err != nil {
		logging.Info(err, cacheResult)
	}

	var userList []models.User
	//json.Unmarshal(cacheResult, &userList)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	json.Unmarshal(cacheResult,&userList)
	if userList != nil {
		return userList, nil
	}
	userList, err =  userService.userRepository.GetUserList(offset, setting.PageSize, requstParams)
	if err != nil {
			logging.Info(err, cacheResult)
		return nil,err
	}

	err = redis.Set(cacheKey,userList,error2.CACHE_EXPIRE_TIME)
	if err != nil {
		logging.Info("set redis cache error",err)
	}
	return userList,nil


}

func (userService *UserService) GetUserTotalNum(requstParams map[string]interface{}) (int, error) {

	return userService.userRepository.GetUserTotalNum(requstParams)

}
func NewUserService() *UserService {

	user := &UserService{}
	user.userRepository = repository.NewUserRepository()
	return user
}
