package common

import (
	"admin/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
)

func InitMiddleware() gin.HandlerFunc {
	return myHandler()
}

func myHandler() gin.HandlerFunc {
	// 定义放行路径列表
	skipAuthPaths := map[string]struct{}{
		"/admin/index": {},
		"/admin/login": {},
	}
	return func(c *gin.Context) {
		// 获取请求路径
		requestPath := c.Request.URL.Path
		// 检查请求路径是否在放行列表中
		if _, exists := skipAuthPaths[requestPath]; exists {
			// 如果在放行列表中，跳过认证检查
			c.Next()
			return
		}

		// 从请求头中获取 jwtToken
		//Authorization := c.Request.Header.Get("Authorization")
		//logrus.Debugf("=================== %s", Authorization)
		//if Authorization == "" {
		//	c.Writer.Header().Set("Authorization", "")
		//	logrus.Warn("未认证！")
		//	c.HTML(http.StatusOK, "401.html", nil)
		//	c.Abort()
		//} else {
		//	c.Set("Authorization", c.Request.Header.Get("Authorization"))
		//	c.Next()
		//}

		// 从 Cookie 中获取 jwtToken
		cookieName := viper.GetString("auth.jwt.cookie.name")
		cookie, err := utils.GetCookie(c, cookieName)
		if err != nil || cookie == "" {
			logrus.Warn("未认证！")
			c.HTML(http.StatusOK, "401.html", nil)
			c.Abort()
			return
		}
		// 验证 Token
		validToken, err := utils.ParseJWT(cookie)
		if validToken {
			c.Next()
		} else {
			logrus.Error("jwtToken解析失败")
			c.HTML(http.StatusOK, "401.html", nil)
			c.Abort()
		}
	}
}
