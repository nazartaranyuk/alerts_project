package repository

import (
	"context"
	"nazartaraniuk/alertsProject/internal/app/db"
	"nazartaraniuk/alertsProject/internal/domain"
)

type AlarmsRepositoryLocal struct {
	Database *db.Database
}

func NewAlarmsRepositoryLocal(db *db.Database) *AlarmsRepositoryLocal {
	return &AlarmsRepositoryLocal{
		Database: db,
	}
}

func (r *AlarmsRepositoryLocal) SaveAlarms(ctx context.Context, alarms []domain.RegionAlarmInfo) error {
	err := r.Database.SaveRegions(ctx, alarms)
	if err != nil {
		return err
	}
	return nil
}
