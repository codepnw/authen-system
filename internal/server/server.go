package server

import (
	"net/http"

	"github.com/codepnw/go-authen-system/config"
	"github.com/codepnw/go-authen-system/internal/db"
	"github.com/gin-gonic/gin"
)

func Run(cfg *config.Config) error {
	_, err := db.NewDatabaseConnection(cfg)
	if err != nil {
		return err
	}

	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	return r.Run(":" + cfg.AppPort)
}
