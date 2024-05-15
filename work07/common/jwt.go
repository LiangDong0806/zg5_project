package common

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
)

const GAOWEIMING = "13243124764"

func SetJwtToken(secretKey string, iat, seconds int64, payload string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["payload"] = payload
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}

func GetJwtToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(GAOWEIMING), nil
	})
	if err != nil {
		return false, err
	}
	if token.Valid {
		return true, nil
	} else {
		return false, fmt.Errorf("errors ")
	}
}
