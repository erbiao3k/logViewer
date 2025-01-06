package logflie

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"logViewerAgent/proxy"
	"logViewerAgent/public/common"
	"logViewerAgent/setting"
	"mime/multipart"
	"net/http"
	"os"
)

// FileUpload 上传文件至服务端
func FileUpload(filePath string) error {

	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)

	fileWriter, _ := bodyWriter.CreateFormFile("file", filePath)

	file, err := os.Open(filePath)
	if err != nil {
		log.Println(common.FuncName(), "文件上传时，文件打开失败：", err)
		return err
	}
	defer file.Close()

	_, err = io.Copy(fileWriter, file)
	if err != nil {
		log.Println(common.FuncName(), "文件上传时，文件读写失败：", err)
		return err
	}

	bodyWriter.Close()

	c := proxy.HttpProxyClient()

	req, err := http.NewRequest(http.MethodPost, setting.Conf.Server+"/api/upload", bodyBuffer)

	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())
	req.Header.Set("User-Agent", "logviewerOpsClient")

	resp, err := c.Do(req)

	respBody, _ := ioutil.ReadAll(resp.Body)

	log.Printf(string(respBody))

	if resp.StatusCode > 399 {
		log.Println(common.FuncName(), "请求错误，err：", resp.Status)
		return http.ErrAbortHandler
	}

	return nil
}
