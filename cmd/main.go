package main

import (
	"fmt"
	"nazartaraniuk/alertsProject/internal/app/client"
	"nazartaraniuk/alertsProject/internal/app/server"
	"nazartaraniuk/alertsProject/internal/config"
	"nazartaraniuk/alertsProject/internal/repository"
	"nazartaraniuk/alertsProject/internal/usecase"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	cfg := config.LoadConfig()
	fmt.Println(cfg.Client.APIBaseURL)

	mainClient := client.NewClient(
		cfg.Client.APIBaseURL,
		time.Second,
		cfg.Client.APIKey,
	)

	alarmsRepository := repository.NewAlarmsRepository(mainClient)
	service := usecase.NewGetAlarmInfoService(*alarmsRepository)

	mainServer, err := server.NewServer(cfg, *service)
	if err != nil {
		log.Fatal(err)
	}

	_ = mainServer.Run()
}
