package statistics

import (
	"github.com/google/uuid"
	"github.com/wealthy-app/wealthy-backend/api/v1/statistics/entities"
	"gorm.io/gorm"
)

type (
	StatisticRepository struct {
		db *gorm.DB
	}

	IStatisticRepository interface {
		SummaryMonthly(IDPersonal uuid.UUID, month, year string) (data entities.StatisticSummaryMonthly, err error)
		Priority(IDPersonal uuid.UUID, month, year string) (data entities.StatisticPriority)
		Category(IDPersonal, category uuid.UUID) (data entities.StatisticTrend)
		expenseWeekly(IDPersonal uuid.UUID, month, year string) (data entities.StatisticExpenseWeekly, err error)
		incomeWeekly(IDPersonal uuid.UUID, month, year string) (data entities.StatisticIncomeWeekly, err error)
		investmentWeekly(IDPersonal uuid.UUID, month, year string) (data entities.StatisticInvestmentWeekly, err error)
		ExpenseDetail(IDPersonal uuid.UUID, month, year string) (data []entities.StatisticDetailExpense, err error)
		SubExpenseDetail(IDPersonal uuid.UUID, IDCategory uuid.UUID, month, year string) (data entities.StatisticExpenseWeekly, err error)
		AnalyticsTrend(IDPersonal uuid.UUID, typeName, period string) (data []entities.StatisticAnalyticsTrends)
		GetProfileByEmail(email string) (data entities.StatisticAccountProfile, err error)
		TopThreeInvestment(IDPersonal uuid.UUID) (data []entities.TopThreeInvestment, err error)
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

func (r *StatisticRepository) Priority(IDPersonal uuid.UUID, month, year string) (data entities.StatisticPriority) {
	if err := r.db.Raw(`SELECT count(tt.id) FILTER ( WHERE  tmtp.priority = 'NEED' OR tmtp.priority = 'WANT' OR tmtp.priority = 'MUST')::numeric as total_transaction,
       count(tt.id)  FILTER (WHERE tmtp.priority = 'NEED')::numeric as priority_need,
       count(tt.id) FILTER (WHERE tmtp.priority = 'WANT')::numeric as priority_want,
       count(tt.id) FILTER (WHERE tmtp.priority = 'MUST')::numeric as priority_must
FROM tbl_transactions tt LEFT JOIN tbl_master_transaction_priorities tmtp ON tt.id_master_transaction_priorities = tmtp.id
WHERE tt.id_personal_account = ? AND to_char(tt.date_time_transaction::DATE, 'MM') = ? AND to_char(tt.date_time_transaction::DATE, 'YYYY') = ?`, IDPersonal, month, year).Scan(&data).Error; err != nil {
		return entities.StatisticPriority{}
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
    COALESCE(SUM(tt.amount) FILTER (WHERE tt.date_time_transaction BETWEEN CONCAT('` + year + `', '-', '` + month + `', '-26') AND  CONCAT('` + year + `', '-', '` + month + `', '-31')), 0)::numeric as date_range_26_30
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
    COALESCE(SUM(tt.amount) FILTER (WHERE tt.date_time_transaction BETWEEN CONCAT('` + year + `', '-', '` + month + `', '-26') AND  CONCAT('` + year + `', '-', '` + month + `', '-31')), 0)::numeric as date_range_26_30
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
    COALESCE(SUM(tt.amount) FILTER (WHERE tt.date_time_transaction BETWEEN CONCAT('` + year + `', '-', '` + month + `', '-26') AND  CONCAT('` + year + `', '-', '` + month + `', '-31')), 0)::numeric as date_range_26_30
FROM tbl_transactions tt LEFT JOIN tbl_master_transaction_types tmtt ON tmtt.id = tt.id_master_transaction_types WHERE tt.id_personal_account=? AND tmtt.type = 'INVEST'`

	if err := r.db.Raw(sql, IDPersonal).Scan(&data).Error; err != nil {
		return entities.StatisticInvestmentWeekly{}, err
	}
	return data, nil
}

func (r *StatisticRepository) ExpenseDetail(IDPersonal uuid.UUID, month, year string) (data []entities.StatisticDetailExpense, err error) {
	if err := r.db.Raw(`SELECT tt.id_master_expense_categories as id,
       tmec.expense_types              as category,
       COALESCE(SUM(tt.amount), 0)     as amount,
       tmec.image_path                 as transaction_category_icon
FROM tbl_transactions tt
         INNER JOIN tbl_master_expense_categories_editable tmec ON tmec.id = tt.id_master_expense_categories
		 INNER JOIN tbl_master_transaction_types tmtt ON tt.id_master_transaction_types = tmtt.id
WHERE tt.id_personal_account = ?
  and tmtt.type != 'TRAVEL'
  AND tmec.id_personal_accounts = ?
  AND to_char(tt.date_time_transaction::DATE, 'MM') = ?
  AND to_char(tt.date_time_transaction::DATE, 'YYYY') = ?
GROUP BY tmec.expense_types, tt.id_master_expense_categories, tmec.image_path`, IDPersonal, IDPersonal, month, year).Scan(&data).Error; err != nil {
		return []entities.StatisticDetailExpense{}, err
	}
	return data, nil
}

func (r *StatisticRepository) SubExpenseDetail(IDPersonal uuid.UUID, IDCategory uuid.UUID, month, year string) (data entities.StatisticExpenseWeekly, err error) {
	if err := r.db.Raw(`SELECT tmec.expense_types   as category_name,
       COALESCE(SUM(tt.amount)
                FILTER (WHERE tt.date_time_transaction BETWEEN CONCAT(to_char(tt.date_time_transaction::DATE, 'YYYY'), '-', to_char(tt.date_time_transaction::DATE, 'MM'), '-01') AND CONCAT(to_char(tt.date_time_transaction::DATE, 'YYYY'), '-', to_char(tt.date_time_transaction::DATE, 'MM'), '-04')),
                0)::numeric as date_range_01_04,
       COALESCE(SUM(tt.amount)
                FILTER (WHERE tt.date_time_transaction BETWEEN CONCAT(to_char(tt.date_time_transaction::DATE, 'YYYY'), '-', to_char(tt.date_time_transaction::DATE, 'MM'), '-05') AND CONCAT(to_char(tt.date_time_transaction::DATE, 'YYYY'), '-', to_char(tt.date_time_transaction::DATE, 'MM'), '-11')),
                0)::numeric as date_range_05_11,
       COALESCE(SUM(tt.amount)
                FILTER (WHERE tt.date_time_transaction BETWEEN CONCAT(to_char(tt.date_time_transaction::DATE, 'YYYY'), '-', to_char(tt.date_time_transaction::DATE, 'MM'), '-12') AND CONCAT(to_char(tt.date_time_transaction::DATE, 'YYYY'), '-', to_char(tt.date_time_transaction::DATE, 'MM'), '-18')),
                0)::numeric as date_range_12_18,
       COALESCE(SUM(tt.amount)
                FILTER (WHERE tt.date_time_transaction BETWEEN CONCAT(to_char(tt.date_time_transaction::DATE, 'YYYY'), '-', to_char(tt.date_time_transaction::DATE, 'MM'), '-19') AND CONCAT(to_char(tt.date_time_transaction::DATE, 'YYYY'), '-', to_char(tt.date_time_transaction::DATE, 'MM'), '-25')),
                0)::numeric as date_range_19_25,
       COALESCE(SUM(tt.amount)
                FILTER (WHERE tt.date_time_transaction BETWEEN CONCAT(to_char(tt.date_time_transaction::DATE, 'YYYY'), '-', to_char(tt.date_time_transaction::DATE, 'MM'), '-26') AND CONCAT(to_char(tt.date_time_transaction::DATE, 'YYYY'), '-', to_char(tt.date_time_transaction::DATE, 'MM'), '-30')),
                0)::numeric as date_range_26_30
FROM tbl_transactions tt
         LEFT JOIN tbl_master_transaction_types tmtt ON tmtt.id = tt.id_master_transaction_types
         INNER JOIN tbl_master_expense_categories tmec ON tmec.id = tt.id_master_expense_categories
WHERE tt.id_personal_account = ?
  AND tmtt.type = 'EXPENSE'
  AND to_char(tt.date_time_transaction::DATE, 'MM') = ?
  AND to_char(tt.date_time_transaction::DATE, 'YYYY') = ?
  AND tt.id_master_expense_categories = ?
GROUP BY tmec.expense_types`, IDPersonal, month, year, IDCategory).Scan(&data).Error; err != nil {
		return entities.StatisticExpenseWeekly{}, err
	}
	return data, nil
}

func (r *StatisticRepository) AnalyticsTrend(IDPersonal uuid.UUID, typeName, period string) (data []entities.StatisticAnalyticsTrends) {
	periods := period + ` MONTHS`

	sql := `SELECT coalesce(SUM(tt.amount),0) as total,
       concat(TO_CHAR(tt.date_time_transaction::date, 'Mon'),'-', TO_CHAR(tt.date_time_transaction::date, 'YYYY')) as period
FROM tbl_transactions tt
INNER JOIN tbl_master_transaction_types tmtt ON tt.id_master_transaction_types = tmtt.id
WHERE tmtt.type='` + typeName + `' AND tt.date_time_transaction::date > CURRENT_DATE - INTERVAL '` + periods + `'
  AND tt.id_personal_account=?
GROUP BY period`
	if err := r.db.Raw(sql, IDPersonal).Scan(&data).Error; err != nil {
		return []entities.StatisticAnalyticsTrends{}
	}
	return data
}

func (r *StatisticRepository) GetProfileByEmail(email string) (data entities.StatisticAccountProfile, err error) {
	if err := r.db.Raw(`SELECT tmg.id as id_gender, pa.file_name ,pa.image_path, pa.id,pa.username, pa.name, pa.dob as date_of_birth, pa.refer_code, pa.email, tmat.account_type, tmg.gender_name as gender, tmr.roles as user_roles
FROM tbl_personal_accounts pa
INNER JOIN tbl_master_account_types tmat ON tmat.id = pa.id_master_account_types
LEFT JOIN tbl_master_genders tmg ON tmg.id = pa.id_master_gender
INNER JOIN tbl_authentications ta ON ta.id_personal_accounts = pa.id
INNER JOIN tbl_master_roles tmr ON tmr.id = ta.id_master_roles
WHERE pa.email=?`, email).Scan(&data).Error; err != nil {
		return entities.StatisticAccountProfile{}, err
	}
	return data, nil
}

func (r *StatisticRepository) TopThreeInvestment(IDPersonal uuid.UUID) (data []entities.TopThreeInvestment, err error) {
	if err := r.db.Raw(`SELECT ti.stock_code, ti.initial_investment FROM tbl_investment ti WHERE ti.id_personal_accounts = ? ORDER BY ti.initial_investment DESC lIMIT 3`, IDPersonal).Scan(&data).Error; err != nil {
		return []entities.TopThreeInvestment{}, err
	}
	return data, nil
}
