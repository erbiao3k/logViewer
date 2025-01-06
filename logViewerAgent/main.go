package main

import (
	"bufio"
	"fmt"
	"log"
	"logViewerAgent/public/common"
	"logViewerAgent/public/file"
	"logViewerAgent/setting"
	file2 "logViewerAgent/todo/logflie"
	"logViewerAgent/todo/project"
	"logViewerAgent/todo/status"
	"os"
	"strings"
	"time"
)

// cleanLogInErrorRuntime 当运行出错时，清空日志内容
func cleanLogInErrorRuntime(storageFile []string) {
	log.Println("正在清理复制的日志内容")
	for _, sf := range storageFile {
		f, err := os.OpenFile(sf, os.O_WRONLY|os.O_TRUNC, 0777)

		defer f.Close()

		if err != nil {
			log.Panicf("文件打开失败，err：%s", err)
		}

		writer := bufio.NewWriter(f)
		writer.WriteString(" ")
		writer.Flush()

	}
}

func init() {
	// 必须指定配置文件
	if len(os.Args) < 2 {
		log.Fatalf("用法：%s conf/config.ini", os.Args[0])
	}

	fmt.Println("\n+-+-+-加载配置文件+-+-+-+-+-+-+-+-加载配置文件+-+-+-+-+-+-+-+-+-")
	if err := setting.Init(os.Args[1]); err != nil {
		log.Fatal("配置文件加载失败, err：", err.Error())
	}

	if !common.IsValueInSlice(setting.Conf.Area, common.AllArea) {
		log.Fatalf("项目所在区域配置错误，仅支持：%s", common.AllArea)
	}

	if !common.IsValueInSlice(setting.Conf.Env, common.AllEnv) {
		log.Fatalf("项目环境配置错误，仅支持：%s", common.AllEnv)
	}

	proxyAddr := setting.Conf.ProxyUrl
	if len(proxyAddr) != 0 {
		log.Println("使用代理地址：", proxyAddr)
	}

	fmt.Println("\n+-+-+-系统权限测试+-+-+-+-+-+-+-+-系统权限测试+-+-+-+-+-+-+-+-+-")

	f, err := os.OpenFile("pem_text", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)

	if err != nil {
		log.Panicf("文件打开失败，err：%s", err)
	}

	_, err = f.Write([]byte("content" + "\n"))
	if err != nil {
		log.Fatalf("当前用户无法执行写入权限，无法部署日志提取工具：%s", err)
	}

	f.Close()

	err = os.Remove(f.Name())
	if err != nil {
		log.Fatalf("当前用户无法执行删除动作，无法部署日志提取工具%s", err)
	}
}

