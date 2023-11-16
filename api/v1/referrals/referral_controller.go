package referrals

type (
	ReferralController struct {
		useCase IReferralUseCase
	}

	IReferralController interface {
	}
)

func NewReferralController(useCase IReferralUseCase) *ReferralController {
	return &ReferralController{useCase: useCase}
}