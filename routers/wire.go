//go:build wireinject
// +build wireinject

package routers

import (
	"github.com/google/wire"
	"github.com/wealthy-app/wealthy-backend/api/v1/accounts"
	"github.com/wealthy-app/wealthy-backend/api/v1/budgets"
	"github.com/wealthy-app/wealthy-backend/api/v1/images"
	"github.com/wealthy-app/wealthy-backend/api/v1/internals"
	"github.com/wealthy-app/wealthy-backend/api/v1/investments"
	"github.com/wealthy-app/wealthy-backend/api/v1/masters"
	"github.com/wealthy-app/wealthy-backend/api/v1/notifications"
	"github.com/wealthy-app/wealthy-backend/api/v1/payments"
	"github.com/wealthy-app/wealthy-backend/api/v1/referrals"
	"github.com/wealthy-app/wealthy-backend/api/v1/statistics"
	"github.com/wealthy-app/wealthy-backend/api/v1/subsriptions"
	"github.com/wealthy-app/wealthy-backend/api/v1/tracks"
	"github.com/wealthy-app/wealthy-backend/api/v1/transactions"
	"github.com/wealthy-app/wealthy-backend/api/v1/wallets"
	"gorm.io/gorm"
)

func Accounts(db *gorm.DB) *accounts.AccountController {
	panic(wire.Build(wire.NewSet(
		accounts.NewAccountRepository,
		accounts.NewAccountUseCase,
		accounts.NewAccountController,
		wire.Bind(new(accounts.IAccountRepository), new(*accounts.AccountRepository)),
		wire.Bind(new(accounts.IAccountUseCase), new(*accounts.AccountUseCase)),
		wire.Bind(new(accounts.IAccountController), new(*accounts.AccountController)),
	)))
	return &accounts.AccountController{}
}

func Wallets(db *gorm.DB) *wallets.WalletController {
	panic(wire.Build(wire.NewSet(
		wallets.NewWalletRepository,
		wallets.NewWalletUseCase,
		wallets.NewWalletController,
		wire.Bind(new(wallets.IWalletRepository), new(*wallets.WalletRepository)),
		wire.Bind(new(wallets.IWalletUseCase), new(*wallets.WalletUseCase)),
		wire.Bind(new(wallets.IWalletController), new(*wallets.WalletController)),
	)))
	return &wallets.WalletController{}
}

func Masters(db *gorm.DB) *masters.MasterController {
	panic(wire.Build(wire.NewSet(
		masters.NewMasterController,
		masters.NewMasterUseCase,
		masters.NewMasterRepository,
		wire.Bind(new(masters.IMasterController), new(*masters.MasterController)),
		wire.Bind(new(masters.IMasterUseCase), new(*masters.MasterUseCase)),
		wire.Bind(new(masters.IMasterRepository), new(*masters.MasterRepository)),
	)))
	return &masters.MasterController{}
}

func Transactions(db *gorm.DB) *transactions.TransactionController {
	panic(wire.Build(wire.NewSet(
		transactions.NewTransactionRepository,
		transactions.NewTransactionUseCase,
		transactions.NewTransactionController,
		wire.Bind(new(transactions.ITransactionController), new(*transactions.TransactionController)),
		wire.Bind(new(transactions.ITransactionUseCase), new(*transactions.TransactionUseCase)),
		wire.Bind(new(transactions.ITransactionRepository), new(*transactions.TransactionRepository)),
	)))
	return &transactions.TransactionController{}
}

func Budgets(db *gorm.DB) *budgets.BudgetController {
	panic(wire.Build(wire.NewSet(
		budgets.NewBudgetRepository,
		budgets.NewBudgetUseCase,
		budgets.NewBudgetController,
		wire.Bind(new(budgets.IBudgetController), new(*budgets.BudgetController)),
		wire.Bind(new(budgets.IBudgetUseCase), new(*budgets.BudgetUseCase)),
		wire.Bind(new(budgets.IBudgetRepository), new(*budgets.BudgetRepository)),
	)))
	return &budgets.BudgetController{}
}

func Statistics(db *gorm.DB) *statistics.StatisticController {
	panic(wire.Build(wire.NewSet(
		statistics.NewStatisticRepository,
		statistics.NewStatisticUseCase,
		statistics.NewStatisticController,
		wire.Bind(new(statistics.IStatisticController), new(*statistics.StatisticController)),
		wire.Bind(new(statistics.IStatisticUseCase), new(*statistics.StatisticUseCase)),
		wire.Bind(new(statistics.IStatisticRepository), new(*statistics.StatisticRepository)),
	)))
	return &statistics.StatisticController{}
}

