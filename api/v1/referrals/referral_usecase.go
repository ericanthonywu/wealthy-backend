package referrals

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/referrals/dtos"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/referrals/entities"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/semicolon-indonesia/wealthy-backend/utils/personalaccounts"
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
		Statistic(ctx *gin.Context) (response dtos.ReferralResponse, httpCode int, errInfo []errorsinfo.Errors)
		List(ctx *gin.Context) (response dtos.ReferralResponseWithCustomer, httpCode int, errInfo []errorsinfo.Errors)
		NormalTier(dataMember []entities.ReferralUserReward) (tier []dtos.TierDetail, tierCustomer []dtos.TierDetailWithCustomer)
		UnusualTier(dataMember []entities.ReferralUserReward, currentLevel int) (tier []dtos.TierDetail, tierCustomer []dtos.TierDetailWithCustomer)
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

func (s *ReferralUseCase) Statistic(ctx *gin.Context) (response dtos.ReferralResponse, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		dataProfile   entities.ReferralAccountProfile
		dataFirstNode entities.ReferralUserReward
		dataMember    []entities.ReferralUserReward
		dtoResponse   dtos.ReferralResponse
		err           error
	)

	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "token contains invalid information")
		return dtoResponse, http.StatusUnauthorized, errInfo
	}

	dataProfile, err = s.AccountProfile(personalAccount.ID)
	if err != nil {
		logrus.Error(err.Error())
	}

	referralCode := dataProfile.ReferType

	// GETTING ROOT NODE INFORMATION
	dataFirstNode, err = s.repo.FirstNode(referralCode)
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "not found data from referral code : "+referralCode)
		return dtos.ReferralResponse{}, http.StatusNotFound, errInfo
	}

	// GETTING MEMBER OF NODE
	dataMember, err = s.repo.MemberNode(referralCode)
	if err != nil || len(dataMember) == 0 {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "not found data from referral code : "+referralCode)
		return dtos.ReferralResponse{}, http.StatusNotFound, errInfo
	}

	// MAPPING FOR TIERS
	if dataFirstNode.Level == 0 {
		dtoResponse.Tier, _ = s.NormalTier(dataMember)
	}

	if dataFirstNode.Level > 0 {
		dtoResponse.Tier, _ = s.UnusualTier(dataMember, dataFirstNode.Level)
	}

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	return dtoResponse, http.StatusOK, errInfo
}

func (s *ReferralUseCase) List(ctx *gin.Context) (response dtos.ReferralResponseWithCustomer, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		dtoResponse   dtos.ReferralResponseWithCustomer
		dataProfile   entities.ReferralAccountProfile
		dataFirstNode entities.ReferralUserReward
		dataMember    []entities.ReferralUserReward
		err           error
	)

	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "token contains invalid information")
		return dtoResponse, http.StatusUnauthorized, errInfo
	}

	dataProfile, err = s.AccountProfile(personalAccount.ID)
	if err != nil {
		logrus.Error(err.Error())
	}

	referralCode := dataProfile.ReferType

	// GETTING ROOT NODE INFORMATION
	dataFirstNode, err = s.repo.FirstNode(referralCode)
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "not found data from referral code : "+referralCode)
		return dtos.ReferralResponseWithCustomer{}, http.StatusNotFound, errInfo
	}

	// GETTING MEMBER OF NODE
	dataMember, err = s.repo.MemberNode(referralCode)
	if err != nil || len(dataMember) == 0 {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "not found data from referral code : "+referralCode)
		return dtos.ReferralResponseWithCustomer{}, http.StatusNotFound, errInfo
	}

	// MAPPING FOR TIERS
	if dataFirstNode.Level == 0 {
		_, dtoResponse.Tier = s.NormalTier(dataMember)
	}

	if dataFirstNode.Level > 0 {
		_, dtoResponse.Tier = s.UnusualTier(dataMember, dataFirstNode.Level)
	}

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