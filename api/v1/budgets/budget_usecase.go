package budgets

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/wealthy-app/wealthy-backend/api/v1/budgets/dtos"
	"github.com/wealthy-app/wealthy-backend/api/v1/budgets/entities"
	"github.com/wealthy-app/wealthy-backend/constants"
	"github.com/wealthy-app/wealthy-backend/utils/datecustoms"
	"github.com/wealthy-app/wealthy-backend/utils/errorsinfo"
	"github.com/wealthy-app/wealthy-backend/utils/personalaccounts"
	"github.com/wealthy-app/wealthy-backend/utils/utilities"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type (
	BudgetUseCase struct {
		repo IBudgetRepository
	}

	IBudgetUseCase interface {
		AllLimit(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		Overview(ctx *gin.Context, month, year string) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		LatestMonths(ctx *gin.Context, categoryID uuid.UUID) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		Limit(ctx *gin.Context, dtoRequest *dtos.BudgetSetRequest, purpose string) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		Trends(ctx *gin.Context, IDCategory uuid.UUID, month, year string) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		Travels(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		UpdateTravelInfo(ctx *gin.Context, IDTravel string, request map[string]interface{}) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
	}
)

func NewBudgetUseCase(repo IBudgetRepository) *BudgetUseCase {
	return &BudgetUseCase{repo: repo}
}

func (s *BudgetUseCase) AllLimit(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		dtoResponse     dtos.AllBudgetLimit
		budgetDetail    []dtos.AllBudgetDetail
		subCategoryInfo []dtos.SubCategoryInfo
		stringBuilder   strings.Builder
	)

	month := fmt.Sprintf("%02s", ctx.Query("month"))
	year := ctx.Query("year")

	// accountUUID
	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	// get all category by accountID
	dataCategory, err := s.repo.CategoryByAccountID(accountUUID)
	if err != nil {
		logrus.Error(err.Error())
	}

	if len(dataCategory) > 0 {
		for _, v := range dataCategory {

			// initial
			totalBudgetCategory := 0.0

			// get sub category by categoryID
			dataSubCategory, err := s.repo.GetSubCategory(accountUUID, v.CategoryID)
			if err != nil {
				logrus.Error(err.Error())
			}

			// if data sub category exist , get amount from tbl_budget by sub_category_ID
			if len(dataSubCategory) > 0 {
				for _, v := range dataSubCategory {
					dataBudgetSubCat, err := s.repo.GetAmountBudgetSubCategory(accountUUID, v.SubCategoryID, month, year)
					if err != nil {
						logrus.Error(err.Error())
					}

					// mapping sub category
					subCategoryInfo = append(subCategoryInfo, dtos.SubCategoryInfo{
						SubCategoryID:   v.SubCategoryID,
						SubCategoryName: v.SubCategoryName,
						SubCategoryIcon: v.SubCategoryIcon,
						BudgetLimit: dtos.Limit{
							CurrencyCode: "IDR",
							Value:        int(dataBudgetSubCat.Amount),
						},
					})

					// override
					totalBudgetCategory += dataBudgetSubCat.Amount
				}
			}

			// if no data sub category exist, get amount from tbl_badget by category_ID
			if len(dataSubCategory) == 0 {
				dataBudgetCategory, err := s.repo.GetAmountBudgetCategory(accountUUID, v.CategoryID, month, year)
				if err != nil {
					logrus.Error(err.Error())
				}

				totalBudgetCategory = totalBudgetCategory + dataBudgetCategory.Amount

				// mapping sub category
				subCategoryInfo = []dtos.SubCategoryInfo{}
			}

			// integration mapping
			budgetDetail = append(budgetDetail, dtos.AllBudgetDetail{
				CategoryID:   v.CategoryID,
				CategoryName: v.CategoryName,
				CategoryIcon: v.CategoryIcon,
				BudgetInfo: dtos.Limit{
					CurrencyCode: "IDR",
					Value:        int(totalBudgetCategory),
				},
				SubCategory: subCategoryInfo,
			})

			// reset
			subCategoryInfo = nil
		}
	}

	monthINT, err := strconv.Atoi(month)
	if err != nil {
		logrus.Error(err.Error())
	}

	stringBuilder.WriteString(datecustoms.IntToMonthName(monthINT))
	stringBuilder.WriteString(" ")
	stringBuilder.WriteString(year)

	dtoResponse.Period = stringBuilder.String()
	dtoResponse.AllBudgetDetail = budgetDetail
	return dtoResponse, http.StatusOK, []errorsinfo.Errors{}
}

func (s *BudgetUseCase) Overview(ctx *gin.Context, month, year string) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		dtoResponse   dtos.BudgetOverview
		dataDetails   []dtos.OverviewDetail
		stringBuilder strings.Builder
	)

	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	// get all category by accountID
	dataCategory, err := s.repo.CategoryByAccountID(accountUUID)
	if err != nil {
		logrus.Error(err.Error())
	}

	if len(dataCategory) > 0 {
		for _, v := range dataCategory {
			// initial
			totalBudgetCategory := 0.0

			// get sub category by categoryID
			dataSubCategory, err := s.repo.GetSubCategory(accountUUID, v.CategoryID)
			if err != nil {
				logrus.Error(err.Error())
			}

			// if data sub category exist , get amount from tbl_budget by sub_category_ID
			if len(dataSubCategory) > 0 {
				for _, v := range dataSubCategory {
					dataBudgetSubCat, err := s.repo.GetAmountBudgetSubCategory(accountUUID, v.SubCategoryID, month, year)
					if err != nil {
						logrus.Error(err.Error())
					}

					// override
					totalBudgetCategory += dataBudgetSubCat.Amount
				}
			}

			// if no data sub category exist, get amount from tbl_badget by category_ID
			if len(dataSubCategory) == 0 {
				dataBudgetCategory, err := s.repo.GetAmountBudgetCategory(accountUUID, v.CategoryID, month, year)
				if err != nil {
					logrus.Error(err.Error())
				}

				totalBudgetCategory += dataBudgetCategory.Amount
			}

			// get transaction by category id
			dataTransaction, err := s.repo.GetTransactionByCategory(accountUUID, v.CategoryID, month, year)
			if err != nil {
				logrus.Error(err.Error())
			}

			// get number of transactions by category id
			dataNumberOfTransactions, err := s.repo.GetNumberOfTransactionByCategory(accountUUID, v.CategoryID, month, year)
			if err != nil {
				logrus.Error(err.Error())
			}

			// integration mapping
			dataDetails = append(dataDetails, dtos.OverviewDetail{
				CategoryName:    v.CategoryName,
				CategoryID:      v.CategoryID,
				TransactionIcon: v.CategoryIcon,
				BudgetLimit: dtos.Limit{
					CurrencyCode: "IDR",
					Value:        int(totalBudgetCategory),
				},
				TransactionSpending: dtos.Transaction{
					CurrencyCode: "IDR",
					Value:        int(dataTransaction.Amount),
				},
				NumberOfTransaction: int(dataNumberOfTransactions.NumberOfTrx),
			})
		}
	}

	monthINT, err := strconv.Atoi(month)
	if err != nil {
		logrus.Error(err.Error())
	}

	stringBuilder.WriteString(datecustoms.IntToMonthName(monthINT))
	stringBuilder.WriteString(" ")
	stringBuilder.WriteString(year)

	dtoResponse.Period = stringBuilder.String()
	dtoResponse.Details = dataDetails

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	return dtoResponse, http.StatusOK, errInfo
}

