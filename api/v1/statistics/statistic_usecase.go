package statistics

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/wealthy-app/wealthy-backend/api/v1/statistics/dtos"
	"github.com/wealthy-app/wealthy-backend/api/v1/statistics/entities"
	"github.com/wealthy-app/wealthy-backend/constants"
	"github.com/wealthy-app/wealthy-backend/utils/datecustoms"
	"github.com/wealthy-app/wealthy-backend/utils/errorsinfo"
	"github.com/wealthy-app/wealthy-backend/utils/personalaccounts"
	"net/http"
	"strconv"
	"strings"
)

type (
	StatisticUseCase struct {
		repo IStatisticRepository
	}

	IStatisticUseCase interface {
		Weekly(ctx *gin.Context, month, year string) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		Summary(ctx *gin.Context, month, year, email string) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		Priority(ctx *gin.Context, month, year string) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		Trend(ctx *gin.Context, month, year string) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		expenseWeekly(IDPersonal uuid.UUID, month, year string) (data []dtos.ExpenseWeekly)
		subExpenseWeekly(IDPersonal uuid.UUID, IDCategory uuid.UUID, month, year string) (categoryName string, data []dtos.WeeklySubExpenseDetail, err error)
		incomeWeekly(IDPersonal uuid.UUID, month, year string) (data []dtos.IncomeWeekly)
		investmentWeekly(IDPersonal uuid.UUID, month, year string) (data []dtos.InvestmentWeekly)
		ExpenseDetail(ctx *gin.Context, month, year, email string) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		SubExpenseDetail(ctx *gin.Context, month, year string, IDCategory uuid.UUID) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		isDataPriorityNotEmpty(data entities.StatisticPriority) bool
		AnalyticsTrend(ctx *gin.Context, period string, typeName string) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		TopThreeInvestment(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
	}
)

func NewStatisticUseCase(repo IStatisticRepository) *StatisticUseCase {
	return &StatisticUseCase{repo: repo}
}

func (s *StatisticUseCase) Weekly(ctx *gin.Context, month, year string) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		stringBuilder strings.Builder
		dtoResponse   dtos.WeeklyData
	)

	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		httpCode = http.StatusUnauthorized
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "token contains invalid information")
		return response, httpCode, errInfo
	}

	dataExpenseWeekly := s.expenseWeekly(personalAccount.ID, month, year)
	dataIncomeWeekly := s.incomeWeekly(personalAccount.ID, month, year)
	dataInvestmentWeekly := s.investmentWeekly(personalAccount.ID, month, year)

	monthINT, err := strconv.Atoi(month)
	if err != nil {
		logrus.Error(err.Error())
	}

	stringBuilder.WriteString(datecustoms.IntToMonthName(monthINT))
	stringBuilder.WriteString(" ")
	stringBuilder.WriteString(year)

	dtoResponse.Period = stringBuilder.String()
	dtoResponse.Expense = dataExpenseWeekly
	dtoResponse.Income = dataIncomeWeekly
	dtoResponse.Investment = dataInvestmentWeekly

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	return dtoResponse, http.StatusOK, errInfo
}

