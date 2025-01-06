package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"logViewerServer/models/mysql"
	"logViewerServer/models/redis"
	"logViewerServer/pubilc"
	"logViewerServer/setting"
	"logViewerServer/webhook"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// RegisterHandler 页面-注册
func RegisterHandler(c *gin.Context) {
	switch c.Request.Method {
	case http.MethodGet:
		c.HTML(http.StatusOK, "register.html", nil)
	case http.MethodPost:
		pm := mysql.PmInfo{
			PmArea:   c.PostForm("area"),
			PmEmail:  c.PostForm("email"),
			PmPhone:  c.PostForm("phone"),
			PmPasswd: c.PostForm("password"),
			Enable:   false,
		}

		if pubilc.IsValueInSlice("", []string{pm.PmEmail, pm.PmArea, pm.PmPhone, pm.PmPasswd}) {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "errmsg": "注册信息不完整"})
			return
		}

		if !strings.HasSuffix(pm.PmEmail, "@myemal.com") {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "errmsg": "请使用myemal.com邮箱登陆系统"})
		}

		// 密码复杂度判断
		if err := pubilc.CheckPasswordLever(pm.PmPasswd); err != nil {
			log.Println("err：", err)
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "errmsg": err.Error()})
			return
		}

		phone, err := mysql.GetWherePmInfo("pm_phone", pm.PmPhone)
		if err != nil {
			log.Println("确认手机号是否已注册时失败，err：", err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "errmsg": "确认手机号是否已注册时失败"})
			return
		}

		email, err := mysql.GetWherePmInfo("pm_email", pm.PmEmail)
		if err != nil {
			log.Println("确认邮箱是否已注册时失败，err：", err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "errmsg": "确认邮箱是否已注册时失败"})
			return
		}

		if len(phone) != 0 || len(email) != 0 {
			log.Printf("请确认手机号【%s】或邮箱【%s】是否已注册", pm.PmPhone, pm.PmEmail)
			c.JSON(http.StatusOK, gin.H{"status": false, "errmsg": fmt.Sprintf("请确认手机号【%s】或邮箱【%s】是否已注册", pm.PmPhone, pm.PmEmail)})
			return
		}

		if err := mysql.CreatePmInfo(&pm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "errmsg": "用户注册失败", "data": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": true, "errmsg": "账号注册成功"})

		// 用户注册成功后立马通知审核人
		info := url.Values{"pm_email": {pm.PmEmail}, "pm_area": {pm.PmArea}, "pm_phone": {pm.PmPhone}}.Encode()
		webhook.NotifyRegisterReview(pm.PmArea, pm.PmEmail, setting.Conf.DownloadAddr+"/ae"+"?"+info)

	default:
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"msg": "MethodNotAllowed",
		})
	}

}

// AccountEnableHandler 页面-账号激活
func AccountEnableHandler(c *gin.Context) {
	pmEmail, emailErr := c.GetQuery("pm_email")
	pmPhone, phoneErr := c.GetQuery("pm_phone")
	pmArea, areaErr := c.GetQuery("pm_area")

	if !emailErr || !areaErr || !phoneErr {
		log.Println("账号信息激活时，传入参数异常，err：", emailErr, areaErr, phoneErr)
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "errmsg": "账号信息激活时，传入参数异常"})
		return
	}

	_, err := mysql.EnablePmInfo(pmEmail, pmPhone, pmArea)
	if err != nil {
		log.Println("账号信息激活异常，err：", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "errmsg": "账号信息激活异常"})
		return
	}

	if err := webhook.NotifyRegisterEnable(pmEmail, pmPhone); err != nil {
		log.Println("NotifyRegisterEnable,err：", err)
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(200, fmt.Sprintf(`<script>alert("账号户【%s】激活成功")</script><p font size="100"></p>`, pmEmail))
}

// IndexHandler 页面-首页
func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

