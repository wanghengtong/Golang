package model

import "time"

type Admin struct {
	Id        int64
	AdminName string `xorm:"varchar(200)"`
	LoginName string
	Password  string    `xorm:"varchar(200)"`
	LoginTime time.Time `xorm:"login_time"`
	Created   time.Time `xorm:"created"`
	Updated   time.Time `xorm:"updated"`
}
