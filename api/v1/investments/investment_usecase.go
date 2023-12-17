package investments

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/investments/dtos"
	"github.com/semicolon-indonesia/wealthy-backend/constants"
	"github.com/semicolon-indonesia/wealthy-backend/utils/datecustoms"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/semicolon-indonesia/wealthy-backend/utils/personalaccounts"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type (
	InvestmentUseCase struct {
		repo IInvestmentRepository
	}

	IInvestmentUseCase interface {
		Portfolio(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		GainLoss(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
	}
)

func NewInvestmentUseCase(repo IInvestmentRepository) *InvestmentUseCase {
	return &InvestmentUseCase{repo: repo}
}

func (s *InvestmentUseCase) Portfolio(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		dtoResponse            dtos.InvestmentResponse
		investmentDetail       []dtos.InvestmentDetails
		investmentInfo         []dtos.InvestmentInfo
		stockCodePrevious      string
		brokerNamePrevious     string
		investmentNamePrevious string
		totalLot               int64
		buy                    float64
		buyCollections         []float64
		potentialReturn        float64
		unreliazeReturn        float64
		initialInvestment      float64
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

	maxTrxData := len(trxData)

	// mapping for response
	for k, v := range trxData {

		// get trading info
		dataTrading, err := s.repo.GetTradingInfo(v.StockCode)
		if err != nil {
			logrus.Error(err.Error())
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
			return response, http.StatusInternalServerError, errInfo
		}

		// get ready for mapping response
		// if previous has data
		if brokerNamePrevious != "" {
			// if broker previous still same new broker
			if brokerNamePrevious == v.BrokerName {
				// if still same stock code
				if stockCodePrevious == v.StockCode {
					// calculate average
					totalBuy := 0.0
					var averageBuy float64

					totalLot += v.Lot

					// buy value
					buy = float64(v.Lot * v.Price * 100)

					// append buy to collection
					buyCollections = append(buyCollections, buy)

					// initial investment
					initialInvestment += float64(v.Lot * v.Price)

					// is latest
					if k == (maxTrxData - 1) {
						// calculate average
						totalBuy = 0.0
						averageBuy = 0

						// add total buy
						for _, buyColl := range buyCollections {
							totalBuy += buyColl
						}

						// renew data
						// average buy
						averageBuy = totalBuy / float64(totalLot)

						// rounding average buy
						averageBuyRounding, err := strconv.ParseFloat(fmt.Sprintf("%.2f", averageBuy), 64)
						if err != nil {
							logrus.Error(err.Error())
						}

						// potential return
						potentialReturn = float64(dataTrading.Close) - averageBuy*float64(totalLot)*100

						unreliazeReturn += potentialReturn

						// replace with new portfolio
						investmentNamePrevious = dataTrading.Name

						// append latest
						investmentInfo = append(investmentInfo, dtos.InvestmentInfo{
							InitialInvestment: initialInvestment,
							Name:              investmentNamePrevious,
							StockCode:         stockCodePrevious,
							Lot:               totalLot,
							AverageBuy:        averageBuyRounding,
							PotentialReturn:   potentialReturn,
						})

						// append to investment detail
						investmentDetail = append(investmentDetail, dtos.InvestmentDetails{
							BokerName:           v.BrokerName,
							UnrealizedPotential: float64(unreliazeReturn),
							Info:                investmentInfo,
						})

						unreliazeReturn = 0.0
						initialInvestment = 0
					}
				}

				// if not same stock code with previous
				if stockCodePrevious != v.StockCode {
					// calculate average
					totalBuy := 0.0
					var averageBuy float64

					//
					for _, buyColl := range buyCollections {
						totalBuy += buyColl
					}

					// renew data
					// average buy
					averageBuy = totalBuy / float64(totalLot)

					// rounding average buy
					averageBuyRounding, err := strconv.ParseFloat(fmt.Sprintf("%.2f", averageBuy), 64)

					if err != nil {
						logrus.Error(err.Error())
					}

					// potential return
					potentialReturn = float64(dataTrading.Close) - averageBuy*float64(totalLot)*100

					// append previous data
					investmentInfo = append(investmentInfo, dtos.InvestmentInfo{
						InitialInvestment: initialInvestment,
						Name:              investmentNamePrevious,
						StockCode:         stockCodePrevious,
						Lot:               totalLot,
						AverageBuy:        averageBuyRounding,
						PotentialReturn:   potentialReturn,
					})

					// clear previous data
					buyCollections = nil
					totalBuy = 0.0
					averageBuy = 0
					initialInvestment = 0

					// override value and renew
					totalLot = v.Lot
					stockCodePrevious = v.StockCode
					investmentNamePrevious = dataTrading.Name

					// buy value
					buy = float64(v.Lot * v.Price * 100)

					// append new buy value
					buyCollections = append(buyCollections, buy)

					// initial investment
					initialInvestment += float64(v.Lot * v.Price)

					// if latest data
					if k == (maxTrxData - 1) {
						// calculate average
						totalBuy = 0.0
						averageBuy = 0

						for _, buyColl := range buyCollections {
							totalBuy += buyColl
						}

						// renew data
						// average buy
						averageBuy = totalBuy / float64(totalLot)

						// rounding average buy
						averageBuyRounding, err = strconv.ParseFloat(fmt.Sprintf("%.2f", averageBuy), 64)
						if err != nil {
							logrus.Error(err.Error())
						}

						// potential return
						potentialReturn = float64(dataTrading.Close) - averageBuy*float64(totalLot)*100
						unreliazeReturn += potentialReturn

						// replace with new portfolio
						investmentNamePrevious = dataTrading.Name

						// append latest
						investmentInfo = append(investmentInfo, dtos.InvestmentInfo{
							Name:              investmentNamePrevious,
							InitialInvestment: initialInvestment,
							StockCode:         stockCodePrevious,
							Lot:               totalLot,
							AverageBuy:        averageBuyRounding,
							PotentialReturn:   potentialReturn,
						})

						// append to investment detail
						investmentDetail = append(investmentDetail, dtos.InvestmentDetails{
							BokerName:           v.BrokerName,
							UnrealizedPotential: unreliazeReturn,
							Info:                investmentInfo,
						})

						unreliazeReturn = 0.0
						initialInvestment = 0
					}
				}
			}

			// if broker previous is not same new broker
			if brokerNamePrevious != v.BrokerName {
				totalBuy := 0.0
				var averageBuy float64

				// calculate average
				for _, buyColl := range buyCollections {
					totalBuy += buyColl
				}

				// renew data

				averageBuy = totalBuy / float64(totalLot)
				averageBuyRounding, err := strconv.ParseFloat(fmt.Sprintf("%.2f", averageBuy), 64)
				if err != nil {
					logrus.Error(err.Error())
				}

				// potential return
				potentialReturn = float64(dataTrading.Close) - averageBuy*float64(totalLot)*100
				unreliazeReturn += potentialReturn

				// append latest
				investmentInfo = append(investmentInfo, dtos.InvestmentInfo{
					InitialInvestment: initialInvestment,
					Name:              investmentNamePrevious,
					StockCode:         stockCodePrevious,
					Lot:               totalLot,
					AverageBuy:        averageBuyRounding,
					PotentialReturn:   potentialReturn,
				})

				// append to investment detail
				investmentDetail = append(investmentDetail, dtos.InvestmentDetails{
					BokerName:           brokerNamePrevious,
					UnrealizedPotential: unreliazeReturn,
					Info:                investmentInfo,
				})

				// clear
				investmentInfo = nil
				unreliazeReturn = 0.0
				initialInvestment = 0

				// replace with new portfolio
				investmentNamePrevious = dataTrading.Name

				// renew data
				buyCollections = nil
				brokerNamePrevious = v.BrokerName
				stockCodePrevious = v.StockCode
				investmentNamePrevious = dataTrading.Name
				totalLot = v.Lot

				// buy value
				buy = float64(v.Lot * v.Price * 100)

				// append buy collections
				buyCollections = append(buyCollections, buy)

				// initial investment
				initialInvestment += float64(v.Lot * v.Price)

				// is latest
				if k == (maxTrxData - 1) {

					// reset value
					investmentInfo = nil
					totalBuy = 0.0
					averageBuy = 0

					// renew value

					totalLot = v.Lot

					for _, buyColl := range buyCollections {
						totalBuy += buyColl
					}

					// average buy
					averageBuy = totalBuy / float64(totalLot)

					// rounding average buy
					averageBuyRounding, err := strconv.ParseFloat(fmt.Sprintf("%.2f", averageBuy), 64)
					if err != nil {
						logrus.Error(err.Error())
					}

					// potential return
					potentialReturn = float64(dataTrading.Close) - averageBuy*float64(totalLot)*100
					unreliazeReturn += potentialReturn

					// replace with new portfolio
					investmentNamePrevious = dataTrading.Name

					// append latest
					investmentInfo = append(investmentInfo, dtos.InvestmentInfo{
						InitialInvestment: initialInvestment,
						Name:              investmentNamePrevious,
						StockCode:         stockCodePrevious,
						Lot:               totalLot,
						AverageBuy:        averageBuyRounding,
						PotentialReturn:   potentialReturn,
					})

					// append to investment detail
					investmentDetail = append(investmentDetail, dtos.InvestmentDetails{
						BokerName:           v.BrokerName,
						UnrealizedPotential: unreliazeReturn,
						Info:                investmentInfo,
					})

					unreliazeReturn = 0.0
					initialInvestment = 0
				}
			}
		}

		// fist time broker name
		if brokerNamePrevious == "" {
			// set data for first time

			brokerNamePrevious = v.BrokerName
			stockCodePrevious = v.StockCode
			investmentNamePrevious = dataTrading.Name
			totalLot = v.Lot

			// buy value
			buy = float64(v.Lot * v.Price * 100)

			// initial investment
			initialInvestment += float64(v.Lot * v.Price)

			// append buy
			buyCollections = append(buyCollections, buy)

			// is latest
			if k == (maxTrxData - 1) {

				// reset value
				investmentInfo = nil
				totalBuy := 0.0
				var averageBuy float64

				totalLot = v.Lot

				for _, buyColl := range buyCollections {
					totalBuy += buyColl
				}

				// renew data
				averageBuy = totalBuy / float64(totalLot)

				averageBuyRounding, err := strconv.ParseFloat(fmt.Sprintf("%.2f", averageBuy), 64)
				if err != nil {
					logrus.Error(err.Error())
				}

				// potential return
				potentialReturn = float64(dataTrading.Close) - averageBuy*float64(totalLot)*100
				unreliazeReturn += potentialReturn

				// replace with new portfolio
				investmentNamePrevious = dataTrading.Name

				// append latest
				investmentInfo = append(investmentInfo, dtos.InvestmentInfo{
					InitialInvestment: initialInvestment,
					Name:              investmentNamePrevious,
					StockCode:         stockCodePrevious,
					Lot:               totalLot,
					AverageBuy:        averageBuyRounding,
					PotentialReturn:   potentialReturn,
				})

				// append to investment detail
				investmentDetail = append(investmentDetail, dtos.InvestmentDetails{
					BokerName:           v.BrokerName,
					UnrealizedPotential: unreliazeReturn,
					Info:                investmentInfo,
				})

				unreliazeReturn = 0.0
				initialInvestment = 0
			}
		}
	}

	dtoResponse.Details = investmentDetail

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

	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", constants.TokenInvalidInformation)
		return response, http.StatusUnauthorized, errInfo
	}

	dataTrx, err := s.repo.TrxInfoSell(personalAccount.ID)
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
		dataTrading, err := s.repo.GetTradingInfo(v.StockCode)
		if err != nil {
			logrus.Error(err.Error())
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
			return response, http.StatusInternalServerError, errInfo
		}

		dtoResponse = append(dtoResponse, dtos.InvestmentGainLoss{
			DataTransaction:   v.DateTransaction,
			StockCode:         v.StockCode,
			Lot:               v.Lot,
			Price:             0,
			Name:              dataTrading.Name,
			InitialInvestment: 0,
			Percentage:        "",
			TotalDays:         datecustoms.TotalDaysBetweenDate(v.DateTransaction),
		})
	}

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	return dtoResponse, http.StatusOK, errInfo
}