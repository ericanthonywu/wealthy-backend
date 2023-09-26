package wallets

import "github.com/gin-gonic/gin"

type (
	WalletController struct {
		useCase IWalletUseCase
	}

	IWalletController interface {
		Add(ctx *gin.Context)
	}
)

func NewWalletController(useCase IWalletUseCase) *WalletController {
	return &WalletController{useCase: useCase}
}

func (c *WalletController) Add(ctx *gin.Context) {

}
