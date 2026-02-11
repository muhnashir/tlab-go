package repository

import (
	"context"
	"database/sql"
	"errors"
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

func (r *MysqlWalletRepo) Create(ctx context.Context, wallet *domain.Wallet) error {
	_, err := r.db.Insert("wallets").
		Rows(goqu.Record{
			"user_id": wallet.UserID,
			"balance": wallet.Balance,
		}).
		Executor().ExecContext(ctx)
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

func (r *MysqlWalletRepo) GetWalletForUpdate(ctx context.Context, tx *goqu.TxDatabase, walletID int64) (*domain.Wallet, error) {
	var wallet domain.Wallet
	found, err := tx.From("wallets").
		Where(goqu.C("id").Eq(walletID)).
		ForUpdate(goqu.Wait).
		ScanStructContext(ctx, &wallet)

	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errors.New("wallet not found")
	}
	return &wallet, nil
}

func (r *MysqlWalletRepo) UpdateBalanceWithTx(ctx context.Context, tx *goqu.TxDatabase, walletID int64, newBalance float64) error {
	_, err := tx.Update("wallets").
		Set(goqu.Record{"balance": newBalance}).
		Where(goqu.C("id").Eq(walletID)).
		Executor().ExecContext(ctx)
	return err
}

func (r *MysqlWalletRepo) UpdateBalance(ctx context.Context, id int64, amount float64) error {
	return nil
}

func (r *MysqlWalletRepo) Delete(ctx context.Context, id int64) error {
	return nil
}
