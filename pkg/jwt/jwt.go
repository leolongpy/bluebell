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
func GetToken(userID int64) (aToken, rToken string, err error) {
	c := MyClaims{
		userID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(
				time.Duration(viper.GetInt("auth.jwt_expire")) * time.Hour).Unix(),

			Issuer: "bluebell",
		},
	}
	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(mySecret)
	rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		Issuer:    "bluebell",
	}).SignedString(mySecret)
	return

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

// RefreshToken 刷新AccessToken
func RefreshToken(aToken, rToken string) (newAToken, newRToken string, err error) {
	if _, err = jwt.Parse(rToken, func(token *jwt.Token) (interface{}, error) {
		return mySecret, nil
	}); err != nil {
		return
	}

	var mc = new(MyClaims)
	_, err = jwt.ParseWithClaims(aToken, mc, func(token *jwt.Token) (interface{}, error) {
		return mySecret, nil
	})
	v, _ := err.(*jwt.ValidationError)

	if v.Errors == jwt.ValidationErrorExpired {
		return GetToken(mc.UserID)
	}
	return

}
