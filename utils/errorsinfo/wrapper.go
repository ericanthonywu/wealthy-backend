package errorsinfo

type Errors struct {
	ErrorCode        string `json:"error_code,omitempty"`
	ErrorDescription string `json:"error_description"`
}

func ErrorWrapper(errInfo []Errors, errorCode, errorDescription string) []Errors {
	errInfo = append(errInfo, Errors{
		ErrorCode:        errorCode,
		ErrorDescription: errorDescription,
	})
	return errInfo
}

func ErrorInfoWrapper(errInfo []string, errorDescription string) []string {
	return append(errInfo, errorDescription)
}

func ErrorWrapperArray(errInfo []string, errorDescription string) []string {
	return append(errInfo, errorDescription)
}