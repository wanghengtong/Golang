package service

import (
	"fmt"
	"iam-auth/model"
	"log"
	"xorm.io/xorm"
)

type UserSrvice struct {
	Engine *xorm.Engine
}

func (userSrvice *UserSrvice) List(engine *xorm.Engine) ([]model.User, error) {
	users := []model.User{}
	err := engine.Find(&users)
	if err != nil {
		log.Println("查询用户列表失败:", err)
		return nil, err
	}
	return users, nil
}

func (userSrvice *UserSrvice) Delete(engine *xorm.Engine, userId int64) (int64, error) {
	count, err := engine.ID(userId).Delete(&model.User{})
	if err != nil {
		fmt.Printf("删除失败: %v\n", err)
	}
	fmt.Printf("删除成功, 影响行数: %v\n", count)
	return count, err
}

func (userSrvice *UserSrvice) Get(engine *xorm.Engine, userId int64) (model.User, error) {
	var user model.User
	_, err := engine.ID(userId).Get(&user)
	return user, err
}

func (userSrvice *UserSrvice) Update(engine *xorm.Engine, user model.User) (int64, error) {
	count, err := engine.ID(user.Id).AllCols().Update(&user)
	return count, err
}
