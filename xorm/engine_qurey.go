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

	user := User{Id: 1000, Name: "wanght", Age: 18, Password: "123456"}
	affected, _ := engine.Insert(&user)
	fmt.Println(affected)
	if affected >= 1 {
		fmt.Println("插入成功")
	}
	fmt.Println()

	user1 := User{Id: 10001, Name: "wanght", Age: 18, Password: "123456"}
	user2 := User{Id: 10002, Name: "wanght", Age: 18, Password: "123456"}
	affected, _ = engine.Insert(&user1, &user2)
	fmt.Println(affected)
	if affected >= 1 {
		fmt.Println("插入成功")
	}
	fmt.Println()

	var users []User
	users = append(users, User{Id: 10003, Name: "wanght", Age: 18, Password: "123456"})
	users = append(users, User{Id: 10004, Name: "wanght", Age: 18, Password: "123456"})
	affected, _ = engine.Insert(&users)
	fmt.Println(affected)
	if affected >= 1 {
		fmt.Println("插入成功")
	}
	fmt.Println()

	var users2 []*User
	if err := engine.SQL("select * from user").Find(&users2); err != nil {
		fmt.Printf("查询失败: %v\n", err)
	} else {
		for _, user := range users2 {
			fmt.Println(user)
		}
	}
	fmt.Println()

	var user3 User
	engine.SQL("select * from user where  id = 1000").Asc("id").Get(&user3)
	fmt.Println("user3:", user3)
	fmt.Println()

	user22 := make([]User, 0)
	err = engine.SQL("select * from user").Asc("id").Find(&user22)
	if err != nil {
		fmt.Printf("查询失败: %v\n", err)
	} else {
		for _, user := range user22 {
			fmt.Println("user22 Asc:", user)
		}
	}
	fmt.Println()

	user23 := make([]User, 0)
	err = engine.SQL("select * from user").Desc("id").Find(&user23)
	if err != nil {
		fmt.Printf("查询失败: %v\n", err)
	} else {
		for _, user := range user23 {
			fmt.Println("user23 Desc:", user)
		}
	}
	fmt.Println()

	var user4 []User
	err = engine.Where("id > ?", 1000).Find(&user4)
	if err != nil {
		fmt.Printf("查询失败: %v\n", err)
	} else {
		for _, user := range user4 {
			fmt.Println("user4::::", user)
		}
	}
	fmt.Println()

	result, _ := engine.Query("select * from user")
	fmt.Println(result)
	fmt.Println()

	result1, _ := engine.QueryString("select * from user")
	fmt.Println(result1)
	fmt.Println()

	result2, _ := engine.QueryInterface("select * from user")
	fmt.Println(result2)
	fmt.Println()

	user11 := User{}
	engine.Get(&user11)
	fmt.Println(user11)
	fmt.Println()

	var name string
	engine.Table(&user).Where("id=10001").Cols("age").Get(&name)
	fmt.Println(name)
	fmt.Println()

	user222 := make([]User, 0)
	engine.Where("id = 10002").And("age =18").Find(&user222)
	fmt.Println(user222)
	fmt.Println()

	usercount := User{Age: 18}
	count1, _ := engine.Count(&usercount)
	fmt.Println(count1)
	fmt.Println()

	engine.Iterate(&User{}, func(i int, bean interface{}) error {
		fmt.Println(i, bean)
		fmt.Println(bean.(*User))
		return nil
	})
	fmt.Println()

	rows, err := engine.Rows(&User{})
	defer rows.Close()
	userBean := new(User)
	for rows.Next() {
		err = rows.Scan(userBean)
		fmt.Println(userBean)
	}
	fmt.Println()

	engine.Close()
}
