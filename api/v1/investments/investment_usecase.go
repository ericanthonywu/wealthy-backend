package investments

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/wealthy-app/wealthy-backend/api/v1/investments/dtos"
	"github.com/wealthy-app/wealthy-backend/api/v1/investments/entities"
	"github.com/wealthy-app/wealthy-backend/utils/datecustoms"
	"github.com/wealthy-app/wealthy-backend/utils/errorsinfo"
	"net/http"
	"sort"
	"strconv"
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
		dtoResponse      dtos.InvestmentResponse
		investmentDetail []dtos.InvestmentDetails
		investmentInfo   []dtos.InvestmentInfo
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

	// mapping response
	walletName := ""
	maxData := len(trxData) - 1
	totalInvestment := 0.0
	totalPotentialReturn := 0.0
	dateTotal := 0

	for k, v := range trxData {

		if v.TotalLot == 0 {
			continue
		}

		dateTotal = datecustoms.TotalDaysBetweenDate(v.DateTransaction)
		if dateTotal < 0 {
			dateTotal = 0
		}

		totalInvestment += v.InitialInvestment

		dataTrading, err := s.repo.GetTradingInfo(v.StockCode)
		if err != nil {
			logrus.Error(err.Error())
		}

		// if previous wallet same with new data
		if walletName == v.WalletName {

			// calculation for
			closePrice := float64(dataTrading.Close)
			potentialReturnString := fmt.Sprintf("%.2f", (closePrice-v.AverageBuy)*float64(v.TotalLot)*100)
			potentialReturn, _ := strconv.ParseFloat(potentialReturnString, 64)
			percentagePotentialReturn := (potentialReturn / v.InitialInvestment) * 100
			totalPotentialReturn += potentialReturn

			investmentInfo = append(investmentInfo, dtos.InvestmentInfo{
				Name:              dataTrading.Name,
				InitialInvestment: v.InitialInvestment,
				StockCode:         v.StockCode,
				Lot:               v.TotalLot,
				AverageBuy:        v.AverageBuy,
				PotentialReturn:   potentialReturn,
				PercentageReturn:  fmt.Sprintf("%.2f", percentagePotentialReturn) + "%",
				TotalDays:         int64(dateTotal),
			})

			// latest data
			if k == maxData {
				investmentDetail = append(investmentDetail, dtos.InvestmentDetails{
					BrokerName:          walletName,
					Info:                investmentInfo,
					UnrealizedPotential: potentialReturn,
				})

				// clear
				investmentInfo = nil
				closePrice = 0
				potentialReturn = 0

			}
		}

		// if previous broker empty
		if walletName == "" {
			// set broker name
			walletName = v.WalletName

			// calculation for
			closePrice := float64(dataTrading.Close)
			potentialReturnString := fmt.Sprintf("%.2f", (closePrice-v.AverageBuy)*float64(v.TotalLot)*100)
			potentialReturn, _ := strconv.ParseFloat(potentialReturnString, 64)
			percentagePotentialReturn := (potentialReturn / v.InitialInvestment) * 100
			totalPotentialReturn += potentialReturn

			investmentInfo = append(investmentInfo, dtos.InvestmentInfo{
				Name:              dataTrading.Name,
				InitialInvestment: v.InitialInvestment,
				StockCode:         v.StockCode,
				Lot:               v.TotalLot,
				AverageBuy:        v.AverageBuy,
				PotentialReturn:   potentialReturn,
				PercentageReturn:  fmt.Sprintf("%.2f", percentagePotentialReturn) + "%",
				TotalDays:         int64(dateTotal),
			})

			if k == maxData {
				investmentDetail = append(investmentDetail, dtos.InvestmentDetails{
					BrokerName:          walletName,
					Info:                investmentInfo,
					UnrealizedPotential: potentialReturn,
				})

				// clear
				investmentInfo = nil
				closePrice = 0
				potentialReturn = 0
			}
		}

		// if previous broker name different with new data
		if walletName != v.WalletName {
			// potential return
			closePrice := float64(dataTrading.Close)
			potentialReturnString := fmt.Sprintf("%.2f", (closePrice-v.AverageBuy)*float64(v.TotalLot)*100)
			potentialReturn, _ := strconv.ParseFloat(potentialReturnString, 64)

			investmentDetail = append(investmentDetail, dtos.InvestmentDetails{
				BrokerName:          walletName,
				Info:                investmentInfo,
				UnrealizedPotential: potentialReturn,
			})

			// clear
			investmentInfo = nil
			closePrice = 0
			potentialReturn = 0

			// renew
			walletName = v.WalletName

			// calculation for
			closePrice = float64(dataTrading.Close)
			potentialReturnString = fmt.Sprintf("%.2f", (closePrice-v.AverageBuy)*float64(v.TotalLot)*100)
			potentialReturn, _ = strconv.ParseFloat(potentialReturnString, 64)
			percentagePotentialReturn := (potentialReturn / v.InitialInvestment) * 100
			totalPotentialReturn += potentialReturn

			investmentInfo = append(investmentInfo, dtos.InvestmentInfo{
				Name:              dataTrading.Name,
				InitialInvestment: v.InitialInvestment,
				StockCode:         v.StockCode,
				Lot:               v.TotalLot,
				AverageBuy:        v.AverageBuy,
				PotentialReturn:   potentialReturn,
				PercentageReturn:  fmt.Sprintf("%.2f", percentagePotentialReturn) + "%",
				TotalDays:         int64(dateTotal),
			})

			// latest data
			if k == maxData {
				investmentDetail = append(investmentDetail, dtos.InvestmentDetails{
					BrokerName:          walletName,
					Info:                investmentInfo,
					UnrealizedPotential: potentialReturn,
				})

				// clear
				investmentInfo = nil

			}
		}
	}

	dtoResponse.TotalInvestment = totalInvestment
	dtoResponse.TotalPotentialReturn = totalPotentialReturn

	if totalPotentialReturn > 0 {
		percentage := (totalPotentialReturn / totalInvestment) * 100
		dtoResponse.PercentagePotentialReturn = fmt.Sprintf("%.2f", percentage) + "%"
	} else {
		dtoResponse.PercentagePotentialReturn = fmt.Sprintf("%.2f", 0.0) + "%"
	}

	if len(investmentDetail) == 0 {
		dtoResponse.Details = []dtos.InvestmentDetails{}
	} else {
		dtoResponse.Details = investmentDetail
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

		totalLot := float64(v.Lot)
		priceSell := float64(v.Price)

		// sell calculation
		valueSell, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", totalLot*priceSell*100), 64)
		feeSell, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", (v.FeeSell/100)*valueSell), 64)
		netSell, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", valueSell-feeSell), 64)

		// buy information
		valueBuy, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", (totalLot*dataInvestment.AverageBuy)*100), 64)
		feeBuy, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", (v.FeeBuy/100)*valueBuy), 64)
		netBuy, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", valueBuy+feeBuy), 64)

		// gain loss
		gainLoss, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", netSell-netBuy), 64)

		// percentage return
		percentageReturn, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", (gainLoss/valueBuy)*100), 64)

		dtoResponse = append(dtoResponse, dtos.InvestmentGainLoss{
			DataTransaction:   v.DateTransaction,
			BrokerName:        v.WalletName,
			StockCode:         v.StockCode,
			Lot:               v.Lot,
			Price:             float64(v.Price),
			Name:              dataTrading.Name,
			InitialInvestment: valueSell,
			Percentage:        fmt.Sprintf("%.2f", percentageReturn) + "%",
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