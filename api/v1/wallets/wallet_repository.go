package wallets

import (
	"errors"
	"github.com/google/uuid"
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
		List(IDPersonal uuid.UUID) (data []entities.WalletEntity, err error)
		UpdateAmount(IDWallet string, amount int64) (data []entities.WalletEntity, httpCode int, err error)
		InitTransaction(trx *entities.WalletInitTransaction, trxDetail *entities.WalletInitTransactionDetail) (err error)
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

func (r *WalletRepository) List(IDPersonal uuid.UUID) (data []entities.WalletEntity, err error) {
	if err := r.db.Where("id_account=?", IDPersonal).Find(&data).Error; err != nil {
		return []entities.WalletEntity{}, err
	}
	return data, nil
}

func (r *WalletRepository) UpdateAmount(IDWallet string, amount int64) (data []entities.WalletEntity, httpCode int, err error) {
	result := r.db.Table("tbl_wallets").Where("id = ?", IDWallet).Update("amount", amount).Scan(&data)

	if result.Error != nil || result.RowsAffected == 0 {
		return []entities.WalletEntity{}, http.StatusInternalServerError, result.Error
	}

	return data, http.StatusOK, nil
}

func (r *WalletRepository) InitTransaction(trx *entities.WalletInitTransaction, trxDetail *entities.WalletInitTransactionDetail) (err error) {
	err = r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&trx).Error; err != nil {
			return err
		}

		if err := tx.Create(&trxDetail).Error; err != nil {
			return err
		}
		return nil
	})

	return err
}