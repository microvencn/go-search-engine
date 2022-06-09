package utils

import (
	"crypto/md5"
	"fmt"
)

func Md5Encrypt(str string) string {
	ans := md5.Sum([]byte(str))
	return fmt.Sprintf("%x", ans)
}
