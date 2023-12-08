package budgets

import (
	"github.com/google/uuid"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/budgets/entities"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type (
	BudgetRepository struct {
		db *gorm.DB
	}

	IBudgetRepository interface {
		SubCategoryBudget(IDPersonal uuid.UUID, month, year string) (data []entities.SubCategoryBudget, err error)
		TotalSpendingAndNumberOfCategory(IDPersonal uuid.UUID, month, year string) (data []entities.BudgetTotalSpendingAndNumberOfCategory)
		BudgetLimit(IDPersonal uuid.UUID, month, year string) (data []entities.BudgetLimit)
		Category(IDPersonal uuid.UUID, month string, year string, category uuid.UUID) (data []entities.BudgetCategory)
		LatestMonths(IDPersonal uuid.UUID, category uuid.UUID) (data []entities.BudgetLatestMonth)
		Limit(model *entities.BudgetSetEntities) (err error)
		isBudgetAlreadyExist(model *entities.BudgetSetEntities) (exist bool, id uuid.UUID)
		PersonalBudget(IDPersonal uuid.UUID, month, year string) (data []entities.PersonalBudget, err error)
		PersonalTransaction(IDPersonal uuid.UUID, month, year string) (data []entities.PersonalTransaction, err error)
		Trends(IDPersonal uuid.UUID, IDCategory uuid.UUID, month, year string) (data entities.TrendsWeekly, err error)
		BudgetEachCategory(IDPersonal uuid.UUID, IDCategory uuid.UUID, month, year string) (data entities.BudgetEachCategory, err error)
		CategoryInfo(IDCategory uuid.UUID) (data entities.CategoryInfo, err error)
		Travels(IDPersonal uuid.UUID) (data []entities.BudgetTravel, err error)
		GetXchangeCurrency(IDMasterExchange uuid.UUID) (data entities.BudgetExistsExchangeExist, err error)
	}
)

func NewBudgetRepository(db *gorm.DB) *BudgetRepository {
	return &BudgetRepository{db: db}
}

func (r *BudgetRepository) SubCategoryBudget(IDPersonal uuid.UUID, month, year string) (data []entities.SubCategoryBudget, err error) {
	if err := r.db.Raw(`SELECT tmec.id as category_id,
       tmec.expense_types as category_name,
       tmes.id as sub_category_id,
       tmes.subcategories as sub_category_name,
       (SELECT b.amount FROM tbl_budgets b WHERE b.id_master_subcategories = tmes.id
        AND b.id_personal_accounts = ?
        AND to_char(b.created_at, 'MM') = ?
        AND to_char(b.created_at, 'YYYY') = ? ) as budget_limit
FROM tbl_master_expense_categories tmec
         LEFT JOIN tbl_master_expense_subcategories tmes ON tmes.id_master_expense_categories = tmec.id
ORDER BY tmec.expense_types ASC`, IDPersonal, month, year).Scan(&data).Error; err != nil {
		return []entities.SubCategoryBudget{}, err
	}
	return data, nil
}

func (r *BudgetRepository) TotalSpendingAndNumberOfCategory(IDPersonal uuid.UUID, month, year string) (data []entities.BudgetTotalSpendingAndNumberOfCategory) {
	if err := r.db.Raw(`SELECT tmec.id,tmec.expense_types as category, COALESCE(SUM(tt.amount),0) as spending, COUNT(tt.id_master_expense_categories) as number_of_category
								FROM tbl_master_expense_categories tmec
    							LEFT JOIN tbl_transactions tt ON tt.id_master_expense_categories = tmec.id
								WHERE tt.id_personal_account= ? AND to_char(tt.date_time_transaction::DATE, 'MM') = ? AND to_char(tt.date_time_transaction::DATE, 'YYYY') = ?
							 GROUP BY tmec.expense_types,tmec.id`, IDPersonal, month, year).Scan(&data).Error; err != nil {
		return []entities.BudgetTotalSpendingAndNumberOfCategory{}
	}
	return data
}

func (r *BudgetRepository) BudgetLimit(IDPersonal uuid.UUID, month, year string) (data []entities.BudgetLimit) {
	if err := r.db.Raw(`SELECT tb.id_master_categories as id_master_expense ,COALESCE(SUM(tb.amount),0) as budget_limit, tmec.expense_types
FROM tbl_budgets tb LEFT JOIN tbl_master_expense_categories tmec ON tb.id_master_categories = tmec.id
WHERE tb.id_personal_accounts= ? AND to_char(tb.created_at, 'MM') = ? AND to_char(tb.created_at, 'YYYY') = ?
GROUP BY tmec.expense_types, tb.id_master_categories`, IDPersonal, month, year).Scan(&data).Error; err != nil {
		return []entities.BudgetLimit{}
	}
	return data
}

