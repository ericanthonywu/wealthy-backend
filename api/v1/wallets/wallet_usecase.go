package wallets

type (
	WalletUseCase struct {
		repo IWalletRepository
	}

	IWalletUseCase interface {
		Add()
	}
)

func NewWalletUseCase(repo IWalletRepository) *WalletUseCase {
	return &WalletUseCase{repo: repo}
}

func (s *WalletUseCase) Add() {
	s.repo.Add("")
}
