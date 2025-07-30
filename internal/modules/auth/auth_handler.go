package auth

import (
	"github.com/codepnw/go-authen-system/internal/middleware"
	"github.com/codepnw/go-authen-system/internal/modules/user"
	"github.com/codepnw/go-authen-system/internal/utils/response"
	"github.com/codepnw/go-authen-system/internal/utils/security"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type authHandler struct {
	uc       AuthUsecase
	validate *validator.Validate
}

func NewAuthHandler(uc AuthUsecase) *authHandler {
	return &authHandler{
		uc:       uc,
		validate: validator.New(),
	}
}

func (h *authHandler) Register(c *gin.Context) {
	req := new(user.CreateUserRequest)

	if err := c.ShouldBindJSON(req); err != nil {
		response.BadRequest(c, "", err)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.BadRequest(c, "", err)
		return
	}

	data, err := h.uc.Register(c, req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.Created(c, data)
}

func (h *authHandler) Login(c *gin.Context) {
	req := new(LoginRequestDTO)

	if err := c.ShouldBindJSON(req); err != nil {
		response.BadRequest(c, "", err)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.BadRequest(c, "", err)
		return
	}

	result, err := h.uc.Login(c, req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, "", result)
}

func (h *authHandler) Profile(c *gin.Context) {
	user, ok := c.Get(middleware.UserContextKey)
	if !ok {
		response.Unauthorized(c, nil)
		return
	}

	response.Success(c, "", user)
}

func (h *authHandler) RefreshToken(c *gin.Context) {
	req := new(RefreshTokenRequestDTO)

	if err := c.ShouldBindJSON(req); err != nil {
		response.BadRequest(c, "", err)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.BadRequest(c, "", err)
		return
	}

	accessToken, refreshToken, err := h.uc.RefreshToken(c, req.RefreshToken)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}
	
	res := &RefreshTokenResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	response.Success(c, "", res)
}

func (h *authHandler) Logout(c *gin.Context) {
	user, ok := c.Get(middleware.UserContextKey)
	if !ok {
		response.Unauthorized(c, nil)
		return
	}

	u := user.(*security.TokenUser)

	if err := h.uc.Logout(c, u.ID); err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, "logout success", nil)
}
