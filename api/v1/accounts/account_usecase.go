package accounts

import (
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/wealthy-app/wealthy-backend/api/v1/accounts/dtos"
	"github.com/wealthy-app/wealthy-backend/api/v1/accounts/entities"
	"github.com/wealthy-app/wealthy-backend/constants"
	"github.com/wealthy-app/wealthy-backend/utils/datecustoms"
	"github.com/wealthy-app/wealthy-backend/utils/errorsinfo"
	"github.com/wealthy-app/wealthy-backend/utils/password"
	"github.com/wealthy-app/wealthy-backend/utils/token"
	"github.com/wealthy-app/wealthy-backend/utils/utilities"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"
)

type (
	AccountUseCase struct {
		repo IAccountRepository
	}

	IAccountUseCase interface {
		SignIn(request *dtos.AccountSignInRequest) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		SignUp(request *dtos.AccountSignUpRequest) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		SignOut()
		GetProfile(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		UpdateProfile(ctx *gin.Context, request map[string]interface{}) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		ChangePassword(ctx *gin.Context, request *dtos.AccountChangePassword) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		ForgotPassword(ctx *gin.Context, request *dtos.AccountForgotPasswordRequest) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		ValidateRefCode(request *dtos.AccountRefCodeValidationRequest) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		SetAvatar(ctx *gin.Context, request *dtos.AccountAvatarRequest) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		RemoveAvatar(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		SearchAccount(ctx *gin.Context, dtoRequest *dtos.AccountGroupSharing) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		InviteSharing(ctx *gin.Context, dtoResponse *dtos.AccountGroupSharing) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		AcceptSharing(ctx *gin.Context, dtoRequest *dtos.AccountGroupSharingAccept) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		RejectSharing(ctx *gin.Context, dtoRequest *dtos.AccountGroupSharingAccept) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		RemoveSharing(ctx *gin.Context, dtoRequest *dtos.AccountGroupSharingRemove) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		GroupSharingAccepted(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		GroupSharingPending(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		VerifyOTP(ctx *gin.Context, request *dtos.AccountOTPVerify) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		ChangePasswordForgot(ctx *gin.Context, request *dtos.AccountChangeForgotPassword) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		DeleteAccount(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
	}
)

func NewAccountUseCase(repo IAccountRepository) *AccountUseCase {
	return &AccountUseCase{repo: repo}
}

func (s *AccountUseCase) SignUp(request *dtos.AccountSignUpRequest) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		dtoResponse       dtos.AccountSignUpResponse
		role              string
		idPersonalAccount uuid.UUID
		level             int
		err               error
	)

	personalAccountEntity := entities.AccountSignUpPersonalAccountEntity{
		Username:  request.Username,
		Name:      request.Name,
		Email:     request.Email,
		ReferCode: request.RefCode,
	}

	authAccountActivity := entities.AccountSignUpAuthenticationsEntity{
		Password: password.Generate(request.Password),
		Active:   true,
	}

	// get all referral code that registered
	listRefCodeExist := s.repo.ListRefCode()

	// is new referral code already exist
	if utilities.Contains(listRefCodeExist, request.RefCode) {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", errors.New("referral code is already registered").Error())
		return response, http.StatusBadRequest, errInfo
	}

	// save data into db
	role, idPersonalAccount, httpCode, errInfo = s.repo.SignUp(&personalAccountEntity, &authAccountActivity)

	// can not save data into db
	if len(errInfo) > 0 {
		resp := struct {
			Message string `json:"message,omitempty"`
		}{}
		return resp, http.StatusBadRequest, errInfo
	}

	// check referral code reference is not or empty
	if request.RefCodeReference != "" {

		// check referral code reference from payload is existed in database
		if !utilities.Contains(listRefCodeExist, request.RefCodeReference) {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", errors.New("unknown referrals code reference").Error())
			return response, http.StatusBadRequest, errInfo
		}

		// set up for reward position ( determined level )
		level, err = s.repo.GetLevelReferenceCode(request.RefCodeReference)
		if err != nil {
			logrus.Error(err.Error())
		}

		if err == nil {
			level = level + 1
		}

		if level == 0 {
			level = level + 1
		}
	}

	// set up for reward level section
	if request.RefCodeReference == "" {
		level = 0
	}

	newID, err := uuid.NewUUID()
	if err != nil {
		logrus.Error(err.Error())
	}

	// writing for user reference table
	model := entities.AccountRewards{
		ID:               newID,
		RefCode:          request.RefCode,
		RefCodeReference: request.RefCodeReference,
		Level:            level,
	}

	err = s.repo.WriteRewardsList(&model)
	if err != nil {
		logrus.Error(err.Error())
	}

	dtoResponse.Customer = dtos.CustomerDetail{
		ID:       idPersonalAccount,
		Name:     request.Name,
		Username: request.Username,
		Email:    request.Email,
	}

	dtoResponse.Account = dtos.Account{
		Role: role,
	}

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	// move contents of duplicate
	err = s.repo.DuplicateExpenseCategory(idPersonalAccount)
	if err != nil {
		logrus.Error(err.Error())
	}

	// move contents of sub-category
	err = s.repo.DuplicateExpenseSUbCategory(idPersonalAccount)
	if err != nil {
		logrus.Error(err.Error())
	}

	// move contents of income category
	err = s.repo.DuplicateIncomeCategory(idPersonalAccount)
	if err != nil {
		logrus.Error(err.Error())
	}

	return dtoResponse, httpCode, errInfo
}

func (s *AccountUseCase) SignIn(request *dtos.AccountSignInRequest) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		dtoResponse dtos.AccountSignInResponse
		err         error
	)

	authentication := entities.AccountSignInAuthenticationEntity{
		Email:    request.Email,
		Password: request.Password,
	}

	// check get data based on email
	data := s.repo.SignInAuth(authentication)

	// if data not found based on email and password
	if data.ID == uuid.Nil || data.Email == "" {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: "email address, " + request.Email + " , not registered",
		}
		return resp, http.StatusBadRequest, []errorsinfo.Errors{}
	}

	// if email found but password not match
	resultOfCompare := password.Compare(data.Password, []byte(request.Password))
	if !resultOfCompare {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: "password invalid for email address : " + request.Email,
		}
		return resp, http.StatusBadRequest, []errorsinfo.Errors{}
	}

	// generate token
	dtoResponse.Token, dtoResponse.TokenExp, err = token.JWTBuilder(request.Email, data.Roles)
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	dtoResponse.Customer.Email = request.Email
	dtoResponse.Customer.CustomerID = data.ID.String()
	dtoResponse.Account.Role = data.Roles
	dtoResponse.AccountType.Type = data.Type

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	return dtoResponse, httpCode, errInfo
}

