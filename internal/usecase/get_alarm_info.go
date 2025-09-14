package usecase

import (
	"nazartaraniuk/alertsProject/internal/domain"
	"nazartaraniuk/alertsProject/internal/repository"
)

type GetAlarmInfoService struct {
	alarmRepository repository.AlarmsRepository
}

func (s *GetAlarmInfoService) GetCurrentAlerts() ([]domain.RegionAlarmInfo, error) {
	return s.alarmRepository.GetCurrentAlerts()
}

func (s *GetAlarmInfoService) GetAlarmInfoByRegion(regionID string) (domain.RegionAlarmInfo, error) {
	return s.alarmRepository.GetAlarmInfoByRegion(regionID)
}

func NewGetAlarmInfoService(repository repository.AlarmsRepository) *GetAlarmInfoService {
	return &GetAlarmInfoService{
		alarmRepository: repository,
	}
}
