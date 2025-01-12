package controller

import (
	"admin/model"
	"admin/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"math"
	"net/http"
	"strconv"
	"xorm.io/xorm"
)

type UserController struct {
	Engine     *xorm.Engine
	UserSrvice *service.UserService
}

func (this *UserController) PageList(ctx *gin.Context) {
	w := ctx.Writer
	// 获取分页参数
	pageIndex, err := strconv.Atoi(ctx.Query("pageIndex"))
	if err != nil || pageIndex < 1 {
		pageIndex = 1
	}
	pageSize, err := strconv.Atoi(ctx.Query("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	searchQuery := ctx.Query("searchQuery")
	fmt.Println(searchQuery)

	// 计算偏移量
	offset := (pageIndex - 1) * pageSize

	// 分页查询用户列表
	var request = model.User{}
	request.Name = searchQuery
	users, totalUsers, err := this.UserSrvice.ListWithPagination(this.Engine, offset, pageSize, request)
	if err != nil {
		model.ReturnError(ctx, 500, "Internal server error")
		return
	}

	// 计算总页数
	totalPages := int(math.Ceil(float64(totalUsers) / float64(pageSize)))

	// 计算上一页和下一页
	prevPage := pageIndex - 1
	if prevPage < 1 {
		prevPage = 1
	}
	nextPage := pageIndex + 1
	if nextPage > totalPages {
		nextPage = totalPages
	}

	// 获取当前页的数据
	start := (pageIndex - 1) * pageSize
	end := start + pageSize
	if end > len(users) {
		end = len(users)
	}

	tmpl, err := template.ParseFiles("resources/templates/list.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	if admin, ok := ctx.Get("CurrentAdmin"); ok {
		data := map[string]interface{}{
			"Users":       users,
			"Admin":       admin,
			"CurrentPage": pageIndex,
			"PrevPage":    prevPage,
			"NextPage":    nextPage,
			"TotalPages":  totalPages,
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Error executing template", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
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

	if admin, ok := ctx.Get("CurrentAdmin"); ok {
		data := map[string]interface{}{
			"Users": users,
			"Admin": admin,
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Error executing template", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
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
	ctx.Redirect(http.StatusSeeOther, "/user/pageList")
}

func (this *UserController) ToEdit(ctx *gin.Context) {
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
	ctx.Redirect(http.StatusSeeOther, "/user/pageList")
}

func (this *UserController) ToAdd(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "addUser.html", nil)
}
func (this *UserController) Add(ctx *gin.Context) {
	// 解析表单数据
	var user model.User
	if err := ctx.ShouldBind(&user); err != nil {
		model.ReturnError(ctx, 5000, "invalid form data")
		return
	}
	// 新增用户信息
	count, err := this.UserSrvice.Add(this.Engine, user)
	if err != nil || count < 1 {
		fmt.Printf("新增失败: %v\n", err)
		model.ReturnError(ctx, 5001, "add failed")
		return
	}
	// 重定向回用户列表页面或其他页面
	ctx.Redirect(http.StatusSeeOther, "/user/pageList")
}
