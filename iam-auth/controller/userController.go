package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"iam-auth/model"
	"iam-auth/service"
	"net/http"
	"strconv"
	"xorm.io/xorm"
)

type UserController struct {
	Engine     *xorm.Engine
	UserSrvice *service.UserSrvice
}

func (this *UserController) List(ctx *gin.Context) {
	w := ctx.Writer
	users, err := this.UserSrvice.List(this.Engine)
	if err != nil {
		model.ReturnError(ctx, 500, "Internal server error")
		return
	}
	tmpl, err := template.ParseFiles("resources/templates/list.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{
		"Users": users,
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}

func (this *UserController) Delete(ctx *gin.Context) {
	userIdString := ctx.Query("id")
	if userIdString == "" {
		model.ReturnError(ctx, 5000, "id is required")
		return
	}
	userId, err := strconv.ParseInt(userIdString, 10, 64)
	if err != nil {
		model.ReturnError(ctx, 5000, "invalid user id")
		return
	}
	count, err := this.UserSrvice.Delete(this.Engine, userId)
	if err != nil || count < 1 {
		fmt.Printf("删除失败: %v\n", err)
		model.ReturnError(ctx, 5001, "delete failed")
		return
	}
	// 重定向回用户列表页面或其他页面
	ctx.Redirect(http.StatusSeeOther, "/user/list")
}

func (this *UserController) Edit(ctx *gin.Context) {
	w := ctx.Writer
	userIdString := ctx.Query("id")
	if userIdString == "" {
		model.ReturnError(ctx, 5000, "id is required")
		return
	}
	userId, err := strconv.ParseInt(userIdString, 10, 64)
	if err != nil {
		model.ReturnError(ctx, 5000, "invalid user id")
		return
	}
	user, err := this.UserSrvice.Get(this.Engine, userId)
	if err != nil {
		fmt.Printf("用户不存在: %v\n", err)
		model.ReturnError(ctx, 4004, "user not exist")
		return
	}
	tmpl, err := template.ParseFiles("resources/templates/editUser.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{
		"User": user,
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}

func (this *UserController) Update(ctx *gin.Context) {
	// 解析表单数据
	var user model.User
	if err := ctx.ShouldBind(&user); err != nil {
		model.ReturnError(ctx, 5000, "invalid form data")
		return
	}
	// 更新用户信息
	count, err := this.UserSrvice.Update(this.Engine, user)
	if err != nil || count < 1 {
		fmt.Printf("更新失败: %v\n", err)
		model.ReturnError(ctx, 5001, "update failed")
		return
	}
	// 重定向回用户列表页面或其他页面
	ctx.Redirect(http.StatusSeeOther, "/user/list")
}