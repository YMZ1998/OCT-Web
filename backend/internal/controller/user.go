package controller

import (
	"net/http"
	"oct-backend/internal/model"
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
	if err != nil || token == "" {
		response.JSON(c, 401, "登录失败", nil)
		return
	}
	response.JSON(c, 0, "登录成功", gin.H{"token": token, "id": user.ID.Hex()})
}

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
	id := c.Param("id")
	if authID, ok := c.Get("userID"); !ok || authID != id {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "无权限修改该用户信息"})
		return
	}

	var req struct {
		Email  *string `json:"email"`
		Gender *string `json:"gender"`
		Age    *int    `json:"age"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, 400, "参数错误", nil)
		return
	}

	user, err := ctl.Service.UpdateUserByID(id, req.Email, req.Gender, req.Age)
	if err != nil {
		response.JSON(c, 500, "更新失败", nil)
		return
	}

	response.JSON(c, 0, "更新成功", user)
}

func (ctl *UserController) UpdateUserProjects(c *gin.Context) {
	id := c.Param("id")
	if authID, ok := c.Get("userID"); !ok || authID != id {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "无权限修改该用户项目数据"})
		return
	}

	var req struct {
		ProjectState model.ProjectState `json:"project_state"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, 400, "参数错误", nil)
		return
	}

	user, err := ctl.Service.UpdateUserProjectState(id, req.ProjectState)
	if err != nil {
		response.JSON(c, 500, "项目数据更新失败", nil)
		return
	}

	response.JSON(c, 0, "更新成功", user)
}

func (ctl *UserController) DeleteUser(c *gin.Context) {
	response.JSON(c, 501, "暂未实现", nil)
}
