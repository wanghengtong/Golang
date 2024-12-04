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
	fmt.Println()

	type User struct {
		Id       int64
		Name     string
		Age      int
		Password string    `xorm:"varchar(200)"`
		Created  time.Time `xorm:"created"`
		Updated  time.Time `xorm:"updated"`
	}
	fmt.Println()

	var user User
	//user.Id = 10001
	//engine.Delete(&user)
	//user.Name = "wanght11"
	// user.Age = 0
	// engine.ID(1000).Cols("age").Update(&user)
	// engine.ID(1000).Update(&user)
	// engine.ID(1000).MustCols("age").Update(&user)
	engine.ID(1000).AllCols().Update(&user)

	engine.Close()
}