func (s *AccountUseCase) SignOut() {

}

func (s *AccountUseCase) GetProfile(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var dtoResponse dtos.AccountProfile

	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	// get profile by customer ID
	dataProfile := s.repo.GetProfile(accountUUID)

	if dataProfile.Email == "" || dataProfile.ID == uuid.Nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "can not give profile info. profile not set properly")
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	dtoResponse.AccountCustomer.ID = dataProfile.ID
	dtoResponse.AccountCustomer.Name = dataProfile.Name
	dtoResponse.AccountCustomer.ReferType = dataProfile.ReferType
	dtoResponse.AccountCustomer.Username = dataProfile.Username
	dtoResponse.AccountCustomer.Email = dataProfile.Email

	// account type, override in beta promotion session
	accountType := ctx.MustGet("accountType").(string)

	if dataProfile.IDGender == uuid.Nil {
		dtoResponse.AccountGender.ID = ""
	} else {
		dtoResponse.AccountGender.ID = dataProfile.IDGender.String()
	}
	dtoResponse.AccountGender.Value = dataProfile.Gender

	dtoResponse.AccountDetail.AccountType = accountType
	dtoResponse.AccountDetail.UserRoles = dataProfile.UserRoles

	if dataProfile.ImagePath == "" {
		dtoResponse.AccountAvatar.URL = ""
		dtoResponse.AccountAvatar.FileName = ""
	} else {
		dtoResponse.AccountAvatar.URL = os.Getenv("APP_HOST") + "/v1/" + dataProfile.ImagePath
		dtoResponse.AccountAvatar.FileName = dataProfile.FileName
	}

	dtoResponse.AccountCustomer.DOB = dataProfile.DOB

	return dtoResponse, http.StatusOK, []errorsinfo.Errors{}
}

