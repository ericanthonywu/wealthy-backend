package statistics

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/statistics/dtos"
	"github.com/semicolon-indonesia/wealthy-backend/utils/datecustoms"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/semicolon-indonesia/wealthy-backend/utils/personalaccounts"
	"github.com/sirupsen/logrus"
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
		Summary(ctx *gin.Context, month, year string) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		Statistic(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		TransactionPriority(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		Trend(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		Category(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		expenseWeekly(IDPersonal uuid.UUID, month, year string) (data []dtos.ExpenseWeekly)
		incomeWeekly(IDPersonal uuid.UUID, month, year string) (data []dtos.IncomeWeekly)
		investmentWeekly(IDPersonal uuid.UUID, month, year string) (data []dtos.InvestmentWeekly)
	}
)

func NewStatisticUseCase(repo IStatisticRepository) *StatisticUseCase {
	return &StatisticUseCase{repo: repo}
}

func (s *StatisticUseCase) Statistic(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		httpCode = http.StatusNotFound
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "not found")
		return response, httpCode, errInfo
	}

	//response = s.repo.Statistic(personalAccount.ID)
	return response, http.StatusOK, []errorsinfo.Errors{}
}

func (s *StatisticUseCase) Weekly(ctx *gin.Context, month, year string) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		stringBuilder strings.Builder
		dtoResponse   dtos.WeeklyData
	)

	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		httpCode = http.StatusNotFound
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
		DateRange: "01-04",
		Amount:    dataExpenseWeekly.DateRange0104,
	})

	expenseWeekly = append(expenseWeekly, dtos.ExpenseWeekly{
		DateRange: "05-11",
		Amount:    dataExpenseWeekly.DateRange0511,
	})

	expenseWeekly = append(expenseWeekly, dtos.ExpenseWeekly{
		DateRange: "12-18",
		Amount:    dataExpenseWeekly.DateRange1218,
	})

	expenseWeekly = append(expenseWeekly, dtos.ExpenseWeekly{
		DateRange: "19-25",
		Amount:    dataExpenseWeekly.DateRange1925,
	})

	expenseWeekly = append(expenseWeekly, dtos.ExpenseWeekly{
		DateRange: "26-30",
		Amount:    dataExpenseWeekly.DateRange2630,
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
		DateRange: "01-04",
		Amount:    dataIncomeWeekly.DateRange0104,
	})

	incomeWeekly = append(incomeWeekly, dtos.IncomeWeekly{
		DateRange: "05-11",
		Amount:    dataIncomeWeekly.DateRange0511,
	})

	incomeWeekly = append(incomeWeekly, dtos.IncomeWeekly{
		DateRange: "12-18",
		Amount:    dataIncomeWeekly.DateRange1218,
	})

	incomeWeekly = append(incomeWeekly, dtos.IncomeWeekly{
		DateRange: "19-25",
		Amount:    dataIncomeWeekly.DateRange1925,
	})

	incomeWeekly = append(incomeWeekly, dtos.IncomeWeekly{
		DateRange: "26-30",
		Amount:    dataIncomeWeekly.DateRange2630,
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
		DateRange: "01-04",
		Amount:    dataInvestmentWeekly.DateRange0104,
	})

	investmentWeekly = append(investmentWeekly, dtos.InvestmentWeekly{
		DateRange: "05-11",
		Amount:    dataInvestmentWeekly.DateRange0511,
	})

	investmentWeekly = append(investmentWeekly, dtos.InvestmentWeekly{
		DateRange: "12-18",
		Amount:    dataInvestmentWeekly.DateRange1218,
	})

	investmentWeekly = append(investmentWeekly, dtos.InvestmentWeekly{
		DateRange: "19-25",
		Amount:    dataInvestmentWeekly.DateRange1925,
	})

	investmentWeekly = append(investmentWeekly, dtos.InvestmentWeekly{
		DateRange: "26-30",
		Amount:    dataInvestmentWeekly.DateRange2630,
	})

	return investmentWeekly
}

