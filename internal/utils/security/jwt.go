package security

import (
	"fmt"
	"time"

	"github.com/codepnw/go-authen-system/config"
	"github.com/golang-jwt/jwt/v5"
)

type TokenConfig struct {
	SecretKey  string
	RefreshKey string
}

type generateTokenParams struct {
	ID       int64
	Email    string
	Role     string
	Key      string
	Duration time.Duration
}

type TokenUser struct {
	ID    int64
	Email string
	Role  string
}

func NewJWTToken(cfg *config.Config) *TokenConfig {
	return &TokenConfig{
		SecretKey:  cfg.JWTSecretKey,
		RefreshKey: cfg.JWTRefreshKey,
	}
}

func (t *TokenConfig) GenerateAccessToken(user *TokenUser) (string, error) {
	duration := time.Hour * 24

	return t.generateToken(&generateTokenParams{
		ID:       user.ID,
		Email:    user.Email,
		Role:     user.Role,
		Key:      t.SecretKey,
		Duration: duration,
	})
}

func (t *TokenConfig) GenerateRefreshToken(user *TokenUser) (string, error) {
	duration := time.Hour * 24 * 7

	return t.generateToken(&generateTokenParams{
		ID:       user.ID,
		Email:    user.Email,
		Role:     user.Role,
		Key:      t.RefreshKey,
		Duration: duration,
	})
}

func (t *TokenConfig) generateToken(input *generateTokenParams) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": input.ID,
		"email":   input.Email,
		"role":    input.Role,
		"exp":     time.Now().Add(input.Duration).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(input.Key))
	if err != nil {
		return "", fmt.Errorf("sign token failed: %w", err)
	}

	return tokenStr, nil
}
