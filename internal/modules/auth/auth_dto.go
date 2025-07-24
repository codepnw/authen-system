package auth

import "github.com/codepnw/go-authen-system/internal/modules/user"

type AuthResponseDTO struct {
	User         *user.User `json:"user"`
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
}
