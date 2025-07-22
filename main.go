package main

import (
	"log"

	"github.com/codepnw/go-authen-system/config"
	"github.com/codepnw/go-authen-system/internal/server"
)

const configFileName = "config"

func main() {
	cfg, err := config.InitConfig(configFileName)
	if err != nil {
		log.Fatal(err)
	}

	if err = server.Run(cfg); err != nil {
		log.Fatal(err)
	}
}
