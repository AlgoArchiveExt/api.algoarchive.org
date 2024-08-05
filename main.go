package main

// Boilerplate is from https://github.com/akmamun/gin-boilerplate-examples/tree/main
import (
	"main/config"
	"main/infra/logger"
	"main/routers"
	"time"

	"github.com/spf13/viper"
)

func main() {
	// Set Timezone
	viper.SetDefault("SERVER_TIMEZONE", "US/Eastern")
	loc, _ := time.LoadLocation(viper.GetString("SERVER_TIMEZONE"))
	time.Local = loc

	if err := config.SetupConfig(); err != nil {
		logger.Fatalf("config SetupConfig() error: %s", err)
	}

	// masterDSN, replicaDSN := config.DbConfiguration()

	// if err := database.DBConnection(masterDSN, replicaDSN); err != nil {
	// 	logger.Fatalf("database DbConnection error: %s", err)
	// }

	router := routers.Routes()

	logger.Fatalf("%v", router.Run(config.ServerConfig()))

}
