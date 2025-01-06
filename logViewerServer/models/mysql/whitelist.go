package mysql

import "logViewerServer/dao"

type WhiteList struct {
	ProjectName string `json:"project_name"`
	ProjectEnv  string `json:"project_env"`
	PublicIp    string `json:"public_ip"`
	*dao.GormModel
}

// GetAllWhiteList 获取表白名单所有记录
func GetAllWhiteList() (WhiteList []WhiteList, err error) {
	if err = dao.DB.Find(&WhiteList).Error; err != nil {
		return nil, err
	}
	return
}

// CreateWhiteListRecord 插入一条白名单记录
func CreateWhiteListRecord(ip string, record *WhiteList) (err error) {
	if err = dao.DB.Where("public_ip = ?", ip).FirstOrCreate(&record).Error; err != nil {
		return err
	}
	return
}