func (s *BudgetUseCase) LatestMonths(ctx *gin.Context, categoryID uuid.UUID) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		dtoResponse dtos.LatestMonth
		details     []dtos.LatestDetails
	)

	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		httpCode = http.StatusUnauthorized
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "token contains invalid information")
		return []dtos.LatestDetails{}, httpCode, errInfo
	}

	dataLatestMonth := s.repo.LatestMonths(personalAccount.ID, categoryID)

	if len(dataLatestMonth) == 0 {
		httpCode = http.StatusOK
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "no data")
		return []dtos.LatestDetails{}, httpCode, errInfo
	}

	for _, v := range dataLatestMonth {
		var percentage string
		var statusComparison string

		// PERCENTAGE
		if v.BudgetLimit <= 0 {
			percentage = "0%"
			statusComparison = "BELOW BUDGET"
		} else {
			percentage = fmt.Sprintf("%.f", (float64(v.TotalSpending)/float64(v.BudgetLimit))*100) + "%"
			if (float64(v.TotalSpending)/float64(v.BudgetLimit))*100 > 100 {
				statusComparison = "OVERSPENT"
			} else if (float64(v.TotalSpending)/float64(v.BudgetLimit))*100 == 100 {
				statusComparison = "RISK Of SPENT ( IN LIMIT )"
			} else {
				statusComparison = "BELOW BUDGET"
			}
		}

		details = append(details, dtos.LatestDetails{
			Period: v.Period,
			BudgetInfo: dtos.Limit{
				CurrencyCode: "IDR",
				Value:        v.BudgetLimit,
			},
			TransactionSpending: dtos.Transaction{
				CurrencyCode: "IDR",
				Value:        v.TotalSpending,
			},
			Percentage:       percentage,
			StatusComparison: statusComparison,
		})
	}

	dtoResponse.Details = details
	return dtoResponse, http.StatusOK, []errorsinfo.Errors{}
}