func projectInit() {
	for {
		fmt.Println("\n+-+-+-获取当前服务器运行的服务+-+-+-+-+-+-+-+-获取当前服务器运行的服务+-+-+-+-+-+-+-+-+-")

		var svcList []string
		conf := setting.Conf
		for _, path := range []string{conf.WebPath, conf.FosPath, conf.BasePath, conf.Op2ContainerLogPath, conf.Op2ConvertPath} {
			svcPath, err := status.GetSvcList(path)
			log.Printf("配置文件指定的路径：%s下，检测到服务目录：%s", path, svcPath)
			if err != nil {
				log.Printf("获取当前服务器部署的服务列表异常，%s", err.Error())
				continue
			}
			svcList = append(svcList, svcPath...)
		}

		fmt.Println("\n+-+-+-提交当前服务器运行的服务信息至服务端+-+-+-+-+-+-+-+-提交当前服务器运行的服务信息至服务端+-+-+-+-+-+-+-+-+-")
		if err := project.CreateProject(svcList); err != nil {
			log.Printf("项目创建失败，err：%s，服务列表：%s", err, svcList)
			break
		}

		break
	}
	time.Sleep(5 * time.Minute)
}
func projectRunner() {
	for {

		fmt.Println("\n+-+-+-从服务端查询当前项目日志提取申请+-+-+-+-+-+-+-+-从服务端查询当前项目日志提取申请+-+-+-+-+-+-+-+-+-")
		statusList, err := status.LogStatus()
		if err != nil {
			log.Println(fmt.Sprintf("服务端查询日志提取申请失败：err：%s", err.Error()))
			break
		}

		for _, stat := range statusList {
			// 获取需要上传的文件，并将文件保存在storage/yyyy-mm-dd_hh-mm-ss下
			fmt.Println("\n+-+-+-正在收集并压缩日志+-+-+-+-+-+-+-+-+-正在收集并压缩日志+-+-+-+-+-+-+-+-")

			log.Println("日志提取申请内容：", stat)

			err := project.UpdateProjectLogCommitStatus(
				stat["SvcName"], stat["LogDate"],
				stat["CreateTime"], common.LogStatusPacking,
				stat["LocalIP"], "",
			)
			if err != nil {
				log.Println("收集日志时，更新服务端状态失败，err：", err)
				break
			}

			logDir, createTime, storageFile, err := file2.CollectFile(stat)

			if err != nil {
				log.Println("日志文件收集异常：", err)
				cleanLogInErrorRuntime(storageFile)
			}
			// 未查询到日志时，不继续执行
			if logDir == "" {
				break
			}

			// 压缩日志所在目录
			log.Println("开始压缩日志：", stat)
			fileZip, err := file.Targz(logDir, stat["SvcName"], createTime)

			storageFile = append(storageFile, fileZip)

			if err != nil {
				cleanLogInErrorRuntime(storageFile)

				err := project.UpdateProjectLogCommitStatus(
					stat["SvcName"], stat["LogDate"],
					stat["CreateTime"], common.LogStatusPackFailed,
					stat["LocalIP"], "",
				)

				if err != nil {
					log.Println("压缩日志时，更新服务端日志状态失败，err：", err)
					break
				}

				log.Println(fmt.Sprintf("打包日志错误，err：%s", err))
				break
			}

			log.Println("已压缩文件：", fileZip)

			if err := project.UpdateProjectLogCommitStatus(
				stat["SvcName"], stat["LogDate"],
				stat["CreateTime"], common.LogStatusPacked,
				stat["LocalIP"], "",
			); err != nil {
				log.Println("压缩文件后，更新服务端状态失败：err：", err, ",压缩文件：", fileZip)
				cleanLogInErrorRuntime(storageFile)
				break
			}

			// 上传日志到服务器
			fmt.Println("\n+-+-+-上传日志中+-+-+-+-+-+-+-+-+-上传日志中+-+-+-+-+-+-+-+-")

			if err := project.UpdateProjectLogCommitStatus(
				stat["SvcName"], stat["LogDate"],
				stat["CreateTime"], common.LogStatusUploading,
				stat["LocalIP"], "",
			); err != nil {
				log.Println("上传日志中，更新服务端状态失败：err：", err, ",压缩文件：", fileZip)
				cleanLogInErrorRuntime(storageFile)
				break
			}

			log.Println("提交当前服务器公网IP至服务端")
			if err := project.AddWhiteList(); err != nil {
				log.Printf(fmt.Sprintf("提交当前服务器公网IP至服务端失败，err：%s", err.Error()))
				cleanLogInErrorRuntime(storageFile)
				break
			}

			err = file2.FileUpload(fileZip)

			if err != nil {
				cleanLogInErrorRuntime(storageFile)

				err := project.UpdateProjectLogCommitStatus(
					stat["SvcName"], stat["LogDate"],
					stat["CreateTime"], common.LogStatusUplodFailed,
					stat["LocalIP"], "",
				)

				if err != nil {
					log.Println("文件上传失败后，更新服务端状态失败，err：", err)
					break
				}

				log.Println(fmt.Sprintf("日志文件上传失败，日志文件包：%s，err：%s", fileZip, err))
				break
			}

			err = project.UpdateProjectLogCommitStatus(
				stat["SvcName"], stat["LogDate"],
				stat["CreateTime"], common.LogStatusUploaded,
				stat["LocalIP"], strings.Split(fileZip, "/")[1],
			)

			if err != nil {
				log.Println("日志文件上传成功后，更新服务端状态失败，err：", err)
				cleanLogInErrorRuntime(storageFile)
				break
			}

			log.Println("清理上传残留文件：", fileZip, " ", logDir)

			if err := os.Remove(fileZip); err != nil {
				log.Printf("清理上传残留文件%s失败，err：%s", fileZip, err)
				cleanLogInErrorRuntime(storageFile)
				break
			}

			if err := os.RemoveAll(logDir); err != nil {
				log.Printf("清理上传残留文件%s失败，err：%s", logDir, err)
				cleanLogInErrorRuntime(storageFile)
				break
			}

		}
		time.Sleep(5 * time.Minute)
	}
}

func main() {

	for {
		projectInit()
		projectRunner()
	}
}
