package logflie

import (
	"log"
	"logViewerAgent/public/common"
	"logViewerAgent/public/file"
	"logViewerAgent/setting"
	"logViewerAgent/todo/project"
	"os"
	"path/filepath"
	"strings"
)

// CollectFile 搜索符合条件的日志文件，并复制到storage/yyyy-mm-dd_hh-mm-ss下
func CollectFile(status map[string]string) (string, string, []string, error) {
	svcName := status["SvcName"]
	logDate := status["LogDate"]
	createTime := status["CreateTime"]
	packageTime := common.FormatTime()

	conf := setting.Conf

	logPath := ""

	op19WebSvcList := []string{common.SvcAdmin, common.SvcApi, common.SvcPortal, common.SvcSchedule}
	op20ContainerSvcList := []string{common.Svc2BaseServer, common.Svc2User, common.Svc2Seal, common.Svc2Contract, common.Svc2SignServer, common.Svc2ScheduleServer, common.Svc2Convert, common.Svc2ApiGateway, common.Svc2Web}

	// 识别op1.9的服务目录
	if common.IsValueInSlice(svcName, op19WebSvcList) {
		logPath = conf.WebPath + "/" + svcName + "/log/"
	}

	if svcName == common.SvcFos {
		logPath = conf.FosPath + "/" + svcName + "/logs/"
	}

	if svcName == common.SvcConvert {
		logPath = conf.BasePath + "/" + "tomcat-convert-8080/logs/"
	}

	if svcName == common.SvcSign {
		logPath = conf.BasePath + "/" + "tomcat-sign-8887/logs/"
	}

	// 识别op2.0的服务目录
	if common.IsValueInSlice(svcName, op20ContainerSvcList) {
		logPath = conf.Op2ContainerLogPath + "/" + svcName + "/"
	}

	if svcName == common.Svc2SignServer {
		logPath = conf.Op2ContainerLogPath + "/base-doc-sign-op2/"
	}

	if svcName == common.Svc2Convert {
		logPath = conf.Op2ConvertPath + "/" + common.Svc2Convert + "/logs/"
	}

	// 获取配置文件中指定位置的所有文件和目录
	filePaths, err := filepath.Glob(filepath.Join(logPath, "*"))
	if err != nil {
		log.Printf(common.FuncName(), "获取配置文件指定位置目录失败, err：", err)
		return "", "", nil, err
	}

	// 获取目标日志文件清单：
	//		1、当天的日志不带日期，也可能带日期，带日期时日志格式为：api-2021-01-01.log
	//		2、非当日的日志带日期，日志格式为：api-2021-01-01.log
	var logList []string
	nowDate := common.FormatDate()

	//按照给定的时间字符串【格式：2021-01-01】去取日志，依据业务不同日志规范，可能取到当天的或者非当天的
	for _, fp := range filePaths {
		if strings.Contains(fp, logDate) {
			logList = append(logList, fp)
		}
	}

	// 依据业务不同日志规范，可能当天的日志不带日期，因此当时间字符串【格式：2021-01-01】为今天时，将不带日期的日志文件也加入打包日志清单
	if nowDate == logDate {
		for _, fp := range filePaths {
			if !common.RegexpDate(fp) {
				logList = append(logList, fp)
			}
		}
	}

	// 收集日志清单时，不存在报错或需要打包的日志清单不为空时，将日志文件集中到起来
	if err != nil || len(logList) == 0 {
		err := project.UpdateProjectLogCommitStatus(
			svcName, logDate, createTime, common.LogStatusNotFound,
			status["LocalIP"], "")
		if err != nil {
			log.Println(common.FuncName(), "收集日志时，日志不存在，", err)
			return "", "", nil, err
		}
		log.Println(common.FuncName(), svcName, "服务未查询到", logDate, "的日志,", "err：", err)
		return "", "", nil, err
	}

	log.Println(common.FuncName(), "匹配到的日志文件：", logList)

	logStorage := "storage/" + svcName + "-" + packageTime + "/"

	// 创建日志打包目录
	if err := os.MkdirAll(logStorage, 0777); err != nil {
		log.Println(common.FuncName(), "创建日志打包目录失败：", logStorage, "，err：", err)
		return "", "", nil, err
	}

	var logStorageFile []string

	//复制日志文件到打包目录
	for _, value := range logList {
		value = strings.ReplaceAll(value, `\`, "/")
		filename := strings.Split(value, "/")[len(strings.Split(value, "/"))-1]
		storageFile := logStorage + "/" + filename
		_, err := file.CopyFile(storageFile, value)
		if err != nil {
			log.Println(common.FuncName(), "复制文件错误：", err)
			return "", "", nil, err
		}
		logStorageFile = append(logStorageFile, storageFile)
	}
	return logStorage, createTime, logStorageFile, nil
}
