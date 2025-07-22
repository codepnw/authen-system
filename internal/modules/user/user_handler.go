package user

import "github.com/gin-gonic/gin"

type userHandler struct {
	uc UserUsecase
}

func NewUserHandler(uc UserUsecase) *userHandler {
	return &userHandler{uc: uc}
}

func (h *userHandler) CreateUser(c *gin.Context) {}

func (h *userHandler) GetProfile(c *gin.Context) {}

func (h *userHandler) GetUsers(c *gin.Context) {}

func (h *userHandler) UpdateUser(c *gin.Context) {}

func (h *userHandler) DeleteUser(c *gin.Context) {}
