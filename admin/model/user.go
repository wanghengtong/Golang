package model

import "time"

type User struct {
	Id      int64
	Name    string
	Age     int
	Sex     int
	Mobile  string
	Mail    string
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
}
