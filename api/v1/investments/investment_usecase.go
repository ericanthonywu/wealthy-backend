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
		averageBuy             float64
		averageBuyCollection   []float64
		potentialReturn        float64
		unreliazeReturn        float64
		initialInvestment      float64
	)

	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", constants.TokenInvalidInformation)
		return response, http.StatusUnauthorized, errInfo
	}

	// fetch transaction data first
	trxData, err := s.repo.TrxInfo(personalAccount.ID)
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
			Message: "no data transaction for investment",
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
					totalAvg := 0.0
					var resultAverageBuy float64

					totalLot += v.Lot
					averageBuy = float64(v.Lot * v.Price / v.Lot)
					averageBuyCollection = append(averageBuyCollection, averageBuy)
					initialInvestment += float64(v.Lot * v.Price)

					// is latest
					if k == (maxTrxData - 1) {
						// calculate average
						totalAvg = 0.0
						resultAverageBuy = 0

						for _, avg := range averageBuyCollection {
							totalAvg += avg
						}

						// renew data
						resultAverageBuy = totalAvg / float64(totalLot)
						averageBuyFinal, err := strconv.ParseFloat(fmt.Sprintf("%.2f", resultAverageBuy), 64)
						if err != nil {
							logrus.Error(err.Error())
						}
						potentialReturn = float64(dataTrading.Close) - resultAverageBuy*float64(totalLot)*100
						unreliazeReturn += potentialReturn

						// replace with new portfolio
						investmentNamePrevious = dataTrading.Name

						// append latest
						investmentInfo = append(investmentInfo, dtos.InvestmentInfo{
							InitialInvestment: initialInvestment,
							Name:              investmentNamePrevious,
							StockCode:         stockCodePrevious,
							Lot:               totalLot,
							AverageBuy:        averageBuyFinal,
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
					totalAvg := 0.0
					var resultAverageBuy float64

					for _, avg := range averageBuyCollection {
						totalAvg += avg
					}

					// renew data
					resultAverageBuy = totalAvg / float64(totalLot)

					// return
					averageBuyFinal, err := strconv.ParseFloat(fmt.Sprintf("%.2f", resultAverageBuy), 64)
					if err != nil {
						logrus.Error(err.Error())
					}
					potentialReturn = float64(dataTrading.Close) - resultAverageBuy*float64(totalLot)*100

					// append previous data
					investmentInfo = append(investmentInfo, dtos.InvestmentInfo{
						InitialInvestment: initialInvestment,
						Name:              investmentNamePrevious,
						StockCode:         stockCodePrevious,
						Lot:               totalLot,
						AverageBuy:        averageBuyFinal,
						PotentialReturn:   potentialReturn,
					})

					// clear arrays
					averageBuyCollection = nil
					totalAvg = 0.0
					resultAverageBuy = 0
					initialInvestment = 0

					// override value and renew
					totalLot = v.Lot
					stockCodePrevious = v.StockCode
					investmentNamePrevious = dataTrading.Name
					averageBuy = float64(v.Lot * v.Price / v.Lot)
					averageBuyCollection = append(averageBuyCollection, averageBuy)
					initialInvestment += float64(v.Lot * v.Price)

					// if latest data
					if k == (maxTrxData - 1) {
						// calculate average
						totalAvg = 0.0
						resultAverageBuy = 0

						for _, avg := range averageBuyCollection {
							totalAvg += avg
						}

						// renew data
						resultAverageBuy = totalAvg / float64(totalLot)
						averageBuyFinal, err = strconv.ParseFloat(fmt.Sprintf("%.2f", resultAverageBuy), 64)
						if err != nil {
							logrus.Error(err.Error())
						}
						potentialReturn = float64(dataTrading.Close) - resultAverageBuy*float64(totalLot)*100
						unreliazeReturn += potentialReturn

						// replace with new portfolio
						investmentNamePrevious = dataTrading.Name

						// append latest
						investmentInfo = append(investmentInfo, dtos.InvestmentInfo{
							Name:              investmentNamePrevious,
							InitialInvestment: initialInvestment,
							StockCode:         stockCodePrevious,
							Lot:               totalLot,
							AverageBuy:        averageBuyFinal,
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
				totalAvg := 0.0
				var resultAverageBuy float64

				// calculate average
				for _, avg := range averageBuyCollection {
					totalAvg += avg
				}

				// renew data
				resultAverageBuy = totalAvg / float64(totalLot)
				averageBuyFinal, err := strconv.ParseFloat(fmt.Sprintf("%.2f", resultAverageBuy), 64)
				if err != nil {
					logrus.Error(err.Error())
				}

				potentialReturn = float64(dataTrading.Close) - resultAverageBuy*float64(totalLot)*100
				unreliazeReturn += potentialReturn

				// append latest
				investmentInfo = append(investmentInfo, dtos.InvestmentInfo{
					InitialInvestment: initialInvestment,
					Name:              investmentNamePrevious,
					StockCode:         stockCodePrevious,
					Lot:               totalLot,
					AverageBuy:        averageBuyFinal,
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
				averageBuyCollection = nil
				brokerNamePrevious = v.BrokerName
				stockCodePrevious = v.StockCode
				investmentNamePrevious = dataTrading.Name
				totalLot = v.Lot
				averageBuy = float64(v.Lot * v.Price / v.Lot)
				averageBuyCollection = append(averageBuyCollection, averageBuy)
				initialInvestment += float64(v.Lot * v.Price)

				// is latest
				if k == (maxTrxData - 1) {

					investmentInfo = nil

					// calculate average
					totalAvg = 0.0
					resultAverageBuy = 0

					totalLot = v.Lot

					for _, avg := range averageBuyCollection {
						totalAvg += avg
					}

					// renew data
					resultAverageBuy = totalAvg / float64(totalLot)
					averageBuyFinal, err := strconv.ParseFloat(fmt.Sprintf("%.2f", resultAverageBuy), 64)
					if err != nil {
						logrus.Error(err.Error())
					}
					potentialReturn = float64(dataTrading.Close) - resultAverageBuy*float64(totalLot)*100
					unreliazeReturn += potentialReturn

					// replace with new portfolio
					investmentNamePrevious = dataTrading.Name

					// append latest
					investmentInfo = append(investmentInfo, dtos.InvestmentInfo{
						InitialInvestment: initialInvestment,
						Name:              investmentNamePrevious,
						StockCode:         stockCodePrevious,
						Lot:               totalLot,
						AverageBuy:        averageBuyFinal,
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
			initialInvestment += float64(v.Lot * v.Price)
			brokerNamePrevious = v.BrokerName
			stockCodePrevious = v.StockCode
			investmentNamePrevious = dataTrading.Name
			totalLot = v.Lot
			averageBuy = float64(v.Lot * v.Price / v.Lot)
			averageBuyCollection = append(averageBuyCollection, averageBuy)

			// is latest
			if k == (maxTrxData - 1) {

				investmentInfo = nil

				// calculate average
				totalAvg := 0.0
				var resultAverageBuy float64

				totalLot = v.Lot

				for _, avg := range averageBuyCollection {
					totalAvg += avg
				}

				// renew data
				resultAverageBuy = totalAvg / float64(totalLot)
				averageBuyFinal, err := strconv.ParseFloat(fmt.Sprintf("%.2f", resultAverageBuy), 64)
				if err != nil {
					logrus.Error(err.Error())
				}
				potentialReturn = float64(dataTrading.Close) - resultAverageBuy*float64(totalLot)*100
				unreliazeReturn += potentialReturn

				// replace with new portfolio
				investmentNamePrevious = dataTrading.Name

				// append latest
				investmentInfo = append(investmentInfo, dtos.InvestmentInfo{
					InitialInvestment: initialInvestment,
					Name:              investmentNamePrevious,
					StockCode:         stockCodePrevious,
					Lot:               totalLot,
					AverageBuy:        averageBuyFinal,
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
			Message: "no data found for gain loss investment",
		}
		return resp, http.StatusInternalServerError, errInfo
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