func (s *StatisticUseCase) expenseWeekly(IDPersonal uuid.UUID, month, year string) (data []dtos.ExpenseWeekly) {
	var expenseWeekly []dtos.ExpenseWeekly
	dataExpenseWeekly, err := s.repo.expenseWeekly(IDPersonal, month, year)

	if err != nil {
		return []dtos.ExpenseWeekly{}
	}

	expenseWeekly = append(expenseWeekly, dtos.ExpenseWeekly{
		StartDate: year + "-" + month + "-01",
		EndDate:   year + "-" + month + "-04",
		Amount: dtos.ExpenseTransaction{
			CurrencyCode: "IDR",
			Value:        dataExpenseWeekly.DateRange0104,
		},
	})

	expenseWeekly = append(expenseWeekly, dtos.ExpenseWeekly{
		StartDate: year + "-" + month + "-05",
		EndDate:   year + "-" + month + "-11",
		Amount: dtos.ExpenseTransaction{
			CurrencyCode: "IDR",
			Value:        dataExpenseWeekly.DateRange0511,
		},
	})

	expenseWeekly = append(expenseWeekly, dtos.ExpenseWeekly{
		StartDate: year + "-" + month + "-12",
		EndDate:   year + "-" + month + "-18",
		Amount: dtos.ExpenseTransaction{
			CurrencyCode: "IDR",
			Value:        dataExpenseWeekly.DateRange1218,
		},
	})

	expenseWeekly = append(expenseWeekly, dtos.ExpenseWeekly{
		StartDate: year + "-" + month + "-19",
		EndDate:   year + "-" + month + "-25",
		Amount: dtos.ExpenseTransaction{
			CurrencyCode: "IDR",
			Value:        dataExpenseWeekly.DateRange1925,
		},
	})

	expenseWeekly = append(expenseWeekly, dtos.ExpenseWeekly{
		StartDate: year + "-" + month + "-26",
		EndDate:   year + "-" + month + "-30",
		Amount: dtos.ExpenseTransaction{
			CurrencyCode: "IDR",
			Value:        dataExpenseWeekly.DateRange2630,
		},
	})

	return expenseWeekly
}

func (s *StatisticUseCase) incomeWeekly(IDPersonal uuid.UUID, month, year string) (data []dtos.IncomeWeekly) {
	var incomeWeekly []dtos.IncomeWeekly
	dataIncomeWeekly, err := s.repo.incomeWeekly(IDPersonal, month, year)

	if err != nil {
		return []dtos.IncomeWeekly{}
	}

	incomeWeekly = append(incomeWeekly, dtos.IncomeWeekly{
		StartDate: year + "-" + month + "-01",
		EndDate:   year + "-" + month + "-04",
		Amount: dtos.IncomeTransaction{
			CurrencyCode: "IDR",
			Value:        dataIncomeWeekly.DateRange0104,
		},
	})

	incomeWeekly = append(incomeWeekly, dtos.IncomeWeekly{
		StartDate: year + "-" + month + "-05",
		EndDate:   year + "-" + month + "-11",
		Amount: dtos.IncomeTransaction{
			CurrencyCode: "IDR",
			Value:        dataIncomeWeekly.DateRange0511,
		},
	})

	incomeWeekly = append(incomeWeekly, dtos.IncomeWeekly{
		StartDate: year + "-" + month + "-12",
		EndDate:   year + "-" + month + "-18",
		Amount: dtos.IncomeTransaction{
			CurrencyCode: "IDR",
			Value:        dataIncomeWeekly.DateRange1218,
		},
	})

	incomeWeekly = append(incomeWeekly, dtos.IncomeWeekly{
		StartDate: year + "-" + month + "-19",
		EndDate:   year + "-" + month + "-25",
		Amount: dtos.IncomeTransaction{
			CurrencyCode: "IDR",
			Value:        dataIncomeWeekly.DateRange1925,
		},
	})

	incomeWeekly = append(incomeWeekly, dtos.IncomeWeekly{
		StartDate: year + "" + month + "-26",
		EndDate:   year + "" + month + "-30",
		Amount: dtos.IncomeTransaction{
			CurrencyCode: "IDR",
			Value:        dataIncomeWeekly.DateRange2630,
		},
	})

	return incomeWeekly
}

