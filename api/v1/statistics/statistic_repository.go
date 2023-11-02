package statistics

import (
	"github.com/google/uuid"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/statistics/entities"
	"gorm.io/gorm"
)

type (
	StatisticRepository struct {
		db *gorm.DB
	}

	IStatisticRepository interface {
		SummaryMonthly(IDPersonal uuid.UUID, month, year string) (data entities.StatisticSummaryMonthly, err error)
		TransactionPriority(IDPersonal uuid.UUID) (data []entities.StatisticTransactionPriority)
		Trend(IDPersonal uuid.UUID) (data entities.StatisticTrend)
		Category(IDPersonal, category uuid.UUID) (data entities.StatisticTrend)
		expenseWeekly(IDPersonal uuid.UUID, month, year string) (data entities.StatisticExpenseWeekly, err error)
		incomeWeekly(IDPersonal uuid.UUID, month, year string) (data entities.StatisticIncomeWeekly, err error)
		investmentWeekly(IDPersonal uuid.UUID, month, year string) (data entities.StatisticInvestmentWeekly, err error)
	}
)

func NewStatisticRepository(db *gorm.DB) *StatisticRepository {
	return &StatisticRepository{db: db}
}

func (r *StatisticRepository) SummaryMonthly(IDPersonal uuid.UUID, month, year string) (data entities.StatisticSummaryMonthly, err error) {
	if err := r.db.Raw(`SELECT COALESCE(SUM(tt.amount) FILTER (WHERE tmtt.type = 'EXPENSE' ), 0)::numeric as total_expense,
       COALESCE(SUM(tt.amount) FILTER (WHERE tmtt.type = 'INCOME' ), 0)::numeric as total_income,
       COALESCE(SUM(tt.amount) FILTER (WHERE tmtt.type = 'TRANSFER' ), 0)::numeric as total_transfer,
       COALESCE(SUM(tt.amount) FILTER (WHERE tmtt.type = 'INVEST' ), 0)::numeric as total_invest
	   FROM tbl_transactions tt
		LEFT JOIN tbl_master_transaction_types tmtt ON tmtt.id = tt.id_master_transaction_types
			WHERE tt.id_personal_account=? AND to_char(tt.date_time_transaction::DATE, 'MM') = ? AND to_char(tt.date_time_transaction::DATE, 'YYYY') = ?`, IDPersonal, month, year).Scan(&data).Error; err != nil {
		return entities.StatisticSummaryMonthly{}, err
	}
	return data, nil
}

func (r *StatisticRepository) TransactionPriority(IDPersonal uuid.UUID) (data []entities.StatisticTransactionPriority) {
	if err := r.db.Raw(`WITH temp AS (SELECT COUNT(tt.id_master_transaction_priorities)::DECIMAL as total_transaction_priorities,
                     COUNT(tt.id_master_transaction_priorities)
                     FILTER ( WHERE tt.id_master_transaction_priorities = 'f05d9cb4-1ae4-4a1a-b566-a906900fdcad' ):: DECIMAL as transaction_need,
                     COUNT(tt.id_master_transaction_priorities)
                     FILTER ( WHERE tt.id_master_transaction_priorities = '9b96cdf8-8173-4d54-9142-e6ebd1f6aea3' ):: DECIMAL as transaction_must,
                     COUNT(tt.id_master_transaction_priorities)
                     FILTER ( WHERE tt.id_master_transaction_priorities = 'd68a049c-7f66-4ab3-a511-c492d3f200c4' ):: DECIMAL as transaction_want
              FROM tbl_transactions tt INNER JOIN tbl_master_transaction_priorities tmtp ON tmtp.id = tt.id_master_transaction_priorities
              WHERE tt.id_personal_account = ? AND to_char(tt.date_time_transaction::DATE, 'MM') = EXTRACT(MONTH FROM current_timestamp)::text
                AND tmtp.active = true)
SELECT (transaction_need / total_transaction_priorities) * 100::DECIMAL AS transaction_need_percentage,
       (transaction_must / total_transaction_priorities) * 100::DECIMAL AS transaction_must_percentage,
       (transaction_want / total_transaction_priorities) * 100::DECIMAL AS transaction_want_percentage
FROM temp`, IDPersonal).Scan(&data).Error; err != nil {
		return []entities.StatisticTransactionPriority{}
	}
	return data
}

