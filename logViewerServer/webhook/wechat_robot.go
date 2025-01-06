package webhook

import (
	"fmt"
	"log"
	"logViewerServer/models/mysql"
	"logViewerServer/pubilc"
	"logViewerServer/setting"
	"net/http"
	"net/url"
	"strings"
)

const contentType = "application/json"

// NotifyRegisterReview 通知审核人激活账号
func NotifyRegisterReview(area, email, reviewAddr string) (err error) {
	data := fmt.Sprintf(`{"msgtype": "text", "text": {"content": "%s的%s提交了一个日志系统注册申请，请点击链接激活账号：%s","mentioned_mobile_list":["@all"]}}`, area, email, reviewAddr)

	if _, err := http.Post(pubilc.WorkwxAdminNotify, contentType, strings.NewReader(data)); err != nil {
		log.Println("账号激活信息发送失败，err：", err)
		return err
	}
	return
}

// NotifyRegisterEnable 通知注册人账号已激活
func NotifyRegisterEnable(email, phone string) (err error) {
	data := fmt.Sprintf(`{"msgtype": "text", "text": {"content": "日志系统账号【%s】已激活，登陆地址：%s","mentioned_mobile_list":["%s"]}}`, email, setting.Conf.DownloadAddr, phone)

	if _, err := http.Post(pubilc.WorkwxPublicNotify, contentType, strings.NewReader(data)); err != nil {
		log.Println("通知注册人账号已激活失败，err：", err)
		return err
	}
	return
}

// SendLogAddr 通知日志下载地址到企业wx机器人
func SendLogAddr(logAddr, projectInfo, phone string) (err error) {
	data := fmt.Sprintf(`{"msgtype": "text", "text": {"content": "[社会社会] 您有一份日志提取完成\n\n[社会社会] 日志信息：%s\n\n[社会社会] 下载地址：%s","mentioned_mobile_list":["%s"]}}`, projectInfo, logAddr, phone)

	if _, err := http.Post(pubilc.WorkwxPublicNotify, contentType, strings.NewReader(data)); err != nil {
		log.Println("日志提醒发送失败，err：", err)
		return err
	}
	log.Println("日志提醒发送成功，日志地址：", logAddr)
	return
}

// NotifyZipfile 通知日志下载地址到企业wx机器人
func NotifyZipfile() {
	notifyInfo, err := mysql.GetProjectLogCommitNotify()

	if err != nil || len(notifyInfo) == 0 {
		log.Println("未查询到需要通知的提取日志，err：", err)
	}

	for _, info := range notifyInfo {

		pmInfo, err := mysql.GetWherePmInfo("pm_email", info.PmEmail)
		if err != nil {
			log.Println("发送消息到企业微信机器人时，通过邮箱查询人员的手机号失败，err：", err)
			return
		}
		zipfileDownloadAddr := url.Values{"filename": {info.LogDownloadAddr}}

		projectInfo := info.LogDate + "-" + info.ProjectName + "-" + info.ProjectEnv + "-" + info.SvcName + "-" + info.SvcAddr

		if err := SendLogAddr(setting.Conf.DownloadAddr+"/log/download"+"?"+zipfileDownloadAddr.Encode(), projectInfo, pmInfo[0].PmPhone); err != nil {
			log.Println("发送消息到企业微信机器人失败，err：", err)
			mysql.UpdateProjectLogCommitLogAddrAfterNotify(info.LogDownloadAddr, info.PmEmail, pubilc.LogStatusSendFailed)
		} else {
			mysql.UpdateProjectLogCommitLogAddrAfterNotify(info.LogDownloadAddr, info.PmEmail, pubilc.LogStatusSent)
		}
	}
}
