package dtos

import "github.com/google/uuid"

type (
	AccountSignUpRequest struct {
		Username         string `json:"username"`
		Name             string `json:"name"`
		Email            string `json:"email"`
		Password         string `json:"password"`
		RefCode          string `json:"referral_code"`
		RefCodeReference string `json:"referral_code_reference,omitempty"`
	}

	AccountSignUpResponse struct {
		Customer CustomerDetail `json:"customer_info,omitempty"`
		Account  Account        `json:"account_info,omitempty"`
	}

	CustomerDetail struct {
		ID       uuid.UUID `json:"id,omitempty"`
		Username string    `json:"username"`
		Name     string    `json:"name"`
		Email    string    `json:"email"`
	}

	Account struct {
		Role string `json:"role,omitempty"`
	}

	AccountSignInRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	AccountSignInResponse struct {
		Customer    CustomerSignIn `json:"customer,omitempty"`
		Account     Account        `json:"account_info,omitempty"`
		AccountType AccountType    `json:"account_type,omitempty"`
		Token       string         `json:"token,omitempty"`
		TokenExp    int64          `json:"token_exp,omitempty"`
	}

	AccountType struct {
		Type string `json:"type,omitempty"`
	}

	CustomerSignIn struct {
		CustomerID string `json:"customer_id,omitempty"`
		Email      string `json:"customer_email,omitempty"`
	}

	AccountSignOutRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	AccountSetProfileRequest struct {
		Name     string `json:"name" copier:"must"`
		Username string `json:"username" copier:"must"`
		DOB      string `json:"date_of_birth" copier:"must"`
		Gender   string `json:"id_master_gender" copier:"must"`
	}

	AccountProfile struct {
		AccountCustomer AccountCustomerInfo `json:"customer_info"`
		AccountGender   AccountGender       `json:"customer_gender"`
		AccountAvatar   AccountAvatar       `json:"customer_avatar"`
		AccountDetail   AccountDetail       `json:"customer_account"`
		AccountGeo      AccountGeo          `json:"customer_geographic"`
	}

	AccountCustomerInfo struct {
		ID        uuid.UUID `json:"customer_id"`
		Email     string    `json:"customer_email"`
		Username  string    `json:"customer_username"`
		Name      string    `json:"customer_name"`
		DOB       string    `json:"customer_dob"`
		ReferType string    `json:"customer_referral_code"`
	}

	AccountGender struct {
		ID    string `json:"gender_id"`
		Value string `json:"gender_value"`
	}

	AccountDetail struct {
		AccountType string `json:"account_type"`
		UserRoles   string `json:"account_roles"`
	}

	AccountAvatar struct {
		URL      string `json:"avatar_url"`
		FileName string `json:"avatar_filename"`
	}

	AccountGeo struct {
		Latitude  string `json:"latitude"`
		Longitude string `json:"longitude"`
	}

	AccountChangePassword struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	AccountRefCodeValidationRequest struct {
		RefCode string `json:"referral_code"`
	}

	AccountRefCodeValidationResponse struct {
		Available bool `json:"available"`
	}

	AccountAvatarRequest struct {
		ImageBase64 string `json:"image_base64"`
	}

	AccountAvatarResponse struct {
		Success bool `json:"success"`
	}

	AccountGroupSharing struct {
		EmailAccount string `json:"email_account"`
	}

	AccountGroupSharingAccept struct {
		IDSender    string `json:"id_group_sender"`
		IDRecipient string `json:"id_group_recipient"`
	}

	AccountGroupSharingRemove struct {
		EmailAccount string `json:"email_account"`
	}

	AccountSearch struct {
		ID uuid.UUID `json:"account_id"`
	}

	AccountForgotPasswordRequest struct {
		EmailAccount string `json:"email_account"`
	}

	AccountOTPVerify struct {
		OTPCode      string `json:"otp_code"`
		EmailAccount string `json:"email_account"`
	}

	AccountShare struct {
		AccountShareDetail AccountShareDetail `json:"account_detail"`
		Status             string             `json:"status"`
	}

	AccountShareDetail struct {
		Name      string `json:"account_name"`
		Email     string `json:"account_email"`
		ImagePath string `json:"account_avatar"`
		Type      string `json:"account_type"`
	}

	AccountChangeForgotPassword struct {
		NewPassword string `json:"new_password"`
	}

	NotificationPending struct {
		ID                      uuid.UUID                 `json:"id"`
		NotificationTitle       string                    `json:"notification_title"`
		NotificationDescription string                    `json:"notification_description"`
		IDPersonalAccounts      uuid.UUID                 `json:"id_personal_accounts"`
		IsRead                  bool                      `json:"is_read"`
		IDGroupSender           uuid.UUID                 `json:"id_group_sender"`
		IDGroupRecipient        uuid.UUID                 `json:"id_group_recipient"`
		AccountDetail           NotificationPendingDetail `json:"account_detail"`
		CreatedAt               string                    `json:"created_at"`
	}

	NotificationPendingDetail struct {
		AccountName  string `json:"account_name"`
		AccountImage string `json:"account_image"`
		AccountType  string `json:"account_type"`
	}
)