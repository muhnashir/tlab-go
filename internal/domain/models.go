package domain

import (
	"context"
	"time"
)

// User represents a user entity
type User struct {
	ID        int64     `json:"id" db:"id" goqu:"skipinsert"`
	Name      string    `json:"name" db:"name" validate:"required"`
	Email     string    `json:"email" db:"email" validate:"required,email"`
	Password  string    `json:"password,omitempty" db:"password" validate:"required,min=6"`
	CreatedAt time.Time `json:"created_at" db:"created_at" goqu:"skipinsert"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" goqu:"skipinsert"`
}

// Wallet represents a digital wallet associated with a user
type Wallet struct {
	ID        int64     `json:"id" db:"id" goqu:"skipinsert"`
	UserID    int64     `json:"user_id" validate:"required" db:"user_id"`
	Balance   float64   `json:"balance" db:"balance"`
	CreatedAt time.Time `json:"created_at" db:"created_at" goqu:"skipinsert"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" goqu:"skipinsert"`
}

// TransactionStatus defines possible statuses for a transaction
type TransactionStatus string

const (
	TransactionStatusPending TransactionStatus = "pending"
	TransactionStatusSuccess TransactionStatus = "success"
	TransactionStatusFailed  TransactionStatus = "failed"
)

// Transaction represents a financial transaction between wallets
type Transaction struct {
	ID               int64             `json:"id" goqu:"skipinsert"`
	SenderWalletID   *int64            `json:"sender_wallet_id"`   // Nullable if system sends money
	ReceiverWalletID *int64            `json:"receiver_wallet_id"` // Nullable if withdrawing to external
	Amount           float64           `json:"amount" validate:"required,gt=0"`
	Status           TransactionStatus `json:"status"`
	CreatedAt        time.Time         `json:"created_at" goqu:"skipinsert"`
	UpdatedAt        time.Time         `json:"updated_at" goqu:"skipinsert"`
}

// UserRepository defines methods for interacting with user data
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error
}

// WalletRepository defines methods for interacting with wallet data
type WalletRepository interface {
	Create(ctx context.Context, wallet *Wallet) error
	CreateWithTx(ctx context.Context, tx interface{}, wallet *Wallet) error // For transaction support
	GetByID(ctx context.Context, id int64) (*Wallet, error)
	GetByUserID(ctx context.Context, userID int64) (*Wallet, error)
	GetWalletForUpdate(ctx context.Context, tx interface{}, userID int64) (*Wallet, error) // For locking
	UpdateBalance(ctx context.Context, id int64, amount float64) error                     // Atomic update
	UpdateBalanceWithTx(ctx context.Context, tx interface{}, walletID int64, newBalance float64) error
	Delete(ctx context.Context, id int64) error
}

// TransactionRepository defines methods for interacting with transaction data
type TransactionRepository interface {
	Create(ctx context.Context, transaction *Transaction) error
	GetByID(ctx context.Context, id int64) (*Transaction, error)
	GetByWalletID(ctx context.Context, walletID int64, limit, offset int) ([]Transaction, error)
	UpdateStatus(ctx context.Context, id int64, status TransactionStatus) error
}

// UserService defines business logic for users
type UserService interface {
	Register(ctx context.Context, user *User) error
	Login(ctx context.Context, email, password string) (string, error) // Returns JWT token
	GetProfile(ctx context.Context, id int64) (*User, error)
}

// TransactionService defines business logic for transactions
type TransactionService interface {
	TopUp(ctx context.Context, userID int64, amount float64) (*Wallet, error)
	Transfer(ctx context.Context, senderID, receiverID int64, amount float64) (*Transaction, error)
	GetHistory(ctx context.Context, userID int64, page, limit int) ([]Transaction, error)
	GetBalance(ctx context.Context, userID int64) (*Wallet, error)
}
