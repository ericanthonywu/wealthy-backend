package masters

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/wealthy-app/wealthy-backend/api/v1/masters/dtos"
	"github.com/wealthy-app/wealthy-backend/api/v1/masters/entities"
	"github.com/wealthy-app/wealthy-backend/utils/errorsinfo"
	"github.com/wealthy-app/wealthy-backend/utils/personalaccounts"
	"github.com/wealthy-app/wealthy-backend/utils/utilities"
	"net/http"
	"strings"
)

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
		TransactionPriority() (data interface{})
		Gender() (data interface{})
		Price(ctx *gin.Context) (data interface{}, httpCode int, errInfo []errorsinfo.Errors)
		SubExpenseCategories(expenseID uuid.UUID) (data interface{})
		Exchange() (data interface{}, httpCode int, errInfo []errorsinfo.Errors)
		PersonalIncomeCategory(ctx *gin.Context) (data interface{}, httpCode int, errInfo []errorsinfo.Errors)
		PersonalExpenseCategory(ctx *gin.Context) (data interface{}, httpCode int, errInfo []errorsinfo.Errors)
		PersonalExpenseSubCategory(ctx *gin.Context, expenseIDUUID uuid.UUID) (data interface{}, httpCode int, errInfo []errorsinfo.Errors)
		RenameIncomeCategory(ctx *gin.Context, id uuid.UUID, request *dtos.RenameCatRequest) (data interface{}, httpCode int, errInfo []errorsinfo.Errors)
		RenameExpenseCategory(ctx *gin.Context, id uuid.UUID, request *dtos.RenameCatRequest) (data interface{}, httpCode int, errInfo []errorsinfo.Errors)
		RenameSubExpenseCategory(ctx *gin.Context, id uuid.UUID, request *dtos.RenameCatRequest) (data interface{}, httpCode int, errInfo []errorsinfo.Errors)
		AddIncomeCategory(ctx *gin.Context, request *dtos.AddCategory) (data interface{}, httpCode int, errInfo []errorsinfo.Errors)
		AddExpenseCategory(ctx *gin.Context, request *dtos.AddCategory) (data interface{}, httpCode int, errInfo []errorsinfo.Errors)
		AddSubExpenseCategory(ctx *gin.Context, request *dtos.AddCategory) (data interface{}, httpCode int, errInfo []errorsinfo.Errors)
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
	var walletResponse []dtos.WalletResponse

	dataWallet := s.repo.WalletType()

	if len(dataWallet) == 0 {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: "no master data for wallet",
		}
		return resp
	}

	if len(dataWallet) > 0 {
		for _, v := range dataWallet {
			walletResponse = append(walletResponse, dtos.WalletResponse{
				ID:         v.ID,
				WalletName: utilities.CapitalizeWords(strings.ReplaceAll(v.WalletType, "_", " ")),
			})
		}
	}

	return walletResponse
}

func (s *MasterUseCase) InvestType() (data interface{}) {
	return s.repo.InvestType()
}

func (s *MasterUseCase) Broker() (data interface{}) {
	return s.repo.Broker()
}

func (s *MasterUseCase) TransactionPriority() (data interface{}) {
	return s.repo.TransactionPriority()
}

func (s *MasterUseCase) Gender() (data interface{}) {
	return s.repo.Gender()
}

func (s *MasterUseCase) SubExpenseCategories(expenseID uuid.UUID) (data interface{}) {
	if s.repo.ExpenseIDExist(expenseID) {
		return s.repo.SubExpenseCategory(expenseID)
	}
	return data
}

func (s *MasterUseCase) Exchange() (data interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	dataExchange, err := s.repo.Exchange()
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return []entities.Exchange{}, http.StatusInternalServerError, errInfo
	}

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	if len(dataExchange) == 0 {
		message := struct {
			Message string `json:"message"`
		}{
			Message: "no data currency exchange",
		}
		return message, http.StatusNotFound, errInfo
	}

	return dataExchange, http.StatusOK, errInfo
}

func (s *MasterUseCase) PersonalIncomeCategory(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {

	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	data, err := s.repo.PersonalIncomeCategory(accountUUID)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	if len(data) == 0 {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: "no data for income category",
		}
		return resp, http.StatusBadRequest, []errorsinfo.Errors{}
	}

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	return data, http.StatusOK, errInfo
}

