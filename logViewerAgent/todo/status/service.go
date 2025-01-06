package status

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"logViewerAgent/proxy"
	"logViewerAgent/public/common"
	"logViewerAgent/public/file"
	"logViewerAgent/setting"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

// SvcInfo 当前服务器运行的服务清单
type SvcInfo struct {
	Status bool     `json:"status"`
	Data   []string `json:"data"`
	Errmsg string   `json:"errmsg"`
}

// GetSvcList 整理该服务器指定路径下运行了哪些服务
func GetSvcList(svcPath string) (svcList []string, err error) {
	// 获取配置文件中指定位置的所有文件和目录
	filepathNames, err := filepath.Glob(filepath.Join(svcPath, "*"))
	if err != nil {
		log.Println(common.FuncName(), "获取配置文件指定位置目录失败, err：", err)
		return nil, err
	}
	// 获取服务端现在支持提取日志的服务
	postData := url.Values{"AuthKey": {common.AuthKey}}
	bodyBuffer := strings.NewReader(postData.Encode())

	c := proxy.HttpProxyClient()

	req, err := http.NewRequest(http.MethodPost, setting.Conf.Server+"/api/svclist", bodyBuffer)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "logviewerOpsClient")

	resp, err := c.Do(req)

	if err != nil {
		log.Println(common.FuncName(), "获取响应信息失败, err：", err, "Body内容：", resp)
		return nil, err
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(common.FuncName(), "获取响应的Body失败, err：", err)
		return nil, err
	}

	if resp.StatusCode > 399 {
		log.Println(common.FuncName(), "请求错误，err：", resp.Status)
		return nil, http.ErrAbortHandler
	}

	var svc SvcInfo

	err = json.Unmarshal(b, &svc)
	if err != nil {
		log.Println(common.FuncName(), "json转struct失败!, err：", err, "Body内容：", string(b))
		return nil, err
	}

	// 通过比对服务端支持提取日志的服务清单
	for _, path := range filepathNames {

		path = strings.ReplaceAll(path, `\`, "/")

		svcName := strings.Split(path, "/")[len(strings.Split(path, "/"))-1]
		if !file.IsDir(path) {
			continue
		}

		if common.IsValueInSlice(svcName, svc.Data) {
			svcList = append(svcList, svcName)
		}
		if strings.Contains(svcName, "tomcat-convert-8080") {
			svcList = append(svcList, common.SvcConvert)
		}
		if strings.Contains(svcName, "tomcat-sign-8887") {
			svcList = append(svcList, common.SvcSign)
		}
		if strings.Contains(svcName, "base-doc-sign-op2") {
			svcList = append(svcList, common.Svc2SignServer)
		}
	}

	if len(svcList) == 0 {
		return nil, fmt.Errorf(common.FuncName(), "未找到任何服务")
	}

	return svcList, nil
}
