package password

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func Generate(pass string) string {
	password := []byte(pass)

	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		logrus.Error(err.Error())
	}
	return string(hashedPassword)
}

func Compare(hashedPassword string, password []byte) bool {
	//return bcrypt.CompareHashAndPassword(hashedPassword, password)
	// Since we'll be getting the hashed password from the DB it
	// will be a utilities so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPassword)
	err := bcrypt.CompareHashAndPassword(byteHash, password)
	if err != nil {
		log.Println(err)
		return false
	}

	return true

}