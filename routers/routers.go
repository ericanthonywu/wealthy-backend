package routers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func API(router *gin.RouterGroup, db *gorm.DB) {
	account := Accounts(db)
	wallet := Wallets(db)

	v1group := router.Group("/v1")
	{
		accountGroup := v1group.Group("/accounts")
		{
			accountGroup.POST("/signup", account.SignUp)
			accountGroup.POST("/signin", account.SignIn)
			// accountGroup.POST("/profile", account.SignIn) // get profile

		}

		walletGroup := v1group.Group("/wallet")
		{
			walletGroup.POST("/", tokenSignature(), wallet.Add)
			walletGroup.GET("/", tokenSignature(), wallet.List)
		}
	}

}
