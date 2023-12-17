package transactions

import (
	"github.com/google/uuid"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/transactions/entities"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type (
	TransactionRepository struct {
		db *gorm.DB
	}

	ITransactionRepository interface {
		Add(trx *entities.TransactionEntity, trxDetail *entities.TransactionDetailEntity) (err error)
		ExpenseDetailHistoryWithoutDate(IDPersonal uuid.UUID) (data []entities.TransactionDetailHistory)
		ExpenseDetailHistoryWithDate(IDPersonal uuid.UUID, startDate, endDate string) (data []entities.TransactionDetailHistory)
		ExpenseTotalHistoryWithoutDate(IDPersonal uuid.UUID) (data entities.TransactionExpenseTotalHistory)
		ExpenseTotalHistoryWithDate(IDPersonal uuid.UUID, startDate, endDate string) (data entities.TransactionExpenseTotalHistory)

		IncomeDetailHistoryWithoutData(IDPersonal uuid.UUID) (data []entities.TransactionDetailHistory)
		IncomeDetailHistoryWithData(IDPersonal uuid.UUID, startDate, endDate string) (data []entities.TransactionDetailHistory)
		IncomeTotalHistoryWithoutDate(IDPersonal uuid.UUID) (data entities.TransactionIncomeTotalHistory)
		IncomeTotalHistoryWithData(IDPersonal uuid.UUID, startDate, endDate string) (data entities.TransactionIncomeTotalHistory)

		TransferMoneyInTotalHistoryWithoutData(IDPersonal uuid.UUID) (data []entities.TransactionIncomeTotalHistory)
		TransferMoneyOutTotalHistoryWithData(IDPersonal uuid.UUID, startDate, endDate string) (data []entities.TransactionIncomeTotalHistory)
		TransferDetailWithoutData(IDPersonal uuid.UUID) (data []entities.TransactionDetailTransfer)
		TransferDetailWithData(IDPersonal uuid.UUID, startDate, endDate string) (data []entities.TransactionDetailTransfer)

		InvestDetailWithoutData(IDPersonal uuid.UUID) (data []entities.TransactionDetailInvest)
		InvestDetailWithData(IDPersonal uuid.UUID, startDate, endDate string) (data []entities.TransactionDetailInvest)
		InvestTotalHistoryWithoutDate(IDPersonal uuid.UUID) (data entities.TransactionInvestTotalHistory)
		InvestTotalHistoryWithData(IDPersonal uuid.UUID, startDate, endDate string) (data entities.TransactionInvestTotalHistory)

		TravelDetailWithData(IDPersonal uuid.UUID, idTravel uuid.UUID, startDate, endDate string) (data []entities.TransactionDetailTravel)
		TravelDetailWithoutData(IDPersonal, idTravel uuid.UUID) (data []entities.TransactionDetailTravel)

		IncomeSpendingMonthlyTotal(IDPersonal uuid.UUID, month, year string) (data entities.TransactionIncomeSpendingTotalMonthly)
		IncomeSpendingMonthlyDetail(IDPersonal uuid.UUID, month, year string) (data []entities.TransactionIncomeSpendingDetailMonthly)
		IncomeSpendingAnnuallyTotal(IDPersonal uuid.UUID, year string) (data entities.TransactionIncomeSpendingTotalAnnually)
		IncomeSpendingAnnuallyDetail(IDPersonal uuid.UUID, year string) (data []entities.TransactionIncomeSpendingDetailAnnually)

		InvestMonthlyTotal(IDPersonal uuid.UUID, month, year string) (data entities.TransactionInvestmentTotals)
		InvestMonthlyDetail(IDPersonal uuid.UUID, month, year string) (data []entities.TransactionInvestmentDetail)
		InvestAnnuallyTotal(IDPersonal uuid.UUID, year string) (data entities.TransactionInvestmentTotals)
		InvestAnnuallyDetail(IDPersonal uuid.UUID, year string) (data []entities.TransactionInvestmentDetail)
		ByNote(IDPersonal uuid.UUID, month, year string) (data []entities.TransactionByNotes)
		Suggestion(IDPersoalAccount uuid.UUID) (data []entities.TransactionSuggestionNotes, err error)

		WalletExist(IDWallet uuid.UUID) bool
		BudgetWithCurrency(IDTravel uuid.UUID) (data entities.TransactionWithCurrency, err error)

		GetTradingInfo(stockCode string) (data entities.InvestmentTreding, err error)
	}
)

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Add(trx *entities.TransactionEntity, trxDetail *entities.TransactionDetailEntity) (err error) {
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

func (r *TransactionRepository) ExpenseDetailHistoryWithoutDate(IDPersonal uuid.UUID) (data []entities.TransactionDetailHistory) {
	if err := r.db.Raw(`SELECT tt.date_time_transaction::text       as transaction_date,
       tmec.expense_types ::text            as transaction_category,
       COALESCE(SUM(tt.amount), 0)::numeric as transaction_amount,
       CASE
           WHEN td.note IS NOT NULL THEN td.note
           WHEN td.note IS NUll THEN ''
           END :: text                      as transaction_note
FROM tbl_transactions tt
         LEFT JOIN tbl_master_expense_categories tmec ON tmec.id = tt.id_master_expense_categories
         INNER JOIN tbl_master_transaction_types tmtt ON tt.id_master_transaction_types = tmtt.id
         LEFT JOIN tbl_transaction_details td ON td.id_transactions = tt.id
WHERE tmtt.type = 'EXPENSE'
  AND tt.id_personal_account = ?
  AND to_char(tt.date_time_transaction::DATE, 'MM') = EXTRACT(
        MONTH FROM current_timestamp)::text
GROUP BY transaction_date, transaction_category, note
ORDER BY transaction_date DESC`, IDPersonal).Scan(&data).Error; err != nil {
		return []entities.TransactionDetailHistory{}
	}
	return data
}

func (r *TransactionRepository) ExpenseDetailHistoryWithDate(IDPersonal uuid.UUID, startDate, endDate string) (data []entities.TransactionDetailHistory) {
	if err := r.db.Raw(`SELECT tt.date_time_transaction    as transaction_date,
             tmec.expense_types          as transaction_category,
             COALESCE(SUM(tt.amount), 0) as transaction_amount,
             CASE
                 WHEN td.note IS NULL THEN ''
                 WHEN td.note IS NOT NULL then td.note
                 END                     as transaction_note
      FROM tbl_transactions tt
               LEFT JOIN tbl_master_expense_categories tmec ON tmec.id = tt.id_master_expense_categories
               INNER JOIN tbl_master_transaction_types tmtt ON tt.id_master_transaction_types = tmtt.id
               LEFT JOIN tbl_transaction_details td ON td.id_transactions = tt.id
      WHERE tmtt.type = 'EXPENSE'
        AND tt.id_personal_account = ?
        AND tt.date_time_transaction BETWEEN ? AND ?
      GROUP BY transaction_date, transaction_category, note
      ORDER BY transaction_date DESC`, IDPersonal, startDate, endDate).Scan(&data).Error; err != nil {
		return []entities.TransactionDetailHistory{}
	}
	return data
}

func (r *TransactionRepository) ExpenseTotalHistoryWithoutDate(IDPersonal uuid.UUID) (data entities.TransactionExpenseTotalHistory) {
	if err := r.db.Raw(`SELECT COALESCE(SUM(tt.amount) FILTER ( WHERE  tt.id_master_expense_categories <> '00000000-0000-0000-0000-000000000000' ), 0) as total_expense
      FROM tbl_transactions tt
      WHERE tt.id_personal_account = ?
        AND to_char(tt.date_time_transaction::DATE, 'MM') = EXTRACT(
              MONTH FROM current_timestamp)::text;`, IDPersonal).Scan(&data).Error; err != nil {
		return entities.TransactionExpenseTotalHistory{}
	}

	return data
}

func (r *TransactionRepository) ExpenseTotalHistoryWithDate(IDPersonal uuid.UUID, startDate, endDate string) (data entities.TransactionExpenseTotalHistory) {
	if err := r.db.Raw(`SELECT COALESCE(SUM(tt.amount) FILTER ( WHERE  tt.id_master_expense_categories <> '00000000-0000-0000-0000-000000000000' ), 0) as total_expense
      FROM tbl_transactions tt
      WHERE tt.id_personal_account = ?
        AND tt.date_time_transaction BETWEEN ? AND ?`, IDPersonal, startDate, endDate).Scan(&data).Error; err != nil {
		return entities.TransactionExpenseTotalHistory{}
	}

	return data
}

func (r *TransactionRepository) IncomeDetailHistoryWithoutData(IDPersonal uuid.UUID) (data []entities.TransactionDetailHistory) {
	if err := r.db.Raw(`SELECT tt.date_time_transaction::text    as transaction_date,
       tmic.income_types::text          as transaction_category,
       COALESCE(SUM(tt.amount), 0)::numeric as transaction_amount,
       CASE
           WHEN td.note IS NOT NULL then td.note
           WHEN td.note IS NULL then ''
           END::text                     as transaction_note
FROM tbl_transactions tt
         LEFT JOIN tbl_master_income_categories tmic ON tmic.id = tt.id_master_income_categories
         INNER JOIN tbl_master_transaction_types tmtt ON tt.id_master_transaction_types = tmtt.id
         LEFT JOIN tbl_transaction_details td ON td.id_transactions = tt.id
WHERE tmtt.type = 'INCOME'
  AND tt.id_personal_account = ?
  AND to_char(tt.date_time_transaction::DATE, 'MM') = EXTRACT(
        MONTH FROM current_timestamp)::text
GROUP BY transaction_date, transaction_category, transaction_note
ORDER BY transaction_date DESC`, IDPersonal).Scan(&data).Error; err != nil {
		return []entities.TransactionDetailHistory{}
	}
	return data
}

func (r *TransactionRepository) IncomeDetailHistoryWithData(IDPersonal uuid.UUID, startDate, endDate string) (data []entities.TransactionDetailHistory) {
	if err := r.db.Raw(`SELECT tt.date_time_transaction::text    as transaction_date,
       tmic.income_types::text           as transaction_category,
       COALESCE(SUM(tt.amount), 0)::numeric as transaction_amount,
       CASE
           WHEN td.note IS NOT NULL then td.note
           WHEN td.note IS NULL then ''
           END::text                     as transaction_note
FROM tbl_transactions tt
         LEFT JOIN tbl_master_income_categories tmic ON tmic.id = tt.id_master_income_categories
         INNER JOIN tbl_master_transaction_types tmtt ON tt.id_master_transaction_types = tmtt.id
         LEFT JOIN tbl_transaction_details td ON td.id_transactions = tt.id
WHERE tmtt.type = 'INCOME'
  AND tt.id_personal_account = ?
  AND tt.date_time_transaction BETWEEN ? AND ?
GROUP BY transaction_date, transaction_category, transaction_note
ORDER BY transaction_date DESC`, IDPersonal, startDate, endDate).Scan(&data).Error; err != nil {
		return []entities.TransactionDetailHistory{}
	}
	return data
}

func (r *TransactionRepository) IncomeTotalHistoryWithoutDate(IDPersonal uuid.UUID) (data entities.TransactionIncomeTotalHistory) {
	if err := r.db.Raw(`SELECT COALESCE(SUM(tt.amount) FILTER ( WHERE tt.id_master_income_categories <> '00000000-0000-0000-0000-000000000000' ), 0) as total_income
FROM tbl_transactions tt
WHERE tt.id_personal_account = ?
  AND to_char(tt.date_time_transaction::DATE, 'MM') = EXTRACT(
        MONTH FROM current_timestamp)::text`, IDPersonal).Scan(&data).Error; err != nil {
		return entities.TransactionIncomeTotalHistory{}
	}
	return data
}

func (r *TransactionRepository) IncomeTotalHistoryWithData(IDPersonal uuid.UUID, startDate, endDate string) (data entities.TransactionIncomeTotalHistory) {
	if err := r.db.Raw(`SELECT COALESCE(SUM(tt.amount) FILTER ( WHERE tt.id_master_income_categories <> '00000000-0000-0000-0000-000000000000' ), 0) as total_income
      FROM tbl_transactions tt
      WHERE tt.id_personal_account = ?
        AND tt.date_time_transaction BETWEEN ? AND ?`, IDPersonal, startDate, endDate).Scan(&data).Error; err != nil {
		return entities.TransactionIncomeTotalHistory{}
	}

	return data
}

func (r *TransactionRepository) InvestTotalHistoryWithoutDate(IDPersonal uuid.UUID) (data entities.TransactionInvestTotalHistory) {
	if err := r.db.Raw(`SELECT COALESCE(SUM(tt.amount) FILTER ( WHERE tt.id_master_invest <> '00000000-0000-0000-0000-000000000000' ), 0) as total_invest
FROM tbl_transactions tt
WHERE tt.id_personal_account = ?
  AND to_char(tt.date_time_transaction::DATE, 'MM') = EXTRACT(
        MONTH FROM current_timestamp)::text`, IDPersonal).Scan(&data).Error; err != nil {
		return entities.TransactionInvestTotalHistory{}
	}
	return data
}

func (r *TransactionRepository) InvestTotalHistoryWithData(IDPersonal uuid.UUID, startDate, endDate string) (data entities.TransactionInvestTotalHistory) {
	if err := r.db.Raw(`SELECT COALESCE(SUM(tt.amount) FILTER ( WHERE tt.id_master_invest <> '00000000-0000-0000-0000-000000000000' ), 0) as total_invest
      FROM tbl_transactions tt
      WHERE tt.id_personal_account = ?
        AND tt.date_time_transaction BETWEEN ? AND ?`, IDPersonal, startDate, endDate).Scan(&data).Error; err != nil {
		return entities.TransactionInvestTotalHistory{}
	}

	return data
}

func (r *TransactionRepository) TransferMoneyInTotalHistoryWithoutData(IDPersonal uuid.UUID) (data []entities.TransactionIncomeTotalHistory) {
	return nil
}

func (r *TransactionRepository) TransferMoneyOutTotalHistoryWithData(IDPersonal uuid.UUID, startDate, endDate string) (data []entities.TransactionIncomeTotalHistory) {
	return nil
}

func (r *TransactionRepository) TransferDetailWithData(IDPersonal uuid.UUID, startDate, endDate string) (data []entities.TransactionDetailTransfer) {
	if err := r.db.Raw(`SELECT tt.date_time_transaction::text    as transaction_date,
             COALESCE(SUM(tt.amount), 0)::numeric as transaction_amount,
             CASE
                 WHEN td.note IS NOT NULL then td.note
                 WHEN td.note IS NULL then ''
                 END ::text                    as transaction_note,
             td.transfer_to ::text                     as transaction_destination,
             td.transfer_from ::text                  as transaction_source
      FROM tbl_transactions tt
               INNER JOIN tbl_master_transaction_types tmtt ON tt.id_master_transaction_types = tmtt.id
               LEFT JOIN tbl_transaction_details td ON td.id_transactions = tt.id
      WHERE tmtt.type = 'TRANSFER'
        AND tt.id_personal_account = ?
        AND tt.date_time_transaction BETWEEN ? AND ?
      GROUP BY transaction_date, transaction_note, td.transfer_to, td.transfer_from
      ORDER BY transaction_date DESC`, IDPersonal, startDate, endDate).Scan(&data).Error; err != nil {
		return []entities.TransactionDetailTransfer{}
	}

	return data
}

func (r *TransactionRepository) TransferDetailWithoutData(IDPersonal uuid.UUID) (data []entities.TransactionDetailTransfer) {
	if err := r.db.Raw(`SELECT tt.date_time_transaction::text    as transaction_date,
       COALESCE(SUM(tt.amount), 0)::numeric as transaction_amount,
       CASE
           WHEN td.note IS NOT NULL then td.note
           WHEN td.note IS NULL then ''
           END::text                    as transaction_note,
       td.transfer_to::text                     as transaction_destination,
       td.transfer_from::text                   as transaction_source
FROM tbl_transactions tt
         INNER JOIN tbl_master_transaction_types tmtt ON tt.id_master_transaction_types = tmtt.id
         LEFT JOIN tbl_transaction_details td ON td.id_transactions = tt.id
WHERE tmtt.type = 'TRANSFER'
  AND tt.id_personal_account = ?
  AND to_char(tt.date_time_transaction::DATE, 'MM') = EXTRACT(
        MONTH FROM current_timestamp)::text
GROUP BY transaction_date, transaction_note, td.transfer_to, td.transfer_from
ORDER BY transaction_date DESC`, IDPersonal).Scan(&data).Error; err != nil {
		return []entities.TransactionDetailTransfer{}
	}

	return data
}

func (r *TransactionRepository) InvestDetailWithoutData(IDPersonal uuid.UUID) (data []entities.TransactionDetailInvest) {
	if err := r.db.Raw(`SELECT tt.date_time_transaction::text as transaction_date,
       td.lot * tt.amount::numeric    as transaction_amount_total,
       tt.amount::numeric             as price,
       CASE
           WHEN td.note IS NULL THEN ''
           ELSE td.note
           END ::text                 as transaction_note,
       td.lot ::int                   as lot,
       td.stock_code::text            as stock_code,
       CASE
           WHEN td.sellbuy = 0 THEN 'SELL'
           WHEN td.sellbuy = 1 THEN 'BUY'
           ELSE ''
           END ::text                 as sell_buy
FROM tbl_transactions tt
         INNER JOIN tbl_transaction_details td ON td.id_transactions = tt.id
WHERE tt.id_personal_account = '?'
GROUP BY transaction_note, lot, stock_code, transaction_date, sell_buy, tt.amount
ORDER BY transaction_date DESC`, IDPersonal).Scan(&data).Error; err != nil {
		return []entities.TransactionDetailInvest{}
	}

	return data
}

func (r *TransactionRepository) InvestDetailWithData(IDPersonal uuid.UUID, startDate, endDate string) (data []entities.TransactionDetailInvest) {
	if err := r.db.Raw(`SELECT tt.date_time_transaction as transaction_date,
       td.lot * tt.amount       as transaction_amount_total,
       tt.amount                as price,
       CASE
           WHEN td.note IS NULL THEN ''
           ELSE td.note
           END                  as transaction_note,
       td.lot                   as lot,
       td.stock_code            as stock_code,
       CASE
           WHEN td.sellbuy = 0 THEN 'SELL'
           WHEN td.sellbuy = 1 THEN 'BUY'
           ELSE ''
           END                  as sell_buy
FROM tbl_transactions tt
         INNER JOIN tbl_master_transaction_types tmtt ON tmtt.id = tt.id_master_transaction_types
         INNER JOIN tbl_transaction_details td ON td.id_transactions = tt.id
WHERE tt.id_personal_account = ?
  AND tmtt.type = 'INVEST'
  AND tt.date_time_transaction BETWEEN ? AND ?
GROUP BY transaction_note, lot, stock_code, transaction_date, sell_buy, tt.amount
ORDER BY transaction_date DESC`, IDPersonal, startDate, endDate).Scan(&data).Error; err != nil {
		return []entities.TransactionDetailInvest{}
	}

	return data
}

func (r *TransactionRepository) TravelDetailWithoutData(IDPersonal, idTravel uuid.UUID) (data []entities.TransactionDetailTravel) {
	if err := r.db.Raw(`SELECT tt.date_time_transaction, tt.id as id_transaction,tt.amount, tmec.expense_types as category,td.note
FROM tbl_transactions tt
         INNER JOIN tbl_master_expense_categories tmec ON tmec.id = tt.id_master_expense_categories
         INNER JOIN tbl_master_transaction_types tmtt ON tt.id_master_transaction_types = tmtt.id
         INNER JOIN tbl_transaction_details td ON td.id_transactions = tt.id
WHERE tmtt.type = 'TRAVEL'
  AND tt.id_personal_account = ?
  AND td.id_travel = ?
ORDER BY tt.date_time_transaction::DATE ASC`, IDPersonal, idTravel).Scan(&data).Error; err != nil {
	}
	return
}

func (r *TransactionRepository) TravelDetailWithData(IDPersonal uuid.UUID, idTravel uuid.UUID, startDate, endDate string) (data []entities.TransactionDetailTravel) {
	if err := r.db.Raw(`SELECT tt.date_time_transaction, tt.id  as id_transaction, tt.amount,tmec.expense_types as category,td.note
FROM tbl_transactions tt
         INNER JOIN tbl_master_expense_categories tmec ON tmec.id = tt.id_master_expense_categories
         INNER JOIN tbl_master_transaction_types tmtt ON tt.id_master_transaction_types = tmtt.id
         INNER JOIN tbl_transaction_details td ON td.id_transactions = tt.id
WHERE tmtt.type = 'TRAVEL'
  AND tt.id_personal_account = ?
  AND td.id_travel = ?
  AND tt.date_time_transaction BETWEEN ? AND ?
ORDER BY tt.date_time_transaction::DATE ASC`, IDPersonal, idTravel, startDate, endDate).Scan(&data).Error; err != nil {

	}
	return
}

func (r *TransactionRepository) IncomeSpendingMonthlyTotal(IDPersonal uuid.UUID, month, year string) (data entities.TransactionIncomeSpendingTotalMonthly) {
	if err := r.db.Raw(`SELECT concat(to_char(to_date(t.date_time_transaction, 'YYYY-MM-DD'), 'Mon'), ' ',
              to_char(t.date_time_transaction::DATE, 'YYYY')) ::text                                  as month,
       to_char(t.date_time_transaction::DATE, 'YYYY')::text                                           as year,
       COALESCE(SUM(t.amount) FILTER ( WHERE t.id_master_income_categories <> '00000000-0000-0000-0000-000000000000'), 0)::numeric  as total_income,
       COALESCE(SUM(t.amount) FILTER ( WHERE t.id_master_expense_categories <> '00000000-0000-0000-0000-000000000000'),
                0)::numeric                                                                           as total_spending,
       COALESCE(SUM(t.amount) FILTER ( WHERE t.id_master_income_categories <> '00000000-0000-0000-0000-000000000000'), 0) -
       COALESCE(SUM(t.amount) FILTER ( WHERE t.id_master_expense_categories <> '00000000-0000-0000-0000-000000000000'), 0)::numeric as net_income
FROM tbl_transactions t
WHERE to_char(t.date_time_transaction::DATE, 'MM') = ?
  AND to_char(t.date_time_transaction::DATE, 'YYYY') = ?
  AND t.id_personal_account = ?
GROUP BY year, month`, month, year, IDPersonal).Scan(&data).Error; err != nil {
		return entities.TransactionIncomeSpendingTotalMonthly{}
	}
	return data
}

func (r *TransactionRepository) IncomeSpendingMonthlyDetail(IDPersonal uuid.UUID, month, year string) (data []entities.TransactionIncomeSpendingDetailMonthly) {
	if err := r.db.Raw(`SELECT concat(to_char(tt.date_time_transaction::DATE, 'YYYY'), '-',
              to_char(tt.date_time_transaction::DATE, 'MM'), '-',
              to_char(tt.date_time_transaction::DATE, 'DD'))::text as date,
       CASE
           WHEN tmec.expense_types IS NOT NULL THEN tmec.expense_types
           WHEN tmic.income_types IS NOT NULL THEN tmic.income_types
           END ::text                                                as transaction_category,
       CASE
           WHEN tmec.expense_types IS NOT NULL THEN 'EXPENSE'
           WHEN tmic.income_types IS NOT NULL THEN 'INCOME'
           END :: text                                               as transaction_type,
       coalesce(SUM(tt.amount), 0)::numeric                         as transaction_amount,
       CASE
           WHEN td.note IS NULL THEN ''
           WHEN td.note IS NOT NULL THEN td.note
           END::text                                            as transaction_note
FROM tbl_transactions tt
         LEFT JOIN tbl_master_expense_categories tmec
                   ON tt.id_master_expense_categories = tmec.id
         LEFT JOIN tbl_master_income_categories tmic
                   ON tt.id_master_income_categories = tmic.id
         LEFT JOIN tbl_transaction_details td ON tt.id = td.id_transactions
         INNER JOIN tbl_master_transaction_types tmtt ON tmtt.id = tt.id_master_transaction_types
WHERE to_char(tt.date_time_transaction::DATE, 'MM') = ?
  AND to_char(tt.date_time_transaction::DATE, 'YYYY') = ?
  AND tt.id_personal_account = ?
  AND (tmtt.type = 'INCOME' OR tmtt.type = 'EXPENSE')
GROUP BY date, transaction_category, tmec.expense_types, tmic.income_types, td.note
ORDER BY date DESC`, month, year, IDPersonal).Scan(&data).Error; err != nil {
		return []entities.TransactionIncomeSpendingDetailMonthly{}
	}

	return data
}

func (r *TransactionRepository) IncomeSpendingAnnuallyTotal(IDPersonal uuid.UUID, year string) (data entities.TransactionIncomeSpendingTotalAnnually) {
	if err := r.db.Raw(`SELECT to_char(t.date_time_transaction::DATE, 'YYYY') ::text  as transaction_period,
       COALESCE(SUM(t.amount) FILTER ( WHERE t.id_master_income_categories <> '00000000-0000-0000-0000-000000000000'), 0)::numeric  as total_income,
       COALESCE(SUM(t.amount) FILTER ( WHERE t.id_master_expense_categories <> '00000000-0000-0000-0000-000000000000'), 0)::numeric as total_spending,
       COALESCE(SUM(t.amount) FILTER ( WHERE t.id_master_income_categories <> '00000000-0000-0000-0000-000000000000'), 0) -
       COALESCE(SUM(t.amount) FILTER ( WHERE t.id_master_expense_categories <> '00000000-0000-0000-0000-000000000000'), 0)::numeric as net_income
FROM tbl_transactions t
WHERE to_char(t.date_time_transaction::DATE, 'YYYY') = ? AND t.id_personal_account = ?
GROUP BY transaction_period`, year, IDPersonal).Scan(&data).Error; err != nil {
		return entities.TransactionIncomeSpendingTotalAnnually{}
	}
	return data
}

func (r *TransactionRepository) IncomeSpendingAnnuallyDetail(IDPersonal uuid.UUID, year string) (data []entities.TransactionIncomeSpendingDetailAnnually) {
	if err := r.db.Raw(`SELECT CONCAT(to_char(tt.date_time_transaction::DATE, 'YYYY'), '-',
              to_char(tt.date_time_transaction::DATE, 'MM'))         as date_origin,
       to_char(tt.date_time_transaction::DATE, 'MM') ::text          as month,
       CONCAT(to_char(to_date(tt.date_time_transaction, 'YYYY-MM-DD'), 'Mon'), ' ',
              to_char(tt.date_time_transaction::DATE, 'YYYY'))::text as month_year,
       date_part('days', (date_trunc('month', tt.date_time_transaction::DATE) +
                          interval '1 month - 1 day')) ::numeric     as total_days_in_month,
       COALESCE(SUM(tt.amount)
                FILTER ( WHERE tt.id_master_income_categories <> '00000000-0000-0000-0000-000000000000' ),
                0) :: numeric                                        as total_income,
       COALESCE(SUM(tt.amount)
                FILTER ( WHERE tt.id_master_expense_categories <> '00000000-0000-0000-0000-000000000000' ),
                0) :: numeric                                        as total_spending,
       COALESCE(SUM(tt.amount)
                FILTER ( WHERE tt.id_master_income_categories <> '00000000-0000-0000-0000-000000000000' ),
                0) - COALESCE(SUM(tt.amount) FILTER ( WHERE tt.id_master_expense_categories <>
                                                            '00000000-0000-0000-0000-000000000000' ),
                              0) :: numeric                          as net_income
FROM tbl_transactions tt
         LEFT JOIN tbl_master_expense_categories tmec ON tt.id_master_expense_categories = tmec.id
         LEFT JOIN tbl_master_income_categories tmic
                   ON tt.id_master_income_categories = tmic.id
WHERE to_char(tt.date_time_transaction::DATE, 'YYYY') = ?
  AND tt.id_personal_account = ?
GROUP BY month_year, month, total_days_in_month, date_origin
ORDER BY month DESC`, year, IDPersonal).Scan(&data).Error; err != nil {
		return []entities.TransactionIncomeSpendingDetailAnnually{}
	}
	return data
}

func (r *TransactionRepository) InvestMonthlyTotal(IDPersonal uuid.UUID, month, year string) (data entities.TransactionInvestmentTotals) {
	if err := r.db.Raw(`SELECT COALESCE(SUM(t.amount) FILTER ( WHERE td.sellbuy = 1 ), 0)::numeric as total_buy,
       COALESCE(SUM(t.amount) FILTER ( WHERE td.sellbuy = 0 ), 0)::numeric as total_sell,
       COALESCE(SUM(t.amount) FILTER ( WHERE td.sellbuy = 0 ), 0) -
       COALESCE(SUM(t.amount) FILTER ( WHERE td.sellbuy = 1 ), 0)::numeric as total_current_portfolio
FROM tbl_transactions t
         INNER JOIN tbl_transaction_details td ON td.id_transactions = t.id
         INNER JOIN tbl_master_transaction_types tmtt ON tmtt.id = t.id_master_transaction_types
WHERE to_char(t.date_time_transaction::DATE, 'MM') = ?
  AND to_char(t.date_time_transaction::DATE, 'YYYY') = ?
  AND t.id_personal_account = ?
  AND tmtt.type = 'INVEST'`, month, year, IDPersonal).Scan(&data).Error; err != nil {
		return entities.TransactionInvestmentTotals{}
	}
	return data
}

func (r *TransactionRepository) InvestMonthlyDetail(IDPersonal uuid.UUID, month, year string) (data []entities.TransactionInvestmentDetail) {
	if err := r.db.Raw(`SELECT concat(to_char(t.date_time_transaction::DATE, 'DD'), ' ',
              to_char(to_date(t.date_time_transaction, 'YYYY-MM-DD'), 'Mon'), ' ',
              to_char(t.date_time_transaction::DATE, 'YYYY'))::text         as date,
       td.lot::numeric                                                      as lot,
       td.stock_code ::text                                                 as stock_code,
	   t.amount as price
FROM tbl_transactions t
         INNER JOIN tbl_transaction_details td ON td.id_transactions = t.id
         INNER JOIN tbl_master_transaction_types tmtt ON tmtt.id = t.id_master_transaction_types
WHERE to_char(t.date_time_transaction::DATE, 'MM') = ?
  AND to_char(t.date_time_transaction::DATE, 'YYYY') = ?
  AND t.id_personal_account = ?
  AND tmtt.type = 'INVEST'
GROUP BY date, lot, stock_code,price
ORDER BY date DESC`, month, year, IDPersonal).Scan(&data).Error; err != nil {
		return []entities.TransactionInvestmentDetail{}
	}
	return data
}

func (r *TransactionRepository) InvestAnnuallyTotal(IDPersonal uuid.UUID, year string) (data entities.TransactionInvestmentTotals) {
	if err := r.db.Raw(`SELECT COALESCE(SUM(t.amount) FILTER ( WHERE td.sellbuy = 1), 0) as total_buy,
       COALESCE(SUM(t.amount) FILTER ( WHERE td.sellbuy = 0), 0) as total_sell,
       COALESCE(SUM(t.amount) FILTER ( WHERE td.sellbuy = 0), 0) -
       COALESCE(SUM(t.amount) FILTER ( WHERE td.sellbuy = 1), 0) as total_current_portfolio
FROM tbl_transactions t
         INNER JOIN tbl_transaction_details td ON td.id_transactions = t.id
         INNER JOIN tbl_master_transaction_types tmtt ON tmtt.id = t.id_master_transaction_types
WHERE to_char(t.date_time_transaction::DATE, 'YYYY') = ?
  AND t.id_personal_account = ?
  AND tmtt.type = 'INVEST'`, year, IDPersonal).Scan(&data).Error; err != nil {
		return entities.TransactionInvestmentTotals{}
	}
	return data
}

func (r *TransactionRepository) InvestAnnuallyDetail(IDPersonal uuid.UUID, year string) (data []entities.TransactionInvestmentDetail) {
	if err := r.db.Raw(`SELECT concat(to_char(t.date_time_transaction::DATE, 'DD'),' ',to_char(to_date(t.date_time_transaction, 'YYYY-MM-DD'), 'Mon'), ' ',
              to_char(t.date_time_transaction::DATE, 'YYYY'))::text         as date,
       COALESCE(SUM(t.amount) FILTER ( WHERE td.sellbuy = 2 ), 0)::numeric  as total_buy,
       COALESCE(SUM(t.amount) FILTER ( WHERE td.sellbuy = 1 ), 0) ::numeric as total_sell,
       td.lot :: numeric                                                    as lot,
       td.stock_code :: text                                                as stock_code,
		t.amount as price
FROM tbl_transactions t
         INNER JOIN tbl_transaction_details td ON td.id_transactions = t.id
         INNER JOIN tbl_master_transaction_types tmtt ON tmtt.id = t.id_master_transaction_types
WHERE to_char(t.date_time_transaction::DATE, 'YYYY') = ?
  AND t.id_personal_account = ?
  AND tmtt.type = 'INVEST'
GROUP BY t.date_time_transaction, date, lot, stock_code, price
ORDER BY t.date_time_transaction::DATE DESC`, year, IDPersonal).Scan(&data).Error; err != nil {
		return []entities.TransactionInvestmentDetail{}
	}
	return data
}

func (r *TransactionRepository) ByNote(IDPersonal uuid.UUID, month, year string) (data []entities.TransactionByNotes) {
	if err := r.db.Raw(`SELECT
    COALESCE(SUM(t.amount),0) as amount,
    td.note as transaction_note,
    tmec.expense_types
FROM tbl_transactions t
INNER JOIN tbl_transaction_details td ON td.id_transactions = t.id
INNER JOIN tbl_master_transaction_types tmtt ON tmtt.id = t.id_master_transaction_types
LEFT JOIN tbl_master_expense_categories tmec ON t.id_master_expense_categories = tmec.id
WHERE to_char(t.date_time_transaction::DATE, 'MM') = ?
  AND to_char(t.date_time_transaction::DATE, 'YYYY') = ?
  AND t.id_personal_account = ?
  AND tmtt.type = 'EXPENSE'
GROUP BY td.note, tmec.expense_types`, month, year, IDPersonal).Scan(&data).Error; err != nil {
		return []entities.TransactionByNotes{}
	}
	return data
}

func (r *TransactionRepository) Suggestion(IDPersoalAccount uuid.UUID) (data []entities.TransactionSuggestionNotes, err error) {
	if err = r.db.Raw(`SELECT DISTINCT td.note FROM tbl_transactions t
INNER JOIN tbl_transaction_details td ON td.id_transactions = t.id
INNER JOIN tbl_master_transaction_types tmtt ON tmtt.id = t.id_master_transaction_types
WHERE t.id_personal_account=? AND tmtt.type='EXPENSE'`, IDPersoalAccount).Scan(&data).Error; err != nil {
		return []entities.TransactionSuggestionNotes{}, err
	}
	return data, nil
}

func (r *TransactionRepository) WalletExist(IDWallet uuid.UUID) bool {
	var model entities.TransactionWalletExist

	if err := r.db.Raw(`SELECT EXISTS (SELECT 1 FROM tbl_wallets tw WHERE tw.id=?)`, IDWallet).Scan(&model).Error; err != nil {
		return model.Exists
	}
	return model.Exists
}

func (r *TransactionRepository) BudgetWithCurrency(IDTravel uuid.UUID) (data entities.TransactionWithCurrency, rr error) {
	if err := r.db.Raw(`SELECT tmec.currency_value FROM tbl_budgets tb INNER JOIN tbl_master_exchange_currency tmec ON tmec.id = tb.id_master_exchance_currency
WHERE tb.id=?`, IDTravel).
		Scan(&data).Error; err != nil {
		logrus.Error(err.Error())
		return entities.TransactionWithCurrency{}, err
	}
	return data, nil
}

func (r *TransactionRepository) GetTradingInfo(stockCode string) (data entities.InvestmentTreding, err error) {
	if err := r.db.Raw(`SELECT tmd.symbol, tmd.name, tmd.close::numeric FROM tbl_master_trading tmd WHERE tmd.symbol=?`, stockCode).
		Scan(&data).Error; err != nil {
		return entities.InvestmentTreding{}, err
	}
	return data, nil
}