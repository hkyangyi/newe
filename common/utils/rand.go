package utils

import (
	"math/rand"
	"time"
)

// 获取随机数字字符串
func GetRandomString(n int) string {
	var letters = []rune("0123456789")
	s := make([]rune, n)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range s {
		s[i] = letters[r.Intn(len(letters))]
	}
	return string(s)
}
