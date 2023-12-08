package investments

type (
	InvestmentUseCase struct {
		repo IInvestmentRepository
	}

	IInvestmentUseCase interface {
	}
)

func NewInvestmentUseCase(repo IInvestmentRepository) *InvestmentUseCase {
	return &InvestmentUseCase{repo: repo}
}