func (s *StatisticUseCase) Summary(ctx *gin.Context, month, year string) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		monthPrevious        string
		yearPrevious         string
		dtoResponse          dtos.Summary
		expensePercentage    int
		investmentPercentage int
		netIncomePercentage  int
	)
	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		httpCode = http.StatusNotFound
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "not found")
		return response, httpCode, errInfo
	}

	dataCurrentSummary, err := s.repo.SummaryMonthly(personalAccount.ID, month, year)
	if err != nil {
		httpCode = http.StatusInternalServerError
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return response, httpCode, errInfo
	}

	monthINT, err := strconv.Atoi(month)
	if err != nil {
		logrus.Error(err.Error())
	}

	yearINT, err := strconv.Atoi(year)
	if err != nil {
		logrus.Error(err.Error())
	}

	if monthINT-1 < 0 {
		monthPrevious = "12"
		yearPrevious = strconv.Itoa(yearINT - 1)

	} else {
		monthPrevious = strconv.Itoa(monthINT - 1)
		yearPrevious = year
	}

	dataPreviousSummary, err := s.repo.SummaryMonthly(personalAccount.ID, monthPrevious, yearPrevious)

	if dataCurrentSummary.TotalExpense > 0 {
		expensePercentage = (dataPreviousSummary.TotalExpense / dataCurrentSummary.TotalExpense) * 100
	}

	if dataCurrentSummary.TotalInvest > 0 {
		investmentPercentage = (dataPreviousSummary.TotalInvest / dataCurrentSummary.TotalInvest) * 100
	}

	netIncome := dataCurrentSummary.TotalIncome - dataCurrentSummary.TotalExpense

	previousNetIncome := dataPreviousSummary.TotalIncome - dataPreviousSummary.TotalExpense
	currentNetIncome := dataCurrentSummary.TotalIncome - dataCurrentSummary.TotalExpense

	if currentNetIncome > 0 {
		netIncomePercentage = (previousNetIncome / currentNetIncome) * 100
	}

	dtoResponse.Expense.TotalAmount = dataCurrentSummary.TotalExpense
	dtoResponse.Expense.Percentage = fmt.Sprintf("%d", expensePercentage) + "%"

	dtoResponse.Investment.TotalAmount = dataCurrentSummary.TotalInvest
	dtoResponse.Investment.Percentage = fmt.Sprintf("%d", investmentPercentage) + "%"

	dtoResponse.NetIncome.TotalAmount = netIncome
	dtoResponse.NetIncome.Percentage = fmt.Sprintf("%d", netIncomePercentage) + "%"

	return dtoResponse, http.StatusOK, errInfo

}

func (s *StatisticUseCase) TransactionPriority(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		httpCode = http.StatusNotFound
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "not found")
		return response, httpCode, errInfo
	}

	response = s.repo.TransactionPriority(personalAccount.ID)
	return response, http.StatusOK, []errorsinfo.Errors{}
}

func (s *StatisticUseCase) Trend(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		httpCode = http.StatusNotFound
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "not found")
		return response, httpCode, errInfo
	}

	response = s.repo.Trend(personalAccount.ID)
	return response, http.StatusOK, []errorsinfo.Errors{}
}

func (s *StatisticUseCase) Category(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var categoryUUID uuid.UUID

	category := ctx.Query("category")
	usrEmail := ctx.MustGet("email").(string)

	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		httpCode = http.StatusNotFound
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "not found")
		return response, httpCode, errInfo
	}

	if category == "" {
		httpCode = http.StatusBadRequest
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "category required")
		return response, httpCode, errInfo
	}

	if category != "" {
		categoryUUID, _ = uuid.Parse(category)
	}

	response = s.repo.Category(personalAccount.ID, categoryUUID)
	return response, http.StatusOK, []errorsinfo.Errors{}
}