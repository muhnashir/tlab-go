package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"wallet-api/internal/domain"
	"wallet-api/internal/utils"
)

type DefaultUserService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) domain.UserService {
	return &DefaultUserService{repo: repo}
}

func (s *DefaultUserService) Register(ctx context.Context, user *domain.User) error {
	// Check if user already exists
	existing, err := s.repo.GetByEmail(ctx, user.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) { // Assuming repo returns error on not found or empty
		// This needs adjustment based on repo implementation. If repo returns nil on not found:
	}
	if existing != nil {
		return errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return s.repo.Create(ctx, user)
}

func (s *DefaultUserService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *DefaultUserService) GetProfile(ctx context.Context, id int64) (*domain.User, error) {
	return s.repo.GetByID(ctx, id)
}
