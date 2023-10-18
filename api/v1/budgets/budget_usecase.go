package budgets

import "github.com/gin-gonic/gin"

type (
	BudgetUseCase struct {
		repo IBudgetRepository
	}

	IBudgetUseCase interface {
		AllCategories(ctx *gin.Context)
		Set()
	}
)

func NewBudgetUseCase(repo IBudgetRepository) *BudgetUseCase {
	return &BudgetUseCase{repo: repo}
}

func (s *BudgetUseCase) AllCategories(ctx *gin.Context) {

	//s.repo.AllCategories()
	return
}

func (s *BudgetUseCase) Set() {
	s.repo.Set()
}
