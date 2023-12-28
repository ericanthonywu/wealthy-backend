package datecustoms

import (
	"github.com/sirupsen/logrus"
	"regexp"
	"time"
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

func TimeRFC3339ToString(timeString time.Time) string {
	return timeString.Format("2006-01-02 15:04:05")
}

func NowTransaction() string {
	timeNow := time.Now()
	return timeNow.Format("2006-01-02")
}

func TotalDaysBetweenDate(dateOrigin string) int {
	// Parse the date string to a time object
	layout := "2006-01-02"
	date, err := time.Parse(layout, dateOrigin)
	if err != nil {
		logrus.Errorf(err.Error())
	}

	// Get the current date
	now := time.Now()

	truncatedNow := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	truncatedDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

	daysDiff := int(truncatedNow.Sub(truncatedDate).Hours() / 24)

	return daysDiff
}

func TotalDaysBetweenDateTimeFormat(dateOrigin time.Time) int {
	// Parse the date string to a time object
	dateOrigin.Format("2006-01-02")

	// Get the current date
	now := time.Now().UTC()

	truncatedNow := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	truncatedDate := time.Date(dateOrigin.Year(), dateOrigin.Month(), dateOrigin.Day(), 0, 0, 0, 0, time.UTC)

	daysDiff := int(truncatedNow.Sub(truncatedDate).Hours() / 24)

	return daysDiff
}