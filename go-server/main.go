package main

// Boilerplate is from https://github.com/akmamun/gin-boilerplate-examples/tree/main
import (
	"main/config"
	"main/infra/logger"
	"main/routers"
)

func main() {
	if err := config.SetupConfig(); err != nil {
		logger.Fatalf("config SetupConfig() error: %s", err)
	}

	router := routers.Routes()

	logger.Fatalf("%v", router.Run(config.ServerConfig()))
}
