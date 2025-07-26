package user

import (
	"context"
	"errors"
	"time"

	"github.com/codepnw/go-authen-system/internal/utils/security"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error)
	GetProfile(ctx context.Context, id int64) (*User, error)
	GetUsers(ctx context.Context) ([]*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, id int64, req *UpdateUserRequest) error
	DeleteUser(ctx context.Context, id int64) error
}

type userUsecase struct {
	repo UserRepository
}

func NewUserUsecase(repo UserRepository) UserUsecase {
	return &userUsecase{repo: repo}
}

func (uc *userUsecase) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
	// Check Email
	found, err := uc.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if found != nil {
		return nil, errors.New("email already exists")
	}

	if req.Password != req.ConfirmPassword {
		return nil, errors.New("password and confirm_password not match")
	}

	// Hash Password
	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	// Create User
	created, err := uc.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (uc *userUsecase) DeleteUser(ctx context.Context, id int64) error {
	if err := uc.repo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

func (uc *userUsecase) GetProfile(ctx context.Context, id int64) (*User, error) {
	user, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (uc *userUsecase) GetUsers(ctx context.Context) ([]*User, error) {
	users, err := uc.repo.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (uc *userUsecase) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return uc.repo.FindByEmail(ctx, email)
}

func (uc *userUsecase) UpdateUser(ctx context.Context, id int64, req *UpdateUserRequest) error {
	user, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if req.Email != nil {
		user.Email = *req.Email
	}

	if req.Username != nil {
		user.Username = *req.Username
	}

	now := time.Now()
	user.UpdatedAt = &now

	if err = uc.repo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}
