package referrals

type (
	ReferralUseCase struct {
		repo IReferralRepository
	}

	IReferralUseCase interface {
	}
)

func NewReferralUseCase(repo IReferralRepository) *ReferralUseCase {
	return &ReferralUseCase{repo: repo}
}