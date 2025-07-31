package server

import (
	"github.com/codepnw/go-authen-system/config"
	"github.com/codepnw/go-authen-system/internal/db"
	"github.com/codepnw/go-authen-system/internal/middleware"
	"github.com/codepnw/go-authen-system/pkg/logger"
	"github.com/gin-gonic/gin"
)

func Run(cfg *config.Config) error {
	// Connect Database
	db, err := db.NewDatabaseConnection(cfg)
	if err != nil {
		return err
	}

	// Init Logger
	log, err := logger.Init()
	if err != nil {
		return err
	}
	defer log.Sync()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(middleware.LoggerMiddleware())
	r.LoadHTMLGlob("templates/*.html")

	// Routes Config
	routes := setupRoutes{
		router: r,
		db:     db,
		cfg:    cfg,
	}
	routes.healthCheck()
	routes.userRoutes()
	routes.authRoutes()

	return r.Run(":" + cfg.AppPort)
}
