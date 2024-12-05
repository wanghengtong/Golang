package common

import (
	"github.com/sirupsen/logrus"
	"os"
)

func init() {
	// 初始化日志
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}
