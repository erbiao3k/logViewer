package common

const (
	// Salt 加密规则中的盐
	Salt = "EASj@.ljOJ7yOt4K%j77"

	// 客户环境级别

	EnvProd = "prod"
	EnvPro  = "pro"
	EnvSit  = "sit"
	EnvUat  = "uat"
	EnvPoc  = "poc"

	// 日志提取申请状态

	LogStatusSubmitted   = "日志申请已提交"
	LogStatusPacking     = "打包日志中"
	LogStatusPackFailed  = "打包日志失败"
	LogStatusPacked      = "日志打包完成"
	LogStatusUploading   = "上传日志中"
	LogStatusNotFound    = "日志不存在"
	LogStatusUplodFailed = "上传日志失败"
	LogStatusUploaded    = "日志上传完成"
	LogStatusSent        = "日志提取通知完成"
	LogStatusSendFailed  = "日志提取通知失败"

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
	// AuthKey 接口认证的key值
	AuthKey = "iqBBCoy6m45DvZU87$kfkcV65A#V2FyW4$iZXjPSh@@#FaYCg39PbAGJZGsdAI8xOWGBz8dZSiPEw^7B6wJ1&kkW$zoZm3984u3wm7"

	// 客户对接区域

	AreaSouth = "南区"
	AreaNorth = "北区"
	AreaEast  = "东区"
	AreaWest  = "西区"
)

var (
	AllArea = []string{AreaEast, AreaNorth, AreaWest, AreaSouth}
	AllEnv  = []string{EnvProd, EnvSit, EnvUat, EnvPoc, EnvPro}
)
