package error

var ErrorMsg = map[int]string{

	SUCCESS:        "ok",
	ERROR:          "fail",
	INVALID_PARAMS: "请求参数错误2",
	EMAIL_HAS_EXISTS:"邮箱已存在",
	REGISTER_USER__FAIL:"注册失败",
	TOKEN_IS_EMPTY:"token为空",
	ERROR_AUTH_CHECK_TOKEN_FAIL:"token认证失败",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT:"token失效",
	ERROR_AUTH:"权限校验失败",
	ERROR_AUTH_TOKEN:"生成token失败",

}

func GetErrorMsg(code int) string {

	if value, ok := ErrorMsg[code]; ok {
		return value
	}
	return ErrorMsg[ERROR]
}
