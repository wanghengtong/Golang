package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/thinkerou/favicon"
	"iam-auth/controller"
	"iam-auth/model"
	"iam-auth/router"
	"net/http"
	"os"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

func main() {
	var (
		userName string = "root"
		password string = "MSFpMmlZ1ZkY6ZLSZrPU"
		host     string = "192.168.20.210"
		port     int    = 3306
		dbName   string = "chiansectest"
		charest  string = "utf8mb4"
	)
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", userName, password, host, port, dbName, charest)
	engine, err := xorm.NewEngine("mysql", dataSourceName)
	if err != nil {
		fmt.Println("数据库连接失败，请检查！")
	}

	// 配置 XORM 日志
	engine.SetLogger(log.NewSimpleLogger(os.Stdout))
	// 显示 SQL 语句
	engine.ShowSQL(true)

	err = engine.Sync(new(model.Admin))
	if err != nil {
		fmt.Println("表结构同步失败，请检查！")
	}

	err = engine.Sync(new(model.User))
	if err != nil {
		fmt.Println("表结构同步失败，请检查！")
	}

	ginServer := gin.Default()
	ginServer.Use(favicon.New("./resources/favicon.ico"))
	ginServer.LoadHTMLGlob("./resources/templates/*")
	// ginServer.Use(myHandler())

	setupRouter(ginServer, engine)
	err = ginServer.Run(":8080")
	if err != nil {
		return
	}
}

func myHandler() gin.HandlerFunc {
	// 定义放行路径列表
	skipAuthPaths := map[string]struct{}{
		"/admin/index": {},
		"/admin/login": {},
	}
	return func(c *gin.Context) {
		// 获取请求路径
		requestPath := c.Request.URL.Path
		// 检查请求路径是否在放行列表中
		if _, exists := skipAuthPaths[requestPath]; exists {
			// 如果在放行列表中，跳过认证检查
			c.Next()
			return
		}
		// 获取 Authorization 头
		Authorization := c.Request.Header.Get("Authorization")
		fmt.Println("===================", Authorization)
		if Authorization == "" {
			c.Writer.Header().Set("Authorization", "")
			fmt.Println("未认证！")
			c.HTML(http.StatusOK, "401.html", nil)
			c.Abort()
		} else {
			c.Set("Authorization", c.Request.Header.Get("Authorization"))
			c.Next()
		}
	}
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
