package investments

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/investments/dtos"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/investments/entities"
	"github.com/semicolon-indonesia/wealthy-backend/utils/datecustoms"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/sirupsen/logrus"
	"net/http"
	"sort"
)

type (
	InvestmentUseCase struct {
		repo IInvestmentRepository
	}

	IInvestmentUseCase interface {
		Portfolio(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		GainLoss(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		sortByAttribute(data []entities.InvestmentTransaction, attribute string)
	}
)

func NewInvestmentUseCase(repo IInvestmentRepository) *InvestmentUseCase {
	return &InvestmentUseCase{repo: repo}
}

func (s *InvestmentUseCase) Portfolio(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		dtoResponse dtos.InvestmentResponse
	)

	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	// fetch transaction data first
	trxData, err := s.repo.TrxInfo(accountUUID)
	if err != nil {
		logrus.Errorf(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return response, http.StatusInternalServerError, errInfo
	}

	// if no transaction data
	if len(trxData) == 0 {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: "no data for portfolio investment",
		}
		return resp, http.StatusNotFound, []errorsinfo.Errors{}
	}

	// clear
	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	return dtoResponse, http.StatusOK, errInfo
}

func (s *InvestmentUseCase) GainLoss(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		dtoResponse []dtos.InvestmentGainLoss
		err         error
	)

	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	dataTrx, err := s.repo.InvestmentTrx(accountUUID)
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	if len(dataTrx) == 0 {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: "no data for gain loss investment",
		}
		return resp, http.StatusNotFound, []errorsinfo.Errors{}
	}

	for _, v := range dataTrx {

		dataInvestment, err := s.repo.GetInvestmentDataHelper(accountUUID, v.StockCode)
		if err != nil {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
			return struct{}{}, http.StatusInternalServerError, errInfo
		}

		// get trading info
		dataTrading, err := s.repo.GetTradingInfo(v.StockCode)
		if err != nil {
			logrus.Error(err.Error())
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
			return response, http.StatusInternalServerError, errInfo
		}

		// sell calculation
		valueSell := float64(v.Lot * v.Price * 100)
		feeSell := v.FeeSell * valueSell
		netSell := valueSell - feeSell

		// buy information
		valueBuy := dataInvestment.ValueBuy
		netBuy := dataInvestment.NetBuy

		// gain loss
		gainLoss := netSell - netBuy

		// percentage return
		percentageReturn := gainLoss / valueBuy

		dtoResponse = append(dtoResponse, dtos.InvestmentGainLoss{
			DataTransaction:   v.DateTransaction,
			BrokerName:        v.BrokerName,
			StockCode:         v.StockCode,
			Lot:               v.Lot,
			Price:             float64(v.Price),
			Name:              dataTrading.Name,
			InitialInvestment: valueSell,
			Percentage:        percentageReturn,
			TotalDays:         datecustoms.TotalDaysBetweenDate(v.DateTransaction),
			GainLoss:          gainLoss,
		})
	}

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	return dtoResponse, http.StatusOK, errInfo
}

func (s *InvestmentUseCase) sortByAttribute(data []entities.InvestmentTransaction, attribute string) {
	switch attribute {
	case "StockCode":
		sort.Slice(data, func(i, j int) bool { return data[i].StockCode < data[j].StockCode })
	default:
		sort.Slice(data, func(i, j int) bool { return data[i].DateTransaction < data[j].DateTransaction })
	}
}