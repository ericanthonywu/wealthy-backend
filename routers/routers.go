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
				typeGroup.GET("/invest", tokenSignature(), master.Invest)
				typeGroup.GET("/broker", tokenSignature(), master.Broker)
			}

			categoriesGroup := masterGroup.Group("/categories")
			{
				categoriesGroup.GET("/income", tokenSignature(), master.IncomeType)
				categoriesGroup.GET("/expense", tokenSignature(), master.ExpenseType)
			}

			transactionGroup := masterGroup.Group("/transactions")
			{
				transactionGroup.GET("/priority", tokenSignature(), master.TransactionPriority)
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
			transactionGroup.GET("/income-spending", tokenSignature(), transaction.IncomeSpending)
			transactionGroup.GET("/investment", tokenSignature(), transaction.Investment)
			transactionGroup.GET("/notes", tokenSignature(), transaction.ByNotes)

			transactionHistory := transactionGroup.Group("/history")
			{
				transactionHistory.GET("/expense", tokenSignature(), transaction.ExpenseTransactionHistory)
				transactionHistory.GET("/income", tokenSignature(), transaction.IncomeTransactionHistory)
				transactionHistory.GET("/transfer", tokenSignature(), transaction.TransferTransactionHistory)
				transactionHistory.GET("/invest", tokenSignature(), transaction.InvestTransactionHistory)
			}
		}
	}
}
