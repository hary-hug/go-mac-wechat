package helper

import (
	"math/rand"
	"time"
)

// GetRandomString return a random string
func  GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// Substr 字符串截取
func Substr(str string, len int) string {
	r := []rune(str)
	return string(r[:len])
}
