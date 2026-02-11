package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
	"wallet-api/internal/domain"

	"github.com/doug-martin/goqu/v9"
)

// DefaultWalletService handles wallet and transaction operations
type DefaultWalletService struct {
	db    *sql.DB // raw DB handle to initiate transactions
	wRepo domain.WalletRepository
	tRepo domain.TransactionRepository
}

// Ensure interface compliance
var _ domain.TransactionService = &DefaultWalletService{}

func NewWalletService(db *sql.DB, wRepo domain.WalletRepository, tRepo domain.TransactionRepository) domain.TransactionService {
	return &DefaultWalletService{
		db:    db,
		wRepo: wRepo,
		tRepo: tRepo,
	}
}

func (s *DefaultWalletService) TopUp(ctx context.Context, userID int64, amount float64) (*domain.Wallet, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	goquDb := goqu.New("mysql", s.db)
	txDb, err := goquDb.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer txDb.Rollback()

	// 1. Try to fetch existing wallet with lock (using repository)
	wallet, err := s.wRepo.GetWalletForUpdate(ctx, txDb, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check/lock wallet: %w", err)
	}

	if wallet == nil {
		// 2a. Create new wallet if not exists (using repository)
		wallet = &domain.Wallet{
			UserID:    userID,
			Balance:   amount,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := s.wRepo.CreateWithTx(ctx, txDb, wallet); err != nil {
			return nil, fmt.Errorf("failed to create wallet: %w", err)
		}
	} else {
		// 2b. Update existing wallet (using repository)
		newBalance := wallet.Balance + amount
		if err := s.wRepo.UpdateBalanceWithTx(ctx, txDb, wallet.ID, newBalance); err != nil {
			return nil, fmt.Errorf("failed to update wallet balance: %w", err)
		}
		wallet.Balance = newBalance
		wallet.UpdatedAt = time.Now()
	}

	if err := txDb.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return wallet, nil
}

func (s *DefaultWalletService) Transfer(ctx context.Context, senderUserID, receiverUserID int64, amount float64) (*domain.Transaction, error) {
	// 1. Basic Validation
	if amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}
	if senderUserID == receiverUserID {
		return nil, errors.New("cannot transfer to self")
	}

	// 2. Begin Database Transaction
	goquDb := goqu.New("mysql", s.db)
	txDb, err := goquDb.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer txDb.Rollback()

	// 3. Get and lock sender wallet (prevents race condition)
	senderWallet, err := s.wRepo.GetWalletForUpdate(ctx, txDb, senderUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sender wallet: %w", err)
	}
	if senderWallet == nil {
		return nil, errors.New("sender wallet not found")
	}

	// 4. Check Balance
	if senderWallet.Balance < amount {
		return nil, errors.New("insufficient balance")
	}

	// 5. Get and lock receiver wallet (prevents race condition)
	receiverWallet, err := s.wRepo.GetWalletForUpdate(ctx, txDb, receiverUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch receiver wallet: %w", err)
	}
	if receiverWallet == nil {
		return nil, errors.New("receiver wallet not found")
	}

	// 6. Update sender balance (deduct)
	newSenderBalance := senderWallet.Balance - amount
	if err := s.wRepo.UpdateBalanceWithTx(ctx, txDb, senderWallet.ID, newSenderBalance); err != nil {
		return nil, fmt.Errorf("failed to update sender balance: %w", err)
	}

	// 7. Update receiver balance (add)
	newReceiverBalance := receiverWallet.Balance + amount
	if err := s.wRepo.UpdateBalanceWithTx(ctx, txDb, receiverWallet.ID, newReceiverBalance); err != nil {
		return nil, fmt.Errorf("failed to update receiver balance: %w", err)
	}

	// 8. Create transaction record
	now := time.Now()
	transaction := &domain.Transaction{
		SenderWalletID:   &senderWallet.ID,
		ReceiverWalletID: &receiverWallet.ID,
		Amount:           amount,
		Status:           domain.TransactionStatusSuccess,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := s.tRepo.CreateWithTx(ctx, txDb, transaction); err != nil {
		return nil, fmt.Errorf("failed to create transaction record: %w", err)
	}

	// 9. Commit transaction
	if err := txDb.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return transaction, nil
}

func (s *DefaultWalletService) GetHistory(ctx context.Context, userID int64, page, limit int) ([]domain.Transaction, error) {
	// Simple implementation delegated to repo
	// We need to resolve UserID to WalletID first, outside of transaction is fine for read-only history
	wallet, err := s.wRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if wallet == nil {
		return nil, errors.New("wallet not found")
	}

	offset := (page - 1) * limit
	return s.tRepo.GetByWalletID(ctx, wallet.ID, limit, offset)
}

func (s *DefaultWalletService) GetBalance(ctx context.Context, userID int64) (*domain.Wallet, error) {
	wallet, err := s.wRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if wallet == nil {
		// Return default wallet with 0 balance
		return &domain.Wallet{
			UserID:  userID,
			Balance: 0,
		}, nil
	}
	return wallet, nil
}
