package auth

import (
	"context"

	"github.com/codepnw/go-authen-system/config"
	"github.com/codepnw/go-authen-system/internal/modules/user"
	"github.com/codepnw/go-authen-system/internal/utils/security"
)

type AuthUsecase interface {
	Register(ctx context.Context, req *user.CreateUserRequest) (*AuthResponseDTO, error)
	// Login()
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

	credential := &security.TokenUser{
		ID:    user.ID,
		Email: user.Email,
		// TODO: change role later
		Role:  "user",
	}

	// Generate Access Token
	accessToken, err := uc.tokenConfig.GenerateAccessToken(credential)
	if err != nil {
		return nil, err
	}

	// Generate Refresh Token
	refreshToken, err := uc.tokenConfig.GenerateRefreshToken(credential)
	if err != nil {
		return nil, err
	}

	response := &AuthResponseDTO{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return response, nil
}
