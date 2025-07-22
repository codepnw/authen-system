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

	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")

	// Routes
	routes := setupRoutes{
		router: r,
		db:     db,
	}
	routes.healthCheck()
	routes.userRoutes()

	return r.Run(":" + cfg.AppPort)
}
