//go:build wireinject
// +build wireinject

package routers

import (
	"github.com/google/wire"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/accounts"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/budgets"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/images"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/internals"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/masters"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/statistics"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/transactions"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/wallets"
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