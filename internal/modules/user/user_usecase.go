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
	UpdateUser(ctx context.Context, id int64, req *UpdateUserRequest) error
	DeleteUser(ctx context.Context, id int64) error
}

type userUsecase struct {
	repo UserRepository
}

func NewUserUsecase(repo UserRepository) UserUsecase {
	return &userUsecase{repo: repo}
}

func (u *userUsecase) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
	// Check Email
	found, err := u.repo.FindByEmail(ctx, req.Email)
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
	created, err := u.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (u *userUsecase) DeleteUser(ctx context.Context, id int64) error {
	if err := u.repo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

func (u *userUsecase) GetProfile(ctx context.Context, id int64) (*User, error) {
	user, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) GetUsers(ctx context.Context) ([]*User, error) {
	users, err := u.repo.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *userUsecase) UpdateUser(ctx context.Context, id int64, req *UpdateUserRequest) error {
	user, err := u.repo.FindByID(ctx, id)
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

	if err = u.repo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}
