package security

import (
	"errors"
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

func (t *TokenConfig) VerifyToken(tokenString, key string) (*TokenUser, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unknow signing method: %v", t.Header)
		}
		return []byte(key), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("verification failed")
	}

	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		return nil, errors.New("token is expired")
	}

	user := new(TokenUser)
	user.ID = int64(claims["user_id"].(float64)) // JSON encode number to float64
	user.Email = claims["email"].(string)
	user.Role = claims["role"].(string)

	return user, nil
}
