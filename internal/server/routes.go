package server

import (
	"net/http"

	"github.com/codepnw/go-authen-system/config"
	"github.com/codepnw/go-authen-system/internal/modules/auth"
	"github.com/codepnw/go-authen-system/internal/modules/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type setupRoutes struct {
	router *gin.Engine
	db     *gorm.DB
	cfg    *config.Config
}

func (r *setupRoutes) healthCheck() {
	r.router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
}

func (r *setupRoutes) userRoutes() {
	repo := user.NewUserRepository(r.db)
	uc := user.NewUserUsecase(repo)
	hdl := user.NewUserHandler(uc)

	user := r.router.Group("/users")

	user.POST("/", hdl.CreateUser)
	user.GET("/", hdl.GetUsers)
	user.GET("/:id", hdl.GetProfile)
	user.PATCH("/:id", hdl.UpdateUser)
	user.DELETE("/:id", hdl.DeleteUser)
}

func (r *setupRoutes) authRoutes() {
	userRepo := user.NewUserRepository(r.db)
	userUsecase := user.NewUserUsecase(userRepo)

	authUsecase := auth.NewAuthUsecase(r.cfg, userUsecase)
	authHandler := auth.NewAuthHandler(authUsecase)

	auth := r.router.Group("/auth")
	auth.POST("/register", authHandler.Register)
}
