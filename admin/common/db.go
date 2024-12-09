package common

import (
	"admin/model"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"time"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

func GetMysqlEngine() (*xorm.Engine, error) {
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
		return nil, fmt.Errorf("数据库连接失败，请检查！: %v", err)
	}

	// 配置 XORM 日志
	engine.SetLogger(log.NewSimpleLogger(os.Stdout))
	// 是否显示 SQL 语句
	showsql := viper.GetBool("logger.showsql")
	engine.ShowSQL(showsql)

	//设置连接池的空闲数大小
	engine.SetMaxIdleConns(viper.GetInt("db.MaxIdleConns"))
	//设置最大打开连接数
	engine.SetMaxOpenConns(viper.GetInt("db.MaxOpenConns"))
	//连接的最大生存时间
	engine.SetConnMaxLifetime(time.Second * time.Duration(viper.GetInt("db.MaxLifetime")))
	if err := engine.Ping(); err != nil {
		logrus.Info(err, engine.DataSourceName())
		return engine, err
	}

	// 同步表结构
	err = syncTables(engine)
	if err != nil {
		return nil, err
	}

	// 启动连接健康检查
	go checkConnectionHealth(engine)
	return engine, nil
}

func syncTables(engine *xorm.Engine) error {
	err := engine.Sync(new(model.Admin))
	if err != nil {
		return fmt.Errorf("表结构同步失败，请检查！: %v", err)
	}
	err = engine.Sync(new(model.User))
	if err != nil {
		return fmt.Errorf("表结构同步失败，请检查！: %v", err)
	}
	return nil
}

func checkConnectionHealth(engine *xorm.Engine) {
	connTimeout := time.Duration(viper.GetInt("db.connTimeout"))
	ticker := time.NewTicker(connTimeout * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if err := engine.Ping(); err != nil {
			logrus.Fatalf("连接健康检查失败: %v\n", err)
			// 重新初始化数据库连接
			newEngine, err := GetMysqlEngine()
			if err != nil {
				logrus.Fatalf("重新初始化数据库连接失败: %v\n", err)
				continue
			}
			engine = newEngine
			logrus.Infof("数据库连接已重新初始化")
		} else {
			logrus.Infof("连接健康检查成功")
		}
	}
}