func (s *AccountUseCase) UpdateProfile(ctx *gin.Context, request map[string]interface{}) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		err            error
		dateOrigin     string
		idMasterGender string
	)

	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	// validate restrict changes
	if request["id"] != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", errors.New("can not change customer ID").Error())
	}

	if request["refer_code"] != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", errors.New("can not change refer code").Error())
	}

	if request["email"] != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", errors.New("can not change email").Error())
	}

	if request["image_path"] != nil || request["filename"] != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", errors.New("can not change avatar. use set avatar API instead").Error())
	}

	if request["id_master_account_type"] != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", errors.New("can not change account type. use subscription API for switch to PRO account").Error())
	}

	// if error occurs
	if len(errInfo) > 0 {
		return struct{}{}, http.StatusBadRequest, errInfo
	}

	// check date of birth ( dob ) value
	value, exists := request["dob"]
	if exists {
		dateOrigin = fmt.Sprintf("%v", value)

		if dateOrigin == "" {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", "date of birth empty value")
			return struct{}{}, http.StatusBadRequest, errInfo
		}

		// validate format dob
		if !datecustoms.ValidDateFormat(dateOrigin) {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", "format date must following YYYY-MM-DD")
			return struct{}{}, http.StatusBadRequest, errInfo
		}
	}

	// check id categories gender value
	value, exists = request["id_master_gender"]
	if exists {
		idMasterGender = fmt.Sprintf("%v", value)

		if idMasterGender == "" {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id categories gender empty value")
			return struct{}{}, http.StatusBadRequest, errInfo
		}

		// check categories gender
		idUUID, err := uuid.Parse(idMasterGender)
		if err != nil {
			logrus.Error(err.Error())
		}

		if !s.repo.GenderData(idUUID) {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id categories gender unregistered")
			return struct{}{}, http.StatusBadRequest, errInfo
		}
	}

	// latitude
	value, exists = request["latitude"]
	if exists {
		latitudeOrigin := fmt.Sprintf("%v", value)

		if latitudeOrigin == "" {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", "latitude empty value")
			return struct{}{}, http.StatusBadRequest, errInfo
		}
	}

	// longitude
	value, exists = request["longitude"]
	if exists {
		longitudeOrigin := fmt.Sprintf("%v", value)

		if longitudeOrigin == "" {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", "longitude empty value")
			return struct{}{}, http.StatusBadRequest, errInfo
		}
	}

	// fcmtoken
	value, exists = request["fcmtoken"]
	if exists {
		fcmToken := fmt.Sprintf("%v", value)

		if fcmToken == "" {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", "fcm token empty value")
			return struct{}{}, http.StatusBadRequest, errInfo
		}
	}

	// update profile
	err = s.repo.UpdateProfile(accountUUID, request)
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	resp := struct {
		Message string `json:"message,omitempty"`
	}{
		Message: "update profile successfully",
	}
	return resp, http.StatusOK, []errorsinfo.Errors{}
}

func (s *AccountUseCase) ChangePassword(ctx *gin.Context, request *dtos.AccountChangePassword) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		update      = make(map[string]interface{})
		dtoResponse = make(map[string]bool)
	)

	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	// get data that stored in db
	dataProfile := s.repo.GetProfilePassword(accountUUID)

	// check if dataProfile existed from database
	if dataProfile.Email == "" || dataProfile.ID == uuid.Nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "system can not change password due to data has deleted improperly")
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// compare old password from payload between old password stored in database
	resultOfCompare := password.Compare(dataProfile.Password, []byte(request.OldPassword))
	if !resultOfCompare {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "old password incorrect")
		return dtoResponse, http.StatusUnprocessableEntity, errInfo
	}

	// generate new password with encryption
	newPassword := password.Generate(request.NewPassword)

	// set new password
	update["password"] = newPassword

	err := s.repo.UpdatePassword(accountUUID, update)
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// reset error info if empty
	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	resp := struct {
		Message string `json:"message,omitempty"`
	}{
		Message: "change password successfully",
	}

	return resp, http.StatusOK, errInfo
}

