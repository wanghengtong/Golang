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
		userGroup.GET("/pageList", ginServer.UserController.PageList)
		userGroup.GET("/list", ginServer.UserController.List)
		userGroup.GET("/delete", ginServer.UserController.Delete)
		userGroup.GET("/toEdit", ginServer.UserController.ToEdit)
		userGroup.POST("/update", ginServer.UserController.Update)
		userGroup.GET("/toAdd", ginServer.UserController.ToAdd)
		userGroup.POST("/add", ginServer.UserController.Add)
	}

	adminGroup := r.Group("/admin")
	{
		adminGroup.GET("/index", ginServer.AdminAuthController.Index)
		adminGroup.POST("/login", ginServer.AdminAuthController.LoginToIndex)
		adminGroup.GET("/logout", ginServer.AdminAuthController.Logout)
	}
}
