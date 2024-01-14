package wallets

import (
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/wealthy-app/wealthy-backend/api/v1/wallets/entities"
	"gorm.io/gorm"
)

type (
	WalletRepository struct {
		db *gorm.DB
	}

	IWalletRepository interface {
		NewWallet(model *entities.WalletEntity) (err error)
		NumberOfWalletsByID(IDPersonal uuid.UUID) (totalWallet int64, err error)
		SetBalanceNonInvestment(trx *entities.WalletInitTransaction, trxDetail *entities.WalletInitTransactionDetail) (err error)
		SetBalanceInvestment(trx *entities.WalletInitTransactionInvestment) (err error)
	}
)

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{
		db: db,
	}
}

func (r *WalletRepository) NewWallet(model *entities.WalletEntity) (err error) {
	result := r.db.Create(&model)
	if result.RowsAffected == 0 {
		logrus.Error(errors.New("can not add wallet").Error())
		return errors.New("can not add wallet")
	}
	return nil
}

func (r *WalletRepository) NumberOfWalletsByID(IDPersonal uuid.UUID) (totalWallet int64, err error) {
	if err := r.db.Model(&entities.WalletEntity{}).
		Where("id_account=?", IDPersonal).
		Count(&totalWallet).Error; err != nil {
		logrus.Error(err.Error())
		return totalWallet, err
	}
	return totalWallet, nil
}

func (r *WalletRepository) SetBalanceNonInvestment(trx *entities.WalletInitTransaction, trxDetail *entities.WalletInitTransactionDetail) (err error) {
	err = r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&trx).Error; err != nil {
			logrus.Error(err.Error())
			return err
		}

		if err := tx.Create(&trxDetail).Error; err != nil {
			logrus.Error(err.Error())
			return err
		}
		return nil
	})

	return nil
}

func (r *WalletRepository) SetBalanceInvestment(trx *entities.WalletInitTransactionInvestment) (err error) {
	return r.db.Create(&trx).Error
}