func (s *StatisticUseCase) investmentWeekly(IDPersonal uuid.UUID, month, year string) (data []dtos.InvestmentWeekly) {
	var investmentWeekly []dtos.InvestmentWeekly
	dataInvestmentWeekly, err := s.repo.investmentWeekly(IDPersonal, month, year)

	if err != nil {
		return []dtos.InvestmentWeekly{}
	}

	investmentWeekly = append(investmentWeekly, dtos.InvestmentWeekly{
		StartDate: year + "-" + month + "-01",
		EndDate:   year + "-" + month + "-04",
		Amount: dtos.InvestTransaction{
			CurrencyCode: "IDR",
			Value:        dataInvestmentWeekly.DateRange0104,
		},
	})

	investmentWeekly = append(investmentWeekly, dtos.InvestmentWeekly{
		StartDate: year + "-" + month + "-05",
		EndDate:   year + "-" + month + "-11",
		Amount: dtos.InvestTransaction{
			CurrencyCode: "IDR",
			Value:        dataInvestmentWeekly.DateRange0511,
		},
	})

	investmentWeekly = append(investmentWeekly, dtos.InvestmentWeekly{
		StartDate: year + "-" + month + "-12",
		EndDate:   year + "-" + month + "-18",
		Amount: dtos.InvestTransaction{
			CurrencyCode: "IDR",
			Value:        dataInvestmentWeekly.DateRange1218,
		},
	})

	investmentWeekly = append(investmentWeekly, dtos.InvestmentWeekly{
		StartDate: year + "-" + month + "-19",
		EndDate:   year + "-" + month + "-25",
		Amount: dtos.InvestTransaction{
			CurrencyCode: "IDR",
			Value:        dataInvestmentWeekly.DateRange1925,
		},
	})

	investmentWeekly = append(investmentWeekly, dtos.InvestmentWeekly{
		StartDate: year + "-" + month + "-26",
		EndDate:   year + "-" + month + "-30",
		Amount: dtos.InvestTransaction{
			CurrencyCode: "IDR",
			Value:        dataInvestmentWeekly.DateRange2630,
		},
	})

	return investmentWeekly
}

func (s *StatisticUseCase) Summary(ctx *gin.Context, month, year, email string) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		monthPrevious        string
		yearPrevious         string
		dtoResponse          dtos.Summary
		expensePercentage    int
		investmentPercentage int
		netIncomePercentage  int
		stringBuilder        strings.Builder
		dataCurrentSummary   entities.StatisticSummaryMonthly
		dataPreviousSummary  entities.StatisticSummaryMonthly
		err                  error
	)

	monthINT, err := strconv.Atoi(month)
	if err != nil {
		logrus.Error(err.Error())
	}

	yearINT, err := strconv.Atoi(year)
	if err != nil {
		logrus.Error(err.Error())
	}

	// if email is empty
	if email == "" {
		usrEmail := ctx.MustGet("email").(string)
		personalAccount := personalaccounts.Informations(ctx, usrEmail)

		if personalAccount.ID == uuid.Nil {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", "token contains invalid information")
			return response, http.StatusNotFound, errInfo
		}

		dataCurrentSummary, err = s.repo.SummaryMonthly(personalAccount.ID, month, year)
		if err != nil {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
			return response, http.StatusInternalServerError, errInfo
		}

		if monthINT-1 < 0 {
			monthPrevious = "12"
			yearPrevious = strconv.Itoa(yearINT - 1)

		} else {
			monthPrevious = strconv.Itoa(monthINT - 1)
			yearPrevious = year
		}

		dataPreviousSummary, err = s.repo.SummaryMonthly(personalAccount.ID, monthPrevious, yearPrevious)
		if err != nil {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
			return response, http.StatusInternalServerError, errInfo
		}
	}

	// if email is not empty
	if email != "" {
		dataProfile, err := s.repo.GetProfileByEmail(email)
		if err != nil {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
			return response, http.StatusInternalServerError, errInfo
		}

		dataCurrentSummary, err = s.repo.SummaryMonthly(dataProfile.ID, month, year)
		if err != nil {
			httpCode = http.StatusInternalServerError
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
			return response, httpCode, errInfo
		}

		if monthINT-1 < 0 {
			monthPrevious = "12"
			yearPrevious = strconv.Itoa(yearINT - 1)
		} else {
			monthPrevious = strconv.Itoa(monthINT - 1)
			yearPrevious = year
		}

		dataPreviousSummary, err = s.repo.SummaryMonthly(dataProfile.ID, monthPrevious, yearPrevious)
		if err != nil {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
			return response, http.StatusInternalServerError, errInfo
		}
	}

	if dataCurrentSummary.TotalExpense > 0 {
		expensePercentage = (dataPreviousSummary.TotalExpense / dataCurrentSummary.TotalExpense) * 100
	}

	if dataCurrentSummary.TotalInvest > 0 {
		investmentPercentage = (dataPreviousSummary.TotalInvest / dataCurrentSummary.TotalInvest) * 100
	}

	// calculate net income, previous and current
	netIncome := dataCurrentSummary.TotalIncome - dataCurrentSummary.TotalExpense
	previousNetIncome := dataPreviousSummary.TotalIncome - dataPreviousSummary.TotalExpense
	currentNetIncome := dataCurrentSummary.TotalIncome - dataCurrentSummary.TotalExpense

	if currentNetIncome > 0 {
		netIncomePercentage = (previousNetIncome / currentNetIncome) * 100
	}

	dtoResponse.Expense.TotalAmount = dtos.SummaryTransaction{
		CurrencyCode: "IDR",
		Value:        dataCurrentSummary.TotalExpense,
	}
	dtoResponse.Expense.Percentage = fmt.Sprintf("%d", expensePercentage) + "%"

	dtoResponse.Investment.TotalAmount = dtos.SummaryTransaction{
		CurrencyCode: "IDR",
		Value:        dataCurrentSummary.TotalInvest,
	}
	dtoResponse.Investment.Percentage = fmt.Sprintf("%d", investmentPercentage) + "%"

	dtoResponse.NetIncome.TotalAmount = dtos.SummaryTransaction{
		CurrencyCode: "IDR",
		Value:        netIncome,
	}
	dtoResponse.NetIncome.Percentage = fmt.Sprintf("%d", netIncomePercentage) + "%"

	stringBuilder.WriteString(datecustoms.IntToMonthName(monthINT))
	stringBuilder.WriteString(" ")
	stringBuilder.WriteString(year)

	dtoResponse.Period = stringBuilder.String()

	// err info empty
	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	return dtoResponse, http.StatusOK, errInfo
}

