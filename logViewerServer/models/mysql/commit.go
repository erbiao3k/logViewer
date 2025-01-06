package mysql

import (
	"fmt"
	"logViewerServer/dao"
	"logViewerServer/pubilc"
	"time"
)

// ProjectLogCommit 项目日志申请记录
type ProjectLogCommit struct {
	LogCommitMd5    string `json:"log_commit_md5"`
	ProjectName     string `json:"project"`
	SvcName         string `json:"svc_name"`
	ProjectEnv      string `json:"project_env"`
	LogDate         string `json:"log_date"`
	PmEmail         string `json:"pm_email"`
	LogStatus       string `json:"log_status"`
	SvcAddr         string `json:"svc_addr"`
	CreateTime      string `json:"create_time"`
	LogDownloadAddr string `json:"log_download_addr"`
	*dao.GormModel
}

// UpdateWhereProjectLogCommit 根据"log_commit_md5"和"create_time"更新日志提取记录的当前状态
func UpdateWhereProjectLogCommit(logCommitMd5, CreateTime, updateValue, SvcAddr string) (ProjectLogCommit *ProjectLogCommit, err error) {

	if err = dao.DB.Model(&ProjectLogCommit).
		Where("log_commit_md5 = ? AND create_time = ? AND svc_addr = ?", logCommitMd5, CreateTime, SvcAddr).
		UpdateColumn("log_status", updateValue).Error; err != nil {
		return nil, err
	}
	return
}

// UpdateProjectLogCommitLogDownloadAddr 根据"log_commit_md5"和"create_time"更新日志提取记录的中文件的下载地址
func UpdateProjectLogCommitLogDownloadAddr(logCommitMd5, CreateTime, logDownloadAddr, SvcAddr string) (ProjectLogCommit *ProjectLogCommit, err error) {
	if err = dao.DB.Model(&ProjectLogCommit).
		Where("log_commit_md5 = ? AND create_time = ? AND svc_addr = ?", logCommitMd5, CreateTime, SvcAddr).
		UpdateColumn("log_download_addr", logDownloadAddr).Error; err != nil {
		return nil, err
	}

	return
}

// GetProjectLogCommitNotify 获取提交记录中已上传完成，但未推送消息的记录
func GetProjectLogCommitNotify() (ProjectLogCommit []ProjectLogCommit, err error) {
	if err = dao.DB.Where("log_status = ?", pubilc.LogStatusUploaded).Find(&ProjectLogCommit).Error; err != nil {
		return nil, err
	}
	return
}

// DelWhereProjectLogCommit 删除提交时间大于3天的记录、删除状态为"日志不存在"的提交记录
func DelWhereProjectLogCommit() (ProjectLogCommit []ProjectLogCommit, err error) {
	//currentTime := time.Now()
	//oldTime := currentTime.AddDate(0, 0, -31).Format("2022-01-11 14:03:37")
	//if err = dao.DB.Where("log_status = ?", middleware.LogStatusNotFound).Or("created_at < ?", oldTime).Delete(&ProjectLogCommit).Error; err != nil {
	if err = dao.DB.Where("log_status = ?", pubilc.LogStatusNotFound).Delete(&ProjectLogCommit).Error; err != nil {
		return nil, err
	}
	return
}

// UpdateProjectLogCommitLogAddrAfterNotify 根据"log_download_addr"、"pm_email"、SvcAddr在通知完成后更新提交记录状态
func UpdateProjectLogCommitLogAddrAfterNotify(LogDownloadAddr, PmEmail, LogStatus string) (ProjectLogCommit *ProjectLogCommit, err error) {
	if err = dao.DB.Model(&ProjectLogCommit).
		Where("log_download_addr = ? AND pm_email = ?", LogDownloadAddr, PmEmail).
		UpdateColumn("log_status", LogStatus).Error; err != nil {
		return nil, err
	}

	return
}

// GetIncompleteProjectLogCommit ProjectLogCommit结构体对应表的where语句
func GetIncompleteProjectLogCommit(field string, symbol string, value string) (ProjectLogCommit []ProjectLogCommit, err error) {
	var whereContext = fmt.Sprintf("%s %s", field, symbol)
	if value != "" {
		whereContext = fmt.Sprintf("%s %s \"%s\"", field, symbol, value)
	}
	if err = dao.DB.Limit(30).Order("created_at desc").Where(whereContext).Find(&ProjectLogCommit).Error; err != nil {
		return nil, err
	}
	return
}

// CreateProjectLogCommit 创建日志提交记录
func CreateProjectLogCommit(commit *ProjectLogCommit) (err error) {
	if err = dao.DB.Create(&commit).Error; err != nil {
		return err
	}
	return
}

// GetWhereProjectLogCommit 查询当日是否提交过日志提取申请
func GetWhereProjectLogCommit(projectEnv, svcName, logDate, projectName string) (logInfo []time.Time, err error) {
	var l []ProjectLogCommit
	if err = dao.DB.
		Where("project_env = ? AND svc_name = ? AND log_date = ? AND project_name = ?",
			projectEnv, svcName, logDate, projectName).
		Select("created_at").Last(&l).Error; err != nil {
		return nil, err
	}

	for _, i := range l {
		logInfo = append(logInfo, i.CreatedAt)
	}
	return
}
