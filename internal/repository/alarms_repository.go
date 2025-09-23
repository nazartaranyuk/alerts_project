package repository

import (
	"nazartaraniuk/alertsProject/internal/app/client"
	"nazartaraniuk/alertsProject/internal/domain"

	"github.com/sirupsen/logrus"
)

type AlarmsRepository struct {
	client client.Client
}

func (rep *AlarmsRepository) GetCurrentAlerts() ([]domain.RegionAlarmInfo, error) {
	response, err := rep.client.GetCurrentAlerts()
	if err != nil {
		logrus.Printf("Service cannot get alarms info: %v", err)
		return []domain.RegionAlarmInfo{}, err
	}

	return response, err
}

func (rep *AlarmsRepository) GetAlarmInfoByRegion(regionID string) (domain.RegionAlarmInfo, error) {
	response, err := rep.client.GetAlarmInfoByRegion(regionID)
	if err != nil {
		logrus.Printf("Service cannot get alarm info by region: %v", err)
	}

	return response, err
}

func NewAlarmsRepository(client *client.Client) *AlarmsRepository {
	return &AlarmsRepository{
		client: *client,
	}
}
