package common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func StartServer(ginEngine *gin.Engine) {
	port := viper.GetInt("server.port")
	if port == 0 {
		port = 8080
	}
	logrus.Infof("启动服务器，监听端口: %d", port)
	err := ginEngine.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		logrus.Fatalf("服务器启动失败: %v", err)
	}
}