func (s *BudgetUseCase) Limit(ctx *gin.Context, dtoRequest *dtos.BudgetSetRequest, purpose string) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		model              entities.BudgetSetEntities
		dtoResponse        dtos.BudgetSetResponse
		filename           string
		targetPath         string
		imagePath          string
		isSubCategoryValid bool
	)

	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	if purpose == constants.Travel {
		imageData, err := base64.StdEncoding.DecodeString(dtoRequest.ImageBase64)
		filename = fmt.Sprintf("%d", time.Now().Unix()) + ".png"
		targetPath = "assets/travel/" + filename
		imagePath = "images/travel/" + filename

		err = utilities.SaveImage(imageData, targetPath)
		if err != nil {
			logrus.Error(err.Error())
		}

		IDMasterExchangeCurrencyUUID, err := uuid.Parse(dtoRequest.IDMasterExchangeCurrency)
		if err != nil {
			logrus.Error(err.Error())
		}

		// category and subcategory injection
		IDCategoryUUID, err := uuid.Parse("d5d3b801-2d10-4d4d-8e16-67599876cbc8")
		if err != nil {
			logrus.Error(err.Error())
		}

		IDSubCategoryUUID, err := uuid.Parse("ff34f563-5d60-45c6-90c3-af08cd100376")
		if err != nil {
			logrus.Error(err.Error())
		}

		// validate id categories exchange currency
		dataCurrency, err := s.repo.GetXchangeCurrency(IDMasterExchangeCurrencyUUID)
		if err != nil {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
			return struct{}{}, http.StatusInternalServerError, errInfo
		}

		if !dataCurrency.Exists {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id categories exchange currency unknown")
			return struct{}{}, http.StatusBadRequest, errInfo
		}

		model.Departure = dtoRequest.Departure
		model.Arrival = dtoRequest.Arrival
		model.Filename = filename
		model.ImagePath = imagePath
		model.TravelStartDate = dtoRequest.TravelStartDate
		model.TravelEndDate = dtoRequest.TravelEndDate
		model.IDMasterTransactionType = dtoRequest.IDMasterTransactionTypes
		model.IDMasterExchangeCurrency = IDMasterExchangeCurrencyUUID
		model.IDCategory = IDCategoryUUID
		model.IDSubCategory = IDSubCategoryUUID
		model.Currency = dtoRequest.ExchangeRate
	}

	if purpose != constants.Travel {
		model.IDCategory = dtoRequest.IDCategory

		// fetch sub category information by categoryID
		dataSubCategory, err := s.repo.GetSubCategory(accountUUID, dtoRequest.IDCategory)
		if err != nil {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
			return struct{}{}, http.StatusInternalServerError, errInfo
		}

		if len(dataSubCategory) > 0 {

			if dtoRequest.IDSubCategory == uuid.Nil {
				resp := struct {
					Message string `json:"message"`
				}{
					Message: "sub category id required",
				}
				return resp, http.StatusBadRequest, []errorsinfo.Errors{}
			}

			// check sub category input same as in database
			for _, v := range dataSubCategory {
				if dtoRequest.IDSubCategory == v.SubCategoryID {
					isSubCategoryValid = true
					break
				}
			}

			// is sub category ID valid
			if !isSubCategoryValid {
				resp := struct {
					Message string `json:"message"`
				}{
					Message: "sub category id invalid for category id",
				}
				return resp, http.StatusBadRequest, []errorsinfo.Errors{}
			}

			// set id sub category
			model.IDSubCategory = dtoRequest.IDSubCategory
		}

		if len(dataSubCategory) == 0 {
			model.IDSubCategory = uuid.Nil
		}

	}

	model.Amount = int64(float64(dtoRequest.Amount) * dtoRequest.ExchangeRate)
	model.IDPersonalAccount = accountUUID
	model.ID = uuid.New()

	err := s.repo.Limit(&model)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "insert new budget issue")
		return response, http.StatusInternalServerError, errInfo
	}

	// if err empty
	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	dtoResponse.ID = model.ID
	dtoResponse.Status = true
	return dtoResponse, httpCode, errInfo
}

