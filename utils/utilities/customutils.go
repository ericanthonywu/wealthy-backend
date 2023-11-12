package utilities

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
	"strings"
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
	defer file.Close()

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