func (s *MasterUseCase) PersonalExpenseCategory(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	data, err := s.repo.PersonalExpenseCategory(accountUUID)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return data, http.StatusInternalServerError, errInfo
	}

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	return data, http.StatusOK, errInfo
}

func (s *MasterUseCase) PersonalExpenseSubCategory(ctx *gin.Context, expenseIDUUID uuid.UUID) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		httpCode = http.StatusUnauthorized
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "token contains invalid information")
		return response, httpCode, errInfo
	}

	data, err := s.repo.PersonalExpenseSubCategory(personalAccount.ID, expenseIDUUID)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return data, http.StatusInternalServerError, errInfo
	}

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	return data, http.StatusOK, errInfo
}

func (s *MasterUseCase) RenameIncomeCategory(ctx *gin.Context, id uuid.UUID, request *dtos.RenameCatRequest) (data interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		err         error
		collections = make(map[string]string)
	)

	// account id uuid format
	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	// get all income categories
	dataAllCategories, err := s.repo.GetAllIncomeCategories(accountUUID)
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// if data empty
	if len(dataAllCategories) == 0 {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: "no data for rename category name base category id : " + id.String(),
		}
		return resp, http.StatusBadRequest, []errorsinfo.Errors{}
	}

	// if data not empty
	if len(dataAllCategories) > 0 {
		// mapping collections
		for _, v := range dataAllCategories {
			collections[v] = v
		}
	}

	// check request with collections
	_, exist := collections[request.NewCategoryName]
	if exist {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: "category name is identically with another category in same group. please another one",
		}
		return resp, http.StatusBadRequest, []errorsinfo.Errors{}
	}

	// rename category
	err = s.repo.RenameIncomeCategory(request.NewCategoryName, id, accountUUID)
	if err != nil {
		logrus.Error()
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// if no error
	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: "rename income category success",
	}
	return response, http.StatusOK, errInfo
}

func (s *MasterUseCase) RenameExpenseCategory(ctx *gin.Context, id uuid.UUID, request *dtos.RenameCatRequest) (data interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		collections = make(map[string]string)
	)

	// account id uuid format
	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	// get all expense categories
	dataAllCategories, err := s.repo.GetAllIExpenseCategories(accountUUID)
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// if data not empty
	if len(dataAllCategories) > 0 {
		// mapping collections
		for _, v := range dataAllCategories {
			collections[v] = v
		}
	}

	// check request with collections
	_, exist := collections[request.NewCategoryName]
	if exist {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: "category name is identically with another category in same group. please another one",
		}
		return resp, http.StatusBadRequest, []errorsinfo.Errors{}
	}

	// rename expense category
	err = s.repo.RenameExpenseCategory(request.NewCategoryName, id, accountUUID)
	if err != nil {
		logrus.Error()
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// if no error
	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: "rename expense category success",
	}
	return response, http.StatusOK, errInfo
}

func (s *MasterUseCase) RenameSubExpenseCategory(ctx *gin.Context, id uuid.UUID, request *dtos.RenameCatRequest) (data interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		collections = make(map[string]string)
	)

	// account uuid format
	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	// get all sub expense categories
	dataAllCategories, err := s.repo.GetAllISubExpenseCategories(accountUUID)
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// if data not empty
	if len(dataAllCategories) > 0 {
		// mapping collections
		for _, v := range dataAllCategories {
			collections[v] = v
		}
	}

	// check request with collections
	_, exist := collections[request.NewCategoryName]
	if exist {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: "category name is identically with another category in same group. please another one",
		}
		return resp, http.StatusBadRequest, []errorsinfo.Errors{}
	}

	// save sub expense category name
	err = s.repo.RenameSubExpenseCategory(request.NewCategoryName, id, accountUUID)
	if err != nil {
		logrus.Error()
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// if no error
	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	response := struct {
		Message string `json:"message"`
	}{
		Message: "rename sub-expense category success",
	}
	return response, http.StatusOK, errInfo
}

