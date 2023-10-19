package statistics

type (
	StatisticController struct {
		useCase IStatisticUseCase
	}

	IStatisticController interface {
	}
)

func NewStatisticController(useCase IStatisticUseCase) *StatisticController {
	return &StatisticController{useCase: useCase}
}
