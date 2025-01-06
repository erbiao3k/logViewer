package setting

import (
	"gopkg.in/ini.v1"
)

var Conf = new(AppConfig)

// AppConfig 应用程序配置
type AppConfig struct {
	Release      bool   `ini:"release"`
	Port         int    `ini:"port"`
	Path         string `ini:"path"`
	DownloadAddr string `ini:"server"`
	*MySQLConfig `ini:"mysql"`
	*MailConfig  `ini:"mail"`
	*RedisConfig `ini:"redis"`
}

// MySQLConfig 数据库配置
type MySQLConfig struct {
	User     string `ini:"user"`
	Password string `ini:"password"`
	DB       string `ini:"db"`
	Host     string `ini:"host"`
	Port     int    `ini:"port"`
}

// RedisConfig 数据库配置
type RedisConfig struct {
	Host     string `ini:"redis_host"`
	Password string `ini:"redis_password"`
	DB       int    `ini:"redis_db"`
	Port     string `ini:"redis_port"`
	PoolSize int    `ini:"redis_pool_size"`
}

// MailConfig 邮箱配置
type MailConfig struct {
	MailSenderAccount  string `ini:"mail_user"`
	MailSenderPassword string `ini:"mail_password"`
	MailServerHost     string `ini:"smtp_host"`
	MailServerPort     int    `ini:"smtp_port"`
}

// Init 初始化函数
func Init(file string) error {
	return ini.MapTo(Conf, file)
}
