package datecustoms

import (
	"regexp"
)

func IntToMonthName(month int) string {
	months := []string{
		"January", "February", "March", "April",
		"May", "June", "July", "August",
		"September", "October", "November", "December",
	}

	if month >= 1 && month <= 12 {
		return months[month-1] // Adjust for 0-based indexing
	}

	return "Invalid Month"
}

func ValidDateFormat(dateOrigin string) bool {
	datePattern := `^\d{4}-\d{2}-\d{2}$`

	r := regexp.MustCompile(datePattern)

	if r.MatchString(dateOrigin) {
		return true
	} else {
		return false
	}
}