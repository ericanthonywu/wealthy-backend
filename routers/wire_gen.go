// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package routers

import (
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/accounts"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/budgets"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/images"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/internals"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/masters"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/notifications"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/payments"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/referrals"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/statistics"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/tracks"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/transactions"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/wallets"
	"github.com/semicolon-indonesia/wealthy-backend/infrastructures/subsriptions"
	"gorm.io/gorm"
)

// Injectors from wire.go:

func Accounts(db *gorm.DB) *accounts.AccountController {
	accountRepository := accounts.NewAccountRepository(db)
	accountUseCase := accounts.NewAccountUseCase(accountRepository)
	accountController := accounts.NewAccountController(accountUseCase)
	return accountController
}

func Wallets(db *gorm.DB) *wallets.WalletController {
	walletRepository := wallets.NewWalletRepository(db)
	walletUseCase := wallets.NewWalletUseCase(walletRepository)
	walletController := wallets.NewWalletController(walletUseCase)
	return walletController
}

func Masters(db *gorm.DB) *masters.MasterController {
	masterRepository := masters.NewMasterRepository(db)
	masterUseCase := masters.NewMasterUseCase(masterRepository)
	masterController := masters.NewMasterController(masterUseCase)
	return masterController
}

func Transactions(db *gorm.DB) *transactions.TransactionController {
	transactionRepository := transactions.NewTransactionRepository(db)
	transactionUseCase := transactions.NewTransactionUseCase(transactionRepository)
	transactionController := transactions.NewTransactionController(transactionUseCase)
	return transactionController
}

func Budgets(db *gorm.DB) *budgets.BudgetController {
	budgetRepository := budgets.NewBudgetRepository(db)
	budgetUseCase := budgets.NewBudgetUseCase(budgetRepository)
	budgetController := budgets.NewBudgetController(budgetUseCase)
	return budgetController
}

func Statistics(db *gorm.DB) *statistics.StatisticController {
	statisticRepository := statistics.NewStatisticRepository(db)
	statisticUseCase := statistics.NewStatisticUseCase(statisticRepository)
	statisticController := statistics.NewStatisticController(statisticUseCase)
	return statisticController
}

func Images(db *gorm.DB) *images.ShowImageController {
	showImageRepository := images.NewShowImageRepository(db)
	showImageUseCase := images.NewShowImageUseCase(showImageRepository)
	showImageController := images.NewShowImageController(showImageUseCase)
	return showImageController
}

func Internals(db *gorm.DB) *internals.InternalController {
	internalRepository := internals.NewInternalRepository(db)
	internalUseCase := internals.NewInternalUseCase(internalRepository)
	internalController := internals.NewInternalController(internalUseCase)
	return internalController
}

func Subscriptions(db *gorm.DB) *subsriptions.SubscriptionController {
	subscriptionRepository := subsriptions.NewSubscriptionRepository(db)
	subscriptionUseCase := subsriptions.NewSubscriptionUseCase(subscriptionRepository)
	subscriptionController := subsriptions.NewSubscriptionController(subscriptionUseCase)
	return subscriptionController
}

func Referrals(db *gorm.DB) *referrals.ReferralController {
	referralRepository := referrals.NewReferralRepository(db)
	referralUseCase := referrals.NewReferralUseCase(referralRepository)
	referralController := referrals.NewReferralController(referralUseCase)
	return referralController
}

func Payments(db *gorm.DB) *payments.PaymentController {
	paymentRepository := payments.NewPaymentRepository(db)
	paymentUseCase := payments.NewPaymentUseCase(paymentRepository)
	paymentController := payments.NewPaymentController(paymentUseCase)
	return paymentController
}

func Tracks(db *gorm.DB) *tracks.TrackController {
	trackRepository := tracks.NewTrackRepository(db)
	trackUseCase := tracks.NewTrackUseCase(trackRepository)
	trackController := tracks.NewTrackController(trackUseCase)
	return trackController
}

func Notifications(db *gorm.DB) *notifications.NotificationController {
	notificationRepository := notifications.NewNotificationRepository(db)
	notificationUseCase := notifications.NewNotificationUseCase(notificationRepository)
	notificationController := notifications.NewNotificationController(notificationUseCase)
	return notificationController
}
