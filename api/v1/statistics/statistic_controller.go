package statistics

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/semicolon-indonesia/wealthy-backend/utils/response"
	"github.com/sirupsen/logrus"
	"net/http"
)

type (
	StatisticController struct {
		useCase IStatisticUseCase
	}

	IStatisticController interface {
		Weekly(ctx *gin.Context)
		Summary(ctx *gin.Context)
		TransactionPriority(ctx *gin.Context)
		Trend(ctx *gin.Context)
		AnalyticsTrend(ctx *gin.Context)
		ExpenseDetail(ctx *gin.Context)
		SubExpenseDetail(ctx *gin.Context)
	}
)

func NewStatisticController(useCase IStatisticUseCase) *StatisticController {
	return &StatisticController{useCase: useCase}
}

func (c *StatisticController) Weekly(ctx *gin.Context) {
	var (
		errInfo    []errorsinfo.Errors
		data       interface{}
		statusCode int
	)
	month := ctx.Query("month")
	year := ctx.Query("year")

	if month == "" || year == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "month or year required in query url")
		response.SendBack(ctx, nil, errInfo, http.StatusBadRequest)
		return
	}

	data, statusCode, errInfo = c.useCase.Weekly(ctx, month, year)

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	response.SendBack(ctx, data, errInfo, statusCode)
	return
}

func (c *StatisticController) Summary(ctx *gin.Context) {
	var (
		errInfo    []errorsinfo.Errors
		data       interface{}
		statusCode int
	)

	// get query parameter
	month := ctx.Query("month")
	year := ctx.Query("year")
	email := ctx.Query("email")

	// validate
	if month == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "month required in query url")
	}

	if year == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "year required in query url")
	}

	if len(errInfo) > 0 {
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	data, statusCode, errInfo = c.useCase.Summary(ctx, month, year, email)
	response.SendBack(ctx, data, errInfo, statusCode)
	return
}

func (c *StatisticController) TransactionPriority(ctx *gin.Context) {
	var (
		errInfo []errorsinfo.Errors
	)

	month := ctx.Query("month")
	year := ctx.Query("year")

	if month == "" || year == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "month or year required in query url")
		response.SendBack(ctx, nil, errInfo, http.StatusBadRequest)
		return
	}

	data, httpCode, errInfo := c.useCase.Priority(ctx, month, year)

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *StatisticController) Trend(ctx *gin.Context) {
	var errInfo []errorsinfo.Errors

	month := ctx.Query("month")
	year := ctx.Query("year")

	if month == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "month required in query url")
	}

	if year == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "year required in query url")
	}

	if len(errInfo) > 0 {
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	data, httpCode, errInfo := c.useCase.Trend(ctx, month, year)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *StatisticController) AnalyticsTrend(ctx *gin.Context) {
	var errInfo []errorsinfo.Errors
	{
	}

	period := ctx.Query("period")
	typeName := ctx.Query("type")

	// validate
	if period == "" || typeName == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "period required in url param")
	}

	if typeName == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "type name required in url param")
	}

	if len(errInfo) > 0 {
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	data, httpCode, errInfo := c.useCase.AnalyticsTrend(ctx, period, typeName)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *StatisticController) ExpenseDetail(ctx *gin.Context) {
	var errInfo []errorsinfo.Errors

	// get value from parameter url
	month := ctx.Query("month")
	year := ctx.Query("year")
	email := ctx.Query("email")

	// validate
	if month == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "month required in query url")
	}

	if year == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "year required in query url")
	}

	if len(errInfo) > 0 {
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	data, httpCode, errInfo := c.useCase.ExpenseDetail(ctx, month, year, email)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *StatisticController) SubExpenseDetail(ctx *gin.Context) {
	var errInfo []errorsinfo.Errors

	month := ctx.Query("month")
	year := ctx.Query("year")
	IDCategory := ctx.Query("categoryid")

	if month == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "month required in query url")
	}

	if year == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "year required in query url")
	}

	if IDCategory == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id category required in query url")
	}

	// if any error, then send response
	if len(errInfo) > 0 {
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	IDCat, err := uuid.Parse(IDCategory)
	if err != nil {
		logrus.Error(err.Error())
	}

	data, httpCode, errInfo := c.useCase.SubExpenseDetail(ctx, month, year, IDCat)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}