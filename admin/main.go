package main

import (
	"admin/common"
	"admin/controller"
	"admin/router"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/thinkerou/favicon"
	"xorm.io/xorm"
)

func main() {
	// 初始化数据库
	engine, err := common.GetMysqlEngine()
	if err != nil {
		logrus.Fatalf("数据库初始化失败: %v", err)
	}
	// 初始化数据
	common.InitData(engine)
	logrus.Info("数据初始化完成")

	// 初始化 Gin 服务器
	ginServer := gin.Default()
	ginServer.Use(favicon.New("./resources/favicon.ico"))
	ginServer.LoadHTMLGlob("./resources/templates/*")
	ginServer.Use(common.InitMiddleware())

	// 注册路由
	setupRouter(ginServer, engine)
	// 启动服务器
	common.StartServer(ginServer)
}

func setupRouter(ginServer *gin.Engine, engine *xorm.Engine) {
	// 初始化控制器
	adminAuthcontroller := &controller.AdminAuthController{
		Engine: engine,
	}
	userController := &controller.UserController{
		Engine: engine,
	}
	// 初始化路由
	authRouter := &router.Router{
		AdminAuthController: adminAuthcontroller,
		UserController:      userController,
	}
	// 注册路由
	authRouter.RegisterRoutes(ginServer)
}
