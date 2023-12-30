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

	month := ctx.Query("month")
	year := ctx.Query("year")

	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", constants.TokenInvalidInformation)
		return response, http.StatusUnauthorized, errInfo
	}

	dataSubCategoryBudget, err := s.repo.SubCategoryBudget(personalAccount.ID, month, year)
	if err != nil {
		logrus.Error(err.Error())
	}

	categoryIDPrevious := uuid.Nil
	categoryNamePrevious := ""
	lengthOfData := len(dataSubCategoryBudget) - 1
	totalBudgetAmount := 0

	if len(dataSubCategoryBudget) > 0 {
		for k, v := range dataSubCategoryBudget {
			// CHECK IF FIRST TIME WITH VALUE NIL
			if categoryIDPrevious == uuid.Nil {

				// MOVE VALUE INTO PREVIOUS
				categoryIDPrevious = v.CategoryID
				categoryNamePrevious = v.CategoryName

				// CHECK IF SUB CATEGORY NOT NIL
				if v.SubCategoryID != uuid.Nil {
					totalBudgetAmount = totalBudgetAmount + int(v.BudgetLimit)

					subCategoryInfo = append(subCategoryInfo, dtos.SubCategoryInfo{
						SubCategoryID:   v.SubCategoryID,
						SubCategoryName: v.SubCategoryName,
						SubCategoryIcon: v.SubCategoryIcon,
						BudgetLimit: dtos.Limit{
							CurrencyCode: "IDR",
							Value:        int(v.BudgetLimit),
						},
					})
				}
				// CHECK IF NOT FIRST TIME WITH VALUE NOT NIL
			} else if categoryIDPrevious != uuid.Nil {

				// CHECK IF PREVIOUS IS SAME AS CURRENT
				if categoryIDPrevious == v.CategoryID {
					totalBudgetAmount = totalBudgetAmount + int(v.BudgetLimit)

					subCategoryInfo = append(subCategoryInfo, dtos.SubCategoryInfo{
						SubCategoryID:   v.SubCategoryID,
						SubCategoryName: v.SubCategoryName,
						SubCategoryIcon: v.SubCategoryIcon,
						BudgetLimit: dtos.Limit{
							CurrencyCode: "IDR",
							Value:        int(v.BudgetLimit),
						},
					})

					// OTHERWISE DIFFERENT
				} else {

					// IF SUB CATEGORY NOT EMPTY
					if len(subCategoryInfo) > 0 {
						budgetDetail = append(budgetDetail, dtos.AllBudgetDetail{
							CategoryID:   categoryIDPrevious,
							CategoryName: categoryNamePrevious,
							CategoryIcon: v.ImagePath,
							SubCategory:  subCategoryInfo,
							BudgetInfo: dtos.Limit{
								CurrencyCode: "IDR",
								Value:        totalBudgetAmount,
							},
						})
						// OTHERWISE EMPTY
					} else {
						budgetDetail = append(budgetDetail, dtos.AllBudgetDetail{
							CategoryID:   categoryIDPrevious,
							CategoryName: categoryNamePrevious,
							CategoryIcon: v.ImagePath,
							SubCategory:  []dtos.SubCategoryInfo{},
							BudgetInfo: dtos.Limit{
								CurrencyCode: "IDR",
								Value:        0,
							},
						})
					}

					// RESET SUB CATEGORY
					subCategoryInfo = []dtos.SubCategoryInfo{}
					totalBudgetAmount = 0

					// RENEW VALUE ID AND NAME
					categoryIDPrevious = v.CategoryID
					categoryNamePrevious = v.CategoryName

					// IF SUB CATEGORY NOT EMPTY
					if v.SubCategoryID != uuid.Nil {
						totalBudgetAmount = totalBudgetAmount + int(v.BudgetLimit)

						subCategoryInfo = append(subCategoryInfo, dtos.SubCategoryInfo{
							SubCategoryID:   v.SubCategoryID,
							SubCategoryName: v.SubCategoryName,
							SubCategoryIcon: v.SubCategoryIcon,
							BudgetLimit: dtos.Limit{
								CurrencyCode: "IDR",
								Value:        int(v.BudgetLimit),
							},
						})

					}

					if k == lengthOfData {
						budgetDetail = append(budgetDetail, dtos.AllBudgetDetail{
							CategoryID:   categoryIDPrevious,
							CategoryName: categoryNamePrevious,
							CategoryIcon: v.ImagePath,
							SubCategory:  []dtos.SubCategoryInfo{},
							BudgetInfo: dtos.Limit{
								CurrencyCode: "IDR",
								Value:        totalBudgetAmount,
							},
						})
					}
				}
			}
		}
	} else {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return dtoResponse, http.StatusNotFound, errInfo
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

	personalDataForSpending := make(map[uuid.UUID]int)
	personalDataForCount := make(map[uuid.UUID]int)

	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	personalBudgetData, err := s.repo.PersonalBudget(accountUUID, month, year)
	if err != nil {
		logrus.Error(err.Error())
	}

	personalTransactionData, err := s.repo.PersonalTransaction(accountUUID, month, year)
	if err != nil {
		logrus.Error(err.Error())
	}

	if len(personalBudgetData) > 0 {

		if len(personalTransactionData) > 0 {
			for _, ptd := range personalTransactionData {
				personalDataForSpending[ptd.ID] = ptd.Amount
				personalDataForCount[ptd.ID] = ptd.Count
			}
		}

		for _, v := range personalBudgetData {
			count := 0
			spendingTrx := 0

			if value, isFound := personalDataForSpending[v.ID]; isFound {
				spendingTrx = value
			}

			if value, isFound := personalDataForCount[v.ID]; isFound {
				count = value
			}

			dataDetails = append(dataDetails, dtos.OverviewDetail{
				CategoryID:      v.ID,
				CategoryName:    v.Category,
				TransactionIcon: v.ImagePath,
				BudgetLimit: dtos.Limit{
					CurrencyCode: "IDR",
					Value:        int(v.BudgetLimit),
				},
				TransactionSpending: dtos.Transaction{
					CurrencyCode: "IDR",
					Value:        spendingTrx,
				},
				NumberOfCategories: count,
			})
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
		return dtoResponse, http.StatusOK, errInfo
	}

	errInfo = errorsinfo.ErrorWrapper(errInfo, "", "data not found")
	return response, http.StatusNotFound, errInfo
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
		model       entities.BudgetSetEntities
		dtoResponse dtos.BudgetSetResponse
		filename    string
		targetPath  string
		imagePath   string
	)

	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", constants.TokenInvalidInformation)
		return struct{}{}, http.StatusUnauthorized, errInfo
	}

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

		// validate id master exchange currency
		dataCurrency, err := s.repo.GetXchangeCurrency(IDMasterExchangeCurrencyUUID)
		if err != nil {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
			return struct{}{}, http.StatusInternalServerError, errInfo
		}

		if !dataCurrency.Exists {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id master exchange currency unknown")
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
	}

	if purpose != constants.Travel {
		model.IDCategory = dtoRequest.IDCategory
		model.IDSubCategory = dtoRequest.IDSubCategory
	}

	model.Amount = dtoRequest.Amount
	model.IDPersonalAccount = personalAccount.ID
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
	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		httpCode = http.StatusUnauthorized
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "token contains invalid information")
		return response, httpCode, errInfo
	}

	dataTrends, err := s.repo.Trends(personalAccount.ID, IDCategory, month, year)
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

	dataBudgetEachCategory, err := s.repo.BudgetEachCategory(personalAccount.ID, IDCategory, month, year)
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