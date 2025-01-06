package auth

import (
	"github.com/gin-gonic/gin"
	"logViewerServer/models/mysql"
	"logViewerServer/pubilc"
	"net/http"
)

// WhiteIpAuth 判断接口在进行数据交换时，是否请求来自白名单
func WhiteIpAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.Request.ParseForm(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "errmsg": "解析Post表单失败", "data": err.Error()})
			return
		}

		// 从请求中获取客户端IP的值，判断是否在白名单中
		clientIp := c.ClientIP()

		whiteList, err := mysql.GetAllWhiteList()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "errmsg": "获取白名单失败", "data": err.Error()})
			return
		}

		var ipList []string
		for _, ip := range whiteList {
			ipList = append(ipList, ip.PublicIp)
		}
		if !pubilc.IsValueInSlice(clientIp, ipList) {
			c.JSON(http.StatusForbidden, gin.H{"status": false, "data": []string{}, "errmsg": "非白名单的网络请求"})
			c.Abort()
			return
		}
	}
}
