package main

import (
	"fmt"
	"nazartaraniuk/alertsProject/internal/app/client"
	"nazartaraniuk/alertsProject/internal/app/server"
	"nazartaraniuk/alertsProject/internal/config"
	"nazartaraniuk/alertsProject/internal/repository"
	"nazartaraniuk/alertsProject/internal/usecase"
	"time"

	"github.com/sirupsen/logrus"
)

// @title Alerts Project API
// @version 2.0
// @description API for getting alarms
// @host localhost:8080
// @BasePath /
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

	mainServer, err := server.NewServer(*cfg, *service)
	if err != nil {
		logrus.Fatal(err)
	}

	_ = mainServer.Run()
}
