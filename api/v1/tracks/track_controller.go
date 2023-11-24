package tracks

type (
	TrackController struct {
		useCase ITrackUseCase
	}

	ITrackController interface {
	}
)

func NewTrackController(useCase ITrackUseCase) *TrackController {
	return &TrackController{useCase: useCase}
}