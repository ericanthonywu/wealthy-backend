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
		GetAllWallets(accountUUID uuid.UUID) (data []entities.WalletEntity, err error)
		GetBalanceInvestment(walletID uuid.UUID) (data entities.WalletInitTransactionInvestment, err error)
		GetBalanceNonInvestment(walletID uuid.UUID) (data entities.WalletInitTransaction, err error)
		GetWalletType(walletUUID uuid.UUID) (data entities.WalletEntity, err error)
		UpdateWalletInformation(UUIDWallet uuid.UUID, request map[string]interface{}) (err error)
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

func (r *WalletRepository) GetAllWallets(accountUUID uuid.UUID) (data []entities.WalletEntity, err error) {
	if err := r.db.Where("id_account=?", accountUUID).Find(&data).Error; err != nil {
		return []entities.WalletEntity{}, err
	}
	return data, nil
}

func (r *WalletRepository) GetBalanceInvestment(walletID uuid.UUID) (data entities.WalletInitTransactionInvestment, err error) {
	if err := r.db.Where("wallet_id = ?", walletID).
		Order("created_at desc").
		First(&data).Error; err != nil {
		logrus.Error(err.Error())
		return entities.WalletInitTransactionInvestment{}, err
	}
	return data, nil
}

func (r *WalletRepository) GetBalanceNonInvestment(walletID uuid.UUID) (data entities.WalletInitTransaction, err error) {
	if err := r.db.Where("id_wallets = ?", walletID).
		Order("created_at desc").
		First(&data).Error; err != nil {
		logrus.Error(err.Error())
		return entities.WalletInitTransaction{}, err
	}
	return data, nil
}

func (r *WalletRepository) GetWalletType(walletUUID uuid.UUID) (data entities.WalletEntity, err error) {
	if err := r.db.Where("id = ?", walletUUID).First(&data).Error; err != nil {
		logrus.Error(err.Error())
		return entities.WalletEntity{}, err
	}
	return data, nil
}

func (r *WalletRepository) UpdateWalletInformation(UUIDWallet uuid.UUID, request map[string]interface{}) (err error) {
	var model entities.WalletEntity

	// set ID
	model.ID = UUIDWallet

	if err := r.db.Model(&model).Updates(request).Error; err != nil {
		logrus.Error(err.Error())
		return err
	}
	return nil
}
