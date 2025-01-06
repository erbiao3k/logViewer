package auth

import (
	"github.com/gin-gonic/gin"
	"logViewerServer/pubilc"
	"net/http"
)

// ApiAuth 判断接口在进行数据交换时，是否带有认证key
func ApiAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.Request.ParseForm(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "errmsg": "解析Post表单失败", "errContent": err})
			return
		}

		// 从请求中获取AuthKey的值，判断与服务端的是否一致
		resp := c.Request
		key := resp.Form.Get("AuthKey")
		if len(key) == 0 || key != pubilc.AuthKey {
			c.JSON(http.StatusUnauthorized, gin.H{"status": false, "data": []string{}, "errmsg": "未认证的请求"})
			c.Abort()
			return
		}
	}
}
