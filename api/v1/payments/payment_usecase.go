package payments

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/payments/dtos"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/payments/entities"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/semicolon-indonesia/wealthy-backend/utils/orderid"
	"github.com/semicolon-indonesia/wealthy-backend/utils/personalaccounts"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
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
		httpCode = http.StatusUnauthorized
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "token contains invalid information")
		return response, httpCode, errInfo
	}

	midTransURL := os.Getenv("MIDTRANS_URL")
	midTransServerKey := os.Getenv("MIDTRANS_SERVER_KEY") + ":"

	dataPrice := s.repo.GetPrice(request.PackageID)

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
			ID:                uuid.New(),
			IDPersonalAccount: personalAccount.ID,
			RedirectURL:       midTransResponse.RedirectUrl,
			Token:             midTransResponse.Token,
			SubscriptionID:    request.PackageID,
			Amount:            dataPrice.Price,
			OrderID:           orderID,
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
	if err := s.repo.MidtransWebhook(request.OrderId); err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return response, http.StatusInternalServerError, errInfo
	}

	// change to basic to pro

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	return response, http.StatusOK, errInfo
}