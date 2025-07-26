package auth

import (
	"github.com/codepnw/go-authen-system/internal/modules/user"
	"github.com/codepnw/go-authen-system/internal/utils/response"
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
