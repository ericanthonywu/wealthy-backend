package token

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func JWTBuilder(email, role string) (tokenString string, err error) {

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["email"] = email
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	appSecret := os.Getenv("APP_SECRET")
	tokenString, err = token.SignedString([]byte(appSecret))

	if err != nil {
		logrus.Error(err.Error())
		return "", err
	}
	return tokenString, nil
}
