package net

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"logViewerAgent/proxy"
	"logViewerAgent/public/common"
	"logViewerAgent/setting"
	"net/http"
	"net/url"
	"strings"
)

type Ip struct {
	Data string `json:"data"`
}

func PublicIP() (string, error) {
	postData := url.Values{"AuthKey": {common.AuthKey}}
	bodyBuffer := strings.NewReader(postData.Encode())

	c := proxy.HttpProxyClient()

	req, err := http.NewRequest(http.MethodPost, setting.Conf.Server+"/api/publicip", bodyBuffer)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "logviewerOpsClient")

	resp, err := c.Do(req)

	if resp.StatusCode > 399 {
		log.Println(common.FuncName(), "请求错误，err：", resp.Status)
		return "", http.ErrAbortHandler
	}

	if err != nil {
		log.Println(common.FuncName(), "获取响应失败, err：", err)
		return "", err
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(common.FuncName(), "获取响应的Body失败, err：", err)
		return "", err
	}

	var ip Ip
	err = json.Unmarshal(b, &ip)
	if err != nil {
		log.Println(common.FuncName(), "json转struct失败!, err：", err, "Body内容：", string(b))
		return "", err
	}

	return ip.Data, nil
}
