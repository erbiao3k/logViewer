package redis

import (
	"logViewerServer/dao"
	"time"
)

// StringSet 写入string
func StringSet(key, value string, expr time.Duration) error {
	r := dao.Rdb()
	defer r.Close()
	err := r.Set(dao.Ctx, key, value, expr).Err()
	if err != nil {
		return err
	}
	return nil
}

// StringGet 查询string
func StringGet(key string) (string, error) {
	r := dao.Rdb()
	defer r.Close()
	value, err := r.Get(dao.Ctx, key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}
