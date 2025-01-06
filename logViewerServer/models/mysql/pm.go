package mysql

import (
	"fmt"
	"logViewerServer/dao"
)

// PmInfo 项目经理信息
type PmInfo struct {
	PmName   string `json:"pm_name"`
	PmArea   string `json:"area"`
	PmEmail  string `json:"email" gorm:"type:varchar(100);unique_index"`
	PmPhone  string `json:"phone"`
	PmPasswd string `json:"login_password"`
	Enable   bool   `json:"enable"`
	*dao.GormModel
}

// CreatePmInfo 创建pm信息
func CreatePmInfo(pmInfo *PmInfo) (err error) {
	if err = dao.DB.Create(&pmInfo).Error; err != nil {
		return err
	}
	return
}

// EnablePmInfo 激活账号
func EnablePmInfo(email, phone, area string) (PmInfo *PmInfo, err error) {
	if err := dao.DB.Model(&PmInfo).
		Where("pm_email = ? AND pm_phone = ? AND pm_area = ?", email, phone, area).
		UpdateColumn("enable", true).Error; err != nil {
		return nil, err
	}
	return
}

// GetWherePmInfo 查询账号信息
func GetWherePmInfo(whereField, whereValue string) (PmInfo []PmInfo, err error) {
	if err = dao.DB.Where(
		fmt.Sprintf("%s = ?", whereField), whereValue).First(&PmInfo).Error; err != nil {
		return nil, err
	}
	return
}
