package investments

type (
	InvestmentController struct {
		useCase IInvestmentUseCase
	}

	IInvestmentController interface {
	}
)

func NewInvestmentController(useCase IInvestmentUseCase) *InvestmentController {
	return &InvestmentController{useCase: useCase}
}