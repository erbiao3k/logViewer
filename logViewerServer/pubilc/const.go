package pubilc

const (
	// 客户环境级别

	EnvProd = "prod"
	EnvPro  = "pro"
	EnvSit  = "sit"
	EnvUat  = "uat"
	EnvPoc  = "poc"

	// Salt 加密规则中的盐
	Salt = "EASj@.lj6J7yOt4K%j77"

	// 客户对接区域

	AreaSouth         = "南区"
	AreaNorthAreaWest = "北区&西区"
	AreaEast          = "东区"

	// 日志提取申请状态

	LogStatusSubmitted   = "日志申请已提交"
	LogStatusNotFound    = "日志不存在"
	LogStatusPacking     = "打包日志中"
	LogStatusPackFailed  = "打包日志失败"
	LogStatusPacked      = "日志打包完成"
	LogStatusUploading   = "上传日志中"
	LogStatusUplodFailed = "上传日志失败"
	LogStatusUploaded    = "日志上传完成"
	LogStatusSent        = "日志提取通知完成"
	LogStatusSendFailed  = "日志提取通知失败"
	LogStatusWaiting     = "正在处理中"

	// 服务列表

	SvcApi             = "api"
	SvcPortal          = "portal"
	SvcSchedule        = "schedule"
	SvcAdmin           = "admin"
	SvcFos             = "fos"
	SvcConvert         = "convert"
	SvcSign            = "sign"
	Svc2BaseServer     = "base-server"
	Svc2Contract       = "contract-service"
	Svc2Seal           = "seal-service"
	Svc2User           = "user-service"
	Svc2ScheduleServer = "schedule-server"
	Svc2SignServer     = "sign-server"
	Svc2Convert        = "convert-op2"
	Svc2ApiGateway     = "api-gateway"
	Svc2Web            = "web"
	// WorkwxPublicNotify 企业微信机器人webhook地址-公共渠道
	WorkwxPublicNotify = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=111"

	// WorkwxAdminNotify 企业微信机器人webhook地址-管理员通知渠道
	// 日志权限申请群
	WorkwxAdminNotify = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=333"

	// AuthKey 接口认证的key值
	AuthKey = "kqBBCoy6m45DvZU87$kfk3V65#V2FyWf$iZXjPSh@@#FaYCg39PbAGJZGsdAI8xOWGBz8dZSiPEw^7B6wJ1&kkW$zozm3284u3wa7"
)

var (
	AllEnv     = []string{EnvProd, EnvSit, EnvUat, EnvPoc, EnvPro}
	AllService = []string{SvcApi, SvcAdmin, SvcPortal, SvcSchedule, SvcFos, SvcConvert, SvcSign, Svc2BaseServer, Svc2User, Svc2Seal, Svc2Contract, Svc2SignServer, Svc2ScheduleServer, Svc2Convert, Svc2ApiGateway, Svc2Web}
	Op2Service = []string{Svc2BaseServer, Svc2User, Svc2Seal, Svc2Contract, Svc2SignServer, Svc2ScheduleServer, Svc2Convert, Svc2ApiGateway, Svc2Web}
)
