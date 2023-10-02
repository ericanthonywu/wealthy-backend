package routers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func API(router *gin.RouterGroup, db *gorm.DB) {
	master := Masters(db)
	account := Accounts(db)
	wallet := Wallets(db)
	transaction := Transactions(db)

	v1group := router.Group("/v1")
	{
		masterGroup := v1group.Group("/masters")
		{
			typeGroup := masterGroup.Group("/types")
			{
				typeGroup.GET("/transaction", tokenSignature(), master.TransactionType)
				typeGroup.GET("/reksadana", tokenSignature(), master.ReksadanaType)
				typeGroup.GET("/wallet", tokenSignature(), master.WalletType)
			}

			categoriesGroup := masterGroup.Group("/categories")
			{
				categoriesGroup.GET("/income", tokenSignature(), master.IncomeType)
				categoriesGroup.GET("/expense", tokenSignature(), master.ExpenseType)
			}
		}

		accountGroup := v1group.Group("/accounts")
		{
			accountGroup.POST("/signup", account.SignUp)
			accountGroup.POST("/signin", account.SignIn)
			// accountGroup.POST("/profile", account.SignIn) // get profile

		}

		walletGroup := v1group.Group("/wallets")
		{
			walletGroup.POST("/", tokenSignature(), wallet.Add)
			walletGroup.GET("/", tokenSignature(), wallet.List)
			walletGroup.PUT("/amount/:id-wallet", tokenSignature(), wallet.UpdateAmount)

		}

		transactionGroup := v1group.Group("/transactions")
		{
			transactionGroup.POST("/", tokenSignature(), transaction.Add)
		}
	}
}