func Images(db *gorm.DB) *images.ShowImageController {
	panic(wire.Build(wire.NewSet(
		images.NewShowImageController,
		images.NewShowImageUseCase,
		images.NewShowImageRepository,
		wire.Bind(new(images.IShowImageController), new(*images.ShowImageController)),
		wire.Bind(new(images.IShowImageUseCase), new(*images.ShowImageUseCase)),
		wire.Bind(new(images.IShowImageRepository), new(*images.ShowImageRepository)),
	)))
	return &images.ShowImageController{}
}

func Internals(db *gorm.DB) *internals.InternalController {
	panic(wire.Build(wire.NewSet(
		internals.NewInternalController,
		internals.NewInternalUseCase,
		internals.NewInternalRepository,
		wire.Bind(new(internals.IInternalController), new(*internals.InternalController)),
		wire.Bind(new(internals.IInternalUseCase), new(*internals.InternalUseCase)),
		wire.Bind(new(internals.IInternalRepository), new(*internals.InternalRepository)),
	)))
	return &internals.InternalController{}
}

func Subscriptions(db *gorm.DB) *subsriptions.SubscriptionController {
	panic(wire.Build(wire.NewSet(
		subsriptions.NewSubscriptionController,
		subsriptions.NewSubscriptionUseCase,
		subsriptions.NewSubscriptionRepository,
		wire.Bind(new(subsriptions.ISubscriptionController), new(*subsriptions.SubscriptionController)),
		wire.Bind(new(subsriptions.ISubscriptionUseCase), new(*subsriptions.SubscriptionUseCase)),
		wire.Bind(new(subsriptions.ISubscriptionRepository), new(*subsriptions.SubscriptionRepository)),
	)))
	return &subsriptions.SubscriptionController{}
}

func Referrals(db *gorm.DB) *referrals.ReferralController {
	panic(wire.Build(wire.NewSet(
		referrals.NewReferralController,
		referrals.NewReferralUseCase,
		referrals.NewReferralRepository,
		wire.Bind(new(referrals.IReferralController), new(*referrals.ReferralController)),
		wire.Bind(new(referrals.IReferralUseCase), new(*referrals.ReferralUseCase)),
		wire.Bind(new(referrals.IReferralRepository), new(*referrals.ReferralRepository)),
	)))
	return &referrals.ReferralController{}
}

func Payments(db *gorm.DB) *payments.PaymentController {
	panic(wire.Build(wire.NewSet(
		payments.NewPaymentController,
		payments.NewPaymentUseCase,
		payments.NewPaymentRepository,
		wire.Bind(new(payments.IPaymentController), new(*payments.PaymentController)),
		wire.Bind(new(payments.IPaymentUseCase), new(*payments.PaymentUseCase)),
		wire.Bind(new(payments.IPaymentRepository), new(*payments.PaymentRepository)),
	)))
	return &payments.PaymentController{}
}

func Tracks(db *gorm.DB) *tracks.TrackController {
	panic(wire.Build(wire.NewSet(
		tracks.NewTrackController,
		tracks.NewTrackUseCase,
		tracks.NewTrackRepository,
		wire.Bind(new(tracks.ITrackController), new(*tracks.TrackController)),
		wire.Bind(new(tracks.ITrackUseCase), new(*tracks.TrackUseCase)),
		wire.Bind(new(tracks.ITrackRepository), new(*tracks.TrackRepository)),
	)))
	return &tracks.TrackController{}
}

func Notifications(db *gorm.DB) *notifications.NotificationController {
	panic(wire.Build(wire.NewSet(
		notifications.NewNotificationController,
		notifications.NewNotificationUseCase,
		notifications.NewNotificationRepository,
		wire.Bind(new(notifications.INotificationController), new(*notifications.NotificationController)),
		wire.Bind(new(notifications.INotificationUseCase), new(*notifications.NotificationUseCase)),
		wire.Bind(new(notifications.INotificationRepository), new(*notifications.NotificationRepository)),
	)))
	return &notifications.NotificationController{}
}

func Investments(db *gorm.DB) *investments.InvestmentController {
	panic(wire.Build(wire.NewSet(
		investments.NewInvestmentController,
		investments.NewInvestmentUseCase,
		investments.NewInvestmentRepository,
		wire.Bind(new(investments.IInvestmentController), new(*investments.InvestmentController)),
		wire.Bind(new(investments.IInvestmentUseCase), new(*investments.InvestmentUseCase)),
		wire.Bind(new(investments.IInvestmentRepository), new(*investments.InvestmentRepository)),
	)))
	return &investments.InvestmentController{}
}