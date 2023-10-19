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
		Statistic(IDPersonal uuid.UUID) (data entities.Statistic)
		TransactionPriority(IDPersonal uuid.UUID) (data []entities.StatisticTransactionPriority)
		Trend(IDPersonal uuid.UUID) (data entities.StatisticTrend)
		Category(IDPersonal, category uuid.UUID) (data entities.StatisticTrend)
	}
)

func NewStatisticRepository(db *gorm.DB) *StatisticRepository {
	return &StatisticRepository{db: db}
}

func (r *StatisticRepository) Statistic(IDPersonal uuid.UUID) (data entities.Statistic) {
	if err := r.db.Raw(`WITH temp AS (SELECT EXTRACT(YEAR FROM current_timestamp)::text  AS year,
                     EXTRACT(MONTH FROM current_timestamp)::text AS month)
SELECT COALESCE(SUM(tt.amount) FILTER
    ( WHERE tmtt.type = 'EXPENSE' AND
            tt.date_time_transaction BETWEEN
                CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-01') AND
                CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-04')),
                0)                                                      AS "expense_01-04",
       COALESCE(SUM(tt.amount) FILTER
           ( WHERE tmtt.type = 'EXPENSE' AND
                   tt.date_time_transaction BETWEEN
                       CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-05') AND
                       CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-11')),
                0)                                                      AS "expense_05-11",
       COALESCE(SUM(tt.amount) FILTER
           ( WHERE tmtt.type = 'EXPENSE' AND
                   tt.date_time_transaction BETWEEN
                       CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-12') AND
                       CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-18')),
                0)                                                      AS "expense_12-18",
       COALESCE(SUM(tt.amount) FILTER
           ( WHERE tmtt.type = 'EXPENSE' AND
                   tt.date_time_transaction BETWEEN
                       CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-19') AND
                       CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-25')),
                0)                                                      AS "expense_19-25",
       COALESCE(SUM(tt.amount) FILTER
           ( WHERE tmtt.type = 'EXPENSE' AND
                   tt.date_time_transaction BETWEEN
                       CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-26') AND
                       CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-30')),
                0)                                                      as "expense_26-30",
       COALESCE(SUM(tt.amount) FILTER
           ( WHERE tmtt.type = 'INCOME' AND
                   tt.date_time_transaction BETWEEN
                       CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-01') AND
                       CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-04')),
                0)                                                      AS "income_01-04",
       COALESCE(SUM(tt.amount) FILTER
           ( WHERE tmtt.type = 'INCOME' AND
                   tt.date_time_transaction BETWEEN
                       CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-05') AND
                       CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-11')),
                0)                                                      AS "income_05-11",
       COALESCE(SUM(tt.amount) FILTER
           ( WHERE tmtt.type = 'INCOME' AND
                   tt.date_time_transaction BETWEEN
                       CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-12') AND
                       CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-18')),
                0)                                                      AS "income_12-18",
       COALESCE(SUM(tt.amount) FILTER
           ( WHERE tmtt.type = 'INCOME' AND
                   tt.date_time_transaction BETWEEN
                       CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-19') AND
                       CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-25')),
                0)                                                      AS "income_19-25",
       COALESCE(SUM(tt.amount) FILTER
           ( WHERE tmtt.type = 'INCOME' AND
                   tt.date_time_transaction BETWEEN
                       CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-26') AND
                       CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-30')),
                0)                                                      as "income_26-30",
       COALESCE(SUM(tt.amount) FILTER
           ( WHERE tmtt.type = 'INVEST' AND
                   tt.date_time_transaction BETWEEN
                       CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-01') AND
                       CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-04')),
                0)                                                      AS "invest_01-04",
       COALESCE(SUM(tt.amount) FILTER
           ( WHERE tmtt.type = 'INVEST' AND
                   tt.date_time_transaction BETWEEN
                       CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-05') AND
                       CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-11')),
                0)                                                      AS "invest_05-11",
       COALESCE(SUM(tt.amount) FILTER
           ( WHERE tmtt.type = 'INVEST' AND
                   tt.date_time_transaction BETWEEN
                       CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-12') AND
                       CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-18')),
                0)                                                      AS "invest_12-18",
       COALESCE(SUM(tt.amount) FILTER
           ( WHERE tmtt.type = 'INVEST' AND
                   tt.date_time_transaction BETWEEN
                       CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-19') AND
                       CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-25')),
                0)                                                      AS "invest_19-25",
       COALESCE(SUM(tt.amount) FILTER
           ( WHERE tmtt.type = 'INVEST' AND
                   tt.date_time_transaction BETWEEN
                       CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-26') AND
                       CONCAT((SELECT temp.year FROM temp), '-', (SELECT temp.month FROM temp), '-30')),
                0)                                                      as "invest_26-30",

       COALESCE(SUM(tt.amount) FILTER ( WHERE  tt.id_master_income_categories IS NOT NULL AND to_char(tt.date_time_transaction::DATE, 'YYYY')='2023' AND to_char(tt.date_time_transaction::DATE, 'MM') = EXTRACT(
               MONTH FROM current_timestamp)::text ), 0)::numeric       as total_income,
        COALESCE(SUM(tt.amount) FILTER ( WHERE  tt.id_master_expense_categories IS NOT NULL  AND to_char(tt.date_time_transaction::DATE, 'YYYY')='2023'  AND to_char(tt.date_time_transaction::DATE, 'MM') = EXTRACT(
               MONTH FROM current_timestamp)::text ), 0)::numeric       as total_expense,
    COALESCE(SUM(tt.amount) FILTER ( WHERE  tt.id_master_income_categories IS NOT NULL AND to_char(tt.date_time_transaction::DATE, 'YYYY')='2023' AND to_char(tt.date_time_transaction::DATE, 'MM') = EXTRACT(
               MONTH FROM current_timestamp)::text ), 0) - COALESCE(SUM(tt.amount) FILTER ( WHERE  tt.id_master_expense_categories IS NOT NULL  AND to_char(tt.date_time_transaction::DATE, 'YYYY')='2023'  AND to_char(tt.date_time_transaction::DATE, 'MM') = EXTRACT(
               MONTH FROM current_timestamp)::text ), 0)::numeric       as total_net_income,
     COALESCE(SUM(tt.amount) FILTER ( WHERE  tmtt.type='INVEST' AND to_char(tt.date_time_transaction::DATE, 'YYYY')='2023' AND to_char(tt.date_time_transaction::DATE, 'MM') = EXTRACT(
               MONTH FROM current_timestamp)::text ), 0)::numeric       as total_invest
FROM tbl_transactions tt
         INNER JOIN tbl_master_transaction_types tmtt ON tt.id_master_transaction_types = tmtt.id
WHERE tt.id_personal_account = ?`, IDPersonal).Scan(&data).Error; err != nil {
		return entities.Statistic{}
	}
	return data
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
