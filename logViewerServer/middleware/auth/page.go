package auth

import (
	"github.com/gin-gonic/gin"
	"logViewerServer/models/redis"
	"net/http"
)

// PageAuth 判断浏览器的loginSession和数据库中通过密码计算的Session是否一致，否则跳转至登陆页面;
func PageAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.FullPath() != "/logviewer/" {
			clientSession, clientSessionErr := c.Cookie("loginSession")
			_, ServerSessionErr := redis.StringGet(clientSession)

			if clientSessionErr != nil || ServerSessionErr != nil {
				c.Redirect(http.StatusMovedPermanently, "/logviewer")
				c.Abort()
				return
			}
		}
	}
}
