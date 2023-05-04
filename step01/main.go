package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-ldap/ldap"
	"net/http"
)

// 中间件进行身份验证
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		fmt.Println("获取到的token为: ", token)
		if token != "twgdh" {
			c.String(403, "身份验证不通过")
			c.Abort() // 终止当前请求，不会将请求转发给路由，所以 处理函数不会去执行
			return
		}
		c.Next()
	}

}
func main() {
	r := gin.Default()

	r.POST("/index", ldapMiddleware(), IndexHandler)

	r.GET("/home", Auth(), HomeHandler)
	fmt.Println("http://127.0.0.1:8002")
	r.Run(":8002")
}

func IndexHandler(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "使用ldap成功"})
}

func HomeHandler(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "欢迎来到管理后台页面"})
}

func ldapMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户名和密码

		username := c.Request.Header.Get("username")
		password := c.Request.Header.Get("password")
		//username := c.PostForm("username")
		//password := c.PostForm("password")

		fmt.Println("username: ", username, "password: ", password)
		// 连接 LDAP 服务器
		l, err := ldap.Dial("tcp", "192.168.3.102:389")
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer l.Close()

		// 绑定 LDAP 服务器
		err = l.Bind("cn="+username+",ou=people,dc=ailieyun,dc=com", password)

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// 身份验证通过，继续处理请求
		c.Next()
	}
}

