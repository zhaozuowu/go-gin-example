package request

type UserListRequst struct {
	Name  string `form:"name" binding:"required"`
	Email string `form:"email" binding:"required"`
}
