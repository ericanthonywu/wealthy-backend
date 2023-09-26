package wallets

import (
	"errors"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/wallets/entities"
	"gorm.io/gorm"
	"net/http"
)

type (
	WalletRepository struct {
		db *gorm.DB
	}

	IWalletRepository interface {
		PersonalAccount(email string) (data entities.WalletPersonalInformationEntity)
		Add(model *entities.WalletEntity) (err error)
		List(email string) (data []entities.WalletEntity, httpCode int, err error)
		UpdateAmount(IDWallet string, amount int64) (data []entities.WalletEntity, httpCode int, err error)
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

func (r *WalletRepository) List(email string) (data []entities.WalletEntity, httpCode int, err error) {
	personalAccountData := r.PersonalAccount(email)

	result := r.db.Where("id_account=?", personalAccountData.ID).Find(&data)
	if result.RowsAffected == 0 {
		return []entities.WalletEntity{}, http.StatusNotFound, errors.New("not found")
	}

	if result.Error != nil {
		return []entities.WalletEntity{}, http.StatusInternalServerError, err
	}
	return data, http.StatusOK, nil
}

func (r *WalletRepository) UpdateAmount(IDWallet string, amount int64) (data []entities.WalletEntity, httpCode int, err error) {
	result := r.db.Table("tbl_wallets").Where("id = ?", IDWallet).Update("amount", amount).Scan(&data)

	if result.Error != nil || result.RowsAffected == 0 {
		return []entities.WalletEntity{}, http.StatusInternalServerError, result.Error
	}

	return data, http.StatusOK, nil
}
