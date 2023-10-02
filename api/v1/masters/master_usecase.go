package masters

type (
	MasterUseCase struct {
		repo IMasterRepository
	}

	IMasterUseCase interface {
		TransactionType() (data interface{})
		IncomeType() (data interface{})
		ExpenseType() (data interface{})
		ReksadanaType() (data interface{})
		WalletType() (data interface{})
		InvestType() (data interface{})
		Broker() (data interface{})
	}
)

func NewMasterUseCase(repo IMasterRepository) *MasterUseCase {
	return &MasterUseCase{repo: repo}
}

func (s *MasterUseCase) TransactionType() (data interface{}) {
	return s.repo.TransactionType()
}

func (s *MasterUseCase) IncomeType() (data interface{}) {
	return s.repo.IncomeType()
}

func (s *MasterUseCase) ExpenseType() (data interface{}) {
	return s.repo.ExpenseType()
}

func (s *MasterUseCase) ReksadanaType() (data interface{}) {
	return s.repo.ReksadanaType()
}

func (s *MasterUseCase) WalletType() (data interface{}) {
	return s.repo.WalletType()
}

func (s *MasterUseCase) InvestType() (data interface{}) {
	return s.repo.InvestType()
}

func (s *MasterUseCase) Broker() (data interface{}) {
	return s.repo.Broker()
}
