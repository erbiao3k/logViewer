package mysql

import (
	"fmt"
	"logViewerServer/dao"
	"logViewerServer/pubilc"
	"strings"
)

// ProjectInfo 项目信息
type ProjectInfo struct {
	Md5Data             string `json:"md5_data"`
	ProjectName         string `json:"project"`
	ProjectEnv          string `json:"env"`
	ProjectArea         string `json:"area"`
	SvcAdmin            string `json:"admin"`
	SvcPortal           string `json:"portal"`
	SvcApi              string `json:"api"`
	SvcSchedule         string `json:"schedule"`
	SvcFos              string `json:"fos"`
	SvcConvert          string `json:"convert"`
	SvcSign             string `json:"sign"`
	Svc2BaseServer      string `json:"op-base"`
	Svc2ContractService string `json:"op2-contract"`
	Svc2SealService     string `json:"op2-seal"`
	Svc2UserService     string `json:"op2-user"`
	Svc2ScheduleServer  string `json:"op2-schedule"`
	Svc2SignServer      string `json:"op2-sign"`
	Svc2Convert         string `json:"op2-convert"`
	Svc2ApiGateway      string `json:"op2-api-gateway"`
	Svc2Web             string `json:"op2-web"`
	*dao.GormModel
}

// GetAllProjectInfo 获取项目信息
func GetAllProjectInfo() (ProjectInfo []ProjectInfo, err error) {
	if err = dao.DB.Find(&ProjectInfo).Error; err != nil {
		return nil, err
	}
	return
}

// FirstOrCreateProjectInfo 判断项目是否已创建,不存在则创建项目，存在则更新项目(强制)
func FirstOrCreateProjectInfo(field string, info *ProjectInfo) (err error) {
	infoCopy := new(ProjectInfo)
	if err = dao.DB.Where(&ProjectInfo{Md5Data: field}).Assign(info).FirstOrCreate(&infoCopy).Error; err != nil {
		return err
	}
	return
}

// GetStructProjectInfo 通过结构体数据以及服务名称查询项目服务运行的服务器IP
func GetStructProjectInfo(info *ProjectInfo, svcName string) (svcAddr []string, err error) {
	var ProjectInfo []ProjectInfo
	sprintfFormat := "svc_%s"

	if pubilc.IsValueInSlice(svcName, pubilc.Op2Service) {
		sprintfFormat = "svc2_%s"
	}

	svc := strings.ReplaceAll(svcName, "-", "_")

	if err = dao.DB.Debug().Where(&info).Where(fmt.Sprintf("%s <> \"\"", fmt.Sprintf(sprintfFormat, svc))).Find(&ProjectInfo).Error; err != nil {
		return
	}

	for _, addr := range ProjectInfo {
		switch svcName {
		case pubilc.SvcFos:
			svcAddr = append(svcAddr, addr.SvcFos)
		case pubilc.SvcApi:
			svcAddr = append(svcAddr, addr.SvcApi)
		case pubilc.SvcPortal:
			svcAddr = append(svcAddr, addr.SvcPortal)
		case pubilc.SvcSchedule:
			svcAddr = append(svcAddr, addr.SvcSchedule)
		case pubilc.SvcAdmin:
			svcAddr = append(svcAddr, addr.SvcAdmin)
		case pubilc.SvcConvert:
			svcAddr = append(svcAddr, addr.SvcConvert)
		case pubilc.SvcSign:
			svcAddr = append(svcAddr, addr.SvcSign)
		case pubilc.Svc2User:
			svcAddr = append(svcAddr, addr.Svc2UserService)
		case pubilc.Svc2Convert:
			svcAddr = append(svcAddr, addr.Svc2Convert)
		case pubilc.Svc2SignServer:
			svcAddr = append(svcAddr, addr.Svc2SignServer)
		case pubilc.Svc2ScheduleServer:
			svcAddr = append(svcAddr, addr.Svc2ScheduleServer)
		case pubilc.Svc2Contract:
			svcAddr = append(svcAddr, addr.Svc2ContractService)
		case pubilc.Svc2Seal:
			svcAddr = append(svcAddr, addr.Svc2SealService)
		case pubilc.Svc2BaseServer:
			svcAddr = append(svcAddr, addr.Svc2BaseServer)
		case pubilc.Svc2ApiGateway:
			svcAddr = append(svcAddr, addr.Svc2ApiGateway)
		case pubilc.Svc2Web:
			svcAddr = append(svcAddr, addr.Svc2Web)
		}
	}
	return

}
