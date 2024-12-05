package router

import (
	"admin/controller"
	"github.com/gin-gonic/gin"
)

type Router struct {
	AdminAuthController *controller.AdminAuthController
	UserController      *controller.UserController
}

func (ginServer *Router) RegisterRoutes(r *gin.Engine) {
	userGroup := r.Group("/user")
	{
		userGroup.GET("/list", ginServer.UserController.List)
		userGroup.GET("/delete", ginServer.UserController.Delete)
		userGroup.GET("/edit", ginServer.UserController.Edit)
		userGroup.POST("/update", ginServer.UserController.Update)
	}

	adminGroup := r.Group("/admin")
	{
		adminGroup.GET("/index", ginServer.AdminAuthController.Index)
		adminGroup.POST("/login", ginServer.AdminAuthController.LoginToIndex)
	}
}
