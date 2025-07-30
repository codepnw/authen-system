package auth

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type AuthRepository interface {
	SaveRefreshToken(ctx context.Context, input *RefreshToken) error
	UpdateRefreshToken(ctx context.Context, input *RefreshToken) error
	IsRefreshToken(ctx context.Context, refreshToken string) bool
	DeleteRefreshToken(ctx context.Context, userID int64) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) SaveRefreshToken(ctx context.Context, input *RefreshToken) error {
	if err := r.db.WithContext(ctx).Create(input).Error; err != nil {
		return err
	}
	return nil
}

func (r *authRepository) UpdateRefreshToken(ctx context.Context, input *RefreshToken) error {
	res := r.db.WithContext(ctx).Where("user_id = ?", input.UserID).Updates(input)
	if res.Error != nil {
		return res.Error
	}

	rows := res.RowsAffected
	if rows == 0 {
		return errors.New("refresh token not found")
	}

	return nil
}

func (r *authRepository) IsRefreshToken(ctx context.Context, refreshToken string) bool {
	err := r.db.WithContext(ctx).First(&RefreshToken{}, "refresh_token = ? AND expires_at > NOW()", refreshToken).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}

	return err == nil
}

func (r *authRepository) DeleteRefreshToken(ctx context.Context, userID int64) error {
	res := r.db.Delete(&RefreshToken{}, "user_id = ?", userID)
	if res.Error != nil {
		return res.Error
	}

	rows := res.RowsAffected
	if rows == 0 {
		return errors.New("user id not found")
	}

	return nil
}
