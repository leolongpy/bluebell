package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"time"
)

var mySecret = []byte("leo")

type MyClaims struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}

// GetToken 生成JWT
func GetToken(userID int64) (string, error) {
	c := MyClaims{
		userID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(
				time.Duration(viper.GetInt("auth.jwt_expire")) * time.Hour).Unix(),

			Issuer: "bluebell",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, c)
	return token.SignedString(mySecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	var mc = new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (interface{}, error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid {
		return mc, nil
	}

	return nil, errors.New("invalid token")
}
