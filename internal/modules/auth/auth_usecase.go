package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/codepnw/go-authen-system/config"
	"github.com/codepnw/go-authen-system/internal/modules/user"
	"github.com/codepnw/go-authen-system/internal/utils/security"
)

type AuthUsecase interface {
	Register(ctx context.Context, req *user.CreateUserRequest) (*AuthResponseDTO, error)
	Login(ctx context.Context, req *LoginRequestDTO) (*AuthResponseDTO, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
	Logout(ctx context.Context, userID int64) error
}

type authUsecase struct {
	authRepo    AuthRepository
	userUsecase user.UserUsecase
	tokenConfig *security.TokenConfig
}

func NewAuthUsecase(cfg *config.Config, authRepo AuthRepository, userUsecase user.UserUsecase) AuthUsecase {
	return &authUsecase{
		authRepo:    authRepo,
		userUsecase: userUsecase,
		tokenConfig: security.NewJWTToken(cfg),
	}
}

func (uc *authUsecase) Register(ctx context.Context, req *user.CreateUserRequest) (*AuthResponseDTO, error) {
	// Create User
	user, err := uc.userUsecase.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	tokenUser := uc.tokenUser(user)

	// Generate Token
	accessToken, refreshToken, err := uc.generateToken(tokenUser)
	if err != nil {
		return nil, err
	}

	// Data Response
	response := uc.authResponse(user, accessToken, refreshToken)

	return response, nil
}

func (uc *authUsecase) Login(ctx context.Context, req *LoginRequestDTO) (*AuthResponseDTO, error) {
	// Check User By Email
	user, err := uc.userUsecase.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	// Check Password
	if ok := security.VerifyPassword(user.Password, req.Password); !ok {
		return nil, errors.New("invalid email or password")
	}

	tokenUser := uc.tokenUser(user)

	// Generate Token
	accessToken, refreshToken, err := uc.generateToken(tokenUser)
	if err != nil {
		return nil, err
	}

	// Save Refresh Token
	err = uc.authRepo.SaveRefreshToken(ctx, &RefreshToken{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(security.RefreshTokenDuration),
	})
	if err != nil {
		return nil, err
	}

	// Data Response
	response := uc.authResponse(user, accessToken, refreshToken)

	return response, nil
}

func (uc *authUsecase) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	// Verify Refresh Token
	user, err := uc.tokenConfig.VerifyRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	// Check Refresh Token in DB
	valid := uc.authRepo.IsRefreshToken(ctx, refreshToken)
	if !valid {
		return "", "", errors.New("invalid refresh token")
	}

	// Generate New Access Token
	newAccessToken, err := uc.tokenConfig.GenerateAccessToken(user)
	if err != nil {
		return "", "", fmt.Errorf("generate access token failed: %w", err)
	}

	// Generate New Refresh Token
	newRefreshToken, err := uc.tokenConfig.GenerateRefreshToken(user)
	if err != nil {
		return "", "", fmt.Errorf("generate refresh token failed: %w", err)
	}

	// Update New Refresh Token
	err = uc.authRepo.UpdateRefreshToken(ctx, &RefreshToken{
		UserID:       user.ID,
		RefreshToken: newRefreshToken,
		ExpiresAt:    time.Now().Add(security.RefreshTokenDuration),
	})
	if err != nil {
		return "", "", fmt.Errorf("save refresh token failed: %w", err)
	}

	return newAccessToken, newRefreshToken, nil
}

func (uc *authUsecase) Logout(ctx context.Context, userID int64) error {
	if err := uc.authRepo.DeleteRefreshToken(ctx, userID); err != nil {
		return fmt.Errorf("delete refresh token failed: %w", err)
	}

	return nil
}

// ------------- Private -------------
func (uc *authUsecase) tokenUser(user *user.User) *security.TokenUser {
	return &security.TokenUser{
		ID:    user.ID,
		Email: user.Email,
		Role:  "user", // TODO: change later
	}
}

func (uc *authUsecase) generateToken(user *security.TokenUser) (string, string, error) {
	// Generate Access Token
	accessToken, err := uc.tokenConfig.GenerateAccessToken(user)
	if err != nil {
		return "", "", err
	}

	// Generate Refresh Token
	refreshToken, err := uc.tokenConfig.GenerateRefreshToken(user)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (uc *authUsecase) authResponse(user *user.User, accessToken, refreshToken string) *AuthResponseDTO {
	return &AuthResponseDTO{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
