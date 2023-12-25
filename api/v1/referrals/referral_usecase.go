package referrals

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/referrals/dtos"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/referrals/entities"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/sirupsen/logrus"
	"net/http"
)

type (
	ReferralUseCase struct {
		repo IReferralRepository
	}

	IReferralUseCase interface {
		AccountProfile(personalID uuid.UUID) (data entities.ReferralAccountProfile, err error)
		AccountProfileByRefCode(refCode string) (data entities.ReferralAccountProfileRefCode, err error)
		Statistic(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		List(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		NormalTier(dataMember []entities.ReferralUserReward) (tier []dtos.TierDetail, tierCustomer []dtos.TierDetailWithCustomer)
		UnusualTier(dataMember []entities.ReferralUserReward, currentLevel int) (tier []dtos.TierDetail, tierCustomer []dtos.TierDetailWithCustomer)
		Earn(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		Withdraw(ctx *gin.Context, request dtos.WithdrawRequest) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
	}
)

func NewReferralUseCase(repo IReferralRepository) *ReferralUseCase {
	return &ReferralUseCase{repo: repo}
}

func (s *ReferralUseCase) AccountProfile(personalID uuid.UUID) (data entities.ReferralAccountProfile, err error) {
	return s.repo.AccountProfile(personalID)
}

func (s *ReferralUseCase) AccountProfileByRefCode(refCode string) (data entities.ReferralAccountProfileRefCode, err error) {
	return s.repo.AccountProfileByRefCode(refCode)
}

func (s *ReferralUseCase) Statistic(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		dataProfile entities.ReferralAccountProfile
		dtoResponse dtos.ReferralResponse
		tierData    []dtos.TierDetail
		tierName    string
		err         error
	)

	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	dataProfile, err = s.AccountProfile(accountUUID)
	if err != nil {
		logrus.Error(err.Error())
	}

	referralCode := dataProfile.ReferType

	dataTierOfRefCode, err := s.repo.GetTierReferralCode(referralCode)
	if err != nil {
		logrus.Error(err.Error())
	}

	if len(dataTierOfRefCode) == 0 {
		res := struct {
			Message string `json:"message"`
		}{
			Message: " have not referral below referral code :" + referralCode,
		}
		return res, http.StatusNotFound, []errorsinfo.Errors{}
	}

	// determine mapping
	var titleCollection = map[int]string{
		1: "1st tier",
		2: "2nd tier",
		3: "3nd tier",
		4: "4nd tier",
		5: "5nd tier",
	}

	// response mapping
	levelPrevious := 0
	totalInTier := 0
	maxData := len(dataTierOfRefCode) - 1

	if len(dataTierOfRefCode) > 0 {
		for k, v := range dataTierOfRefCode {

			if v.Level == 0 {
				continue
			}

			if levelPrevious == v.Level {

				// continue increment
				totalInTier++

				if k == maxData {
					// get tier name
					if value, exist := titleCollection[levelPrevious]; exist {
						tierName = value
					}

					tierData = append(tierData, dtos.TierDetail{
						Name:  tierName,
						Value: totalInTier,
					})

					totalInTier = 0
				}
			}

			if levelPrevious == 0 {
				levelPrevious = v.Level
				totalInTier++

				if k == maxData {
					totalInTier++

					// get tier name
					if value, exist := titleCollection[levelPrevious]; exist {
						tierName = value
					}

					tierData = append(tierData, dtos.TierDetail{
						Name:  tierName,
						Value: totalInTier,
					})

					// clear
					totalInTier = 0
				}
			}

			if levelPrevious != v.Level {
				// save previous
				if value, exist := titleCollection[levelPrevious]; exist {
					tierName = value
				}

				tierData = append(tierData, dtos.TierDetail{
					Name:  tierName,
					Value: totalInTier,
				})

				// reset
				totalInTier = 0

				// renew
				totalInTier++
				levelPrevious = v.Level

				if k == maxData {
					// get tier name
					if value, exist := titleCollection[levelPrevious]; exist {
						tierName = value
					}

					tierData = append(tierData, dtos.TierDetail{
						Name:  tierName,
						Value: totalInTier,
					})

					// clear
					totalInTier = 0
				}

			}
		}
	}

	dtoResponse.Tier = tierData

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	return dtoResponse, http.StatusOK, errInfo
}

func (s *ReferralUseCase) List(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		dtoResponse  dtos.ReferralResponseWithCustomer
		tierInfo     []dtos.TierDetailWithCustomer
		customerInfo []dtos.CustomerDetail
		dataProfile  entities.ReferralAccountProfile
		tierName     string
		err          error
	)

	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	dataProfile, err = s.AccountProfile(accountUUID)
	if err != nil {
		logrus.Error(err.Error())
	}

	referralCode := dataProfile.ReferType

	// get tier of referral code :
	dataTierOfRefCode, err := s.repo.GetTierReferralCode(referralCode)
	if err != nil {
		logrus.Error(err.Error())
	}

	if len(dataTierOfRefCode) == 0 {
		res := struct {
			Message string `json:"message"`
		}{
			Message: " have not referral below referral code :" + referralCode,
		}
		return res, http.StatusNotFound, []errorsinfo.Errors{}
	}

	// determine mapping
	var titleCollection = map[int]string{
		1: "1st tier",
		2: "2nd tier",
		3: "3nd tier",
		4: "4nd tier",
		5: "5nd tier",
	}

	// response mapping
	levelPrevious := 0
	totalInTier := 0
	maxData := len(dataTierOfRefCode) - 1

	if len(dataTierOfRefCode) > 0 {
		for k, v := range dataTierOfRefCode {

			dataAccount, err := s.repo.GetAccountInfoFromRefCode(v.RefCode)
			if err != nil {
				logrus.Error(err.Error())
			}

			if v.Level == 0 {
				continue
			}

			if levelPrevious == v.Level {
				// continue to increment
				totalInTier++

				// store into customer info
				customerInfo = append(customerInfo, dtos.CustomerDetail{
					Name:        dataAccount.Name,
					AccountType: dataAccount.Type,
				})

				// latest data
				if k == maxData {
					// get tier name
					if value, exist := titleCollection[levelPrevious]; exist {
						tierName = value
					}

					// append to tier
					tierInfo = append(tierInfo, dtos.TierDetailWithCustomer{
						Name:           tierName,
						Value:          totalInTier,
						CustomerDetail: customerInfo,
					})
				}
			}

			if levelPrevious == 0 {

				// setup information
				levelPrevious = v.Level
				totalInTier++

				// store into customer info
				customerInfo = append(customerInfo, dtos.CustomerDetail{
					Name:        dataAccount.Name,
					AccountType: dataAccount.Type,
				})

				// latest data
				if k == maxData {
					// get tier name
					if value, exist := titleCollection[levelPrevious]; exist {
						tierName = value
					}

					//append to tier
					tierInfo = append(tierInfo, dtos.TierDetailWithCustomer{
						Name:           tierName,
						Value:          totalInTier,
						CustomerDetail: customerInfo,
					})
				}
			}

			if levelPrevious != v.Level {
				// get tier name
				if value, exist := titleCollection[levelPrevious]; exist {
					tierName = value
				}

				//append to tier
				tierInfo = append(tierInfo, dtos.TierDetailWithCustomer{
					Name:           tierName,
					Value:          totalInTier,
					CustomerDetail: customerInfo,
				})

				// clear previous
				customerInfo = nil
				totalInTier = 0

				// renew data
				levelPrevious = v.Level
				totalInTier++

				// store into customer info
				customerInfo = append(customerInfo, dtos.CustomerDetail{
					Name:        dataAccount.Name,
					AccountType: dataAccount.Type,
				})

				// latest data
				if k == maxData {
					// get tier name
					if value, exist := titleCollection[levelPrevious]; exist {
						tierName = value
					}

					//append to tier
					tierInfo = append(tierInfo, dtos.TierDetailWithCustomer{
						Name:           tierName,
						Value:          totalInTier,
						CustomerDetail: customerInfo,
					})
				}
			}
		}
	}

	dtoResponse.Tier = tierInfo

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	return dtoResponse, http.StatusOK, errInfo
}

func (s *ReferralUseCase) NormalTier(dataMember []entities.ReferralUserReward) (tier []dtos.TierDetail, tierCustomer []dtos.TierDetailWithCustomer) {
	var (
		tier1 int
		tier2 int
		tier3 int
		tier4 int
		tier5 int

		name1 []string
		name2 []string
		name3 []string
		name4 []string
		name5 []string

		accType1 []string
		accType2 []string
		accType3 []string
		accType4 []string
		accType5 []string

		customer1 []dtos.CustomerDetail
		customer2 []dtos.CustomerDetail
		customer3 []dtos.CustomerDetail
		customer4 []dtos.CustomerDetail
		customer5 []dtos.CustomerDetail
	)

	for _, v := range dataMember {

		refCode := v.RefCode
		dataProfileByCode, err := s.repo.AccountProfileByRefCode(refCode)

		if err != nil {
			logrus.Error(err.Error())
		}

		if v.Level == 1 {
			tier1 = tier1 + 1
			name1 = append(name1, dataProfileByCode.Name)
			accType1 = append(accType1, dataProfileByCode.AccountType)
		}

		if v.Level == 2 {
			tier2 = tier2 + 1
			name2 = append(name2, dataProfileByCode.Name)
			accType2 = append(accType2, dataProfileByCode.AccountType)
		}

		if v.Level == 3 {
			tier3 = tier3 + 1
			name3 = append(name3, dataProfileByCode.Name)
			accType3 = append(accType3, dataProfileByCode.AccountType)
		}

		if v.Level == 4 {
			tier4 = tier4 + 1
			name4 = append(name4, dataProfileByCode.Name)
			accType4 = append(accType4, dataProfileByCode.AccountType)
		}

		if v.Level == 5 {
			tier5 = tier5 + 1
			name5 = append(name5, dataProfileByCode.Name)
			accType5 = append(accType1, dataProfileByCode.AccountType)
		}
	}

	tier = []dtos.TierDetail{
		{
			Name:  "tier_1",
			Value: tier1,
		},
		{
			Name:  "tier_2",
			Value: tier2,
		},
		{
			Name:  "tier_3",
			Value: tier3,
		},
		{
			Name:  "tier_4",
			Value: tier4,
		},
		{
			Name:  "tier_5",
			Value: tier5,
		},
	}

	for k, _ := range name1 {
		customer1 = append(customer1, dtos.CustomerDetail{
			Name:        name1[k],
			AccountType: accType1[k],
		})
	}

	for k, _ := range name2 {
		customer2 = append(customer2, dtos.CustomerDetail{
			Name:        name2[k],
			AccountType: accType2[k],
		})
	}

	for k, _ := range name3 {
		customer3 = append(customer3, dtos.CustomerDetail{
			Name:        name3[k],
			AccountType: accType3[k],
		})
	}

	for k, _ := range name4 {
		customer4 = append(customer4, dtos.CustomerDetail{
			Name:        name4[k],
			AccountType: accType4[k],
		})
	}

	for k, _ := range name5 {
		customer5 = append(customer5, dtos.CustomerDetail{
			Name:        name5[k],
			AccountType: accType5[k],
		})
	}

	tierCustomer = []dtos.TierDetailWithCustomer{
		{
			Name:           "tier_1",
			Value:          tier1,
			CustomerDetail: customer1,
		},
		{
			Name:           "tier_2",
			Value:          tier2,
			CustomerDetail: customer2,
		},
		{
			Name:           "tier_3",
			Value:          tier3,
			CustomerDetail: customer3,
		},
		{
			Name:           "tier_4",
			Value:          tier4,
			CustomerDetail: customer4,
		},
		{
			Name:           "tier_5",
			Value:          tier5,
			CustomerDetail: customer5,
		},
	}

	return tier, tierCustomer
}

func (s *ReferralUseCase) UnusualTier(dataMember []entities.ReferralUserReward, currentLevel int) (tier []dtos.TierDetail, tierCustomer []dtos.TierDetailWithCustomer) {
	var (
		tier1 int
		tier2 int
		tier3 int
		tier4 int
		tier5 int

		name1 []string
		name2 []string
		name3 []string
		name4 []string
		name5 []string

		accType1 []string
		accType2 []string
		accType3 []string
		accType4 []string
		accType5 []string

		customer1 []dtos.CustomerDetail
		customer2 []dtos.CustomerDetail
		customer3 []dtos.CustomerDetail
		customer4 []dtos.CustomerDetail
		customer5 []dtos.CustomerDetail
	)

	targetTier := 0
	deduction := targetTier - currentLevel

	for _, v := range dataMember {

		refCode := v.RefCode
		dataProfileByCode, err := s.repo.AccountProfileByRefCode(refCode)
		if err != nil {
			logrus.Error(err.Error())
		}

		if v.Level+deduction == 1 {
			tier1 = tier1 + 1
			name1 = append(name1, dataProfileByCode.Name)
			accType1 = append(accType1, dataProfileByCode.AccountType)
		}

		if v.Level+deduction == 2 {
			tier2 = tier2 + 1
			name2 = append(name2, dataProfileByCode.Name)
			accType2 = append(accType2, dataProfileByCode.AccountType)
		}

		if v.Level+deduction == 3 {
			tier3 = tier3 + 1
			name3 = append(name3, dataProfileByCode.Name)
			accType3 = append(accType3, dataProfileByCode.AccountType)
		}

		if v.Level+deduction == 4 {
			tier4 = tier4 + 1
			name4 = append(name4, dataProfileByCode.Name)
			accType4 = append(accType4, dataProfileByCode.AccountType)
		}

		if v.Level+deduction == 5 {
			tier5 = tier5 + 1
			name5 = append(name5, dataProfileByCode.Name)
			accType5 = append(accType5, dataProfileByCode.AccountType)
		}

	}

	tier = []dtos.TierDetail{
		{
			Name:  "tier_1",
			Value: tier1,
		},
		{
			Name:  "tier_2",
			Value: tier2,
		},
		{
			Name:  "tier_3",
			Value: tier3,
		},
		{
			Name:  "tier_4",
			Value: tier4,
		},
		{
			Name:  "tier_5",
			Value: tier5,
		},
	}

	for k, _ := range name1 {
		customer1 = append(customer1, dtos.CustomerDetail{
			Name:        name1[k],
			AccountType: accType1[k],
		})
	}

	for k, _ := range name2 {
		customer2 = append(customer2, dtos.CustomerDetail{
			Name:        name2[k],
			AccountType: accType2[k],
		})
	}

	for k, _ := range name3 {
		customer3 = append(customer3, dtos.CustomerDetail{
			Name:        name3[k],
			AccountType: accType3[k],
		})
	}

	for k, _ := range name4 {
		customer4 = append(customer4, dtos.CustomerDetail{
			Name:        name4[k],
			AccountType: accType4[k],
		})
	}

	for k, _ := range name5 {
		customer5 = append(customer5, dtos.CustomerDetail{
			Name:        name5[k],
			AccountType: accType5[k],
		})
	}

	tierCustomer = []dtos.TierDetailWithCustomer{
		{
			Name:           "tier_1",
			Value:          tier1,
			CustomerDetail: customer1,
		},
		{
			Name:           "tier_2",
			Value:          tier2,
			CustomerDetail: customer2,
		},
		{
			Name:           "tier_3",
			Value:          tier3,
			CustomerDetail: customer3,
		},
		{
			Name:           "tier_4",
			Value:          tier4,
			CustomerDetail: customer4,
		},
		{
			Name:           "tier_5",
			Value:          tier5,
			CustomerDetail: customer5,
		},
	}

	return tier, tierCustomer
}

func (s *ReferralUseCase) Earn(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	dataAccount, err := s.repo.AccountProfile(accountUUID)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	dataCommission, err := s.repo.GetPreviousCommission(dataAccount.ReferType)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	resp := struct {
		Commission float64 `json:"commission_amount"`
	}{
		Commission: dataCommission.Commission,
	}

	return resp, http.StatusOK, errInfo
}

func (s *ReferralUseCase) Withdraw(ctx *gin.Context, request dtos.WithdrawRequest) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	// get id personal account from token
	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	// determine trx id
	trxID := uuid.New()

	// mapping to model
	model := entities.WithdrawEntities{
		ID:                 trxID,
		IDPersonalAccounts: accountUUID,
		AccountNumber:      request.AccountNumber,
		AccountName:        request.AccountName,
		BankIssue:          request.BankIssue,
		Amount:             float64(request.WithdrawAmount),
		Status:             0,
	}

	// save withdraw request
	_, err := s.repo.SaveWithdraws(&model)
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// no error
	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	resp := struct {
		WithdrawID uuid.UUID `json:"withdraw_id"`
		Message    string    `json:"message"`
	}{
		WithdrawID: trxID,
		Message:    "withdraw request is processing, now",
	}
	return resp, http.StatusOK, errInfo

}