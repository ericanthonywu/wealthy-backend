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
		Customer CustomerSignIn `json:"customer,omitempty"`
		Account  Account        `json:"account_info,omitempty"`
		Token    string         `json:"token,omitempty"`
		TokenExp int64          `json:"token_exp,omitempty"`
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
		AccountCustomer AccountCustomerInfo `json:"customer_info,omitempty"`
		AccountAvatar   AccountAvatar       `json:"avatar,omitempty"`
		AccountDetail   AccountDetail       `json:"account_detail,omitempty"`
	}

	AccountCustomerInfo struct {
		Email     string    `json:"email"`
		Username  string    `json:"username"`
		Name      string    `json:"name"`
		DOB       string    `json:"date_of_birth"`
		ReferType string    `json:"refer_type"`
		ID        uuid.UUID `json:"id"`
		Gender    string    `json:"gender"`
	}

	AccountDetail struct {
		AccountType string `json:"account_type"`
		UserRoles   string `json:"user_roles"`
	}

	AccountAvatar struct {
		URL      string `json:"url"`
		FileName string `json:"filename"`
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
)