func (s *BudgetUseCase) Trends(ctx *gin.Context, IDCategory uuid.UUID, month, year string) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		dtoResponse              dtos.Trends
		trendsInfo               []dtos.TrendsInfo
		totalSpendingTransaction int
		totalRemains             int
		stringBuilder            strings.Builder
	)

	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	dataTrends, err := s.repo.Trends(accountUUID, IDCategory, month, year)
	if err != nil {
		httpCode = http.StatusInternalServerError
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return entities.TrendsWeekly{}, httpCode, errInfo
	}

	// APPEND PROCESS
	trendsInfo = append(trendsInfo, dtos.TrendsInfo{
		StartDate: year + "-" + month + "-01",
		EndDate:   year + "-" + month + "-04",
		TransactionAmount: dtos.Transaction{
			CurrencyCode: "IDR",
			Value:        dataTrends.DateRange0104,
		},
	})

	trendsInfo = append(trendsInfo, dtos.TrendsInfo{
		StartDate: year + "-" + month + "-05",
		EndDate:   year + "-" + month + "-11",
		TransactionAmount: dtos.Transaction{
			CurrencyCode: "IDR",
			Value:        dataTrends.DateRange0511,
		},
	})

	trendsInfo = append(trendsInfo, dtos.TrendsInfo{
		StartDate: year + "-" + month + "-12",
		EndDate:   year + "-" + month + "-18",
		TransactionAmount: dtos.Transaction{
			CurrencyCode: "IDR",
			Value:        dataTrends.DateRange1218,
		},
	})

	trendsInfo = append(trendsInfo, dtos.TrendsInfo{
		StartDate: year + "-" + month + "-19",
		EndDate:   year + "-" + month + "-25",
		TransactionAmount: dtos.Transaction{
			CurrencyCode: "IDR",
			Value:        dataTrends.DateRange1925,
		},
	})

	trendsInfo = append(trendsInfo, dtos.TrendsInfo{
		StartDate: year + "-" + month + "-26",
		EndDate:   year + "-" + month + "-30",
		TransactionAmount: dtos.Transaction{
			CurrencyCode: "IDR",
			Value:        dataTrends.DateRange2630,
		},
	})

	dataBudgetEachCategory, err := s.repo.BudgetEachCategory(accountUUID, IDCategory, month, year)
	if err != nil {
		logrus.Error(err.Error())
	}

	dataCategoryInfo, err := s.repo.CategoryInfo(IDCategory)
	if err != nil {
		logrus.Error(err.Error())
	}

	// PERIOD
	monthINT, err := strconv.Atoi(month)
	if err != nil {
		logrus.Error(err.Error())
	}

	stringBuilder.WriteString(datecustoms.IntToMonthName(monthINT))
	stringBuilder.WriteString(" ")
	stringBuilder.WriteString(year)

	// TOTAL EXPENSE TRANSACTION
	totalSpendingTransaction = dataTrends.DateRange0104 + dataTrends.DateRange0511 + dataTrends.DateRange1218 + dataTrends.DateRange1925 + dataTrends.DateRange2630

	// REMAINS BUDGET
	totalRemains = dataBudgetEachCategory.BudgetLimit - totalSpendingTransaction

	dtoResponse.Period = stringBuilder.String()
	dtoResponse.CategoryID = dataCategoryInfo.CategoryID
	dtoResponse.CategoryName = dataCategoryInfo.CategoryName
	dtoResponse.BudgetInfo.CurrencyCode = "IDR"
	dtoResponse.BudgetInfo.Value = dataBudgetEachCategory.BudgetLimit
	dtoResponse.TrendsInfo = trendsInfo
	dtoResponse.Expense.TransactionSpending.CurrencyCode = "IDR"
	dtoResponse.Expense.TransactionSpending.Value = totalSpendingTransaction
	dtoResponse.Expense.BudgetRemains.CurrencyCode = "IDR"
	dtoResponse.Expense.BudgetRemains.Value = totalRemains
	dtoResponse.Expense.AverageDailySpending.CurrencyCode = "IDR"
	dtoResponse.Expense.AverageDailySpending.Value = totalSpendingTransaction / 30
	dtoResponse.Expense.AverageDailySpendingRecommended.CurrencyCode = "IDR"
	dtoResponse.Expense.AverageDailySpendingRecommended.Value = dataBudgetEachCategory.BudgetLimit / 30

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	return dtoResponse, http.StatusOK, errInfo
}

