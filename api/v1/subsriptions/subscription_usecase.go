package subsriptions

import (
	"github.com/gin-gonic/gin"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"net/http"
)

type (
	SubscriptionUseCase struct {
		repo ISubscriptionRepository
	}

	ISubscriptionUseCase interface {
		Plan(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
	}
)

func NewSubscriptionUseCase(repo ISubscriptionRepository) *SubscriptionUseCase {
	return &SubscriptionUseCase{repo: repo}
}

func (s *SubscriptionUseCase) Plan(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	data := s.repo.Plan()

	if len(data) == 0 {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "no data")
		return data, http.StatusNotFound, errInfo
	}

	return data, http.StatusOK, errInfo
}