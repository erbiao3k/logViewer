package common

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
)

func Encrypt(str string) string {
	m := md5.New()
	io.WriteString(m, str)
	io.WriteString(m, Salt)
	data := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%x", m.Sum(nil))))
	return data
}
