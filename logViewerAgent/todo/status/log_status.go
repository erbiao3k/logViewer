package status

import (
	"encoding/json"
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

type RespCheckLogStatus struct {
	Status bool                `json:"status"`
	Msg    string              `json:"msg"`
	Data   []map[string]string `json:"data"`
}

// LogStatus 查询当前项目是否有日志提取申请
func LogStatus() ([]map[string]string, error) {

	// 准备要提交给服务端查询是否有日志提取申请的信息
	localIp, err := net.GetInnerIP()
	if err != nil {
		log.Println(common.FuncName(), "获取当前服务器内网IP失败：", err)
		return nil, err
	}

	postData := url.Values{
		"AuthKey":     {common.AuthKey},
		"ProjectEnv":  {setting.Conf.Env},
		"ProjectName": {setting.Conf.Project},
		"LocalIP":     {localIp},
	}

	bodyBuffer := strings.NewReader(postData.Encode())

	c := proxy.HttpProxyClient()

	req, err := http.NewRequest(http.MethodPost, setting.Conf.Server+"/api/project/logstatus/check", bodyBuffer)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "logviewerOpsClient")

	resp, err := c.Do(req)
	if err != nil {
		log.Println(common.FuncName(), "get resp failed, err：", err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode > 399 {
		log.Println(common.FuncName(), "请求错误，err：", resp.Status)
		return nil, http.ErrAbortHandler
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(common.FuncName(), "get resp.Body failed, err：", err)
		return nil, err
	}

	var svc RespCheckLogStatus

	err = json.Unmarshal(b, &svc)
	if err != nil {
		log.Println(common.FuncName(), "json转struct失败!, err：", err)
		return nil, err
	}

	if len(svc.Data) == 0 {
		log.Println(common.FuncName(), "未收到日志提取申请", "日志提取查询结果：", svc.Data)
	}

	return svc.Data, nil
}
