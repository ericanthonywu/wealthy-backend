package budgets

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/budgets/dtos"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/budgets/entities"
	"github.com/semicolon-indonesia/wealthy-backend/constants"
	"github.com/semicolon-indonesia/wealthy-backend/utils/datecustoms"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/semicolon-indonesia/wealthy-backend/utils/personalaccounts"
	"github.com/semicolon-indonesia/wealthy-backend/utils/utilities"
	"github.com/sirupsen/logrus"
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
		httpCode = http.StatusUnauthorized
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "token contains invalid information")
		return response, httpCode, errInfo
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
					totalBudgetAmount = totalBudgetAmount + v.BudgetLimit

					subCategoryInfo = append(subCategoryInfo, dtos.SubCategoryInfo{
						SubCategoryID:   v.SubCategoryID,
						SubCategoryName: v.SubCategoryName,
						BudgetLimit: dtos.Limit{
							CurrencyCode: "IDR",
							Value:        v.BudgetLimit,
						},
					})
				}
				// CHECK IF NOT FIRST TIME WITH VALUE NOT NIL
			} else if categoryIDPrevious != uuid.Nil {

				// CHECK IF PREVIOUS IS SAME AS CURRENT
				if categoryIDPrevious == v.CategoryID {
					totalBudgetAmount = totalBudgetAmount + v.BudgetLimit

					subCategoryInfo = append(subCategoryInfo, dtos.SubCategoryInfo{
						SubCategoryID:   v.SubCategoryID,
						SubCategoryName: v.SubCategoryName,
						BudgetLimit: dtos.Limit{
							CurrencyCode: "IDR",
							Value:        v.BudgetLimit,
						},
					})

					// OTHERWISE DIFFERENT
				} else {

					// IF SUB CATEGORY NOT EMPTY
					if len(subCategoryInfo) > 0 {
						budgetDetail = append(budgetDetail, dtos.AllBudgetDetail{
							CategoryID:   categoryIDPrevious,
							CategoryName: categoryNamePrevious,
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
						totalBudgetAmount = totalBudgetAmount + v.BudgetLimit

						subCategoryInfo = append(subCategoryInfo, dtos.SubCategoryInfo{
							SubCategoryID:   v.SubCategoryID,
							SubCategoryName: v.SubCategoryName,
							BudgetLimit: dtos.Limit{
								CurrencyCode: "IDR",
								Value:        v.BudgetLimit,
							},
						})

					}

					if k == lengthOfData {
						budgetDetail = append(budgetDetail, dtos.AllBudgetDetail{
							CategoryID:   categoryIDPrevious,
							CategoryName: categoryNamePrevious,
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

	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		httpCode = http.StatusUnauthorized
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "token contains invalid information")
		return response, httpCode, errInfo
	}

	personalBudgetData, err := s.repo.PersonalBudget(personalAccount.ID, month, year)
	if err != nil {
		logrus.Error(err.Error())
	}

	personalTransactionData, err := s.repo.PersonalTransaction(personalAccount.ID, month, year)
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
				CategoryID:   v.ID,
				CategoryName: v.Category,
				BudgetLimit: dtos.Limit{
					CurrencyCode: "IDR",
					Value:        v.BudgetLimit,
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
		httpCode = http.StatusUnauthorized
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "token contains invalid information")
		return response, httpCode, errInfo
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

		model.Departure = dtoRequest.Departure
		model.Arrival = dtoRequest.Arrival
		model.Filename = filename
		model.ImagePath = imagePath
		model.TravelStartDate = dtoRequest.TravelStartDate
		model.TravelEndDate = dtoRequest.TravelEndDate
		model.IDMasterTransactionType = dtoRequest.IDMasterTransactionTypes
	}

	model.Amount = dtoRequest.Amount
	model.IDPersonalAccount = personalAccount.ID
	model.IDCategory = dtoRequest.IDCategory
	model.IDSubCategory = dtoRequest.IDSubCategory
	model.ID = uuid.New()

	err := s.repo.Limit(&model)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "insert new budget issue")
		return response, http.StatusInternalServerError, errInfo
	}

	dtoResponse.ID = model.ID
	dtoResponse.Status = true
	return dtoResponse, httpCode, []errorsinfo.Errors{}
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

	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		httpCode = http.StatusUnauthorized
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "token contains invalid information")
		return response, httpCode, errInfo
	}

	dataTravel, err := s.repo.Travels(personalAccount.ID)
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return entities.TrendsWeekly{}, http.StatusInternalServerError, errInfo
	}

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	if len(dataTravel) == 0 {
		response := struct {
			Message string `json:"message"`
		}{
			Message: "no travel data using the token",
		}
		return response, http.StatusNotFound, errInfo
	}

	if len(dataTravel) > 0 {
		for _, v := range dataTravel {

			budgetValue, _ := strconv.ParseInt(v.Budget, 10, 64)

			travelDetails = append(travelDetails, dtos.TravelDetails{
				ID:        v.ID,
				Departure: v.Departure,
				Arrival:   v.Arrival,
				ImagePath: os.Getenv("APP_HOST") + "/v1/" + v.ImagePath,
				Filename:  v.Filename,
				Budget: dtos.Amount{
					CurrencyCode: "IDR",
					Value:        budgetValue,
				},
				TravelStartDate: v.TravelStartDate,
				TravelEndDate:   v.TravelEndDate,
			})
		}

		dtoResponse.Details = travelDetails
	}

	return dtoResponse, http.StatusOK, errInfo
}