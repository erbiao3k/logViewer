package proxy

import (
	"log"
	"logViewerAgent/setting"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// tcpPortCheck 检测tcp地址端口连通性
func tcpPortCheck(IpPort string) error {
	_, err := net.DialTimeout("tcp", IpPort, 10*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// HttpProxyClient 给http请求设置代理
func HttpProxyClient() *http.Client {

	proxyAddr := setting.Conf.ProxyUrl

	if len(proxyAddr) != 0 {
		err := tcpPortCheck(strings.Split(proxyAddr, "//")[1])
		if err != nil {
			log.Println("代理地址：", proxyAddr, "，访问异常：", err)
		}
		proxy := func(_ *http.Request) (*url.URL, error) {
			return url.Parse(proxyAddr)
		}

		transport := &http.Transport{Proxy: proxy}
		client := &http.Client{Transport: transport}
		return client
	}
	return &http.Client{}
}