func (s *StatisticUseCase) Priority(ctx *gin.Context, month, year string) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		stringBuilder  strings.Builder
		dtoResponse    dtos.Priority
		monthINT       int
		percentageMust string
		percentageWant string
		percentageNeed string
	)

	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		httpCode = http.StatusNotFound
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "token contains invalid information")
		return response, httpCode, errInfo
	}

	dataPriority := s.repo.Priority(personalAccount.ID, month, year)

	monthINT, err := strconv.Atoi(month)
	if err != nil {
		logrus.Error(err.Error())
	}

	if s.isDataPriorityNotEmpty(dataPriority) {
		percentageMust = fmt.Sprintf("%.2f", (float64(dataPriority.PriorityMust)/float64(dataPriority.TotalTransaction))*100) + "%"
		percentageWant = fmt.Sprintf("%.2f", (float64(dataPriority.PriorityWant)/float64(dataPriority.TotalTransaction))*100) + "%"
		percentageNeed = fmt.Sprintf("%.2f", (float64(dataPriority.PriorityNeed)/float64(dataPriority.TotalTransaction))*100) + "%"
	} else {
		percentageMust = "0%"
		percentageWant = percentageMust
		percentageNeed = percentageWant
	}

	dtoResponse.Info = append(dtoResponse.Info, dtos.PriorityInfo{
		Type:       strings.ToUpper("must"),
		Percentage: percentageMust,
	})

	dtoResponse.Info = append(dtoResponse.Info, dtos.PriorityInfo{
		Type:       strings.ToUpper("want"),
		Percentage: percentageWant,
	})

	dtoResponse.Info = append(dtoResponse.Info, dtos.PriorityInfo{
		Type:       strings.ToUpper("need"),
		Percentage: percentageNeed,
	})

	stringBuilder.WriteString(datecustoms.IntToMonthName(monthINT))
	stringBuilder.WriteString(" ")
	stringBuilder.WriteString(year)

	dtoResponse.Period = stringBuilder.String()

	// if err info is empty
	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	return dtoResponse, http.StatusOK, errInfo
}

