package masters

type (
	MasterController struct {
		useCase IMasterUseCase
	}

	IMasterController interface {
	}
)

func NewMasterController(useCase IMasterUseCase) *MasterController {
	return &MasterController{useCase: useCase}
}
