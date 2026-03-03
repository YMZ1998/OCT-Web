package router

import (
	"github.com/gin-gonic/gin"
	"oct-backend/internal/controller"
	"oct-backend/internal/middleware"
)

func SetupRouter(userCtl *controller.UserController) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.ZapLogger())
	r.POST("/api/v1/register", userCtl.Register)
	r.POST("/api/v1/login", userCtl.Login)

	auth := r.Group("/api/v1/user")
	auth.Use(middleware.JWTAuth())
	auth.GET("/:id", userCtl.GetUser)
	auth.PUT("/:id", userCtl.UpdateUser)
	auth.PUT("/:id/projects", userCtl.UpdateUserProjects)
	auth.DELETE("/:id", userCtl.DeleteUser)
	return r
}
