package wallets

import "gorm.io/gorm"

type (
	WalletRepository struct {
		db *gorm.DB
	}

	IWalletRepository interface {
		PersonalAccount(email string)
		Add(email string)
	}
)

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) PersonalAccount(email string) {
	r.db.Raw("")
}

func (r *WalletRepository) Add(email string) {
	r.PersonalAccount(email)
}
