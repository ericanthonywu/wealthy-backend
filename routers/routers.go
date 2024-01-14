package routers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func API(router *gin.RouterGroup, db *gorm.DB) {

	// version 1
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
	investment := Investments(db)

	// version 2
	categoriesV2 := CategoriesV2(db)
	walletsV2 := WalletsV2(db)

	v1group := router.Group("/v1")
	{
		masterGroup := v1group.Group("/masters", tokenSignature(), accountType(), betaVersion())
		{
			typeGroup := masterGroup.Group("/types")
			{
				typeGroup.GET("/transaction", master.TransactionType)
				typeGroup.GET("/reksadana", master.ReksadanaType)
				typeGroup.GET("/wallet", master.WalletType)
				typeGroup.GET("/invest", master.Invest)
				typeGroup.GET("/broker", master.Broker)
				typeGroup.GET("/stock-code", master.StockCode)
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

			subscriptionGroup := masterGroup.Group("/subscriptions")
			{
				subscriptionGroup.GET("/price", master.Price)
			}
		}

		accountGroup := v1group.Group("/accounts")
		{
			accountGroup.POST("/signup", account.SignUp)
			accountGroup.POST("/signin", account.SignIn)
			accountGroup.PATCH("/profiles", tokenSignature(), accountType(), account.UpdateProfile)
			accountGroup.GET("/profiles", tokenSignature(), accountType(), betaVersion(), account.GetProfile)
			accountGroup.DELETE("", tokenSignature(), accountType(), account.DeleteAccount)

			accountProfileGroup := accountGroup.Group("/profiles", tokenSignature(), accountType())
			{
				accountProfileGroup.POST("/avatar", account.SetAvatar)
				accountProfileGroup.DELETE("/avatar", account.RemoveAvatar)
			}

			accountPasswordGroup := accountGroup.Group("/password")
			{
				forgotGroup := accountPasswordGroup.Group("/forgot")
				{
					forgotGroup.POST("", account.ForgotPassword)
					forgotGroup.POST("/verify", account.VerifyOTP)
					forgotGroup.POST("/change-password", tokenSignature(), accountType(), account.ChangePasswordForgot)
				}
				accountPasswordGroup.POST("/change", tokenSignature(), accountType(), account.ChangePassword)
			}

			accountReferral := accountGroup.Group("/referrals")
			{
				accountReferral.POST("/validate", account.ValidateRefCode)
			}

			sharingGroup := accountGroup.Group("/shares", tokenSignature(), accountType(), betaVersion())
			{
				sharingGroup.POST("/search", account.SearchAccount)
				sharingGroup.POST("/invite", account.InviteSharing)
				sharingGroup.POST("/accept", account.AcceptSharing)
				sharingGroup.POST("/reject", account.RejectSharing)
				sharingGroup.POST("/remove", account.RemoveSharing)

				listGroup := sharingGroup.Group("/list")
				{
					listGroup.GET("/accepted", account.GroupSharingAccepted)
					listGroup.GET("/pending", account.GroupSharingPending)
				}
			}
		}

		budgetGroup := v1group.Group("/budgets", tokenSignature(), accountType(), betaVersion())
		{
			limitGroup := budgetGroup.Group("/limits")
			{
				limitGroup.POST("", budget.Limit)
				limitGroup.GET("", budget.AllLimit)
				limitGroup.PATCH("/travels/:id-travel", budget.UpdateTravelInfo)
			}

			detailsGroup := budgetGroup.Group("/details")
			{
				detailsGroup.GET("/category", budget.Trends)
				detailsGroup.GET("/category-latest", budget.LatestMonths)
				detailsGroup.GET("/travel", budget.Travels)
			}
			budgetGroup.GET("/overview", budget.Overview)
		}

		statisticGroup := v1group.Group("/statistics", tokenSignature(), accountType(), betaVersion())
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

			investmentGroup := statisticGroup.Group("/investments")
			{
				investmentGroup.GET("/top-three", statistic.TopThreeInvestment)
			}
		}

		transactionGroup := v1group.Group("/transactions", tokenSignature(), accountType(), betaVersion())
		{
			transactionGroup.POST("", transaction.Add)
			transactionGroup.GET("/income-spending", transaction.IncomeSpending)
			transactionGroup.GET("/investments", transaction.Investment)
			transactionGroup.GET("/notes", transaction.ByNotes)
			transactionGroup.GET("/cash-flow", transaction.CashFlow)

			transactionHistory := transactionGroup.Group("/history")
			{
				transactionHistory.GET("/expense", transaction.ExpenseTransactionHistory)
				transactionHistory.GET("/income", transaction.IncomeTransactionHistory)
				transactionHistory.GET("/transfer", transaction.TransferTransactionHistory)
				transactionHistory.GET("/invest", transaction.InvestTransactionHistory)
				transactionHistory.GET("/travel", transaction.TravelTransactionHistory)
			}

			investmentTransaction := transactionGroup.Group("/records")
			{
				investmentTransaction.POST("/investments", transaction.AddInvestmentTransaction)
			}

			suggestionGroup := transactionGroup.Group("/suggestions")
			{
				suggestionGroup.GET("/notes", transaction.Suggestion)
			}

			walletGroup := transactionGroup.Group("/wallets")
			{
				walletGroup.GET("/non-investment", transaction.WalletNonInvestment)
				walletGroup.GET("/investment", transaction.WalletInvestment)
			}
		}

		investmentGroup := v1group.Group("/investments", tokenSignature(), accountType(), betaVersion())
		{
			investmentGroup.GET("/portfolio", investment.Portfolio)
			investmentGroup.GET("/gain-loss", investment.GainLoss)
		}

		walletGroup := v1group.Group("/wallets", tokenSignature(), accountType(), betaVersion())
		{
			walletGroup.POST("", wallet.Add)
			walletGroup.GET("", wallet.List)
			walletGroup.PATCH("/amount/:id-wallet", wallet.UpdateAmount)
		}

		subscriptionGroup := v1group.Group("/subscriptions", tokenSignature())
		{
			subscriptionGroup.GET("/faq", subscriptions.FAQ)
		}

		internalGroup := v1group.Group("/internals")
		{
			internalTransactionGroup := internalGroup.Group("/transactions")
			{
				internalTransactionGroup.GET("/notes", internals.TransactionNotes)
			}
		}

		referralGroup := v1group.Group("/referrals", tokenSignature(), accountType(), betaVersion())
		{
			referralGroup.GET("/statistics", referrals.Statistic)
			referralGroup.GET("/list", referrals.List)

			earnGroup := referralGroup.Group("/earns")
			{
				earnGroup.GET("", referrals.Earn)
				earnGroup.POST("/withdraw", referrals.Withdraw)
			}
		}

		paymentGroup := v1group.Group("/payments")
		{
			paymentGroup.POST("/subscriptions", tokenSignature(), payments.Subscriptions)
			webhookGroup := paymentGroup.Group("/webhooks")
			{
				webhookGroup.POST("/midtrans", payments.MidtransWebhook)
			}
		}

		trackGroup := v1group.Group("/tracks", tokenSignature())
		{
			trackGroup.POST("/screen-time", tracks.ScreenTime)
		}

		notificationGroup := v1group.Group("/notifications", tokenSignature(), accountType(), betaVersion())
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

	v2group := router.Group("/v2")
	{
		categoriesGroup := v2group.Group("/categories")
		{
			categoriesGroup.GET("/expense", tokenSignature(), accountType(), categoriesV2.GetCategoriesExpenseList)
			categoriesGroup.GET("/income", tokenSignature(), accountType(), categoriesV2.GetCategoriesIncomeList)
		}

		walletGroup := v2group.Group("/wallets", tokenSignature(), accountType(), betaVersion())
		{
			walletGroup.POST("", walletsV2.NewWallet)
			walletGroup.GET("", walletsV2.GetAllWallets)
			walletGroup.PATCH("/amount/:id-wallet", walletsV2.UpdateWallet)
		}

	}
}