// LoginHandler 页面-登陆
func LoginHandler(c *gin.Context) {
	switch c.Request.Method {
	case http.MethodGet:
		c.HTML(http.StatusOK, "login.html", nil)
	case http.MethodPost:
		userEmail := c.PostForm("username")
		inputPassword := c.PostForm("password")

		if !strings.HasSuffix(userEmail, "@myemal.com") {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "errmsg": "请使用myemal.com邮箱登陆系统"})
			return
		}
		// 密码复杂度判断
		if err := pubilc.CheckPasswordLever(inputPassword); err != nil {
			log.Println("err：", err)
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "errmsg": err.Error()})
			return
		}

		password, err := mysql.GetWherePmInfo("pm_email", userEmail)
		if err != nil {
			log.Println("登陆时，获取账号信息失败，err：", err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "errmsg": "登陆时，获取账号信息失败"})
			return
		}

		if !password[0].Enable {
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.String(http.StatusOK, fmt.Sprintf(`<script>alert("账号户%s未激活，请联系管理员激活账号")</script><p font size="100">账号户%s未激活，请联系管理员激活账号</p>`, userEmail, userEmail))
			c.Abort()
			return
		}

		if len(password) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"status": false, "errmsg": fmt.Sprintf("账号：%s不存在！！！", userEmail)})
			return
		}

		if inputPassword != password[0].PmPasswd {
			c.JSON(http.StatusUnauthorized, gin.H{"status": false, "errmsg": "密码错误！！！"})
			return
		}
		loginSession := pubilc.GenerateSession(inputPassword)
		c.SetCookie("loginName", userEmail, 86400*30, "", "", false, false)
		c.SetCookie("loginSession", loginSession, 86400*30, "", "", false, false)

		if userEmail == "adminEngineer@myemal.com" {
			c.SetCookie("loginSessionAdmin", loginSession, 86400*30, "", "", false, false)
			c.SetCookie("loginNameAdmin", userEmail, 86400*30, "", "", false, false)
		}

		if err = redis.StringSet(loginSession, userEmail, time.Hour*24*30); err != nil {
			log.Println("登录失故障，redis会话信息写入失败：", err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "errmsg": "登陆失败！！！！"})
		}

		c.Redirect(http.StatusMovedPermanently, "/logviewer/log")
	default:
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"msg": "MethodNotAllowed",
		})
	}
}

// ProjectListHandler 页面-返回所有项目信息
func ProjectListHandler(c *gin.Context) {

	var projectList map[string]map[string][]string
	projectList = make(map[string]map[string][]string)

	// 获取所有项目信息原始数据
	projectInfo, err := mysql.GetAllProjectInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "errmsg": "项目信息获取失败"})
	}

	// 按照projectList数据类型，整理项目信息
	for _, p := range projectInfo {
		pEnv := p.ProjectEnv
		pName := p.ProjectName
		_, ok := projectList[pName]

		if !ok {
			projectList[pName] = map[string][]string{pEnv: {}}
		}

		temp := func(svcName string, svcName2 string) {
			if len(svcName) != 0 && !pubilc.IsValueInSlice(svcName2, projectList[pName][pEnv]) {
				projectList[pName][pEnv] = append(projectList[pName][pEnv], svcName2)
			}
		}

		temp(p.SvcPortal, pubilc.SvcPortal)

		temp(p.SvcSchedule, pubilc.SvcSchedule)

		temp(p.SvcAdmin, pubilc.SvcAdmin)

		temp(p.SvcApi, pubilc.SvcApi)

		temp(p.SvcSign, pubilc.SvcSign)

		temp(p.SvcConvert, pubilc.SvcConvert)

		temp(p.SvcFos, pubilc.SvcFos)

		temp(p.Svc2BaseServer, pubilc.Svc2BaseServer)

		temp(p.Svc2ContractService, pubilc.Svc2Contract)

		temp(p.Svc2SealService, pubilc.Svc2Seal)

		temp(p.Svc2UserService, pubilc.Svc2User)

		temp(p.Svc2ScheduleServer, pubilc.Svc2ScheduleServer)

		temp(p.Svc2SignServer, pubilc.Svc2SignServer)

		temp(p.Svc2Convert, pubilc.Svc2Convert)

		temp(p.Svc2ApiGateway, pubilc.Svc2ApiGateway)

		temp(p.Svc2Web, pubilc.Svc2Web)
	}

	// 清理无效项目信息
	delete(projectList, "None")

	c.JSON(http.StatusOK, gin.H{"projectList": projectList})
}