func (s *AccountUseCase) ForgotPassword(ctx *gin.Context, request *dtos.AccountForgotPasswordRequest) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {

	dataProfile, err := s.repo.GetProfileByEmail(request.EmailAccount)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusBadRequest, errInfo
	}

	if dataProfile.Email == "" || dataProfile.ID == uuid.Nil {
		resp := struct {
			Message string `json:"message,omitempty"`
		}{
			Message: "email account not registered in system",
		}
		return resp, http.StatusBadRequest, errInfo
	}

	// get six digit random
	otpCode := utilities.GenerateRandomSixDigitNumber()

	// static path for email template
	templatePath := "./assets/files/reset-pass.html"
	logoPath := constants.LogoPrimary

	// keys that will be used in HTML template
	contents := map[string]string{
		"username": dataProfile.Username,
		"otp":      fmt.Sprintf("%d", otpCode),
		"logo":     logoPath,
	}

	// reading email template file based on static path before
	htmlContent, err := os.ReadFile(templatePath)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "can not find email template. reason : "+err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// replace variable in origin html
	modifiedHTML := utilities.HTMLContentReplacer(string(htmlContent), contents)

	// credential email server
	portNumber, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	smtpServer := os.Getenv("SMTP_SERVER")
	smtpPort := portNumber
	senderEmail := os.Getenv("SMTP_SENDER")
	senderPassword := os.Getenv("SMTP_SENDER_PASS")
	recipientEmail := request.EmailAccount

	// email subject
	subject := constants.EmailSubject

	// SMTP authentication
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpServer)

	// tls configuration
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         smtpServer,
	}

	// try to connect to SMTP server
	smtpAddress := fmt.Sprintf("%s:%d", smtpServer, smtpPort)
	conn, err := tls.Dial("tcp", smtpAddress, tlsConfig)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "can not connecting to SMTP. reason : "+err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// Create an SMTP client
	client, err := smtp.NewClient(conn, smtpServer)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "can not creating SMTP client. reason : "+err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// Authenticate with the Gmail SMTP server
	err = client.Auth(auth)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "can not authenticating with SMTP server. reason : "+err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// setting from address
	from := senderEmail

	// compose the email message with MIME headers
	message := fmt.Sprintf("From: %s\r\n", from)
	message += fmt.Sprintf("To: %s\r\n", recipientEmail)
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += "MIME-version: 1.0;\r\n"
	message += "Content-Type: text/html; charset=\"UTF-8\";\r\n"
	message += "Content-Transfer-Encoding: 7bit;\r\n"
	message += "\r\n" + modifiedHTML

	// set the sender and recipient
	err = client.Mail(senderEmail)
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "can not setting sender. reason : "+err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	err = client.Rcpt(recipientEmail)
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "can not setting recipient. reason : "+err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// send the email message
	wc, err := client.Data()
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "can not opening data connection. reason : "+err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	_, err = wc.Write([]byte(message))
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "can not opening writing email message. reason : "+err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	err = wc.Close()
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "can not closing data connection. reason : "+err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	err = client.Quit()
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "can not quitting SMTP session. reason : "+err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// store into db with preparation model
	currentTime := time.Now()
	expiredTime := currentTime.Add(3 * time.Minute)

	model := entities.AccountForgotPassword{
		ID:                uuid.New(),
		OTPCode:           fmt.Sprintf("%d", otpCode),
		IDPersonalAccount: dataProfile.ID,
		IsVerified:        false,
		Expired:           expiredTime,
		CreatedAt:         currentTime,
	}

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	// store into db
	err = s.repo.ForgotPassword(&model)
	if err != nil {
		logrus.Error(err.Error())
	}

	resp := struct {
		Message string `json:"message,omitempty"`
	}{
		Message: "otp has been sent to email",
	}

	return resp, http.StatusOK, errInfo
}

func (s *AccountUseCase) ValidateRefCode(request *dtos.AccountRefCodeValidationRequest) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var dtoResponse dtos.AccountRefCodeValidationResponse

	// get list ref code
	RefCodeList := s.repo.ListRefCode()

	// check
	if utilities.Contains(RefCodeList, request.RefCode) {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "referrals code already exist on system. please use another code to be registered")
		return struct{}{}, http.StatusOK, errInfo
	}

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	dtoResponse.Available = true
	return dtoResponse, http.StatusOK, errInfo
}