func (r *BudgetRepository) Category(IDPersonal uuid.UUID, month string, year string, category uuid.UUID) (data []entities.BudgetCategory) {
	if err := r.db.Raw(`SELECT tmec.expense_types::text                              as transaction_category,
       (SELECT COALESCE(SUM(b.amount), 0)
        FROM tbl_budgets b
        WHERE b.id_master_categories = tmec.id
          AND b.id_personal_accounts = ?
          AND to_char(b.created_at, 'MM') = ?
          AND to_char(b.created_at, 'YYYY') = ?)::numeric  as budget_limit,
       (SELECT COALESCE(SUM(tt.amount) FILTER ( WHERE tt.id_master_expense_categories = tmec.id ), 0)
        FROM tbl_transactions tt
        WHERE tt.id_personal_account = ?
          AND tt.id_master_expense_categories = tmec.id
          AND to_char(tt.created_at, 'MM') = ?
          AND to_char(tt.created_at, 'YYYY') = ?)::numeric as total_spending,
       (
               (SELECT COALESCE(SUM(b.amount), 0)
                FROM tbl_budgets b
                WHERE b.id_master_categories = tmec.id
                  AND b.id_personal_accounts = ?
                  AND to_char(b.created_at, 'MM') = ?
                  AND to_char(b.created_at, 'YYYY') = ?) -
               (SELECT COALESCE(SUM(tt.amount) FILTER ( WHERE tt.id_master_expense_categories = tmec.id ), 0)
                FROM tbl_transactions tt
                WHERE tt.id_personal_account = ?
                  AND tt.id_master_expense_categories = tmec.id
                  AND to_char(tt.created_at, 'MM') = ?
                  AND to_char(tt.created_at, 'YYYY') = ?))::numeric as total_remaining FROM tbl_master_expense_categories tmec
         LEFT JOIN tbl_transactions tt ON tt.id_master_expense_categories = tmec.id WHERE tmec.id = ?`, IDPersonal, month, year, IDPersonal, month, year, IDPersonal, month, year, IDPersonal, month, year, category).Scan(&data).Error; err != nil {
		return []entities.BudgetCategory{}
	}
	return data
}

func (r *BudgetRepository) LatestMonths(IDPersonal uuid.UUID, category uuid.UUID) (data []entities.BudgetLatestMonth) {
	if err := r.db.Raw(`SELECT concat(to_char(to_date(tt.date_time_transaction, 'YYYY-MM-DD'), 'Mon'), ' ', to_char(tt.date_time_transaction::DATE, 'YYYY'))::text as period,
       coalesce(sum(tt.amount), 0) ::numeric as total_spending,
       (SELECT coalesce(SUM(b.amount), 0) 
        FROM tbl_budgets b 
        WHERE b.id_personal_accounts = ?)::numeric as budget_limit 
FROM tbl_transactions tt 
WHERE tt.date_time_transaction::date > CURRENT_DATE - INTERVAL '6 MONTHS'
  AND tt.id_personal_account = ?
  AND tt.id_master_expense_categories = ?
GROUP BY period 
ORDER BY period DESC`, IDPersonal, IDPersonal, category).Scan(&data).Error; err != nil {
		return []entities.BudgetLatestMonth{}
	}
	return data
}

func (r *BudgetRepository) Limit(model *entities.BudgetSetEntities) (err error) {
	if err = r.db.Create(&model).Error; err != nil {
		return err
	}

	return nil
}

func (r *BudgetRepository) isBudgetAlreadyExist(model *entities.BudgetSetEntities) (exist bool, id uuid.UUID) {
	var m entities.BudgetExistEntities

	if err := r.db.Raw(`SELECT * FROM tbl_budgets b
WHERE b.id_master_categories=? AND b.id_master_subcategories=? AND b.id_personal_accounts=?
AND to_char(b.created_at, 'MM') = EXTRACT(MONTH FROM current_timestamp)::text
AND to_char(b.created_at, 'YYYY') = EXTRACT(YEAR FROM current_timestamp)::text`, model.IDCategory, model.IDSubCategory, model.IDPersonalAccount).Scan(&m).Error; err != nil {
		return false, uuid.Nil
	}

	if m.ID == uuid.Nil {
		return false, uuid.Nil
	}

	return true, m.ID
}

func (r *BudgetRepository) PersonalBudget(IDPersonal uuid.UUID, month, year string) (data []entities.PersonalBudget, err error) {
	if err := r.db.Raw(`SELECT tmec.id, tmec.expense_types as category, (SELECT b.amount FROM tbl_budgets b
                        WHERE b.id_master_categories = tmec.id AND b.id_personal_accounts=?
                          AND to_char(b.created_at, 'MM') = ? AND to_char(b.created_at, 'YYYY') = ?) as budget FROM tbl_master_expense_categories tmec WHERE tmec.active=true`, IDPersonal, month, year).Scan(&data).Error; err != nil {
		return []entities.PersonalBudget{}, err
	}
	return data, nil
}

