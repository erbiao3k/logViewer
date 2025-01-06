package common

import (
	"strings"
	"time"
)

func FormatTime() string {
	now := time.Now().String()
	nowDate := strings.Split(now, " ")[0]
	nowTime := strings.Replace(strings.Split(strings.Split(now, " ")[1], ".")[0], ":", "", -1)
	return nowDate + "-" + nowTime
}

func FormatDate() string {
	now := time.Now().String()
	nowDate := strings.Split(now, " ")[0]
	return nowDate
}
