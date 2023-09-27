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
			masterGroup.GET("/transaction-type", tokenSignature(), master.TransactionType)
			masterGroup.GET("/income-type", tokenSignature(), master.IncomeType)
			masterGroup.GET("/expense-type", tokenSignature(), master.ExpenseType)
			masterGroup.GET("/reksadana-type", tokenSignature(), master.ReksadanaType)
			masterGroup.GET("/wallet-type", tokenSignature(), master.WalletType)
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

		transactionGroup := v1group.GET("/transactions")
		{
			transactionGroup.POST("/expense", tokenSignature(), transaction.Expense)
			transactionGroup.POST("/income", tokenSignature(), transaction.Income)
			transactionGroup.POST("/transfer", tokenSignature(), transaction.Transfer)
			transactionGroup.POST("/invest", tokenSignature(), transaction.Invest)
		}
	}
}
