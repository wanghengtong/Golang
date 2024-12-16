package service

import (
	"admin/model"
	"fmt"
	"log"
	"xorm.io/xorm"
)

type UserService struct {
	Engine *xorm.Engine
}

func (userSrvice *UserService) List(engine *xorm.Engine) ([]model.User, error) {
	users := []model.User{}
	err := engine.Find(&users)
	if err != nil {
		log.Println("查询用户列表失败:", err)
		return nil, err
	}
	return users, nil
}

func (userSrvice *UserService) ListWithPagination(engine *xorm.Engine, offset, pageSize int, request model.User) ([]model.User, int64, error) {
	var users []model.User
	var totalUsers int64
	// 查询总记录数
	count, err := engine.Count(&model.User{})
	if err != nil {
		return nil, 0, err
	}
	totalUsers = count

	// 分页查询用户列表
	if err := engine.Where("name LIKE ?", "%"+request.Name+"%").Limit(pageSize, offset).Find(&users); err != nil {
		return nil, 0, err
	}
	return users, totalUsers, nil
}

func (userSrvice *UserService) Delete(engine *xorm.Engine, userId int64) (int64, error) {
	count, err := engine.ID(userId).Delete(&model.User{})
	if err != nil {
		fmt.Printf("删除失败: %v\n", err)
	}
	fmt.Printf("删除成功, 影响行数: %v\n", count)
	return count, err
}

func (userSrvice *UserService) Get(engine *xorm.Engine, userId int64) (model.User, error) {
	var user model.User
	_, err := engine.ID(userId).Get(&user)
	return user, err
}

func (userSrvice *UserService) Update(engine *xorm.Engine, user model.User) (int64, error) {
	count, err := engine.ID(user.Id).AllCols().Update(&user)
	return count, err
}

func (userSrvice *UserService) Add(engine *xorm.Engine, user model.User) (int64, error) {
	count, err := engine.Insert(&user)
	return count, err
}
