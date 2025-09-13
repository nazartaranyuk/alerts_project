package repository

import (
	"log"
	"nazartaraniuk/alertsProject/internal/app/client"
	"nazartaraniuk/alertsProject/internal/domain"
)

type AlarmsRepository struct {
	client client.Client
}

func (rep *AlarmsRepository) GetCurrentAlerts() ([]domain.RegionAlarmInfo, error) {
	response, err := rep.client.GetCurrentAlerts()
	if err != nil {
		log.Printf("Service cannot get alarms info: %v", err)
		return nil, err
	}

	return response, nil
}

func NewAlarmsRepository(client *client.Client) *AlarmsRepository {
	return &AlarmsRepository{
		client: *client,
	}
}
