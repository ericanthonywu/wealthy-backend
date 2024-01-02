package budgets

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/wealthy-app/wealthy-backend/api/v1/budgets/dtos"
	"github.com/wealthy-app/wealthy-backend/constants"
	"github.com/wealthy-app/wealthy-backend/utils/errorsinfo"
	"github.com/wealthy-app/wealthy-backend/utils/response"
	"github.com/wealthy-app/wealthy-backend/utils/utilities"
	"net/http"
)

type (
	BudgetController struct {
		useCase IBudgetUseCase
	}

	IBudgetController interface {
		AllLimit(ctx *gin.Context)
		Overview(ctx *gin.Context)
		LatestMonths(ctx *gin.Context)
		Limit(ctx *gin.Context)
		Trends(ctx *gin.Context)
		Travels(ctx *gin.Context)
		UpdateTravelInfo(ctx *gin.Context)
	}
)

func NewBudgetController(useCase IBudgetUseCase) *BudgetController {
	return &BudgetController{useCase: useCase}
}

func (c *BudgetController) AllLimit(ctx *gin.Context) {
	data, httpCode, errInfo := c.useCase.AllLimit(ctx)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *BudgetController) Overview(ctx *gin.Context) {
	var errInfo []errorsinfo.Errors

	month := ctx.Query("month")
	year := ctx.Query("year")

	if month == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "month mandatory required in query url")
	}

	if year == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "year mandatory required in query url")
	}

	if len(errInfo) > 0 {
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	data, httpCode, errInfo := c.useCase.Overview(ctx, month, year)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *BudgetController) LatestMonths(ctx *gin.Context) {
	var (
		errInfo  []errorsinfo.Errors
		httpCode int
	)

	categoryID := ctx.Query("categoryid")

	if categoryID == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "category ID required in query url")
		response.SendBack(ctx, nil, errInfo, http.StatusBadRequest)
		return
	}

	catID, err := uuid.Parse(categoryID)
	if err != nil {
		logrus.Error(err.Error())
	}

	data, httpCode, errInfo := c.useCase.LatestMonths(ctx, catID)

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *BudgetController) Limit(ctx *gin.Context) {
	var (
		dtoRequest  dtos.BudgetSetRequest
		dtoResponse interface{}
		errInfo     []errorsinfo.Errors
		httpCode    int
		purpose     string
	)

	purpose = constants.NonTravel

	// bind
	if err := ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "body payload required")
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	// validation sections
	if !utilities.IsEmptyString(dtoRequest.TravelStartDate) {

		// get account type
		accountType := ctx.MustGet("accountType").(string)

		// if basic account
		if accountType == constants.AccountBasic {
			resp := struct {
				Message string `json:"message"`
			}{
				Message: constants.ProPlan,
			}
			response.SendBack(ctx, resp, []errorsinfo.Errors{}, http.StatusUpgradeRequired)
			return
		}

		purpose = constants.Travel
		if utilities.IsEmptyString(dtoRequest.TravelEndDate) {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", "travel_end_date attribute needed in body payload")
		}

		if utilities.IsEmptyString(dtoRequest.ImageBase64) {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", "image_base54 attribute needed in body payload")
		}

		if utilities.IsEmptyString(dtoRequest.IDMasterTransactionTypes.String()) {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id_master_transaction_types attribute needed in body payload")
		}

		if !utilities.ValidateBetweenTwoDateRange(dtoRequest.TravelStartDate, dtoRequest.TravelEndDate) {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", "travel_end_date attribute must greater than travel_start_date attribute in body payload")
		}

		if utilities.IsEmptyString(dtoRequest.IDMasterExchangeCurrency) {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id categories exchange currency empty value")
		}
	}

	// show err info
	if len(errInfo) > 0 {
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	dtoResponse, httpCode, errInfo = c.useCase.Limit(ctx, &dtoRequest, purpose)
	response.SendBack(ctx, dtoResponse, errInfo, httpCode)
	return
}

func (c *BudgetController) Trends(ctx *gin.Context) {
	var (
		errInfo     []errorsinfo.Errors
		httpCode    int
		dtoResponse interface{}
	)

	month := ctx.Query("month")
	year := ctx.Query("year")
	IDCategory := ctx.Query("categoryid")

	if month == "" || year == "" || IDCategory == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "month or year required in query url")
		response.SendBack(ctx, nil, errInfo, http.StatusBadRequest)
		return
	}

	IDCat, err := uuid.Parse(IDCategory)
	if err != nil {
		logrus.Error(err.Error())
	}

	dtoResponse, httpCode, errInfo = c.useCase.Trends(ctx, IDCat, month, year)

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}
	response.SendBack(ctx, dtoResponse, errInfo, httpCode)
	return
}

func (c *BudgetController) Travels(ctx *gin.Context) {
	// get account type
	accountType := ctx.MustGet("accountType").(string)

	// if basic account
	if accountType == constants.AccountBasic {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: constants.ProPlan,
		}
		response.SendBack(ctx, resp, []errorsinfo.Errors{}, http.StatusUpgradeRequired)
		return
	}

	dtoResponse, httpCode, errInfo := c.useCase.Travels(ctx)
	response.SendBack(ctx, dtoResponse, errInfo, httpCode)
	return
}

func (c *BudgetController) UpdateTravelInfo(ctx *gin.Context) {
	var (
		dtoRequest map[string]interface{}
		errInfo    []errorsinfo.Errors
	)

	// bind
	if err := ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "body payload required")
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	IDTravel := ctx.Param("id-travel")

	if IDTravel == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id travel required in query url")
	}

	dtoResponse, httpCode, errInfo := c.useCase.UpdateTravelInfo(ctx, IDTravel, dtoRequest)
	response.SendBack(ctx, dtoResponse, errInfo, httpCode)
	return
}