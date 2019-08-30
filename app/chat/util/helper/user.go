package helper

import (
	"crypto/md5"
	"fmt"
)


// SetPassword return a md5 password
func SetPassword(password string, salt string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(password + salt)))
}