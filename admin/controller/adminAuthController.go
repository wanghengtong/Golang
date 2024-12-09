package controller

import (
	"admin/model"
	"admin/service"
	"admin/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"time"
	"xorm.io/xorm"
)

type AdminAuthController struct {
	Engine     *xorm.Engine
	UserSrvice *service.UserSrvice
}

func (adminAuthcontroller *AdminAuthController) Index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.html", nil)
}

func (adminAuthcontroller *AdminAuthController) Logout(ctx *gin.Context) {
	// 清除认证 Cookie
	cookieName := viper.GetString("auth.jwt.cookie.name")
	utils.ClearCookie(ctx, cookieName)
	ctx.HTML(http.StatusOK, "login.html", nil)
}

func (this *AdminAuthController) LoginToIndex(ctx *gin.Context) {
	loginName := ctx.DefaultPostForm("loginName", "")
	fmt.Println(loginName)
	password := ctx.DefaultPostForm("password", "")
	fmt.Println(password)
	if loginName == "" || password == "" {
		model.ReturnError(ctx, 500, "loginName and password are required")
		return
	}
	var admin model.Admin
	has, err := this.Engine.Where("login_name = ? AND password = ?", loginName, password).Get(&admin)
	if err != nil {
		model.ReturnError(ctx, 500, "Internal server error")
		return
	}
	if !has {
		model.ReturnError(ctx, 2000, "Invalid login credentials")
		return
	}
	admin.LoginTime = time.Now()
	this.Engine.Update(&admin)

	// 生成 JWT Token
	token := utils.GenerateJWT(admin)

	// 设置 Cookie
	cookieName := viper.GetString("auth.jwt.cookie.name")
	utils.SetCookie(ctx, cookieName, token)

	// 重定向回用户列表页面或其他页面
	ctx.Redirect(http.StatusSeeOther, "/user/list")
}
