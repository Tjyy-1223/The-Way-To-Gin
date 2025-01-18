package utils

import (
	"math/rand"
	"time"
)

// RandString 用于生成锁标识，防止任何客户端都能解锁
func RandString(len int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}
