package dtos

import "github.com/google/uuid"

type (
	TransactionExpenseRequest struct {
		Date            string    `json:"date"`
		Amount          int64     `json:"amount"`
		ExpenseCategory uuid.UUID `json:"expense_category"`
		WalletID        uuid.UUID `json:"wallet_id"`
		Repeat          bool      `json:"repeat"`
		Note            string    `json:"note"`
		ExpenseTypeID   uuid.UUID `json:"expense_type_id"`
	}
)
