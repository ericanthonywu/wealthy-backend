package subsriptions

import (
	"github.com/gin-gonic/gin"
	"github.com/wealthy-app/wealthy-backend/utils/errorsinfo"
	"net/http"
)

type (
	SubscriptionUseCase struct {
		repo ISubscriptionRepository
	}

	ISubscriptionUseCase interface {
		FAQ(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
	}
)

func NewSubscriptionUseCase(repo ISubscriptionRepository) *SubscriptionUseCase {
	return &SubscriptionUseCase{repo: repo}
}

func (s *SubscriptionUseCase) FAQ(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	data := s.repo.FAQ()

	if len(data) == 0 {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "no data")
		return data, http.StatusNotFound, errInfo
	}

	return data, http.StatusOK, errInfo
}