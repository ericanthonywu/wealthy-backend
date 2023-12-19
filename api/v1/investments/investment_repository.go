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
		InvestmentTrx(IDPersonal uuid.UUID) (data []entities.InvestmentTransaction, err error)
		TrxInfo(IDPersonal uuid.UUID) (data []entities.InvestmentTransaction, err error)
		GetTradingInfo(stockCode string) (data entities.InvestmentTreding, err error)
		GetBrokerInfo(IDMasterBroker uuid.UUID) (data entities.BrokerInfo, err error)
		GetInvestmentDataHelper(IDPersonal uuid.UUID, stockCode string) (data entities.InvestmentDataHelper, err error)
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

func (r *InvestmentRepository) InvestmentTrx(IDPersonal uuid.UUID) (data []entities.InvestmentTransaction, err error) {
	if err := r.db.Raw(`SELECT tt.date_time_transaction as date_transaction, ttd.stock_code, tt.amount as price, ttd.lot, tmb.broker_name, tw.fee_invest_sell as fee_sell
FROM tbl_transactions tt
         INNER JOIN tbl_transaction_details ttd ON tt.id = ttd.id_transactions
         INNER JOIN tbl_master_broker tmb ON tt.id_master_broker = tmb.id
         INNER JOIN tbl_wallets tw ON tw.id = tt.id_wallets
WHERE tt.id_personal_account = ?
  AND tt.id_master_broker <> '00000000-0000-0000-0000-000000000000'
  AND tt.id_master_invest <> '00000000-0000-0000-0000-000000000000'
  AND ttd.sellbuy = 0
ORDER BY tt.date_time_transaction::DATE DESC`, IDPersonal).
		Scan(&data).Error; err != nil {
		logrus.Error(err.Error())
		return []entities.InvestmentTransaction{}, err
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

func (r *InvestmentRepository) GetInvestmentDataHelper(IDPersonal uuid.UUID, stockCode string) (data entities.InvestmentDataHelper, err error) {
	if err := r.db.Raw(`SELECT * FROM tbl_investment WHERE id_personal_accounts=? AND stock_code=?`, IDPersonal, stockCode).
		Scan(&data).Error; err != nil {
		return entities.InvestmentDataHelper{}, err
	}
	return data, nil
}