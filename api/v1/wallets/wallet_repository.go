package wallets

import (
	"errors"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/wallets/entities"
	"gorm.io/gorm"
)

type (
	WalletRepository struct {
		db *gorm.DB
	}

	IWalletRepository interface {
		PersonalAccount(email string) (data entities.WalletPersonalInformationEntity)
		Add(model *entities.WalletEntity) (err error)
	}
)

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) PersonalAccount(email string) (data entities.WalletPersonalInformationEntity) {
	r.db.Raw(" SELECT pa.id,pa.id_master_account_types, mat.account_type, (SELECT COUNT(id) FROM tbl_wallets WHERE id_account = pa.id) as total_wallet "+
		"FROM tbl_personal_accounts pa "+
		"INNER JOIN tbl_master_account_types mat ON mat.id = pa.id_master_account_types "+
		"WHERE pa.email = ?", email).Scan(&data)
	return data
}

func (r *WalletRepository) Add(model *entities.WalletEntity) (err error) {
	result := r.db.Create(&model)
	if result.RowsAffected == 0 {
		return errors.New("can not add wallet")
	}
	return nil
}
