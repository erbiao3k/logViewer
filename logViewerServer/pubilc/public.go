package pubilc

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"
)

// GenerateSession session生成方法
func GenerateSession(str string) string {
	m := md5.New()
	io.WriteString(m, str)
	io.WriteString(m, Salt)
	data := hex.EncodeToString(m.Sum(nil))
	return strings.ToUpper(data)
}

// OnlyRequest 客户端请求唯一标识字符串生成
func OnlyRequest(str string) string {
	m := md5.New()
	io.WriteString(m, str)
	io.WriteString(m, Salt)
	data := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%x", m.Sum(nil))))
	return data
}

// IsValueInSlice 判断字符串切片中是否包含某个value
func IsValueInSlice(value string, slice []string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func FormatTime() string {
	now := time.Now().String()
	nowDate := strings.Split(now, " ")[0]
	nowTime := strings.Replace(strings.Split(strings.Split(now, " ")[1], ".")[0], ":", "", -1)
	return nowDate + "-" + nowTime
}

// CheckPasswordLever 密码复杂度检查
func CheckPasswordLever(ps string) error {
	if len(ps) < 8 {
		return fmt.Errorf("密码长度至少8位")
	}
	num := `[0-9]{1}`
	az := `[a-z]{1}`
	//A_Z := `[A-Z]{1}`
	symbol := `[!@#~$%^&*()+=|_]{1}`
	if b, err := regexp.MatchString(num, ps); !b || err != nil {
		return fmt.Errorf("密码中未包含数字 :%v", err)
	}
	if b, err := regexp.MatchString(az, ps); !b || err != nil {
		return fmt.Errorf("密码中未包含字母 :%v", err)
	}
	//if b, err := regexp.MatchString(A_Z, ps); !b || err != nil {
	//	return fmt.Errorf("password need A_Z :%v", err)
	//}
	if b, err := regexp.MatchString(symbol, ps); !b || err != nil {
		return fmt.Errorf("密码中未包含符号 :%v", err)
	}
	return nil
}
