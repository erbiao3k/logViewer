package file

import (
	"archive/zip"
	"github.com/Chain-Zhang/pinyin"
	"io"
	"log"
	"logViewerAgent/public/common"
	"logViewerAgent/public/net"
	"logViewerAgent/setting"
	"os"
)

func compress(files []*os.File, dest string) error {
	d, _ := os.Create(dest)
	defer d.Close()
	w := zip.NewWriter(d)
	defer w.Close()
	for _, file := range files {
		err := zipFile(file, "", w)
		if err != nil {
			return err
		}
	}
	return nil
}

func zipFile(file *os.File, prefix string, zw *zip.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = zipFile(f, prefix, zw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func ZipFile(path string, svcName string, packageTime string) (string, error) {

	// 准备要提交给服务端创建项目的信息
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

	dest := "storage/【" + projectNamePinyin + "】-【" + svcName + "】-【" + setting.Conf.Env + "】-【" + packageTime + "】-【" + localIP + "】.zip"

	dir, err := os.Open(path)
	if err != nil {
		log.Fatal(common.FuncName(), "打开文件失败：", path, "", err)
		return "", err
	}
	defer dir.Close()

	var files = []*os.File{dir}
	log.Println("开始压缩")
	err = compress(files, dest)
	if err != nil {
		log.Fatal(common.FuncName(), "压缩文件失败：", err)
		return "", err
	}
	log.Println("压缩完成")
	return dest, nil
}

//
//func DeCompress(zipFile, dest string) error {
//	reader, err := zip.OpenReader(zipFile)
//	if err != nil {
//		return err
//	}
//	defer reader.Close()
//	for _, logflie := range reader.File {
//		rc, err := logflie.Open()
//		if err != nil {
//			return err
//		}
//		defer rc.Close()
//		filename := dest + logflie.Name
//		err = os.MkdirAll(getDir(filename), 0755)
//		if err != nil {
//			return err
//		}
//		w, err := os.Create(filename)
//		if err != nil {
//			return err
//		}
//		defer w.Close()
//		_, err = io.Copy(w, rc)
//		if err != nil {
//			return err
//		}
//		w.Close()
//		rc.Close()
//	}
//	return nil
//}

//func getDir(path string) string {
//	return subString(path, 0, strings.LastIndex(path, "/"))
//}
//
//func subString(str string, start, end int) string {
//	rs := []rune(str)
//	length := len(rs)
//
//	if start < 0 || start > length {
//		panic("start is wrong")
//	}
//
//	if end < start || end > length {
//		panic("end is wrong")
//	}
//
//	return string(rs[start:end])
//}
