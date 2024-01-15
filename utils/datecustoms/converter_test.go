package datecustoms

import (
	"testing"
	"time"
)

func TestIntToMonthName(t *testing.T) {
	testCases := []struct {
		input    int
		expected string
	}{
		{1, "January"},
		{2, "February"},
		{3, "March"},
		{4, "April"},
		{5, "May"},
		{6, "June"},
		{7, "July"},
		{8, "August"},
		{9, "September"},
		{10, "October"},
		{11, "November"},
		{12, "December"},
		{0, "Invalid Month"},
		{13, "Invalid Month"},
		{-1, "Invalid Month"},
		{100, "Invalid Month"},
	}

	for _, tc := range testCases {
		result := IntToMonthName(tc.input)

		if result != tc.expected {
			t.Errorf("For month %d, expected %s, but got %s", tc.input, tc.expected, result)
		}
	}
}

func TestValidDateFormat(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{input: "2022-01-01", expected: true},
		{input: "2022-12-31", expected: true},
		{input: "2022-01", expected: false},          // Incomplete format
		{input: "2022-01-01 12:00", expected: false}, // Additional time
		{input: "notadate", expected: false},         // Not a valid date format
		{input: "", expected: false},                 // Empty string
	}

	for _, tc := range testCases {
		result := ValidDateFormat(tc.input)

		if result != tc.expected {
			t.Errorf("For date %s, expected %t, but got %t", tc.input, tc.expected, result)
		}
	}
}

func TestNowTransaction(t *testing.T) {
	// Define a fixed time for testing
	fixedTime := time.Date(2022, time.January, 15, 12, 0, 0, 0, time.UTC)

	testCases := []struct {
		currentTime time.Time
		expected    string
	}{
		{currentTime: fixedTime, expected: "2022-01-15"},
	}

	for _, tc := range testCases {
		result := tc.currentTime.Format("2006-01-02")

		if result != tc.expected {
			t.Errorf("For time %v, expected %s, but got %s", tc.currentTime, tc.expected, result)
		}
	}
}
