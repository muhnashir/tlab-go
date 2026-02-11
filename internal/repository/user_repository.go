package repository

import (
	"context"
	"database/sql"
	"wallet-api/internal/domain"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
)

const (
	UserTable = `users`
)

type MysqlUserRepository struct {
	db *goqu.Database
}

func NewMysqlUserRepository(db *sql.DB) domain.UserRepository {
	// Setup MySQL dialect for goqu
	return &MysqlUserRepository{
		db: goqu.New("mysql", db),
	}
}

func (r *MysqlUserRepository) Create(ctx context.Context, user *domain.User) error {

	ds := r.db.Insert(UserTable).Rows(goqu.Record{
		"name":       user.Name,
		"email":      user.Email,
		"password":   user.Password,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})

	// 2. Eksekusi langsung tanpa ToSQL manual jika ingin praktis
	// Method ini aman dari nil pointer karena r.db sudah di-init di NewMysqlUserRepository
	_, err := ds.Executor().ExecContext(ctx)
	return err
}

func (r *MysqlUserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	var user domain.User
	found, err := r.db.From(UserTable).
		Where(goqu.C("id").Eq(id)).
		ScanStructContext(ctx, &user)

	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil // Or domain specific error like sql.ErrNoRows
	}
	return &user, nil
}

// GetByEmail retrieves a user by their email address
func (r *MysqlUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	found, err := r.db.From(UserTable).
		Where(goqu.C("email").Eq(email)).
		ScanStructContext(ctx, &user)

	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return &user, nil
}

func (r *MysqlUserRepository) Update(ctx context.Context, user *domain.User) error {
	_, err := r.db.Update(UserTable).
		Set(goqu.Record{
			"name":     user.Name,
			"password": user.Password,
		}).
		Where(goqu.C("id").Eq(user.ID)).
		Executor().ExecContext(ctx)
	return err
}

func (r *MysqlUserRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Delete(UserTable).
		Where(goqu.C("id").Eq(id)).
		Executor().ExecContext(ctx)
	return err
}
