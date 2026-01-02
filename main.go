package main

import (
	"alfdwirhmn/inventory/database"
	"alfdwirhmn/inventory/handler"
	"alfdwirhmn/inventory/repository"
	"alfdwirhmn/inventory/router"
	"alfdwirhmn/inventory/service"
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

	db, err := database.InitDB(config.DB)
	if err != nil {
		panic(err)
	}

	// validate := validator.New()

	logger, err = utils.InitLogger(config.PathLogg, config.Debug)

	repo := repository.NewContainer(db, logger)
	svc := service.NewContainer(repo)
	h := handler.NewContainer(svc, repo, logger, config)

	r := router.NewRouter(h, logger)

	// run server with port from config
	log.Printf("Server running on port %s\n", config.Port)
	if err := http.ListenAndServe(":"+config.Port, r); err != nil {
		log.Fatal("Server error: ", err)
	}
}
