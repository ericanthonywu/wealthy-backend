package password

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func Generate(pass string) string {
	password := []byte(pass)

	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		logrus.Error(err.Error())
	}
	return string(hashedPassword)
}

func Compare(hashedPassword, password []byte) (err error) {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
