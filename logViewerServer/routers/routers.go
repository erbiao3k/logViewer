package routers

import (
	"github.com/gin-gonic/gin"
	"logViewerServer/controller"
	"logViewerServer/middleware/auth"
	"logViewerServer/setting"
	"net/http"
)

func SetupRouter() *gin.Engine {
	if setting.Conf.Release {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// 告诉gin框架去哪里找模板文件
	r.LoadHTMLGlob("templates/*")

	// 告诉gin框架模板文件引用的静态文件去哪里找
	r.Static("/logviewer/static", "static")

	// 路由前缀
	prefixGroup := r.Group("/logviewer")

	// 登陆
	prefixGroup.Any("/", controller.LoginHandler)

	// 注册
	prefixGroup.Any("/register", controller.RegisterHandler)

	// 激活
	prefixGroup.GET("/ae", auth.AdminPageAuth(), controller.AccountEnableHandler)

	// 路由组-浏览器页面
	logGroup := prefixGroup.Group("/log")
	logGroup.Use(auth.PageAuth())
	{
		// 首页
		logGroup.GET("/", controller.IndexHandler)

		// 日志提取申请
		logGroup.Any("/commit", controller.LogCommitHandler)

		// 查询所有项目信息
		logGroup.GET("/pi", controller.ProjectListHandler)

		// 日提取申请的状态查询
		logGroup.GET("/status", controller.LogStatusHandler)

		// 日志下载
		logGroup.GET("/download", controller.FileDownloadHandler)

	}

	//日志文件上传接口
	r.MaxMultipartMemory = 4096 << 20
	prefixGroup.POST("/api/upload", auth.WhiteIpAuth(), controller.FileUploadHandler)

	// 路由组-api接口
	apiGroup := prefixGroup.Group("/api")
	apiGroup.Use(auth.ApiAuth())
	{
		// 返回客户端公网IP
		apiGroup.POST("/publicip", controller.PublicIp)

		// 添加客户端公网IP至白名单
		apiGroup.POST("/whiteip", controller.AddWhiteList)

		// 获取支持上传日志的服务的接口
		apiGroup.POST("/svclist", controller.SvcListHandler)

		// 路由组-项目接口
		apiProjectGroup := apiGroup.Group("/project")
		{
			// 客户端启动时，刷新客户端运行的服务信息
			apiProjectGroup.POST("/create", controller.CreateProject)

			// 路由组-项目状态管理
			apiProjectLogStatus := apiProjectGroup.Group("/logstatus")
			{
				// 查询-日志提取状态
				apiProjectLogStatus.POST("/check", controller.CheckLogStatusHandler)

				// 更新-日志提取状态
				apiProjectLogStatus.POST("/update", controller.UpdateLogStatusHandler)
			}

		}

	}

	// 统一的404错误页面
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"errmsg": "访问的资源不存在",
		})
	})

	// 未知的方法页面
	r.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"errmsg": "MethodNotAllowed",
		})
	})

	return r
}