func (s *StatisticUseCase) Trend(ctx *gin.Context, month, year string) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		stringBuilder strings.Builder
		dtoResponse   dtos.TrendsData
		totalWeekly   int
		totalDaily    int
		looping       int
	)
	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", constants.TokenInvalidInformation)
		return struct{}{}, http.StatusUnauthorized, errInfo
	}

	monthINT, err := strconv.Atoi(month)
	if err != nil {
		logrus.Error(err.Error())
	}

	stringBuilder.WriteString(datecustoms.IntToMonthName(monthINT))
	stringBuilder.WriteString(" ")
	stringBuilder.WriteString(year)

	dataExpenseWeekly := s.expenseWeekly(personalAccount.ID, month, year)

	isData := dataExpenseWeekly[0].Amount.Value
	if isData == 0 {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: "no data for trends statistic",
		}
		return resp, http.StatusBadRequest, []errorsinfo.Errors{}
	}

	looping = 1
	for _, v := range dataExpenseWeekly {
		totalWeekly = totalWeekly + v.Amount.Value
		looping++
	}

	looping = 1
	for _, v := range dataExpenseWeekly {
		if looping == 1 {
			totalDaily += v.Amount.Value / 4
		}

		if looping == 2 || looping == 3 || looping == 4 {
			totalDaily += v.Amount.Value / 7
		}

		if looping == 5 {
			totalDaily += v.Amount.Value / 5
		}
		looping++
	}

	dtoResponse.Period = stringBuilder.String()
	dtoResponse.Expense = dataExpenseWeekly
	dtoResponse.AverageWeekly = totalWeekly / looping
	dtoResponse.AverageDaily = totalDaily / 30

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	return dtoResponse, http.StatusOK, errInfo
}

func (s *StatisticUseCase) ExpenseDetail(ctx *gin.Context, month, year, email string) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		dtoResponse   dtos.ExpenseDetail
		stringBuilder strings.Builder
		totalExpense  int64
		IDUser        uuid.UUID
	)

	// if email empty
	if email == "" {
		usrEmail := ctx.MustGet("email").(string)
		personalAccount := personalaccounts.Informations(ctx, usrEmail)

		if personalAccount.ID == uuid.Nil {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", constants.TokenInvalidInformation)
			return struct{}{}, http.StatusUnauthorized, errInfo
		}
		IDUser = personalAccount.ID
	}

	// if email not empty
	if email != "" {
		dataProfile, err := s.repo.GetProfileByEmail(email)
		if err != nil {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
			return struct{}{}, http.StatusInternalServerError, errInfo
		}
		IDUser = dataProfile.ID
	}

	dataExpenseDetail, err := s.repo.ExpenseDetail(IDUser, month, year)
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	if len(dataExpenseDetail) == 0 {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: "no data for expense detail statistic",
		}
		return resp, http.StatusNotFound, []errorsinfo.Errors{}
	}

	for _, v := range dataExpenseDetail {
		dtoResponse.Expense = append(dtoResponse.Expense, dtos.ExpDetail{
			ID:           v.ID,
			Category:     v.Category,
			CategoryIcon: v.CategoryIcon,
			Amount: dtos.ExpDetailTransaction{
				CurrencyCode: "IDR",
				Value:        int64(v.Amount),
			},
		})

		totalExpense += int64(v.Amount)
	}

	monthINT, err := strconv.Atoi(month)
	if err != nil {
		logrus.Error(err.Error())
	}

	stringBuilder.WriteString(datecustoms.IntToMonthName(monthINT))
	stringBuilder.WriteString(" ")
	stringBuilder.WriteString(year)

	dtoResponse.Period = stringBuilder.String()
	dtoResponse.TotalExpense = totalExpense

	// errInfo empty
	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	return dtoResponse, http.StatusOK, errInfo
}

func (s *StatisticUseCase) SubExpenseDetail(ctx *gin.Context, month, year string, IDCategory uuid.UUID) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		stringBuilder strings.Builder
		dtoResponse   dtos.WeeklySubExpense
	)

	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", constants.TokenInvalidInformation)
		return struct{}{}, http.StatusUnauthorized, errInfo
	}

	categoryName, dataExpenseWeekly, err := s.subExpenseWeekly(personalAccount.ID, IDCategory, month, year)
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// if no data transaction
	if categoryName == "" {
		resp := struct {
			Message string `json:""`
		}{
			Message: "no data for sub expense detail statistic",
		}
		return resp, http.StatusNotFound, []errorsinfo.Errors{}
	}

	monthINT, err := strconv.Atoi(month)
	if err != nil {
		logrus.Error(err.Error())
	}

	stringBuilder.WriteString(datecustoms.IntToMonthName(monthINT))
	stringBuilder.WriteString(" ")
	stringBuilder.WriteString(year)

	dtoResponse.CategoryName = categoryName
	dtoResponse.Period = stringBuilder.String()
	dtoResponse.CategoryID = IDCategory.String()
	dtoResponse.Expense = dataExpenseWeekly

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	return dtoResponse, http.StatusOK, errInfo
}