func (s *AccountUseCase) SetAvatar(ctx *gin.Context, request *dtos.AccountAvatarRequest) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	updateProfile := make(map[string]interface{})

	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	// decode base64 image from string
	imageData, err := base64.StdEncoding.DecodeString(request.ImageBase64)
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// get profile from customer ID
	dataProfile := s.repo.GetProfile(accountUUID)

	// remove old image from storage
	if dataProfile.ImagePath != "" || dataProfile.FileName != "" {
		target := "assets/avatar/" + dataProfile.FileName
		err = os.Remove(target)
		if err != nil {
			logrus.Error(err.Error())
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
			return struct{}{}, http.StatusInternalServerError, errInfo
		}
	}

	// remove empty column after remove file
	updateProfile["image_path"] = ""
	updateProfile["file_name"] = ""
	err = s.repo.UpdateProfile(accountUUID, updateProfile)
	if err != nil {
		logrus.Error(err.Error())
	}

	// setup new file name and target path
	filename := fmt.Sprintf("%d", time.Now().Unix()) + ".png"
	targetPath := "assets/avatar/" + filename

	// save image into storage
	err = utilities.SaveImage(imageData, targetPath)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// setup for new data
	updateProfile["image_path"] = "images/avatar/" + filename
	updateProfile["file_name"] = filename

	// update table with new data
	err = s.repo.UpdateProfile(accountUUID, updateProfile)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	resp := struct {
		Message string `json:"message,omitempty"`
	}{
		Message: "set avatar for profile successfully",
	}

	return resp, http.StatusOK, errInfo
}

func (s *AccountUseCase) RemoveAvatar(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		err         error
		dtoResponse dtos.AccountAvatarResponse
	)

	updateProfile := make(map[string]interface{})

	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	// get profile by customer id
	dataProfile := s.repo.GetProfile(accountUUID)

	// remove previous image
	if dataProfile.ImagePath != "" || dataProfile.FileName != "" {
		target := "assets/avatar/" + dataProfile.FileName
		err = os.Remove(target)
		if err != nil {
			logrus.Error(err.Error())
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
			return struct{}{}, http.StatusInternalServerError, errInfo
		}
	}

	// setup empty value
	updateProfile["image_path"] = ""
	updateProfile["file_name"] = ""

	// update column with empty value
	err = s.repo.UpdateProfile(accountUUID, updateProfile)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	resp := struct {
		Message string `json:"message,omitempty"`
	}{
		Message: "delete avatar from profile successfully",
	}
	dtoResponse.Success = true
	return resp, http.StatusOK, errInfo
}

func (s *AccountUseCase) SearchAccount(ctx *gin.Context, dtoRequest *dtos.AccountGroupSharing) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var dtoResponse dtos.AccountSearch

	dataPersonalAccounts, err := s.repo.SearchAccount(dtoRequest.EmailAccount)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	if dataPersonalAccounts.ID == uuid.Nil {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: "account email not existed",
		}
		return resp, http.StatusNotFound, errInfo
	}

	dtoResponse.ID = dataPersonalAccounts.ID
	return dtoResponse, http.StatusOK, errInfo
}

