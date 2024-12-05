package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// 设置 Cookie
func SetCookie(ctx *gin.Context, key string, value string) {
	// 定义 Cookie 的过期时间
	expirationTime := time.Now().Add(24 * time.Hour)

	// 创建 Cookie
	cookie := &http.Cookie{
		Name:     key,
		Value:    value,
		Expires:  expirationTime,
		HttpOnly: true,  // 设置 HttpOnly 属性以防止 JavaScript 访问
		Secure:   false, // 如果使用 HTTPS，请设置为 true
		Path:     "/",
	}

	// 设置 Cookie 到响应中
	http.SetCookie(ctx.Writer, cookie)
}

// 从 Cookie 中获取值
func GetCookie(ctx *gin.Context, cookieName string) (string, error) {
	cookie, err := ctx.Cookie(cookieName)
	if err != nil {
		logrus.Warnf("在 Cookie 中未找到%s", cookieName)
		return "", err
	}
	return cookie, err
}
