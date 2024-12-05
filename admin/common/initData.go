package common

import (
	"admin/model"
	"github.com/sirupsen/logrus"
	"xorm.io/xorm"
)

func InitData(engine *xorm.Engine) {
	initAdminData(engine)
	initUserData(engine)
}

func initAdminData(engine *xorm.Engine) {
	var admins []model.Admin
	admins = append(admins, model.Admin{Id: 1, AdminName: "超级管理员", LoginName: "superadmin", Password: "superadmin"})
	admins = append(admins, model.Admin{Id: 2, AdminName: "系统管理员", LoginName: "sysadmin", Password: "sysadmin"})
	admins = append(admins, model.Admin{Id: 4, AdminName: "安全管理员", LoginName: "secadmin", Password: "secadmin"})
	admins = append(admins, model.Admin{Id: 3, AdminName: "审计管理员", LoginName: "logadmin", Password: "logadmin"})
	affected, err := engine.Insert(&admins)
	if err != nil {
		logrus.Errorf("插入管理员数据失败: %v", err)
		return
	}
	logrus.Infof("插入管理员数据成功，受影响行数: %d", affected)
}

func initUserData(engine *xorm.Engine) {
	var users []model.User
	users = append(users, model.User{Id: 1, Name: "张三", Age: 18, Sex: 1, Mobile: "15011112222", Mail: "zhangsan@qq.com"})
	users = append(users, model.User{Id: 2, Name: "李四", Age: 18, Sex: 0, Mobile: "15011112222", Mail: "lisi@qq.com"})
	users = append(users, model.User{Id: 3, Name: "王五", Age: 18, Sex: 1, Mobile: "15011112222", Mail: "wangwu@qq.com"})
	affected, err := engine.Insert(&users)
	if err != nil {
		logrus.Errorf("插入用户数据失败: %v", err)
		return
	}
	logrus.Infof("插入用户数据成功，受影响行数: %d", affected)
}
