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
	referrals := Referrals(db)
	payments := Payments(db)
	tracks := Tracks(db)
	notification := Notifications(db)

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

				personalGroup := categoriesGroup.Group("/personals")
				{
					personalGroup.GET("/income", master.PersonalIncomeCategory)
					personalGroup.GET("/expense", master.PersonalExpenseCategory)
					personalGroup.GET("/sub-expense/:expense-id", master.PersonalExpenseSubCategory)

					renameGroup := personalGroup.Group("/renames")
					{
						renameGroup.PUT("/income/:id", master.RenameIncomeCategory)
						renameGroup.PUT("/expense/:id", master.RenameExpenseCategory)
						renameGroup.PUT("/sub-expense/:id", master.RenameSubExpenseCategory)
					}

					addGroup := personalGroup.Group("/adds")
					{
						addGroup.POST("/income", master.AddIncomeCategory)
						addGroup.POST("/expense", master.AddExpenseCategory)
						addGroup.POST("/sub-expense", master.AddSubExpenseCategory)
					}
				}
			}

			transactionGroup := masterGroup.Group("/transactions")
			{
				transactionGroup.GET("/priority", master.TransactionPriority)
			}

			genderGroup := masterGroup.Group("/genders")
			{
				genderGroup.GET("", master.Gender)
			}

			currencyGroup := masterGroup.Group("/currency")
			{
				currencyGroup.GET("/exchange", master.Exchange)
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

			accountPasswordGroup := accountGroup.Group("/password")
			{
				accountPasswordGroup.POST("/change", tokenSignature(), account.ChangePassword)
				accountPasswordGroup.POST("/forgot", account.ForgotPassword)
			}

			otpGroup := accountGroup.Group("/otp")
			{
				otpGroup.POST("/verify", account.VerifyOTP)
			}

			accountReferral := accountGroup.Group("/referrals")
			{
				accountReferral.POST("/validate", account.ValidateRefCode)
			}

			sharingGroup := accountGroup.Group("/shares", tokenSignature())
			{
				sharingGroup.POST("/search", account.SearchAccount)
				sharingGroup.POST("/invite", account.InviteSharing)
				sharingGroup.POST("/accept", account.AcceptSharing)
				sharingGroup.POST("/reject", account.RejectSharing)
				sharingGroup.POST("/remove", account.RemoveSharing)
				sharingGroup.GET("/list", account.ListGroupSharing)
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
				detailsGroup.GET("/travel", budget.Travels)
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

			suggestionGroup := transactionGroup.Group("/suggestions")
			{
				suggestionGroup.GET("/notes", transaction.Suggestion)
			}
		}

		walletGroup := v1group.Group("/wallets", tokenSignature())
		{
			walletGroup.POST("/", wallet.Add)
			walletGroup.GET("", wallet.List)
			walletGroup.PATCH("/amount/:id-wallet", wallet.UpdateAmount)

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

		referralGroup := v1group.Group("/referrals", tokenSignature())
		{
			referralGroup.GET("/statistics", referrals.Statistic)
			referralGroup.GET("/list", referrals.List)
		}

		paymentGroup := v1group.Group("/payments", tokenSignature())
		{
			paymentGroup.POST("/subscriptions", payments.Subscriptions)
			webhookGroup := paymentGroup.Group("/webhooks")
			{
				webhookGroup.POST("/midtrans", payments.MidtransWebhook)
			}
		}

		trackGroup := v1group.Group("/tracks", tokenSignature())
		{
			trackGroup.POST("/screen-time", tracks.ScreenTime)
		}

		notificationGroup := v1group.Group("/notifications", tokenSignature())
		{
			notificationGroup.GET("", notification.GetNotification)
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