package helper

import (
	"github.com/dgrijalva/jwt-go"
	"go-chat/app/chat/conf"
	"time"
)

// define a custom claims
type Claims struct {
	Uid int
	jwt.StandardClaims
}

type Token struct {
	AccessToken  string  `json:"access_token"`
	ExpireAt     int64   `json:"expire_at"`
}

var jwtSecret []byte


func init()  {
	jwtSecret = []byte(conf.Cfg.Common.JwtSecret)
}

// creating a token using a custom claims type
func GenerateToken(uid int, timeout int) (res Token, err error) {

	// expire time for token
	expire := time.Now().Add(time.Duration(timeout) * time.Second)

	claims := Claims{
		uid,
		jwt.StandardClaims{
			ExpiresAt: expire.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecret)

	res.AccessToken = tokenString
	res.ExpireAt = expire.Unix()

	return

}

// parse token using a custom claims
// return claims if success
func ParseToken(token string) (*Claims, error) {

	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
