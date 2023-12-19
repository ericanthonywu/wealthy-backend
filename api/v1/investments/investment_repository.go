package investments

import (
	"github.com/google/uuid"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/investments/entities"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type (
	InvestmentRepository struct {
		db *gorm.DB
	}

	IInvestmentRepository interface {
		InvestmentTrx(IDPersonal uuid.UUID) (data []entities.InvestmentGainLossEntity, err error)
		TrxInfo(IDPersonal uuid.UUID) (data []entities.InvestmentTransaction, err error)
		GetTradingInfo(stockCode string) (data entities.InvestmentTreding, err error)
		GetBrokerInfo(IDMasterBroker uuid.UUID) (data entities.BrokerInfo, err error)
	}
)

func NewInvestmentRepository(db *gorm.DB) *InvestmentRepository {
	return &InvestmentRepository{db: db}
}

func (r *InvestmentRepository) TrxInfo(IDPersonal uuid.UUID) (data []entities.InvestmentTransaction, err error) {
	if err := r.db.Raw(``, IDPersonal).Scan(&data).Error; err != nil {
		return []entities.InvestmentTransaction{}, err
	}
	return data, nil
}

func (r *InvestmentRepository) GetTradingInfo(stockCode string) (data entities.InvestmentTreding, err error) {
	if err := r.db.Raw(`SELECT tmd.symbol, tmd.name, tmd.close::numeric FROM tbl_master_trading tmd WHERE tmd.symbol=?`, stockCode).
		Scan(&data).Error; err != nil {
		return entities.InvestmentTreding{}, err
	}
	return data, nil
}

func (r *InvestmentRepository) InvestmentTrx(IDPersonal uuid.UUID) (data []entities.InvestmentGainLossEntity, err error) {
	if err := r.db.Raw(`SELECT * FROM tbl_investment ti WHERE ti.id_personal_accounts=? ORDER BY ti.created_at DESC`, IDPersonal).
		Scan(&data).Error; err != nil {
		logrus.Error(err.Error())
		return []entities.InvestmentGainLossEntity{}, err
	}
	return data, nil
}

func (r *InvestmentRepository) GetBrokerInfo(IDMasterBroker uuid.UUID) (data entities.BrokerInfo, err error) {
	if err := r.db.Raw(`SELECT * FROM tbl_master_broker WHERE id=?`, IDMasterBroker).
		Scan(&data).Error; err != nil {
		return entities.BrokerInfo{}, err
	}
	return data, nil
}