package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	AppPort       string
	DBUser        string
	DBPass        string
	DBHost        string
	DBName        string
	DBSSLMode     string
	JWTSecretKey  string
	JWTRefreshKey string
}

func InitConfig(fileName string) (*Config, error) {
	viper.SetConfigName(fileName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetDefault("app.port", 8080)
	viper.SetDefault("db.user", "postgres")
	viper.SetDefault("db.password", "")
	viper.SetDefault("db.host", "localhost:5432")
	viper.SetDefault("db.name", "auth_system")
	viper.SetDefault("db.ssl_mode", "disable")
	viper.SetDefault("jwt.secret_key", "secret_key")
	viper.SetDefault("jwt.refresh_key", "refresh_key")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("reading config failed: %w", err)
	}

	return &Config{
		AppPort:       viper.GetString("app.port"),
		DBUser:        viper.GetString("db.user"),
		DBPass:        viper.GetString("db.password"),
		DBHost:        viper.GetString("db.host"),
		DBName:        viper.GetString("db.name"),
		DBSSLMode:     viper.GetString("db.ssl_mode"),
		JWTSecretKey:  viper.GetString("jwt.secret_key"),
		JWTRefreshKey: viper.GetString("jwt.refresh_key"),
	}, nil
}
