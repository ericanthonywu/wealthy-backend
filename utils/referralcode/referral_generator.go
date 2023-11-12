package referralcode

import (
	"github.com/semicolon-indonesia/wealthy-backend/utils/utilities"
	"math/rand"
	"time"
)

func ReferralCodeGenerator(username string, startIndex, endIndex int) string {
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	rand.Seed(time.Now().UnixNano())
	code := make([]byte, 8)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return utilities.CustomSubstring(username, startIndex, endIndex) + "-" + string(code)
}