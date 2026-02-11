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
	// Not implemented in this step request
	return nil, nil
}

func (s *DefaultWalletService) Transfer(ctx context.Context, senderUserID, receiverUserID int64, amount float64) (*domain.Transaction, error) {
	// 1. Basic Validation
	if amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}
	if senderUserID == receiverUserID {
		return nil, errors.New("cannot transfer to self")
	}

	// 2. Begin Database Transaction using Goqu
	goquDb := goqu.New("mysql", s.db)
	txDb, err := goquDb.BeginTx(ctx, nil) // Returns *goqu.TxDatabase, Error
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	// We must defer Rollback on the goqu TxDatabase
	defer txDb.Rollback()

	// 3. Get User Wallets (Resolve UserID to WalletID)
	// Now `txDb` is a *goqu.TxDatabase which has all methods (From, Insert, Update)
	// and executes them on the transaction.

	var senderWallet domain.Wallet
	var receiverWallet domain.Wallet

	// Fetch Sender Wallet with FOR UPDATE to lock and prevent race conditions
	foundSender, err := txDb.From("wallets").
		Where(goqu.C("user_id").Eq(senderUserID)).
		ForUpdate(goqu.Wait).
		ScanStructContext(ctx, &senderWallet)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch sender wallet: %w", err)
	}
	if !foundSender {
		return nil, errors.New("sender wallet not found")
	}

	// 4. Check Balance
	if senderWallet.Balance < amount {
		return nil, errors.New("insufficient balance")
	}

	// Fetch Receiver Wallet
	foundReceiver, err := txDb.From("wallets").
		Where(goqu.C("user_id").Eq(receiverUserID)).
		ForUpdate(goqu.Wait).
		ScanStructContext(ctx, &receiverWallet)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch receiver wallet: %w", err)
	}
	if !foundReceiver {
		return nil, errors.New("receiver wallet not found")
	}

	// 5. Update Balances (Atomicity)
	// Subtract from Sender
	_, err = txDb.Update("wallets").
		Set(goqu.Record{"balance": senderWallet.Balance - amount}).
		Where(goqu.C("id").Eq(senderWallet.ID)).
		Executor().ExecContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to update sender balance: %w", err)
	}

	// Add to Receiver
	_, err = txDb.Update("wallets").
		Set(goqu.Record{"balance": receiverWallet.Balance + amount}).
		Where(goqu.C("id").Eq(receiverWallet.ID)).
		Executor().ExecContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to update receiver balance: %w", err)
	}

	// 6. Insert Transaction Record
	now := time.Now()
	transaction := &domain.Transaction{
		SenderWalletID:   &senderWallet.ID,
		ReceiverWalletID: &receiverWallet.ID,
		Amount:           amount,
		Status:           domain.TransactionStatusSuccess,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	result, err := txDb.Insert("transactions").
		Rows(goqu.Record{
			"sender_wallet_id":   transaction.SenderWalletID,
			"receiver_wallet_id": transaction.ReceiverWalletID,
			"amount":             transaction.Amount,
			"status":             transaction.Status,
			"created_at":         transaction.CreatedAt,
			"updated_at":         transaction.UpdatedAt,
		}).
		Executor().ExecContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction record: %w", err)
	}

	// Get inserted ID
	id, err := result.LastInsertId()
	if err == nil {
		transaction.ID = id
	}

	// 7. Commit Transaction
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
