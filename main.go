package main

import (
	"alfdwirhmn/inventory/router"
	"alfdwirhmn/inventory/utils"
	"log"
	"net/http"

	"go.uber.org/zap"
)

func main() {
	config, err := utils.ReadConfigurationEnv()
	if err != nil {
		log.Fatal("failed loaded configuration:", err)
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	r := router.NewRouter(logger)

	// run server with port from config
	log.Printf("Server running on port %s\n", config.Port)
	if err := http.ListenAndServe(":"+config.Port, r); err != nil {
		log.Fatal("Server error: ", err)
	}
}
