package wallets

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/wealthy-app/wealthy-backend/api/v1/wallets/entities"
	"gorm.io/gorm"
)

type (
	WalletRepository struct {
		db *gorm.DB
	}

	IWalletRepository interface {
		NewWallet(model *entities.WalletEntity) (httpCode int, err error)
		NumberOfWalletsByID(IDPersonal uuid.UUID) (totalWallet int64, httpCode int, err error)
	}
)

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{
		db: db,
	}
}

func (r *WalletRepository) NewWallet(model *entities.WalletEntity) (httpCode int, err error) {
	result := r.db.Create(&model)
	if result.RowsAffected == 0 {
		return http.StatusInternalServerError, errors.New("can not add wallet")
	}
	return http.StatusOK, nil
}

func (r *WalletRepository) NumberOfWalletsByID(IDPersonal uuid.UUID) (totalWallet int64, httpCode int, err error) {
	if err := r.db.Model(&entities.WalletEntity{}).Where("id_account=?", IDPersonal).Count(&totalWallet).Error; err != nil {
		return totalWallet, http.StatusInternalServerError, err
	}
	return totalWallet, http.StatusOK, nil
}