func (r *StatisticRepository) Trend(IDPersonal uuid.UUID) (data entities.StatisticTrend) {
	if err := r.db.Raw(`WITH temp AS (SELECT EXTRACT(YEAR FROM current_timestamp)::text  AS year,
                           EXTRACT(MONTH FROM current_timestamp)::text AS month)
      SELECT COALESCE(SUM(tt.amount) FILTER
          ( WHERE tt.id_master_expense_categories IS NOT NULL AND
                  tt.date_time_transaction BETWEEN
                      CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-01') AND
                      CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-04')),
                      0)::text                                             AS "01-04",
             COALESCE(SUM(tt.amount) FILTER
                 ( WHERE tt.id_master_expense_categories IS NOT NULL AND
                         tt.date_time_transaction BETWEEN
                             CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-05') AND
                             CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-11')),
                      0)::text                                             AS "05-11",
             COALESCE(SUM(tt.amount) FILTER
                 ( WHERE tt.id_master_expense_categories IS NOT NULL AND
                         tt.date_time_transaction BETWEEN
                             CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-12') AND
                             CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-18')),
                      0)::text                                             AS "12-18",
             COALESCE(SUM(tt.amount) FILTER
                 ( WHERE tt.id_master_expense_categories IS NOT NULL AND
                         tt.date_time_transaction BETWEEN
                             CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-19') AND
                             CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-25')),
                      0)::text                                             AS "19-25",
             COALESCE(SUM(tt.amount) FILTER
                 ( WHERE tt.id_master_expense_categories IS NOT NULL AND
                         tt.date_time_transaction BETWEEN
                             CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-26') AND
                             CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-30')),
                      0)::text                                             as "26-30",
             COALESCE(SUM(tt.amount) FILTER ( WHERE to_char(tt.date_time_transaction::DATE, 'MM') = EXTRACT(
                     MONTH FROM current_timestamp)::text ), 0)::numeric       as total_average_weekly,
             ROUND(COALESCE(SUM(tt.amount) FILTER ( WHERE to_char(tt.date_time_transaction::DATE, 'MM') = EXTRACT(
                     MONTH FROM current_timestamp)::text ) / 30, 0))::numeric as total_average_daily
      FROM tbl_transactions tt
      WHERE tt.id_personal_account = ?
        AND tt.id_master_expense_categories IS NOT NULL`, IDPersonal).Scan(&data).Error; err != nil {
		return entities.StatisticTrend{}
	}
	return data
}

func (r *StatisticRepository) Category(IDPersonal, category uuid.UUID) (data entities.StatisticTrend) {
	if err := r.db.Raw(`WITH temp AS (SELECT EXTRACT(YEAR FROM current_timestamp)::text  AS year,
                     EXTRACT(MONTH FROM current_timestamp)::text AS month)
SELECT COALESCE(SUM(tt.amount) FILTER
    ( WHERE tt.id_master_expense_categories IS NOT NULL AND
            tt.date_time_transaction BETWEEN
                CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-01') AND
                CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-04')),
                0) ::text                                               AS "01-04",
       COALESCE(SUM(tt.amount) FILTER
           ( WHERE tt.id_master_expense_categories IS NOT NULL AND
                   tt.date_time_transaction BETWEEN
                       CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-05') AND
                       CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-11')),
                0) ::text                                               AS "05-11",
       COALESCE(SUM(tt.amount) FILTER
           ( WHERE tt.id_master_expense_categories IS NOT NULL AND
                   tt.date_time_transaction BETWEEN
                       CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-12') AND
                       CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-18')),
                0) ::text                                               AS "12-18",
       COALESCE(SUM(tt.amount) FILTER
           ( WHERE tt.id_master_expense_categories IS NOT NULL AND
                   tt.date_time_transaction BETWEEN
                       CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-19') AND
                       CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-25')),
                0) ::text                                               AS "19-25",
       COALESCE(SUM(tt.amount) FILTER
           ( WHERE tt.id_master_expense_categories IS NOT NULL AND
                   tt.date_time_transaction BETWEEN
                       CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-26') AND
                       CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-30')),
                0) ::text                                               as "26-30",
       COALESCE(SUM(tt.amount) FILTER ( WHERE to_char(tt.date_time_transaction::DATE, 'MM') = EXTRACT(
               MONTH FROM current_timestamp)::text ), 0)::numeric       as total_average_weekly,
       ROUND(COALESCE(SUM(tt.amount) FILTER ( WHERE to_char(tt.date_time_transaction::DATE, 'MM') = EXTRACT(
               MONTH FROM current_timestamp)::text ) / 30, 0))::numeric as total_average_daily
FROM tbl_transactions tt
WHERE tt.id_personal_account = ?
  AND tt.id_master_expense_categories = ?
  AND tt.id_master_expense_categories IS NOT NULL;`, IDPersonal, category).Scan(&data).Error; err != nil {

	}

	return data
}

