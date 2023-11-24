package tracks

type (
	TrackUseCase struct {
		repo ITrackRepository
	}

	ITrackUseCase interface {
	}
)

func NewTrackUseCase(repo ITrackRepository) *TrackUseCase {
	return &TrackUseCase{repo: repo}
}