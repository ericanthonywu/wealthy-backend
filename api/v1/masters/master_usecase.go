package masters

type (
	MasterUseCase struct {
		repo IMasterRepository
	}

	IMasterUseCase interface {
	}
)

func NewMasterUseCase(repo IMasterRepository) *MasterUseCase {
	return &MasterUseCase{repo: repo}
}