func (s *MasterUseCase) AddIncomeCategory(ctx *gin.Context, request *dtos.AddCategory) (data interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		dataID      entities.AddEntities
		err         error
		collections = make(map[string]string)
	)

	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	// get all income categories
	dataAllCategories, err := s.repo.GetAllIncomeCategories(accountUUID)
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// if data empty
	if len(dataAllCategories) == 0 {
		// save new income category
		dataID, err = s.repo.AddIncomeCategory(request.CategoryName, accountUUID)
		if err != nil {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
			return struct{}{}, http.StatusInternalServerError, errInfo
		}
	}

	// if data not empty
	if len(dataAllCategories) > 0 {
		// mapping collections
		for _, v := range dataAllCategories {
			collections[v] = v
		}
	}

	// check request with collections
	_, exist := collections[request.CategoryName]
	if exist {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: "new category name already existed. try another one",
		}
		return resp, http.StatusBadRequest, []errorsinfo.Errors{}
	}

	// save new income category
	dataID, err = s.repo.AddIncomeCategory(request.CategoryName, accountUUID)
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// if no error
	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	return dataID, http.StatusOK, errInfo
}

func (s *MasterUseCase) AddExpenseCategory(ctx *gin.Context, request *dtos.AddCategory) (data interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		dataID      entities.AddEntities
		err         error
		collections = make(map[string]string)
	)

	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	// get all income categories
	dataAllCategories, err := s.repo.GetAllIExpenseCategories(accountUUID)
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// if data empty
	if len(dataAllCategories) == 0 {
		// save new expense category
		dataID, err = s.repo.AddExpenseCategory(request.CategoryName, accountUUID)
		if err != nil {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
			return struct{}{}, http.StatusInternalServerError, errInfo
		}
	}

	// if data not empty
	if len(dataAllCategories) > 0 {
		// mapping collections
		for _, v := range dataAllCategories {
			collections[v] = v
		}
	}

	// check request with collections
	_, exist := collections[request.CategoryName]
	if exist {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: "new category name already existed. try another one",
		}
		return resp, http.StatusBadRequest, []errorsinfo.Errors{}
	}

	// save new outcome category
	dataID, err = s.repo.AddExpenseCategory(request.CategoryName, accountUUID)
	if err != nil {
		return data, http.StatusInternalServerError, errInfo
	}

	// if no error
	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	return dataID, http.StatusOK, errInfo
}

func (s *MasterUseCase) AddSubExpenseCategory(ctx *gin.Context, request *dtos.AddCategory) (data interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		dataID      entities.AddEntities
		err         error
		collections = make(map[string]string)
	)

	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	// get all income categories
	dataAllCategories, err := s.repo.GetAllISubExpenseCategories(accountUUID)
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// if data empty
	if len(dataAllCategories) == 0 {
		// save new sub expense category
		dataID, err = s.repo.AddSubExpenseCategory(request.CategoryName, request.ExpenseID, accountUUID)
		if err != nil {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
			return struct{}{}, http.StatusInternalServerError, errInfo
		}
	}

	// if data not empty
	if len(dataAllCategories) > 0 {
		// mapping collections
		for _, v := range dataAllCategories {
			collections[v] = v
		}
	}

	// check request with collections
	_, exist := collections[request.CategoryName]
	if exist {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: "new category name already existed. try another one",
		}
		return resp, http.StatusBadRequest, []errorsinfo.Errors{}
	}

	// save sub category expense category
	dataID, err = s.repo.AddSubExpenseCategory(request.CategoryName, request.ExpenseID, accountUUID)
	if err != nil {
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// if no error
	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	return dataID, http.StatusOK, errInfo
}

func (s *MasterUseCase) Price(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		dtoResponse    []dtos.PriceResponse
		idSubscription uuid.UUID
		isCurrent      bool
	)

	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	dataSubscription, err := s.repo.UserSubscriptionInfo(accountUUID)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	if dataSubscription.ID != uuid.Nil {
		idSubscription = dataSubscription.ID
	} else {
		dataSubscription.ID = uuid.Nil
	}

	dataPrice := s.repo.Price()
	if len(dataPrice) == 0 {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: "code for price not available",
		}
		return resp, http.StatusBadRequest, []errorsinfo.Errors{}
	}

	if len(dataPrice) > 0 {
		for _, v := range dataPrice {
			isCurrent = false

			if idSubscription == v.ID {
				isCurrent = true
			}

			dtoResponse = append(dtoResponse, dtos.PriceResponse{
				ID:            v.ID,
				Title:         v.Title,
				Price:         v.Price,
				ActualPrice:   v.ActualPrice,
				Description:   v.Description,
				IsCurrent:     isCurrent,
				IsRecommended: v.IsRecommended,
			})
		}
	}

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	return dtoResponse, http.StatusOK, errInfo
}