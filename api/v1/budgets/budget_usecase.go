package budgets

type (
	BudgetUseCase struct {
		repo IBudgetRepository
	}

	IBudgetUseCase interface {
		AllCategories()
		Set()
	}
)

func NewBudgetUseCase(repo IBudgetRepository) *BudgetUseCase {
	return &BudgetUseCase{repo: repo}
}

func (s *BudgetUseCase) AllCategories() {
	s.repo.AllCategories()
}

func (s *BudgetUseCase) Set() {
	s.repo.Set()
}
