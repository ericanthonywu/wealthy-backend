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
	budget := Budgets(db)
	statistic := Statistics(db)
	images := Images(db)
	internals := Internals(db)
	subscriptions := Subscriptions(db)

	v1group := router.Group("/v1")
	{
		masterGroup := v1group.Group("/masters", tokenSignature())
		{
			typeGroup := masterGroup.Group("/types")
			{
				typeGroup.GET("/transaction", master.TransactionType)
				typeGroup.GET("/reksadana", master.ReksadanaType)
				typeGroup.GET("/wallet", master.WalletType)
				typeGroup.GET("/invest", master.Invest)
				typeGroup.GET("/broker", master.Broker)
			}

			categoriesGroup := masterGroup.Group("/categories")
			{
				categoriesGroup.GET("/income", master.IncomeType)
				categoriesGroup.GET("/expense", master.ExpenseType)
				categoriesGroup.GET("/sub-expense/:expense-id", master.SubExpenseCategories)
			}

			transactionGroup := masterGroup.Group("/transactions")
			{
				transactionGroup.GET("/priority", master.TransactionPriority)
			}

			genderGroup := masterGroup.Group("/genders")
			{
				genderGroup.GET("", master.Gender)
			}

		}

		accountGroup := v1group.Group("/accounts")
		{
			accountGroup.POST("/signup", account.SignUp)
			accountGroup.POST("/signin", account.SignIn)
			accountGroup.PATCH("/profiles/:id", tokenSignature(), account.UpdateProfile)
			accountGroup.GET("/profiles", tokenSignature(), account.GetProfile)

			accountProfileGroup := accountGroup.Group("/profiles", tokenSignature())
			{
				accountProfileGroup.POST("/avatar", account.SetAvatar)
				accountProfileGroup.DELETE("/avatar/:customer-id", account.RemoveAvatar)
			}

			accountPasswordGroup := accountGroup.Group("/password", tokenSignature())
			{
				accountPasswordGroup.POST("/change/:id", account.ChangePassword)
			}

			accountReferral := accountGroup.Group("/referrals", tokenSignature())
			{
				accountReferral.POST("/validate", account.ValidateRefCode)
			}
		}

		budgetGroup := v1group.Group("/budgets", tokenSignature())
		{
			limitGroup := budgetGroup.Group("/limits")
			{
				limitGroup.POST("", budget.Limit)
				limitGroup.GET("", budget.AllLimit)
			}

			detailsGroup := budgetGroup.Group("/details")
			{
				detailsGroup.GET("/category", budget.Trends)
				detailsGroup.GET("/category-latest", budget.LatestMonths)
			}
			budgetGroup.GET("/overview", budget.Overview)
		}

		statisticGroup := v1group.Group("/statistics", tokenSignature())
		{
			statisticGroup.GET("/trends", statistic.Trend)

			analyticsGroup := statisticGroup.Group("/analytics")
			{
				analyticsGroup.GET("/trend", statistic.AnalyticsTrend)
			}

			transactionStatisticGroup := statisticGroup.Group("/transactions")
			{
				transactionStatisticGroup.GET("/weekly", statistic.Weekly)
				transactionStatisticGroup.GET("/summary", statistic.Summary)
				transactionStatisticGroup.GET("/priority", statistic.TransactionPriority)

				transactionsDetailGroup := transactionStatisticGroup.Group("/details")
				{
					transactionsDetailGroup.GET("/expense", statistic.ExpenseDetail)
					transactionsDetailGroup.GET("/sub-expense", statistic.SubExpenseDetail)
				}
			}
		}

		transactionGroup := v1group.Group("/transactions", tokenSignature())
		{
			transactionGroup.POST("", transaction.Add)
			transactionGroup.GET("/income-spending", transaction.IncomeSpending)
			transactionGroup.GET("/investment", transaction.Investment)
			transactionGroup.GET("/notes", transaction.ByNotes)

			transactionHistory := transactionGroup.Group("/history")
			{
				transactionHistory.GET("/expense", transaction.ExpenseTransactionHistory)
				transactionHistory.GET("/income", transaction.IncomeTransactionHistory)
				transactionHistory.GET("/transfer", transaction.TransferTransactionHistory)
				transactionHistory.GET("/invest", transaction.InvestTransactionHistory)
				transactionHistory.GET("/travel", transaction.TravelTransactionHistory)
			}
		}

		walletGroup := v1group.Group("/wallets", tokenSignature())
		{
			walletGroup.POST("", wallet.Add)
			walletGroup.GET("", wallet.List)
			walletGroup.PUT("/amount/:id-wallet", wallet.UpdateAmount)

		}

		subscriptionGroup := v1group.Group("/subscriptions", tokenSignature())
		{
			subscriptionGroup.GET("/plan", subscriptions.Plan)
			subscriptionGroup.GET("/faq", subscriptions.FAQ)
		}

		internalGroup := v1group.Group("/internals")
		{
			internalTransactionGroup := internalGroup.Group("/transactions")
			{
				internalTransactionGroup.GET("/notes", internals.TransactionNotes)
			}
		}

		imageGroup := v1group.Group("/images")
		{
			avatarGroup := imageGroup.Group("/avatar")
			{
				avatarGroup.GET("/:filename", images.Avatar)
			}

			travelGroup := imageGroup.Group("/travel")
			{
				travelGroup.GET("/:filename", images.Travel)
			}
		}
	}

}