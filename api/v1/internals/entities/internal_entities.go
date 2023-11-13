package entities

type (
	TransactionByNotes struct {
		MonthYear           string `gorm:"column:month_year" json:"month_year"`
		Quantity            int    `gorm:"column:quantity" json:"quantity"`
		TransactionNote     string `gorm:"column:transaction_note" json:"transaction_note"`
		TransactionCategory string `gorm:"column:transaction_category" json:"transaction_category"`
	}
)