package auth

import (
	"context"
	"errors"

	"github.com/codepnw/go-authen-system/config"
	"github.com/codepnw/go-authen-system/internal/modules/user"
	"github.com/codepnw/go-authen-system/internal/utils/security"
)

type AuthUsecase interface {
	Register(ctx context.Context, req *user.CreateUserRequest) (*AuthResponseDTO, error)
	Login(ctx context.Context, req *LoginRequestDTO) (*AuthResponseDTO, error)
}

type authUsecase struct {
	userUsecase user.UserUsecase
	tokenConfig *security.TokenConfig
}

func NewAuthUsecase(cfg *config.Config, userUsecase user.UserUsecase) AuthUsecase {
	return &authUsecase{
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

	// Data Response
	response := uc.authResponse(user, accessToken, refreshToken)

	return response, nil
}

// Private Method
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
