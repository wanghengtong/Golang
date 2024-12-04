package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)
import "github.com/thinkerou/favicon"

func myHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		Authorization := c.Request.Header.Get("Authorization")
		log.Println("===================", Authorization)
		if Authorization == "" {
			c.Writer.Header().Set("Authorization", "")
			log.Println("未认证！")
			c.HTML(http.StatusOK, "401.html", nil)
			c.Abort()
		} else {
			c.Set("Authorization", c.Request.Header.Get("Authorization"))
			c.Next()
		}
	}
}

func main() {
	ginServer := gin.Default()
	ginServer.Use(myHandler())

	ginServer.Use(favicon.New("./favicon.ico"))

	ginServer.LoadHTMLGlob("templates/*")

	ginServer.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})

	ginServer.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"msg": "Hello World",
		})
	})

	ginServer.GET("/user/info", func(c *gin.Context) {
		userid := c.Query("userid")
		username := c.Query("username")
		c.JSON(http.StatusOK, gin.H{
			"userid":   userid,
			"username": username,
		})
	})

	ginServer.GET("/user/info/:userid/:username", func(c *gin.Context) {
		userid := c.Param("userid")
		username := c.Param("username")
		c.JSON(http.StatusOK, gin.H{
			"userid":   userid,
			"username": username,
		})
	})

	ginServer.POST("/json", func(c *gin.Context) {
		data, _ := c.GetRawData()
		var m = make(map[string]interface{})
		_ = json.Unmarshal(data, &m)
		c.JSON(http.StatusOK, m)
	})

	ginServer.POST("/user/add/PostForm", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		c.JSON(http.StatusOK, gin.H{
			"username": username,
			"password": password,
			"msg":      "success",
		})
	})

	ginServer.GET("/redirect", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "https://www.baidu.com")
	})

	ginServer.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusOK, "404.html", nil)
	})

	userGroup := ginServer.Group("/user")
	{
		userGroup.GET("/get", myHandler(), func(c *gin.Context) {
			Authorization := c.MustGet("Authorization").(string)
			c.Writer.Header().Set("Authorization", Authorization)
			c.JSON(http.StatusOK, gin.H{
				"username": "王恒通",
				"password": 123456,
				"msg":      "success",
			})
		})
		userGroup.POST("/add", func(c *gin.Context) {

		})
		userGroup.DELETE("delete/:id")
	}

	err := ginServer.Run(":8080")
	if err != nil {
		return
	}

}
