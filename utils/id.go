package utils

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

// GenerateSessionID 生成会话ID（时间戳+随机）
func GenerateSessionID() string {
	return time.Now().Format("20060102150405") + "_" + RandomString(10)
}

// RandomString 生成指定长度的URL安全随机字符串
func RandomString(length int) string {
	if length <= 0 {
		length = 8
	}
	// 使用随机字节再做URL安全编码，截断到需要长度
	byteLen := (length*6 + 7) / 8 // 近似换算
	b := make([]byte, byteLen)
	_, err := rand.Read(b)
	if err != nil {
		// 回退到时间种子（极端情况）
		return time.Now().Format("150405")
	}
	// URL安全，不含+
	s := base64.RawURLEncoding.EncodeToString(b)
	if len(s) < length {
		return s
	}
	return s[:length]
}
