package project

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"logViewerAgent/proxy"
	"logViewerAgent/public/common"
	"logViewerAgent/public/net"
	"logViewerAgent/setting"
	"net/http"
	"net/url"
	"strings"
)

type RespCreateProject struct {
	Status bool                `json:"status"`
	Msg    string              `json:"msg"`
	Data   []map[string]string `json:"data"`
}

// CreateProject 程序启动时，向服务器提交创建项目申请
func CreateProject(svcList []string) (err error) {

	// 准备要提交给服务端创建项目的信息
	localIP, err := net.GetInnerIP()
	if err != nil {
		log.Println(common.FuncName(), "获取当前服务器内网IP失败：", err)
		return err
	}

	projectName := setting.Conf.Project
	projectEnv := setting.Conf.Env
	md5Data := common.Encrypt(projectName + projectEnv + localIP)
	projectArea := setting.Conf.Area

	postData := url.Values{
		"AuthKey":     {common.AuthKey},
		"ProjectEnv":  {projectEnv},
		"ProjectName": {projectName},
		"ProjectArea": {projectArea},
		"Md5Data":     {md5Data},
	}

	for _, svc := range svcList {
		postData.Add("Svc"+common.Capitalize(svc), localIP)
	}

	bodyBuffer := strings.NewReader(postData.Encode())

	c := proxy.HttpProxyClient()

	req, err := http.NewRequest(http.MethodPost, setting.Conf.Server+"/api/project/create", bodyBuffer)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "logviewerOpsClient")

	resp, err := c.Do(req)

	if err != nil {
		log.Println(common.FuncName(), "获取响应失败, err：", err)
		return err
	}

	if resp.StatusCode > 399 {
		log.Println(common.FuncName(), "请求错误，err：", resp.Status)
		return errors.New("请求错误，状态码大于399")
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(common.FuncName(), "获取响应的Body失败, err：", err)
		return err
	}

	log.Println(string(b))
	return
}

// UpdateProjectLogCommitStatus 更新日志提取状态至服务端
func UpdateProjectLogCommitStatus(service, date, createTime, logStatus, localIP, logDownloadAddr string) (err error) {

	// 准备更新服务端日志提取状态所需的信息
	projectEnv := setting.Conf.Env
	projectName := setting.Conf.Project

	LogCommitMd5 := common.Encrypt(projectEnv + service + date + projectName + createTime)

	postData := url.Values{
		"AuthKey":         {common.AuthKey},
		"LogCommitMd5":    {LogCommitMd5},
		"LogStatus":       {logStatus},
		"CreateTime":      {createTime},
		"LogDownloadAddr": {logDownloadAddr},
		"LocalIP":         {localIP},
	}

	bodyBuffer := strings.NewReader(postData.Encode())

	c := proxy.HttpProxyClient()

	req, err := http.NewRequest(http.MethodPost, setting.Conf.Server+"/api/project/logstatus/update", bodyBuffer)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "logviewerOpsClient")

	resp, err := c.Do(req)

	if resp.StatusCode > 399 {
		log.Println(common.FuncName(), "请求错误，err：", resp.Status)
		return http.ErrAbortHandler
	}

	if err != nil {
		log.Println(common.FuncName(), "获取响应失败, err：", err)
		return err
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(common.FuncName(), "获取响应的Body失败, err：", err)
		return err
	}

	var svc RespCreateProject
	err = json.Unmarshal(b, &svc)
	if err != nil {
		log.Println(common.FuncName(), "json转struct失败, err：", err)
		return err
	}

	return
}