func (s *AccountUseCase) InviteSharing(ctx *gin.Context, dtoResponse *dtos.AccountGroupSharing) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		model              utilities.NotificationEntities
		modelInviteSharing entities.AccountGroupSharing
		err                error
	)

	accountUUID := ctx.MustGet("accountID").(uuid.UUID)
	usrEmail := ctx.MustGet("email").(string)

	// get profile for email target
	dataProfileSender, err := s.repo.GetProfileByEmail(usrEmail)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// get profile for email target
	dataProfileReceipt, err := s.repo.GetProfileByEmail(dtoResponse.EmailAccount)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// if target email same as email sender in token
	if dataProfileReceipt.Email == usrEmail {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "can not share with the same email")
		return struct{}{}, http.StatusBadRequest, errInfo
	}

	// if email target empty from db
	if dataProfileReceipt.Email == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "can not share with the unregistered account")
		return struct{}{}, http.StatusBadRequest, errInfo
	}

	// if already in group sharing
	if s.repo.IsAlreadySharing(accountUUID, dataProfileReceipt.ID) {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: "already invited before",
		}
		return resp, http.StatusBadRequest, []errorsinfo.Errors{}
	}

	// set ID
	IDGroupSharing := uuid.New()

	// save invitation
	modelInviteSharing.ID = IDGroupSharing
	modelInviteSharing.ShareFrom = accountUUID
	modelInviteSharing.ShareTo = dataProfileReceipt.ID
	modelInviteSharing.IsAccepted = false

	err = s.repo.InviteSharing(&modelInviteSharing)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// set notifications
	model.ID = uuid.New()
	model.IDPersonalAccounts = dataProfileReceipt.ID
	model.IsRead = false
	model.NotificationTitle = constants.NotificationTitle
	model.NotificationDescription = dataProfileSender.Name + " has invite you to become group sharing member"
	model.IDGroupSharing = IDGroupSharing.String()

	err = utilities.SetNotifications(ctx, model)
	if err != nil {
		logrus.Error(err.Error())
	}

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	resp := struct {
		Message string `json:"message"`
	}{
		Message: "invitation has been sent successfully",
	}

	return resp, http.StatusOK, errInfo
}

func (s *AccountUseCase) AcceptSharing(ctx *gin.Context, dtoRequest *dtos.AccountGroupSharingAccept) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	// get customer id
	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	// translate string to uuid value
	IDGroupSharingUUID, err := uuid.Parse(dtoRequest.IDGroupSharing)
	if err != nil {
		logrus.Error(err.Error())
	}

	// get information group sharing by id
	dataGroupSharing, err := s.repo.GroupSharingInfoByID(IDGroupSharingUUID)
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// if not eligible
	if dataGroupSharing.ShareFrom == accountUUID {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: "tokens are not entitled to receive invitations",
		}

		return resp, http.StatusBadRequest, []errorsinfo.Errors{}
	}

	// accept sharing
	err = s.repo.AcceptSharing(IDGroupSharingUUID)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	resp := struct {
		Message string `json:"message"`
	}{
		Message: "accept group sharing successfully",
	}

	return resp, http.StatusOK, errInfo
}

func (s *AccountUseCase) RejectSharing(ctx *gin.Context, dtoRequest *dtos.AccountGroupSharingAccept) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	// translate string to uuid value
	IDGroupSharingUUID, err := uuid.Parse(dtoRequest.IDGroupSharing)
	if err != nil {
		logrus.Error(err.Error())
	}

	// accept sharing
	err = s.repo.RejectSharing(IDGroupSharingUUID)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	resp := struct {
		Message string `json:"message"`
	}{
		Message: "reject group sharing successfully",
	}

	return resp, http.StatusOK, errInfo
}

func (s *AccountUseCase) RemoveSharing(ctx *gin.Context, dtoRequest *dtos.AccountGroupSharingRemove) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	// get email user
	usrEmail := ctx.MustGet("email").(string)

	// get customer id
	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	// email token is same as email target
	if usrEmail == dtoRequest.EmailAccount {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "target email same as user token")
		return struct{}{}, http.StatusBadRequest, errInfo
	}

	// get profile by email target for get account id
	dataProfile, err := s.repo.GetProfileByEmail(dtoRequest.EmailAccount)
	if err != nil {
		logrus.Error(err.Error())
	}

	// make sure invitation has accepted before
	dataFirst, dataSecond := s.repo.GroupSharingInfoByIDPersonalAccount(accountUUID, dataProfile.ID)

	if !dataFirst.IsAccepted || !dataSecond.IsAccepted {
		resp := struct {
			Message string `json:"message,omitempty"`
		}{
			Message: "no data group sharing for email account : " + dtoRequest.EmailAccount,
		}
		return resp, http.StatusBadRequest, []errorsinfo.Errors{}
	}

	// remove first data. it has been accepted
	if dataFirst.IsAccepted {
		err = s.repo.RemoveGroupSharingByID(dataFirst.ID)
		if err != nil {
			logrus.Error(err.Error())
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
			return struct{}{}, http.StatusInternalServerError, errInfo
		}
	}

	// remove second data. it has been accepted
	if dataSecond.IsAccepted {
		err = s.repo.RemoveGroupSharingByID(dataSecond.ID)
		if err != nil {
			logrus.Error(err.Error())
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
			return struct{}{}, http.StatusInternalServerError, errInfo
		}
	}

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	resp := struct {
		Message string `json:"message,omitempty"`
	}{
		Message: "remove account group sharing successfully",
	}

	return resp, http.StatusOK, errInfo
}