// LogCommitHandler 页面-日志提取申请提交
func LogCommitHandler(c *gin.Context) {
	projectInfo, err := mysql.GetAllProjectInfo()
	if err != nil {
		log.Println("提交日志提取申请时，查询项目信息失败，err：", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "errmsg": "提交日志提取申请时，查询项目清单失败"})
		return
	}
	var projectList []string
	for _, p := range projectInfo {
		if !pubilc.IsValueInSlice(p.ProjectName, projectList) {
			projectList = append(projectList, p.ProjectName)
		}
	}

	switch c.Request.Method {
	case http.MethodGet:
		c.HTML(http.StatusOK, "logCommit.tmpl", nil)
		return

	case http.MethodPost:
		env := c.PostForm("env")
		service := c.PostForm("service")
		item := c.PostForm("project")
		date := c.PostForm("date")
		username, _ := c.Cookie("loginName")
		reqTime := pubilc.FormatTime()

		data := map[string]interface{}{"环境类型": env, "服务名称": service, "日志时间": date, "项目名称": item, "日志申请人": username}

		log.Println("提交信息：", data)

		if !pubilc.IsValueInSlice(env, pubilc.AllEnv) ||
			!pubilc.IsValueInSlice(service, pubilc.AllService) ||
			!pubilc.IsValueInSlice(item, projectList) {
			c.JSON(http.StatusOK, gin.H{"status": false, "errmsg": "信息填写错误", "提交信息": data})
			return
		}

		md5Data := pubilc.OnlyRequest(env + service + date + item + reqTime)
		pTimeList, err := mysql.GetWhereProjectLogCommit(env, service, date, item)

		if err != nil {
			log.Println("提交日志提取申请时，查询最近一次提交记录的时间失败，err：", err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "errmsg": "提交日志提取申请时，查询最近一次提交记录的时间失败"})
			return
		}

		pi := mysql.ProjectInfo{
			ProjectName: item,
			ProjectEnv:  env,
		}

		svcAddrList, err := mysql.GetStructProjectInfo(&pi, service)
		if err != nil {
			log.Println("提交日志提取申请时，获取服务运行节点IP失败，err：", err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "errmsg": "提交日志提取申请时，获取服务运行节点IP失败"})
			return
		}

		for _, addr := range svcAddrList {
			p := mysql.ProjectLogCommit{
				LogCommitMd5: md5Data,
				CreateTime:   reqTime,
				ProjectName:  item,
				ProjectEnv:   env,
				LogDate:      date,
				SvcName:      service,
				SvcAddr:      addr,
				PmEmail:      username,
				LogStatus:    pubilc.LogStatusSubmitted,
			}

			// 当库中不存在日志提取申请入库
			if len(pTimeList) == 0 {
				err := mysql.CreateProjectLogCommit(&p)
				if err != nil {
					log.Println("日志提取申请时，提取记录入库失败，err：", err)
					c.JSON(http.StatusInternalServerError, gin.H{"status": false, "errmsg": "日志提取申请时，提取记录入库失败"})
					return
				}
			} else {
				// 比较最近一次匹配md5Data的项目日志的提取时间，两次日志提取时间间隔不得小于5min
				if time.Now().Unix()-pTimeList[0].Unix() > 300 {
					err := mysql.CreateProjectLogCommit(&p)
					if err != nil {
						log.Println("日志提取申请时，提取记录入库失败2，err：", err)
						c.JSON(http.StatusInternalServerError, gin.H{"status": false, "errmsg": "日志提取申请时，提取记录入库失败2"})
						return
					}
				} else {
					c.JSON(http.StatusOK, gin.H{"errmsg": "5min内已提交过一次，请稍后再提交", "最近一次提交时间": pTimeList[0]})
					return
				}
			}
		}

		c.Redirect(http.StatusMovedPermanently, "/logviewer/log/status")

	default:
		c.JSON(http.StatusMethodNotAllowed, gin.H{"errmsg": "MethodNotAllowed"})
	}
}

