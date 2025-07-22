package user

import "context"

type UserUsecase interface {
	CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error)
	GetProfile(ctx context.Context, id int64) (User, error)
	GetUsers(ctx context.Context) ([]User, error)
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
	panic("unimplemented")
}

func (u *userUsecase) DeleteUser(ctx context.Context, id int64) error {
	panic("unimplemented")
}

func (u *userUsecase) GetProfile(ctx context.Context, id int64) (User, error) {
	panic("unimplemented")
}

func (u *userUsecase) GetUsers(ctx context.Context) ([]User, error) {
	panic("unimplemented")
}

func (u *userUsecase) UpdateUser(ctx context.Context, id int64, req *UpdateUserRequest) error {
	panic("unimplemented")
}
