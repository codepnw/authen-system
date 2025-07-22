package user

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, input *User) (*User, error)
	FindByID(ctx context.Context, id int64) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	ListUsers(ctx context.Context) ([]*User, error)
	Update(ctx context.Context, input *User) error
	Delete(ctx context.Context, id int64) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (u *userRepository) Create(ctx context.Context, user *User) (*User, error) {
	err := u.db.Create(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userRepository) Delete(ctx context.Context, id int64) error {
	res := u.db.Delete(&User{}, id)
	if res.Error != nil {
		return res.Error
	}

	rows := res.RowsAffected
	if rows == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (u *userRepository) FindByEmail(ctx context.Context, email string) (user *User, err error) {
	res := u.db.First(&user, "email = ?", email)
	if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}

func (u *userRepository) ListUsers(ctx context.Context) (users []*User, err error) {
	if err = u.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (u *userRepository) FindByID(ctx context.Context, id int64) (user *User, err error) {
	res := u.db.First(&user, "id = ?", id)
	if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}

func (u *userRepository) Update(ctx context.Context, input *User) error {
	res := u.db.Save(&input)
	if res.Error != nil {
		return res.Error
	}

	rows := res.RowsAffected
	if rows == 0 {
		return errors.New("user not found")
	}

	return nil
}
