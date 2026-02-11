package repository

import (
	"context"
	"database/sql"
	"wallet-api/internal/domain"

	"github.com/doug-martin/goqu/v9"
)

// MysqlWalletRepo handles database operations for wallets
type MysqlWalletRepo struct {
	db *goqu.Database
}

// NewMysqlWalletRepository creates a new wallet repository
func NewMysqlWalletRepository(db *sql.DB) domain.WalletRepository {
	dialect := goqu.Dialect("mysql")
	return &MysqlWalletRepo{db: dialect.DB(db)}
}

// GoquExecutor defines the common methods between *goqu.Database and *goqu.TxDatabase that we use
type GoquExecutor interface {
	Insert(table interface{}) *goqu.InsertDataset
	From(table ...interface{}) *goqu.SelectDataset
	Update(table interface{}) *goqu.UpdateDataset
	Delete(table interface{}) *goqu.DeleteDataset
}

// Helper to get goqu executor
func getDb(r *MysqlWalletRepo, tx interface{}) GoquExecutor {
	if tx != nil {
		if txDb, ok := tx.(*goqu.TxDatabase); ok {
			return txDb
		}
	}
	return r.db
}

func (r *MysqlWalletRepo) Create(ctx context.Context, wallet *domain.Wallet) error {
	return r.CreateWithTx(ctx, nil, wallet)
}

func (r *MysqlWalletRepo) CreateWithTx(ctx context.Context, tx interface{}, wallet *domain.Wallet) error {
	db := getDb(r, tx)
	result, err := db.Insert("wallets").
		Rows(goqu.Record{
			"user_id":    wallet.UserID,
			"balance":    wallet.Balance,
			"created_at": wallet.CreatedAt,
			"updated_at": wallet.UpdatedAt,
		}).
		Executor().ExecContext(ctx)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err == nil {
		wallet.ID = id
	}
	return err
}

func (r *MysqlWalletRepo) GetByID(ctx context.Context, id int64) (*domain.Wallet, error) {
	var wallet domain.Wallet
	found, err := r.db.From("wallets").
		Where(goqu.C("id").Eq(id)).
		ScanStructContext(ctx, &wallet)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return &wallet, nil
}

func (r *MysqlWalletRepo) GetByUserID(ctx context.Context, userID int64) (*domain.Wallet, error) {
	var wallet domain.Wallet
	found, err := r.db.From("wallets").
		Where(goqu.C("user_id").Eq(userID)).
		ScanStructContext(ctx, &wallet)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return &wallet, nil
}

func (r *MysqlWalletRepo) GetWalletForUpdate(ctx context.Context, tx interface{}, userID int64) (*domain.Wallet, error) {
	db := getDb(r, tx)
	var wallet domain.Wallet

	// Ensure we are using a transaction for locking
	// If generic Interface doesn't support ForUpdate directly in check (it does in goqu), we use From

	found, err := db.From("wallets").
		Where(goqu.C("user_id").Eq(userID)).
		ForUpdate(goqu.Wait).
		ScanStructContext(ctx, &wallet)

	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil // Return nil if not found, let service handle
	}
	return &wallet, nil
}

func (r *MysqlWalletRepo) UpdateBalanceWithTx(ctx context.Context, tx interface{}, walletID int64, newBalance float64) error {
	db := getDb(r, tx)
	_, err := db.Update("wallets").
		Set(goqu.Record{
			"balance":    newBalance,
			"updated_at": goqu.L("NOW()"),
		}).
		Where(goqu.C("id").Eq(walletID)).
		Executor().ExecContext(ctx)
	return err
}

func (r *MysqlWalletRepo) UpdateBalance(ctx context.Context, id int64, amount float64) error {
	// Simple non-transactional update (not recommended for financial ops usually, but implemented for interface)
	_, err := r.db.Update("wallets").
		Set(goqu.Record{"balance": amount}).
		Where(goqu.C("id").Eq(id)).
		Executor().ExecContext(ctx)
	return err
}

func (r *MysqlWalletRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Delete("wallets").
		Where(goqu.C("id").Eq(id)).
		Executor().ExecContext(ctx)
	return err
}
