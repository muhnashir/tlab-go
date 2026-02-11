package handler

import (
	"database/sql"
	"wallet-api/internal/repository"
	"wallet-api/internal/service"
)

type AllHandlers struct {
	AuthHandler   *AuthHandler
	WalletHandler *WalletHandler
}

// InitHandlers assembles the dependency graph: Repository -> Service -> Handler
func InitHandlers(db *sql.DB) *AllHandlers {
	// 1. Initialize Repositories
	userRepo := repository.NewMysqlUserRepository(db)
	walletRepo := repository.NewMysqlWalletRepository(db)
	transactionRepo := repository.NewMysqlTransactionRepository(db)

	// 2. Initialize Services
	authService := service.NewUserService(userRepo)
	walletService := service.NewWalletService(db, walletRepo, transactionRepo)

	// 3. Initialize Handlers
	authHandler := NewAuthHandler(authService)
	// Make sure NewWalletHandler accepts the concrete interface returned by NewWalletService
	walletHandler := NewWalletHandler(walletService)

	return &AllHandlers{
		AuthHandler:   authHandler,
		WalletHandler: walletHandler,
	}
}
