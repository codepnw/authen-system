package auth

import (
	"context"
	"time"

	"github.com/codepnw/go-authen-system/config"
	"github.com/codepnw/go-authen-system/internal/modules/user"
	"github.com/codepnw/go-authen-system/internal/utils/errs"
	"github.com/codepnw/go-authen-system/internal/utils/security"
	"github.com/codepnw/go-authen-system/pkg/logger"
)

const queryTimeout = time.Second * 5

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
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	// Create User
	user, err := uc.userUsecase.CreateUser(ctx, req)
	if err != nil {
		logger.Error("REGIS-001", "create user failed", err)
		return nil, err
	}

	tokenUser := uc.tokenUser(user)

	// Generate Token
	accessToken, refreshToken, err := uc.generateToken(tokenUser)
	if err != nil {
		logger.Error("REGIS-002", "generate token failed", err)
		return nil, err
	}

	// Data Response
	response := uc.authResponse(user, accessToken, refreshToken)
	logger.Info("REGIS-003", "register success", response)

	return response, nil
}

func (uc *authUsecase) Login(ctx context.Context, req *LoginRequestDTO) (*AuthResponseDTO, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	// Check User By Email
	user, err := uc.userUsecase.GetUserByEmail(ctx, req.Email)
	if err != nil {
		logger.Error("LOGIN-001", "get user email failed", err)
		return nil, errs.ErrInvalidEmailOrPassword
	}

	// Check Password
	if ok := security.VerifyPassword(user.Password, req.Password); !ok {
		logger.Error("LOGIN-002", "verify password failed", err)
		return nil, errs.ErrInvalidEmailOrPassword
	}

	tokenUser := uc.tokenUser(user)

	// Generate Token
	accessToken, refreshToken, err := uc.generateToken(tokenUser)
	if err != nil {
		logger.Error("LOGIN-003", "generate token failed", err)
		return nil, errs.ErrGenerateToken
	}

	// Save Refresh Token
	err = uc.authRepo.SaveRefreshToken(ctx, &RefreshToken{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(security.RefreshTokenDuration),
	})
	if err != nil {
		logger.Error("LOGIN-004", "save refresh token failed", err)
		return nil, errs.ErrSaveToken
	}

	// Data Response
	response := uc.authResponse(user, accessToken, refreshToken)
	logger.Info("LOGIN-005", "login success", response)

	return response, nil
}

func (uc *authUsecase) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	// Verify Refresh Token
	user, err := uc.tokenConfig.VerifyRefreshToken(refreshToken)
	if err != nil {
		logger.Error("REFRESH-001", "verify token failed", err)
		return "", "", errs.ErrInvalidToken
	}

	// Check Refresh Token in DB
	valid := uc.authRepo.IsRefreshToken(ctx, refreshToken)
	if !valid {
		logger.Error("REFRESH-002", "check token failed", err)
		return "", "", errs.ErrInvalidToken
	}

	// Generate New Access Token
	newAccessToken, err := uc.tokenConfig.GenerateAccessToken(user)
	if err != nil {
		logger.Error("REFRESH-003", "generate access token failed", err)
		return "", "", errs.ErrGenerateToken
	}

	// Generate New Refresh Token
	newRefreshToken, err := uc.tokenConfig.GenerateRefreshToken(user)
	if err != nil {
		logger.Error("REFRESH-004", "generate refresh token failed", err)
		return "", "", errs.ErrGenerateToken
	}

	// Update New Refresh Token
	err = uc.authRepo.UpdateRefreshToken(ctx, &RefreshToken{
		UserID:       user.ID,
		RefreshToken: newRefreshToken,
		ExpiresAt:    time.Now().Add(security.RefreshTokenDuration),
	})
	if err != nil {
		logger.Error("REFRESH-005", "update token failed", err)
		return "", "", errs.ErrSaveToken
	}

	logger.Info("REFRESH-006", "refresh token success", nil)
	return newAccessToken, newRefreshToken, nil
}

func (uc *authUsecase) Logout(ctx context.Context, userID int64) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	if err := uc.authRepo.DeleteRefreshToken(ctx, userID); err != nil {
		logger.Error("LOGOUT-001", "delete token failed", err)
		return errs.ErrInvalidToken
	}

	logger.Info("LOGOUT-002", "logout success", nil)
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
