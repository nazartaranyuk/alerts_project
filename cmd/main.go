package main

import (
	"context"
	"fmt"
	"nazartaraniuk/alertsProject/internal/app/client"
	"nazartaraniuk/alertsProject/internal/app/db"
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

	database, err := db.NewDatabase(cfg.Database.AdminDSN, cfg.Database.AppDSN, cfg.Database.DbName)

	alarmsRepositoryLocal := repository.NewAlarmsRepositoryLocal(database)
	alarmsRepository := repository.NewAlarmsRepository(mainClient)
	userRepository := repository.NewUserRepository(database)

	alarmsInfoService := usecase.NewGetAlarmInfoService(*alarmsRepository)
	saveAlarmsService := usecase.NewSaveAlarmsService(*alarmsRepository, *alarmsRepositoryLocal)
	userService := usecase.NewUserService(userRepository)

	err = updateDatabase(context.Background(), 2*time.Minute, *saveAlarmsService)
	if err != nil {
		logrus.Println(err)
	}

	if err != nil {
		logrus.Println(err)
	}

	mainServer, err := server.NewServer(*cfg, alarmsInfoService, userService)
	if err != nil {
		logrus.Println(err)
	}

	_ = mainServer.Run()
}

func updateDatabase(ctx context.Context, timeout time.Duration, service usecase.SaveAlarmsService) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, err := service.SaveAlarmsInfo(ctx)
	if err != nil {
		return err
	}
	return nil
}