// LogStatusHandler 页面-日志提取状态查询
func LogStatusHandler(c *gin.Context) {
	username, _ := c.Cookie("loginName")

	commitStatus, err := mysql.GetIncompleteProjectLogCommit("log_status", "IS NOT NULL", "")
	if err != nil {
		log.Println("日志提取状态查询时，查询未完成的记录失败，err：", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "errmsg": "日志提取状态查询时，查询未完成的记录失败"})
		return
	}

	var htmlCode string
	var logZipAddr string
	for _, i := range commitStatus {
		if username == i.PmEmail {

			if pubilc.IsValueInSlice(i.LogStatus, []string{pubilc.LogStatusUploaded, pubilc.LogStatusSent, pubilc.LogStatusSendFailed}) {
				logZipAddr = fmt.Sprintf(`<a href="%s/log/download/?filename=%s">下载链接</a>`, setting.Conf.DownloadAddr, i.LogDownloadAddr)
			} else {
				logZipAddr = pubilc.LogStatusWaiting
			}
			htmlCode += fmt.Sprintf(`<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>`,
				i.CreateTime, i.LogStatus, i.ProjectName, i.ProjectEnv, i.SvcName, i.LogDate, i.SvcAddr, logZipAddr, i.PmEmail)
		}
	}
	c.HTML(http.StatusOK, "logStatus.tmpl", gin.H{"htmlCode": template.HTML(htmlCode)})
}

// FileDownloadHandler 页面-文件下载
func FileDownloadHandler(c *gin.Context) {
	filename, err := c.GetQuery("filename")
	if !err {
		log.Println("文件下载错误，err：", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "errmsg": "文件下载错误"})
		return
	}

	filePath := setting.Conf.Path + "/" + filename
	_, errByOpenFile := os.Open(filePath)

	if errByOpenFile != nil {
		log.Println("要下载的文件未找到，err：", err)
		c.JSON(http.StatusNotFound, gin.H{"status": false, "errmsg": fmt.Sprintf("文件%s不存在", filePath)})
		return
	}
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.File(filePath)
	return
}

// FileUploadHandler 接口-文件上传
func FileUploadHandler(c *gin.Context) {
	file, _ := c.FormFile("file")

	// 上传文件至指定目录
	if err := c.SaveUploadedFile(file, setting.Conf.Path+"/"+file.Filename); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":    fmt.Sprintf("'%s'上传失败！！!", file.Filename),
			"errmsg": err.Error(),
		})
		return
	}
	c.String(http.StatusOK, fmt.Sprintf("'%s'上传成功！！!", file.Filename))

}

// SvcListHandler 接口-返回服务清单
func SvcListHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": true, "data": pubilc.AllService, "errmsg": ""})
}

// CreateProject 接口-创建项目
func CreateProject(c *gin.Context) {
	if err := c.Request.ParseForm(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "errmsg": "解析Post表单失败", "errContent": err})
		return
	}

	// 接收客户端上传的项目信息,存在项目md5值时，不重新创建项目，仅更新项目
	resp := c.Request.Form
	p := mysql.ProjectInfo{
		Md5Data:             resp.Get("Md5Data"),
		ProjectName:         resp.Get("ProjectName"),
		ProjectEnv:          resp.Get("ProjectEnv"),
		ProjectArea:         resp.Get("ProjectArea"),
		SvcAdmin:            resp.Get("SvcAdmin"),
		SvcApi:              resp.Get("SvcApi"),
		SvcPortal:           resp.Get("SvcPortal"),
		SvcSchedule:         resp.Get("SvcSchedule"),
		SvcSign:             resp.Get("SvcSign"),
		SvcConvert:          resp.Get("SvcConvert"),
		SvcFos:              resp.Get("SvcFos"),
		Svc2BaseServer:      resp.Get("SvcBase-server"),
		Svc2ContractService: resp.Get("SvcContract-service"),
		Svc2SealService:     resp.Get("SvcSeal-service"),
		Svc2UserService:     resp.Get("SvcUser-service"),
		Svc2ScheduleServer:  resp.Get("SvcSchedule-server"),
		Svc2SignServer:      resp.Get("SvcSign-server"),
		Svc2Convert:         resp.Get("Svc2Convert"),
		Svc2ApiGateway:      resp.Get("SvcApi-gateway"),
		Svc2Web:             resp.Get("SvcWeb"),
	}

	if p.Md5Data == "" || p.ProjectEnv == "" || p.ProjectArea == "" || p.ProjectName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "errmsg": "创建项目接口参数不完整", "Content": resp})
		return
	}

	// 项目存在则强制更新项目，不存在则创建项目
	if err := mysql.FirstOrCreateProjectInfo(p.Md5Data, &p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "errmsg": "创建项目失败", "errContent": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true, "msg": "项目创建成功",
		"data": map[string]string{
			"ProjectName": p.ProjectName,
			"ProjectEnv":  p.ProjectEnv,
			"ProjectArea": p.ProjectArea,
		}})
}