func (r *BudgetRepository) PersonalTransaction(IDPersonal uuid.UUID, month, year string) (data []entities.PersonalTransaction, err error) {
	if err := r.db.Raw(`SELECT tmec.id, tmec.expense_types as category, coalesce(SUM(tt.amount),0) as amount, COUNT(tt.id_master_expense_categories)
FROM tbl_transactions tt
LEFT JOIN tbl_master_expense_categories tmec ON tmec.id = tt.id_master_expense_categories
LEFT JOIN tbl_master_transaction_types tmtt ON tmtt.id = tt.id_master_transaction_types
WHERE tt.id_personal_account=?
  AND to_char(tt.date_time_transaction::DATE, 'MM') = ? AND to_char(tt.date_time_transaction::DATE, 'YYYY') = ? AND tmtt.type = 'EXPENSE'
group by tmec.id, tmec.expense_types`, IDPersonal, month, year).Scan(&data).Error; err != nil {
		return []entities.PersonalTransaction{}, err
	}
	return data, nil
}

func (r *BudgetRepository) Trends(IDPersonal uuid.UUID, IDCategory uuid.UUID, month, year string) (data entities.TrendsWeekly, err error) {
	sql := `SELECT COALESCE(SUM(tt.amount) FILTER (WHERE tt.date_time_transaction BETWEEN CONCAT('` + year + `', '-', '` + month + ` ', '-01') AND  CONCAT('` + year + `', '-', '` + month + `', '-04')), 0)::numeric as date_range_01_04,
    COALESCE(SUM(tt.amount) FILTER (WHERE tt.date_time_transaction BETWEEN CONCAT('` + year + `', '-', '` + month + `', '-05') AND  CONCAT('` + year + `', '-', '` + month + `', '-11')), 0)::numeric as date_range_05_11,
    COALESCE(SUM(tt.amount) FILTER (WHERE tt.date_time_transaction BETWEEN CONCAT('` + year + `', '-', '` + month + `', '-12') AND  CONCAT('` + year + `', '-', '` + month + `', '-18')), 0)::numeric as date_range_12_18,
    COALESCE(SUM(tt.amount) FILTER (WHERE tt.date_time_transaction BETWEEN CONCAT('` + year + `', '-', '` + month + `', '-19') AND  CONCAT('` + year + `', '-', '` + month + `', '-25')), 0)::numeric as date_range_19_25,
    COALESCE(SUM(tt.amount) FILTER (WHERE tt.date_time_transaction BETWEEN CONCAT('` + year + `', '-', '` + month + `', '-26') AND  CONCAT('` + year + `', '-', '` + month + `', '-30')), 0)::numeric as date_range_26_30
	FROM tbl_transactions tt LEFT JOIN tbl_master_transaction_types tmtt ON tmtt.id = tt.id_master_transaction_types WHERE tt.id_personal_account=?  AND tt.id_master_expense_categories=? AND tmtt.type = 'EXPENSE'`

	if err := r.db.Raw(sql, IDPersonal, IDCategory).Scan(&data).Error; err != nil {
		return entities.TrendsWeekly{}, err
	}
	return data, nil
}

func (r *BudgetRepository) BudgetEachCategory(IDPersonal uuid.UUID, IDCategory uuid.UUID, month, year string) (data entities.BudgetEachCategory, err error) {
	if err := r.db.Raw(`SELECT tmec.expense_types as category ,COALESCE(SUM(b.amount), 0) as budget_limit
	FROM tbl_budgets b INNER JOIN tbl_master_expense_categories tmec ON tmec.id = b.id_master_categories
	WHERE id_master_categories = ? AND b.id_personal_accounts = ? AND to_char(b.created_at, 'MM') = ? AND to_char(b.created_at, 'YYYY') = ? 
	GROUP BY tmec.expense_types`, IDCategory, IDPersonal, month, year).Scan(&data).Error; err != nil {
		return entities.BudgetEachCategory{}, nil
	}
	return data, nil
}

func (r *BudgetRepository) CategoryInfo(IDCategory uuid.UUID) (data entities.CategoryInfo, err error) {
	if err := r.db.Raw(`SELECT tmec.id as category_id, tmec.expense_types as category_name FROM tbl_master_expense_categories tmec WHERE tmec.id=?`, IDCategory).Scan(&data).Error; err != nil {
		return entities.CategoryInfo{}, err
	}
	return data, nil
}

func (r *BudgetRepository) Travels(IDPersonal uuid.UUID) (data []entities.BudgetTravel, err error) {
	if err := r.db.Raw(`SELECT tb.id, tb.departure,tb.arrival,tb.image_path,tb.filename,tb.amount as budget,tb.travel_start_date,tb.travel_end_date
FROM tbl_budgets tb WHERE tb.id_personal_accounts = ? AND id_master_transaction_types = 'd969fb78-1370-4238-adf0-f143d8a662ef'`, IDPersonal).
		Scan(&data).Error; err != nil {
		logrus.Error(err.Error())
		return []entities.BudgetTravel{}, err
	}
	return data, nil
}

func (r *BudgetRepository) GetXchangeCurrency(IDMasterExchange uuid.UUID) (data entities.BudgetExistsExchangeExist, err error) {
	if err := r.db.Raw(`SELECT EXISTS( SELECT 1 FROM tbl_master_exchange_currency tmec WHERE tmec.id= ?)`, IDMasterExchange).
		Scan(&data).Error; err != nil {
		logrus.Error(err.Error())
		return entities.BudgetExistsExchangeExist{}, err
	}
	return data, nil
}