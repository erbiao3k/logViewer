package project

import (
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

// AddWhiteList 添加当前服务器公网IP至服务端白名单
func AddWhiteList() (err error) {

	// 准备提交当前公网IP所需的信息
	publicIP, err := net.PublicIP()
	if err != nil {
		log.Println(common.FuncName(), "获取当前服务器公网IP失败：", err)
		return err
	}

	postData := url.Values{
		"AuthKey":     {common.AuthKey},
		"ProjectEnv":  {setting.Conf.Env},
		"ProjectName": {setting.Conf.Project},
		"PublicIP":    {publicIP},
	}

	body := strings.NewReader(postData.Encode())

	c := proxy.HttpProxyClient()

	req, err := http.NewRequest(http.MethodPost, setting.Conf.Server+"/api/whiteip", body)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "logviewerOpsClient")

	resp, err := c.Do(req)

	b, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode > 399 {
		log.Println(common.FuncName(), "请求错误，err：", string(b))
		return http.ErrAbortHandler
	}

	defer resp.Body.Close()
	log.Println("当前服务器公网IP为：", publicIP)
	return
}