func (r *StatisticRepository) expenseWeekly(IDPersonal uuid.UUID, month, year string) (data entities.StatisticExpenseWeekly, err error) {

	sql := `SELECT COALESCE(SUM(tt.amount) FILTER (WHERE tt.date_time_transaction BETWEEN CONCAT('` + year + `', '-', '` + month + ` ', '-01') AND  CONCAT('` + year + `', '-', '` + month + `', '-04')), 0)::numeric as date_range_01_04,
    COALESCE(SUM(tt.amount) FILTER (WHERE tt.date_time_transaction BETWEEN CONCAT('` + year + `', '-', '` + month + `', '-05') AND  CONCAT('` + year + `', '-', '` + month + `', '-11')), 0)::numeric as date_range_05_11,
    COALESCE(SUM(tt.amount) FILTER (WHERE tt.date_time_transaction BETWEEN CONCAT('` + year + `', '-', '` + month + `', '-12') AND  CONCAT('` + year + `', '-', '` + month + `', '-18')), 0)::numeric as date_range_12_18,
    COALESCE(SUM(tt.amount) FILTER (WHERE tt.date_time_transaction BETWEEN CONCAT('` + year + `', '-', '` + month + `', '-19') AND  CONCAT('` + year + `', '-', '` + month + `', '-25')), 0)::numeric as date_range_19_25,
    COALESCE(SUM(tt.amount) FILTER (WHERE tt.date_time_transaction BETWEEN CONCAT('` + year + `', '-', '` + month + `', '-26') AND  CONCAT('` + year + `', '-', '` + month + `', '-30')), 0)::numeric as date_range_26_30
FROM tbl_transactions tt LEFT JOIN tbl_master_transaction_types tmtt ON tmtt.id = tt.id_master_transaction_types WHERE tt.id_personal_account=? AND tmtt.type = 'EXPENSE'`

	if err := r.db.Raw(sql, IDPersonal).Scan(&data).Error; err != nil {
		return entities.StatisticExpenseWeekly{}, err
	}
	return data, nil
}

func (r *StatisticRepository) incomeWeekly(IDPersonal uuid.UUID, month, year string) (data entities.StatisticIncomeWeekly, err error) {

	sql := `SELECT COALESCE(SUM(tt.amount) FILTER (WHERE tt.date_time_transaction BETWEEN CONCAT('` + year + `', '-', '` + month + ` ', '-01') AND  CONCAT('` + year + `', '-', '` + month + `', '-04')), 0)::numeric as date_range_01_04,
    COALESCE(SUM(tt.amount) FILTER (WHERE tt.date_time_transaction BETWEEN CONCAT('` + year + `', '-', '` + month + `', '-05') AND  CONCAT('` + year + `', '-', '` + month + `', '-11')), 0)::numeric as date_range_05_11,
    COALESCE(SUM(tt.amount) FILTER (WHERE tt.date_time_transaction BETWEEN CONCAT('` + year + `', '-', '` + month + `', '-12') AND  CONCAT('` + year + `', '-', '` + month + `', '-18')), 0)::numeric as date_range_12_18,
    COALESCE(SUM(tt.amount) FILTER (WHERE tt.date_time_transaction BETWEEN CONCAT('` + year + `', '-', '` + month + `', '-19') AND  CONCAT('` + year + `', '-', '` + month + `', '-25')), 0)::numeric as date_range_19_25,
    COALESCE(SUM(tt.amount) FILTER (WHERE tt.date_time_transaction BETWEEN CONCAT('` + year + `', '-', '` + month + `', '-26') AND  CONCAT('` + year + `', '-', '` + month + `', '-30')), 0)::numeric as date_range_26_30
FROM tbl_transactions tt LEFT JOIN tbl_master_transaction_types tmtt ON tmtt.id = tt.id_master_transaction_types WHERE tt.id_personal_account=? AND tmtt.type = 'INCOME'`

	if err := r.db.Raw(sql, IDPersonal).Scan(&data).Error; err != nil {
		return entities.StatisticIncomeWeekly{}, err
	}
	return data, nil
}

func (r *StatisticRepository) investmentWeekly(IDPersonal uuid.UUID, month, year string) (data entities.StatisticInvestmentWeekly, err error) {

	sql := `SELECT COALESCE(SUM(tt.amount) FILTER (WHERE tt.date_time_transaction BETWEEN CONCAT('` + year + `', '-', '` + month + ` ', '-01') AND  CONCAT('` + year + `', '-', '` + month + `', '-04')), 0)::numeric as date_range_01_04,
    COALESCE(SUM(tt.amount) FILTER (WHERE tt.date_time_transaction BETWEEN CONCAT('` + year + `', '-', '` + month + `', '-05') AND  CONCAT('` + year + `', '-', '` + month + `', '-11')), 0)::numeric as date_range_05_11,
    COALESCE(SUM(tt.amount) FILTER (WHERE tt.date_time_transaction BETWEEN CONCAT('` + year + `', '-', '` + month + `', '-12') AND  CONCAT('` + year + `', '-', '` + month + `', '-18')), 0)::numeric as date_range_12_18,
    COALESCE(SUM(tt.amount) FILTER (WHERE tt.date_time_transaction BETWEEN CONCAT('` + year + `', '-', '` + month + `', '-19') AND  CONCAT('` + year + `', '-', '` + month + `', '-25')), 0)::numeric as date_range_19_25,
    COALESCE(SUM(tt.amount) FILTER (WHERE tt.date_time_transaction BETWEEN CONCAT('` + year + `', '-', '` + month + `', '-26') AND  CONCAT('` + year + `', '-', '` + month + `', '-30')), 0)::numeric as date_range_26_30
FROM tbl_transactions tt LEFT JOIN tbl_master_transaction_types tmtt ON tmtt.id = tt.id_master_transaction_types WHERE tt.id_personal_account=? AND tmtt.type = 'INVEST'`

	if err := r.db.Raw(sql, IDPersonal).Scan(&data).Error; err != nil {
		return entities.StatisticInvestmentWeekly{}, err
	}
	return data, nil
}