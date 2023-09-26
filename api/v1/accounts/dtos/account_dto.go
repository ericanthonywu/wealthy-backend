package dtos

type (
	AccountSignUpRequest struct {
		Username string `json:"username"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		RefCode  string `json:"ref-code,omitempty"`
	}

	AccountSignUpResponse struct {
		Username string `json:"username"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Role     string `json:"role"`
	}

	AccountSignInRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	AccountSignInResponse struct {
		Email string `json:"email"`
		Role  string `json:"role"`
		Token string `json:"token"`
	}

	AccountSignOutRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
)
