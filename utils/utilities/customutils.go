package utilities

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func CustomSubstring(originalString string, startIndex, endIndex int) string {
	if startIndex >= 0 && endIndex <= len(originalString) && startIndex <= endIndex {
		substring := originalString[startIndex:endIndex]
		return substring
	} else {
		return ""
	}
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func SaveImage(data []byte, targetPath string) error {
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}

	file, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	fmt.Printf("Image saved to: %s\n", targetPath)
	return nil
}

func GetExtensionFromBase64(base64String string) string {
	mimeType := strings.Split(base64String, ";")[0]
	extension := strings.Split(mimeType, "/")[1]
	return extension
}

func GetBaseURL(c *gin.Context) string {
	url := c.Request.URL.String()
	baseURL := strings.Split(url, "/")[0]
	return baseURL
}

func IsEmptyString(variableContent string) bool {
	if variableContent == "" {
		return true
	}
	return false
}

func ValidateBetweenTwoDateRange(startDateOrigin, endDateOrigin string) bool {
	var result int

	if startDateOrigin != "" || endDateOrigin != "" {
		startDate, err := time.Parse("2006-01-02", startDateOrigin)
		if err != nil {
			logrus.Error(err.Error())
		}

		endDate, err := time.Parse("2006-01-02", endDateOrigin)
		if err != nil {
			logrus.Error(err.Error())
		}

		diff := endDate.Sub(startDate)
		result = int(diff.Hours() / 24)

		if result < 0 {
			return false
		}

		if result >= 0 {
			return true
		}
	}

	return false
}

func GetLastDay(dateOrigin string) int {
	dateNow, _ := time.Parse("2006-01", dateOrigin)
	year, month, _ := dateNow.Date()

	// Calculate the last day of the month
	lastDay := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC)
	return lastDay.Day()
}

func CalculateSHA512(input string) (string, error) {
	hash := sha512.New()
	_, err := hash.Write([]byte(input))
	if err != nil {
		return "", err
	}
	hashInBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashInBytes)
	return hashString, nil
}

func SetNotifications(ctx *gin.Context, data NotificationEntities) (err error) {
	db := ctx.MustGet("db").(*gorm.DB)

	if err := db.Create(&data).Error; err != nil {
		return err
	}

	return nil
}

func GenerateRandomSixDigitNumber() int {
	rand.NewSource(time.Now().UnixNano())
	return rand.Intn(900000) + 100000
}

func HTMLContentReplacer(htmlContent string, variables map[string]string) string {
	for key, value := range variables {
		placeholder := fmt.Sprintf("{{ %s }}", key)
		htmlContent = strings.Replace(htmlContent, placeholder, value, -1)
	}
	return htmlContent
}

func CapitalizeWords(s string) string {
	words := strings.Fields(s) // Split the string into words
	var capitalizedWords []string
	for _, word := range words {
		capitalizedWords = append(capitalizedWords, cases.Title(language.Und).String(word))
	}
	return strings.Join(capitalizedWords, " ") // Join the words back with spaces
}

func RemoveEmptyStrings(arr []string) []string {
	var result []string

	for _, str := range arr {
		if str != "" {
			result = append(result, str)
		}
	}

	return result
}