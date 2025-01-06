package setting

import (
	"gopkg.in/ini.v1"
)

var Conf = new(AppConfig)

// AppConfig 应用程序配置
type AppConfig struct {
	Project             string `ini:"project"`
	Env                 string `ini:"env"`
	WebPath             string `ini:"web_path"`
	FosPath             string `ini:"fos_path"`
	BasePath            string `ini:"base_path"`
	Op2ContainerLogPath string `ini:"op2_container_log_path"`
	Op2ConvertPath      string `ini:"op2_convert_path"`
	Area                string `ini:"area"`
	Server              string `ini:"server"`
	ProxyUrl            string `ini:"proxy_url"`
}

func Init(file string) error {
	return ini.MapTo(Conf, file)
}
