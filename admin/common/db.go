package common

import (
	"admin/model"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

func InitDB() (*xorm.Engine, error) {
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

	// 同步表结构
	err = syncTables(engine)
	if err != nil {
		return nil, err
	}
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
