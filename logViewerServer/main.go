package main

import (
	"fmt"
	"log"
	"logViewerServer/dao"
	"logViewerServer/models/mysql"
	"logViewerServer/routers"
	"logViewerServer/setting"
	"logViewerServer/webhook"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		log.Printf("用法：%s conf/config-dev.ini", os.Args[0])
		return
	}

	// 加载配置文件
	if err := setting.Init(os.Args[1]); err != nil {
		log.Printf("配置文件加载失败, err：%s", err)
		return
	}

	// 连接数据库
	err := dao.InitMySQL(setting.Conf.MySQLConfig)
	if err != nil {
		log.Printf("mysql 连接初始化失败, err:%v\n", err)
		return
	}

	// 程序退出关闭数据库连接
	defer dao.CloseMySQL()

	// 数据库模型绑定
	dao.DB.AutoMigrate(
		&mysql.ProjectLogCommit{},
		&mysql.ProjectInfo{},
		&mysql.PmInfo{},
		&mysql.WhiteList{},
	)

	// 每30s检测一次是否有日志打包完成，有则在群通知提取申请人
	go func() {
		for {
			webhook.NotifyZipfile()
			time.Sleep(5 * time.Second)
		}
	}()

	//每半小时清理一次数据库提交记录
	//	1、删除提交时间大于3天的记录
	//  2、删除状态为"日志不存在"的提交记录
	go func() {
		for {
			mysql.DelWhereProjectLogCommit()
			time.Sleep(30 * time.Minute)
		}
	}()

	// 路由信息
	r := routers.SetupRouter()
	if err := r.Run(fmt.Sprintf(":%d", setting.Conf.Port)); err != nil {
		log.Printf("server startup failed, err:%v\n", err)
	}
}
