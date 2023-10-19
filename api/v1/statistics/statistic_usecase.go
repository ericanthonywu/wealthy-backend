package statistics

type (
	StatisticUseCase struct {
		repo IStatisticRepository
	}

	IStatisticUseCase interface {
	}
)

func NewStatisticUseCase(repo IStatisticRepository) *StatisticUseCase {
	return &StatisticUseCase{repo: repo}
}
