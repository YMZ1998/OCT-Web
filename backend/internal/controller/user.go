package controller

import (
	"oct-backend/internal/response"
	"oct-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	Service *service.UserService
}

func (ctl *UserController) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Gender   string `json:"gender"`
		Age      int    `json:"age"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, 400, "参数错误", nil)
		return
	}
	err := ctl.Service.Register(req.Username, req.Password, req.Email, req.Gender, req.Age)
	if err != nil {
		response.JSON(c, 500, "注册失败", nil)
		return
	}
	response.JSON(c, 0, "注册成功", nil)
}

func (ctl *UserController) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, 400, "参数错误", nil)
		return
	}
	user, err := ctl.Service.GetUserByUsername(req.Username)
	if err != nil {
		response.JSON(c, 401, "用户不存在", nil)
		return
	}
	token, err := ctl.Service.Login(req.Username, req.Password)
	if err != nil {
		response.JSON(c, 401, "登录失败", nil)
		return
	}
	response.JSON(c, 0, "登录成功", gin.H{"token": token, "id": user.ID.Hex()})
}

// CRUD 示例
func (ctl *UserController) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := ctl.Service.GetUserByID(id)
	if err != nil {
		response.JSON(c, 404, "用户不存在", nil)
		return
	}
	response.JSON(c, 0, "success", user)
}
func (ctl *UserController) UpdateUser(c *gin.Context) {
	// ...更新用户信息...
}
func (ctl *UserController) DeleteUser(c *gin.Context) {
	// ...删除用户...
}
