package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
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
	// 初始化 viper 读取配置文件
	// 配置文件名（不带扩展名）
	viper.SetConfigName("/config/config")
	// 配置文件类型
	viper.SetConfigType("yaml")
	// 配置文件搜索路径
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("读取配置文件失败: %v\n", err)
		return
	}

	// 从配置文件中读取数据库连接信息
	var (
		userName = viper.GetString("db.username")
		password = viper.GetString("db.password")
		host     = viper.GetString("db.host")
		port     = viper.GetInt("db.port")
		dbName   = viper.GetString("db.dbname")
		charset  = viper.GetString("db.charset")
	)

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", userName, password, host, port, dbName, charset)
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

	// 初始化插入数据函数调用
	initAdminData(engine)
	initUserData(engine)
	fmt.Println()

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

func initAdminData(engine *xorm.Engine) {
	var admins []model.Admin
	admins = append(admins, model.Admin{Id: 1, AdminName: "超级管理员", LoginName: "superadmin", Password: "superadmin"})
	admins = append(admins, model.Admin{Id: 2, AdminName: "系统管理员", LoginName: "sysadmin", Password: "sysadmin"})
	admins = append(admins, model.Admin{Id: 4, AdminName: "安全管理员", LoginName: "secadmin", Password: "secadmin"})
	admins = append(admins, model.Admin{Id: 3, AdminName: "审计管理员", LoginName: "logadmin", Password: "logadmin"})
	// ... 其他管理员数据 ...
	affected, err := engine.Insert(&admins)
	if err != nil {
		fmt.Printf("插入管理员数据失败: %v\n", err)
		return
	}
	fmt.Printf("插入管理员数据成功，受影响行数: %d\n", affected)
}

func initUserData(engine *xorm.Engine) {
	var users []model.User
	users = append(users, model.User{Id: 1, Name: "张三", Age: 18, Sex: 1, Mobile: "15011112222", Mail: "zhangsan@qq.com"})
	users = append(users, model.User{Id: 2, Name: "李四", Age: 18, Sex: 0, Mobile: "15011112222", Mail: "lisi@qq.com"})
	users = append(users, model.User{Id: 3, Name: "王五", Age: 18, Sex: 1, Mobile: "15011112222", Mail: "wangwu@qq.com"})
	// ... 其他用户数据 ...
	affected, err := engine.Insert(&users)
	if err != nil {
		fmt.Printf("插入用户数据失败: %v\n", err)
		return
	}
	fmt.Printf("插入用户数据成功，受影响行数: %d\n", affected)
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