func (s *AccountUseCase) GroupSharingAccepted(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var dtoResponse []dtos.AccountShare

	// get customer id
	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	// get profile by email in token
	dataGroupSharingWithProfile, err := s.repo.GroupSharingList(accountUUID)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// reserve
	dataGroupSharingWithProfileReserve, err := s.repo.GroupSharingListReserve(accountUUID)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// if not found
	if len(dataGroupSharingWithProfile) == 0 && len(dataGroupSharingWithProfileReserve) == 0 {
		resp := struct {
			Message string `json:"message,omitempty"`
		}{
			Message: "this token has not shared with other accounts",
		}
		return resp, http.StatusNotFound, []errorsinfo.Errors{}
	}

	// clear error info
	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	// response for invited
	if len(dataGroupSharingWithProfileReserve) > 0 {
		// append
		for _, v := range dataGroupSharingWithProfileReserve {
			ImagePath := ""
			if v.ImagePath != "" {
				ImagePath = os.Getenv("APP_HOST") + "/v1/" + v.ImagePath
			}

			dtoResponse = append(dtoResponse, dtos.AccountShare{
				AccountShareDetail: dtos.AccountShareDetail{
					Name:      v.Name,
					Email:     v.Email,
					ImagePath: ImagePath,
					Type:      v.Type,
				},
				Status: strings.ToUpper(v.Status),
			})
		}

		return dtoResponse, http.StatusOK, errInfo
	}

	// response for invites
	if len(dataGroupSharingWithProfile) > 0 {
		// append
		for _, v := range dataGroupSharingWithProfile {
			ImagePath := ""
			if v.ImagePath != "" {
				ImagePath = os.Getenv("APP_HOST") + "/v1/" + v.ImagePath
			}

			dtoResponse = append(dtoResponse, dtos.AccountShare{
				AccountShareDetail: dtos.AccountShareDetail{
					Name:      v.Name,
					Email:     v.Email,
					ImagePath: ImagePath,
					Type:      v.Type,
				},
				Status: strings.ToUpper(v.Status),
			})
		}

		return dtoResponse, http.StatusOK, errInfo
	}

	return
}

func (s *AccountUseCase) GroupSharingPending(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var dtoResponse []dtos.AccountShare

	// get customer id
	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	if accountUUID == uuid.Nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "token contains invalid information")
		return struct{}{}, http.StatusUnauthorized, errInfo
	}

	// get profile by email in token
	dataGroupSharingWithProfile, err := s.repo.GroupSharingListPending(accountUUID)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// reverse
	dataGroupSharingWithProfileReverse, err := s.repo.GroupSharingListPendingReverse(accountUUID)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// if not found
	if len(dataGroupSharingWithProfile) == 0 && len(dataGroupSharingWithProfileReverse) == 0 {
		resp := struct {
			Message string `json:"message,omitempty"`
		}{
			Message: "no pending group sharing invitation",
		}
		return resp, http.StatusNotFound, []errorsinfo.Errors{}
	}

	// clear error info
	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	// response for sender
	if len(dataGroupSharingWithProfile) > 0 {
		for _, v := range dataGroupSharingWithProfile {
			// setup image
			ImagePath := ""
			if v.ImagePath != "" {
				ImagePath = os.Getenv("APP_HOST") + "/v1/" + v.ImagePath
			}

			dtoResponse = append(dtoResponse, dtos.AccountShare{
				AccountShareDetail: dtos.AccountShareDetail{
					Name:      v.Name,
					Email:     v.Email,
					ImagePath: ImagePath,
					Type:      v.Type,
				},
				Status:   strings.ToUpper(v.Status),
				ActionID: dtos.AccountActionID{},
			})
		}
		return dtoResponse, http.StatusOK, errInfo
	}

	// response for receipt
	if len(dataGroupSharingWithProfileReverse) > 0 {
		for _, v := range dataGroupSharingWithProfileReverse {
			// setup image
			ImagePath := ""
			if v.ImagePath != "" {
				ImagePath = os.Getenv("APP_HOST") + "/v1/" + v.ImagePath
			}

			dtoResponse = append(dtoResponse, dtos.AccountShare{
				AccountShareDetail: dtos.AccountShareDetail{
					Name:      v.Name,
					Email:     v.Email,
					ImagePath: ImagePath,
					Type:      v.Type,
				},
				Status: strings.ToUpper(v.Status),
				ActionID: dtos.AccountActionID{
					IDGroupSharing: v.ID.String(),
				},
			})

		}
		return dtoResponse, http.StatusOK, errInfo
	}

	return
}

