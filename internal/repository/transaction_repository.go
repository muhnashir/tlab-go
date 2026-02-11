package repository

import (
	"context"
	"database/sql"
	"errors"
	"wallet-api/internal/domain"

	"github.com/doug-martin/goqu/v9"
)

// MysqlTransactionRepo handles transaction records
type MysqlTransactionRepo struct {
	db *goqu.Database
}

// NewMysqlTransactionRepository creates a new transaction repository
func NewMysqlTransactionRepository(db *sql.DB) domain.TransactionRepository {
	dialect := goqu.Dialect("mysql")
	return &MysqlTransactionRepo{db: dialect.DB(db)}
}

// CreateWithTx creates a transaction record within a database transaction context
func (r *MysqlTransactionRepo) CreateWithTx(ctx context.Context, tx interface{}, transaction *domain.Transaction) error {
	// Cast tx to *goqu.TxDatabase
	txDb, ok := tx.(*goqu.TxDatabase)
	if !ok {
		return errors.New("invalid transaction type")
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
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	transaction.ID = id
	return nil
}

// Create (outside transaction)
func (r *MysqlTransactionRepo) Create(ctx context.Context, transaction *domain.Transaction) error {
	result, err := r.db.Insert("transactions").
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
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	transaction.ID = id
	return nil
}

func (r *MysqlTransactionRepo) GetByID(ctx context.Context, id int64) (*domain.Transaction, error) {
	var transaction domain.Transaction
	found, err := r.db.From("transactions").
		Where(goqu.C("id").Eq(id)).
		ScanStructContext(ctx, &transaction)

	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return &transaction, nil
}

func (r *MysqlTransactionRepo) GetByWalletID(ctx context.Context, walletID int64, limit, offset int) ([]domain.Transaction, error) {
	var transactions []domain.Transaction
	err := r.db.From("transactions").
		Where(goqu.Or(
			goqu.C("sender_wallet_id").Eq(walletID),
			goqu.C("receiver_wallet_id").Eq(walletID),
		)).
		Order(goqu.C("created_at").Desc()).
		Limit(uint(limit)).
		Offset(uint(offset)).
		ScanStructsContext(ctx, &transactions)

	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *MysqlTransactionRepo) UpdateStatus(ctx context.Context, id int64, status domain.TransactionStatus) error {
	_, err := r.db.Update("transactions").
		Set(goqu.Record{"status": status}).
		Where(goqu.C("id").Eq(id)).
		Executor().ExecContext(ctx)
	return err
}
