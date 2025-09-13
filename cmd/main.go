package main

import (
	"fmt"
	"nazartaraniuk/alertsProject/internal/app/client"
	"nazartaraniuk/alertsProject/internal/config"
	"nazartaraniuk/alertsProject/internal/repository"
	"nazartaraniuk/alertsProject/internal/service"
	"time"

	log "github.com/sirupsen/logrus"
)

const(
	CONFIG_PATH = "configs/config.dev.yaml"
)

func main() {
	cfg := config.LoadConfig(CONFIG_PATH)
	fmt.Println(cfg.Client.APIBaseURL)

	client := client.NewClient(
		cfg.Client.APIBaseURL,
		time.Second,
		cfg.Client.APIKey,
	)
	
	repository := repository.NewAlarmsRepository(client)
	service := usecase.NewGetAlarmInfoService(*repository)

	response, _ := service.GetCurrentAlerts()
	log.SetFormatter(&log.JSONFormatter{})
	log.Info(response)

}
