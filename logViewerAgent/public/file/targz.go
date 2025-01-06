package file

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/Chain-Zhang/pinyin"
	"io"
	"log"
	"logViewerAgent/public/common"
	"logViewerAgent/public/net"
	"logViewerAgent/setting"
	"os"
)

func Targz(srcPath, svcName, packageTime string) (string, error) {

	localIP, err := net.GetInnerIP()
	if err != nil {
		log.Println(common.FuncName(), "获取当前服务器内网IP失败：", err)
		return "", err
	}
	projectNameChinese := setting.Conf.Project
	projectNamePinyin, err := pinyin.New(projectNameChinese).Split("").Convert()
	if err != nil {
		log.Println(common.FuncName(), "中文转拼音失败，错误信息：", err, "项目原名称为：", projectNameChinese)
		return "", err
	}

	dist := "storage/【" + projectNamePinyin + "】-【" + svcName + "】-【" + setting.Conf.Env + "】-【" + packageTime + "】-【" + localIP + "】.tar.gz"

	// logflie write
	fw, err := os.Create(dist)
	if err != nil {
		panic(err)
	}
	defer fw.Close()

	// gzip write
	gw := gzip.NewWriter(fw)
	defer gw.Close()

	// tar write
	tw := tar.NewWriter(gw)
	defer tw.Close()

	// 打开文件夹
	dir, err := os.Open(srcPath)
	if err != nil {
		panic(nil)
	}
	defer dir.Close()

	// 读取文件列表
	fis, err := dir.Readdir(0)
	if err != nil {
		panic(err)
	}

	// 遍历文件列表
	for _, fi := range fis {
		// 逃过文件夹, 我这里就不递归了
		if fi.IsDir() {
			continue
		}

		// 打印文件名称
		fmt.Println(fi.Name())

		// 打开文件
		fr, err := os.Open(dir.Name() + "/" + fi.Name())
		if err != nil {
			panic(err)
		}
		defer fr.Close()

		// 信息头
		h := new(tar.Header)
		h.Name = fi.Name()
		h.Size = fi.Size()
		h.Mode = int64(fi.Mode())
		h.ModTime = fi.ModTime()

		// 写信息头
		err = tw.WriteHeader(h)
		if err != nil {
			panic(err)
		}

		// 写文件
		_, err = io.Copy(tw, fr)
		if err != nil {
			panic(err)
		}
	}
	return dist, nil
}
