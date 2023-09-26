//go:build wireinject
// +build wireinject

package routers

import (
	"github.com/google/wire"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/accounts"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/masters"
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
