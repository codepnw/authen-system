package server

import (
	"github.com/codepnw/go-authen-system/config"
	"github.com/codepnw/go-authen-system/internal/db"
	"github.com/gin-gonic/gin"
)

func Run(cfg *config.Config) error {
	db, err := db.NewDatabaseConnection(cfg)
	if err != nil {
		return err
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")

	// Routes
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
