package user

import (
	"strconv"

	"github.com/codepnw/go-authen-system/internal/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type userHandler struct {
	validate *validator.Validate
	uc       UserUsecase
}

func NewUserHandler(uc UserUsecase) *userHandler {
	return &userHandler{
		validate: validator.New(),
		uc:       uc,
	}
}

func (h *userHandler) CreateUser(c *gin.Context) {
	req := new(CreateUserRequest)

	if err := c.ShouldBindJSON(req); err != nil {
		response.BadRequest(c, "", err)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.BadRequest(c, "", err)
		return
	}

	// Create User
	user, err := h.uc.CreateUser(c, req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.Created(c, user)
}

func (h *userHandler) GetProfile(c *gin.Context) {
	// TODO: get id from context
	id, err := getIntParamID(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "", err)
		return
	}

	user, err := h.uc.GetProfile(c, id)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, "", user)
}

func (h *userHandler) GetUsers(c *gin.Context) {
	users, err := h.uc.GetUsers(c)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, "", users)
}

func (h *userHandler) UpdateUser(c *gin.Context) {
	id, err := getIntParamID(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id", err)
		return
	}

	req := new(UpdateUserRequest)

	if err := c.ShouldBindJSON(req); err != nil {
		response.BadRequest(c, "", err)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.BadRequest(c, "", err)
		return
	}

	// Update User
	if err = h.uc.UpdateUser(c, id, req); err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, "user updated", nil)
}

func (h *userHandler) DeleteUser(c *gin.Context) {
	id, err := getIntParamID(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id", err)
		return
	}

	if err = h.uc.DeleteUser(c, id); err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, "user deleted", nil)
}

func getIntParamID(key string) (int64, error) {
	return strconv.ParseInt(key, 10, 64)
}
