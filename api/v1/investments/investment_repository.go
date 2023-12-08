package investments

import (
	"github.com/google/uuid"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/investments/entities"
	"gorm.io/gorm"
)

type (
	InvestmentRepository struct {
		db *gorm.DB
	}

	IInvestmentRepository interface {
		TrxInfo(IDPersonal uuid.UUID) (data []entities.InvestmentTransaction, err error)
		GetTradingInfo(stockCode string) (data entities.InvestmentTreding, err error)
	}
)

func NewInvestmentRepository(db *gorm.DB) *InvestmentRepository {
	return &InvestmentRepository{db: db}
}

func (r *InvestmentRepository) TrxInfo(IDPersonal uuid.UUID) (data []entities.InvestmentTransaction, err error) {
	if err := r.db.Raw(`SELECT tt.amount as price,
       ttd.lot,
       ttd.stock_code,
       tmb.broker_name,
       ttd.sellbuy,
       tw.fee_invest_buy  as fee_buy,
       tw.fee_invest_sell as fee_sell
FROM tbl_transactions tt
         INNER JOIN tbl_transaction_details ttd ON ttd.id_transactions = tt.id
         INNER JOIN tbl_master_broker tmb ON tmb.id = tt.id_master_broker
         INNER JOIN tbl_wallets tw ON tw.id = tt.id_wallets
WHERE tt.id_master_invest <> '00000000-0000-0000-0000-000000000000'
  AND tt.id_personal_account = ?
ORDER BY ttd.stock_code ASC`, IDPersonal).Scan(&data).Error; err != nil {
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