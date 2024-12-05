package common

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// 定义常量
const (
	defaultPort = 8080
	configPath  = "./config/config.yaml"
)

func init() {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalf("读取配置文件失败: %v", err)
	}
}