func (s *BudgetUseCase) Travels(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		dtoResponse   dtos.Travel
		travelDetails []dtos.TravelDetails
	)

	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	dataTravel, err := s.repo.Travels(accountUUID)
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	if len(dataTravel) == 0 {
		response := struct {
			Message string `json:"message"`
		}{
			Message: "no data for travel budget. please set first",
		}
		return response, http.StatusNotFound, errInfo
	}

	if len(dataTravel) > 0 {
		for _, v := range dataTravel {

			var (
				dataXchangeCurrency entities.BudgetExistsExchangeValue
				err                 error
			)

			if v.CurrencyOrigin != "" {
				IDUUID, err := uuid.Parse(v.CurrencyOrigin)
				if err != nil {
					logrus.Error(err.Error())
				}

				dataXchangeCurrency, err = s.repo.GetXchangeCurrencyValue(IDUUID)
				if err != nil {
					logrus.Error(err.Error())
				}
			}

			budgetValue, err := strconv.ParseFloat(v.Budget, 64)
			if err != nil {
				logrus.Error(err.Error())
			}

			travelDetails = append(travelDetails, dtos.TravelDetails{
				ID:        v.ID,
				Departure: v.Departure,
				Arrival:   v.Arrival,
				ImagePath: os.Getenv("APP_HOST") + "/v1/" + v.ImagePath,
				Filename:  v.Filename,
				Budget: dtos.Amount{
					CurrencyCode: "IDR",
					Value:        int64(budgetValue),
				},
				TravelStartDate: v.TravelStartDate,
				TravelEndDate:   v.TravelEndDate,
				CurrencyOrigin:  dataXchangeCurrency.Code,
			})
		}

		dtoResponse.Details = travelDetails
	}

	return dtoResponse, http.StatusOK, errInfo
}

func (s *BudgetUseCase) UpdateTravelInfo(ctx *gin.Context, IDWallet string, request map[string]interface{}) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		amount float64
	)

	// check amount exist from payload
	value, exists := request["amount"]
	if exists {
		amount = value.(float64)

		if amount <= 0 {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", "amount must greater than 0")
			return struct{}{}, http.StatusBadRequest, errInfo
		}
	}

	// convert IDWwallet to UUID
	IDWalletUUID, err := uuid.Parse(IDWallet)
	if err != nil {
		logrus.Error(err.Error())
	}

	// update amount
	err = s.repo.UpdateAmountTravel(IDWalletUUID, request)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	resp := struct {
		Message string `json:"message"`
	}{
		Message: "change amount travel successfully",
	}

	return resp, http.StatusOK, errInfo
}