func (s *AccountUseCase) VerifyOTP(ctx *gin.Context, request *dtos.AccountOTPVerify) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	// getting profile
	dataProfile, err := s.repo.GetProfileByEmail(request.EmailAccount)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusBadRequest, errInfo
	}

	// validate result
	if dataProfile.Email == "" || dataProfile.ID == uuid.Nil {
		resp := struct {
			Message string `json:"message,omitempty"`
		}{
			Message: "email account not registered in system",
		}
		return resp, http.StatusBadRequest, errInfo
	}

	// getting forgot pass data
	dataForgotPassword, err := s.repo.ForgotPasswordData(dataProfile.ID)
	if err != nil {
		logrus.Error(err.Error())
	}

	// validate result
	if dataForgotPassword.OTPCode == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "no data found for this otp : "+request.OTPCode)
		return struct{}{}, http.StatusBadRequest, errInfo
	}

	if dataForgotPassword.OTPCode != request.OTPCode {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: "otp code in valid. please resend new otp code",
		}
		return resp, http.StatusBadRequest, errInfo
	}

	// duration expired
	duration := dataForgotPassword.Expired.Sub(dataForgotPassword.CreatedAt)
	if duration <= 0 {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: "otp code has expired",
		}

		return resp, http.StatusBadRequest, []errorsinfo.Errors{}
	}

	// update for verified
	err = s.repo.UpdateForgotPassword(dataForgotPassword.ID)
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// generate token jwt
	token, _, err := token.JWTBuilder(dataProfile.Email, dataProfile.UserRoles)
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// clear if empty
	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	resp := struct {
		Message string `json:"message"`
		Token   string `json:"token"`
	}{
		Message: "otp has been verified",
		Token:   token,
	}
	return resp, http.StatusOK, errInfo
}

func (s *AccountUseCase) ChangePasswordForgot(ctx *gin.Context, request *dtos.AccountChangeForgotPassword) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	// get customer id
	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	// save
	hashPassword := password.Generate(request.NewPassword)

	err := s.repo.ChangePassword(accountUUID, hashPassword)
	if err != nil {
		logrus.Error(err.Error())
	}

	resp := struct {
		Message string `json:"message"`
	}{
		Message: "change password success",
	}

	// if empty
	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	return resp, http.StatusOK, errInfo
}

func (s *AccountUseCase) DeleteAccount(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	// delete wallet
	go s.repo.DeleteAccountWallet(accountUUID)

	// delete transaction
	go s.repo.DeleteAccountTransaction(accountUUID)

	// delete categories category editable
	go s.repo.DeleteAccountMasterExpenseCategory(accountUUID)

	// delete categories sub category editable
	go s.repo.DeleteAccountMasterSubExpenseCategory(accountUUID)

	// delete categories income editable
	go s.repo.DeleteAccountMasterIncomeCategory(accountUUID)

	// delete budget
	go s.repo.DeleteAccountBudget(accountUUID)

	// delete group sharing
	go s.repo.DeleteAccountGroupSharing(accountUUID)

	// delete account subscription
	go s.repo.DeleteAccountSubscription(accountUUID)

	// delete account withdraw
	go s.repo.DeleteAccountWithdraw(accountUUID)

	// delete account authorization
	go s.repo.DeleteAccountAuthorization(accountUUID)

	// delete account for personal account
	go s.repo.DeleteAccountPersonalAccount(accountUUID)

	// if no error
	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	resp := struct {
		Message string `json:"message"`
	}{
		Message: "delete account successfully",
	}

	return resp, http.StatusOK, errInfo
}
