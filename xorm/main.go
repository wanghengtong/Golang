package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"xorm.io/xorm"
)

func main() {
	var (
		userName string = "root"
		password string = "root"
		host     string = "localhost"
		port     int    = 3306
		dbName   string = "test"
		charest  string = "utf8mb4"
	)
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", userName, password, host, port, dbName, charest)
	engine, err := xorm.NewEngine("mysql", dataSourceName)
	if err != nil {
		fmt.Println("数据库连接失败，请检查！")
	}

	type User struct {
		Id       int64
		Name     string
		Age      int
		Password string    `xorm:"varchar(200)"`
		Created  time.Time `xorm:"created"`
		Updated  time.Time `xorm:"updated"`
	}

	err = engine.Sync(new(User))
	if err != nil {
		fmt.Println("表结构同步失败，请检查！")
	}

}
