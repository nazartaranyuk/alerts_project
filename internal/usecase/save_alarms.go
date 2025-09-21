package usecase

import (
	"context"
	"nazartaraniuk/alertsProject/internal/repository"
)

type SaveAlarmsService struct {
	alarmsRepository      repository.AlarmsRepository
	alarmsRepositoryLocal repository.AlarmsRepositoryLocal
}

func (s *SaveAlarmsService) SaveAlarmsInfo(ctx context.Context) (bool, error) {
	currentAlerts, err := s.alarmsRepository.GetCurrentAlerts()
	if err != nil {
		return false, err
	}
	err = s.alarmsRepositoryLocal.SaveAlarms(ctx, currentAlerts)
	if err != nil {
		return false, err
	}
	return true, nil
}

func NewSaveAlarmsService(
	alarmsRepository repository.AlarmsRepository,
	alarmsRepositoryLocal repository.AlarmsRepositoryLocal,
) *SaveAlarmsService {
	return &SaveAlarmsService{
		alarmsRepository:      alarmsRepository,
		alarmsRepositoryLocal: alarmsRepositoryLocal,
	}
}