func (s *StatisticUseCase) subExpenseWeekly(IDPersonal uuid.UUID, IDCategory uuid.UUID, month, year string) (categoryName string, data []dtos.WeeklySubExpenseDetail, err error) {
	dataSubExpenseWeekly, errs := s.repo.SubExpenseDetail(IDPersonal, IDCategory, month, year)
	if errs != nil {
		return "", []dtos.WeeklySubExpenseDetail{}, errs
	}

	data = append(data, dtos.WeeklySubExpenseDetail{
		StartDate: year + "-" + month + "-10",
		EndDate:   year + "-" + month + "-04",
		Amount: dtos.WeeklySubExpenseDetailTransaction{
			CurrencyCode: "IDR",
			Value:        dataSubExpenseWeekly.DateRange0104,
		},
	})

	data = append(data, dtos.WeeklySubExpenseDetail{
		StartDate: year + "-" + month + "-05",
		EndDate:   year + "-" + month + "-11",
		Amount: dtos.WeeklySubExpenseDetailTransaction{
			CurrencyCode: "IDR",
			Value:        dataSubExpenseWeekly.DateRange0511,
		},
	})

	data = append(data, dtos.WeeklySubExpenseDetail{
		StartDate: year + "-" + month + "-12",
		EndDate:   year + "-" + month + "-18",
		Amount: dtos.WeeklySubExpenseDetailTransaction{
			CurrencyCode: "IDR",
			Value:        dataSubExpenseWeekly.DateRange1218,
		},
	})

	data = append(data, dtos.WeeklySubExpenseDetail{
		StartDate: year + "-" + month + "-19",
		EndDate:   year + "-" + month + "-25",
		Amount: dtos.WeeklySubExpenseDetailTransaction{
			CurrencyCode: "IDR",
			Value:        dataSubExpenseWeekly.DateRange1925,
		},
	})

	data = append(data, dtos.WeeklySubExpenseDetail{
		StartDate: year + "" + month + "-26",
		EndDate:   year + "-" + month + "-30",
		Amount: dtos.WeeklySubExpenseDetailTransaction{
			CurrencyCode: "IDR",
			Value:        dataSubExpenseWeekly.DateRange2630,
		},
	})

	return dataSubExpenseWeekly.CategoryName, data, nil
}

func (s *StatisticUseCase) isDataPriorityNotEmpty(data entities.StatisticPriority) bool {
	return data != entities.StatisticPriority{}
}

func (s *StatisticUseCase) AnalyticsTrend(ctx *gin.Context, period string, typeName string) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	typeName = strings.ToUpper(typeName)
	accountID := ctx.MustGet("accountID").(uuid.UUID)

	dataRepo := s.repo.AnalyticsTrend(accountID, typeName, period)
	if len(dataRepo) == 0 {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: "no data for trends analytics statistic",
		}
		return resp, http.StatusNotFound, []errorsinfo.Errors{}
	}

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	return dataRepo, http.StatusOK, errInfo
}

func (s *StatisticUseCase) TopThreeInvestment(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {

	accountUUID := ctx.MustGet("accountID").(uuid.UUID)
	dataTopThreeInvestment, err := s.repo.TopThreeInvestment(accountUUID)
	if err != nil {
		logrus.Error()
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	if len(dataTopThreeInvestment) == 0 {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: "no data for investment",
		}
		return resp, http.StatusNotFound, []errorsinfo.Errors{}
	}

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	return dataTopThreeInvestment, http.StatusOK, errInfo
}