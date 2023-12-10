package payments

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/payments/dtos"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/payments/entities"
	"github.com/semicolon-indonesia/wealthy-backend/constants"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/semicolon-indonesia/wealthy-backend/utils/orderid"
	"github.com/semicolon-indonesia/wealthy-backend/utils/personalaccounts"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"time"
)

type (
	PaymentUseCase struct {
		repo IPaymentRepository
	}

	IPaymentUseCase interface {
		Subscriptions(ctx *gin.Context, request *dtos.PaymentSubscription) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		MidtransWebhook(ctx *gin.Context, request *dtos.MidTransWebhook) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
	}
)

func NewPaymentUseCase(repo IPaymentRepository) *PaymentUseCase {
	return &PaymentUseCase{repo: repo}
}

func (s *PaymentUseCase) Subscriptions(ctx *gin.Context, request *dtos.PaymentSubscription) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var midTransResponse dtos.MidTansResponse

	orderID := orderid.Generator()

	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "token contains invalid information")
		return struct{}{}, http.StatusUnauthorized, errInfo
	}

	// convert package id ( string to uuid )
	PackageIDUUID, err := uuid.Parse(request.PackageID)
	if err != nil {
		logrus.Error(err.Error())
	}

	// check if package id correct
	checkPackage := s.repo.CheckPackageID(PackageIDUUID)

	if !checkPackage.Exists {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: "package id unknown",
		}
		return resp, http.StatusBadRequest, []errorsinfo.Errors{}
	}

	// check if already subscription
	dataSubs, err := s.repo.GetSubscriptionInformation(personalAccount.ID)
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	if dataSubs.ID != uuid.Nil {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: "you are in subscription period",
		}
		return resp, http.StatusBadRequest, []errorsinfo.Errors{}
	}

	midTransURL := os.Getenv("MIDTRANS_URL")
	midTransServerKey := os.Getenv("MIDTRANS_SERVER_KEY") + ":"

	dataPrice := s.repo.GetPrice(PackageIDUUID)

	bodyPayload := dtos.PaymentSnapRequest{
		Details: dtos.PaymentSnapDetails{
			OrderId:     orderID,
			GrossAmount: dataPrice.Price,
		},
	}

	payload, err := json.MarshalIndent(bodyPayload, "", "\t")
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return midTransResponse, http.StatusInternalServerError, errInfo
	}

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, midTransURL, bytes.NewBuffer(payload))
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return midTransResponse, http.StatusInternalServerError, errInfo
	}

	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(midTransServerKey))
	req.Header.Set("Authorization", basicAuth)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return midTransResponse, http.StatusInternalServerError, errInfo
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logrus.Error(err.Error())
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
			return midTransResponse, http.StatusInternalServerError, errInfo
		}

		err = json.Unmarshal(body, &midTransResponse)
		if err != nil {
			logrus.Error(err.Error())
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
			return midTransResponse, http.StatusInternalServerError, errInfo
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				logrus.Error(err.Error())
			}
		}(resp.Body)

		model := entities.SubsTransaction{
			ID:                 uuid.New(),
			IDPersonalAccount:  personalAccount.ID,
			RedirectURL:        midTransResponse.RedirectUrl,
			Token:              midTransResponse.Token,
			SubscriptionID:     PackageIDUUID,
			Amount:             dataPrice.Price,
			OrderID:            orderID,
			IDMasterSubsPeriod: dataPrice.IDMasterSubsPeriod,
		}

		result, err := s.repo.SaveSubscriptionPayment(&model)
		if !result {
			logrus.Error(err.Error())
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
			return midTransResponse, http.StatusInternalServerError, errInfo
		}

		if len(errInfo) == 0 {
			errInfo = []errorsinfo.Errors{}
		}
		return midTransResponse, resp.StatusCode, errInfo
	} else {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "payment gateway problem")
		return midTransResponse, resp.StatusCode, errInfo
	}
}

func (s *PaymentUseCase) MidtransWebhook(ctx *gin.Context, request *dtos.MidTransWebhook) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		model           entities.SubsInfo
		additionalMonth int
	)

	// get transaction information by order id
	dataTransaction, err := s.repo.GetTransactionInfoByOrderID(request.OrderId)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// check if data transaction not found
	if dataTransaction.ID == uuid.Nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "invalid order id")
		return struct{}{}, http.StatusBadRequest, errInfo
	}

	// get price information
	dataMasterPrice, err := s.repo.GetPeriodName(dataTransaction.IDMasterSubsPeriod)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// update transaction
	if err = s.repo.UpdateStatusTransaction(request.OrderId); err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// get additional month
	switch dataMasterPrice.PeriodName {
	case constants.Monthly:
		additionalMonth = 1
	case constants.SixMonthly:
		additionalMonth = 6
	case constants.Annual:
		additionalMonth = 12
	}

	// writes to user_subscription
	currentTime := time.Now()
	addMonthFinal := currentTime.AddDate(0, additionalMonth, 0)
	model.ID = uuid.New()
	model.IDPersonalAccounts = dataTransaction.IDPersonalAccount
	model.IDSubsTransaction = dataTransaction.ID
	model.PeriodEExpired = addMonthFinal

	err = s.repo.WriteUserSubscription(model)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// change to basic to pro
	IDProAccountUUID, err := uuid.Parse("826ed2f2-7dad-49c4-846b-131fa4e55161")
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	err = s.repo.ChangeAccountUser(dataTransaction.IDPersonalAccount, IDProAccountUUID)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// if empty
	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	return struct{}{}, http.StatusOK, errInfo
}