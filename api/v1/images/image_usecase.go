package images

type (
	ShowImageUseCase struct {
		repo IShowImageRepository
	}

	IShowImageUseCase interface {
	}
)

func NewShowImageUseCase(repo IShowImageRepository) *ShowImageUseCase {
	return &ShowImageUseCase{repo: repo}
}