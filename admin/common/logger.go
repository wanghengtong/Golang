package common

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

func init() {
	// 初始化日志
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)
	// 设置日志级别
	logLevel := viper.GetString("logger.level")
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.Fatalf("解析日志级别失败: %v", err)
	}
	logrus.SetLevel(level)
}