// CheckLogStatusHandler 接口-查询日志提取状态
func CheckLogStatusHandler(c *gin.Context) {
	if err := c.Request.ParseForm(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "errmsg": "解析Post表单失败", "errContent": err})
		return
	}

	// 接收客户端上传的项目信息
	resp := c.Request.Form
	projectName := resp.Get("ProjectName")
	projectEnv := resp.Get("ProjectEnv")
	localIP := resp.Get("LocalIP")

	if projectName == "" || projectEnv == "" || localIP == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "errmsg": "日志提取状态接口参数不完整", "Content": resp})
		return
	}
	// 查询项目日志提取信息
	commitStatus, err := mysql.GetIncompleteProjectLogCommit("log_status", "<>", pubilc.LogStatusSent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errmsg": err})
		return
	}

	var data []map[string]string
	for _, status := range commitStatus {
		if status.LogStatus != pubilc.LogStatusUploaded && status.ProjectName == projectName && status.ProjectEnv == projectEnv && status.SvcAddr == localIP {
			data = append(data, map[string]string{"SvcName": status.SvcName, "LogDate": status.LogDate, "CreateTime": status.CreateTime, "LocalIP": localIP})
		}
	}

	if len(data) == 0 {
		c.JSON(http.StatusOK, gin.H{"status": false, "msg": "无日志申请", "data": data})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true, "msg": "已匹配到提取申请", "data": data})
}

// UpdateLogStatusHandler 接口-更新日志提取状态
func UpdateLogStatusHandler(c *gin.Context) {

	if err := c.Request.ParseForm(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "errmsg": "解析Post表单失败", "errContent": err})
		return
	}

	// 接收客户端上传的项目信息,依据Md5Data更新信息
	resp := c.Request.Form
	p := mysql.ProjectLogCommit{
		LogCommitMd5:    resp.Get("LogCommitMd5"),
		LogStatus:       resp.Get("LogStatus"),
		CreateTime:      resp.Get("CreateTime"),
		LogDownloadAddr: resp.Get("LogDownloadAddr"),
		SvcAddr:         resp.Get("LocalIP"),
	}
	if p.LogCommitMd5 == "" || p.LogStatus == "" || p.CreateTime == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "errmsg": "更新项目接口参数不完整", "Content": resp})
		return
	}

	if _, err := mysql.UpdateWhereProjectLogCommit(p.LogCommitMd5, p.CreateTime, p.LogStatus, p.SvcAddr); err != nil {
		log.Println("Md5为：", p.LogCommitMd5, ",创建时间为：", p.CreateTime, "的提交记录状态更新失败，err：", err)
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "errmsg": "更新日志提交记录状态失败", "Content": resp})
		return
	}
	if p.LogDownloadAddr != "" {
		if _, err := mysql.UpdateProjectLogCommitLogDownloadAddr(p.LogCommitMd5, p.CreateTime, p.LogDownloadAddr, p.SvcAddr); err != nil {
			log.Println("Md5为：", p.LogCommitMd5, ",创建时间为：", p.CreateTime, "的日志下载地址回填失败，err：", err)
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "errmsg": "日志下载地址回填失败", "Content": resp})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": true, "errmsg": "更新日志提交记录状态成功", "Content": resp})
}

// PublicIp 接口-返回客户端公网IP
func PublicIp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": c.ClientIP()})
}

// AddWhiteList 接口-添加白名单
func AddWhiteList(c *gin.Context) {

	if err := c.Request.ParseForm(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "errmsg": "解析Post表单失败", "errContent": err})
		return
	}

	// 接收客户端上传的项目信息
	resp := c.Request.Form
	l := mysql.WhiteList{
		ProjectName: resp.Get("ProjectName"),
		ProjectEnv:  resp.Get("ProjectEnv"),
		PublicIp:    resp.Get("PublicIP"),
	}

	if l.PublicIp == "" || l.ProjectEnv == "" || l.ProjectName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "errmsg": "写入白名单记录参数不完整", "Content": resp})
		return
	}

	if err := mysql.CreateWhiteListRecord(l.PublicIp, &l); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "errmsg": "插入白名单记录失败", "errContent": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": true, "errmsg": "插入白名单记录成功"})
}
