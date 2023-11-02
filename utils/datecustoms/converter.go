package datecustoms

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