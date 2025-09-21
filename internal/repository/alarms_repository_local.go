package repository

import (
	"context"
	"nazartaraniuk/alertsProject/internal/app/db"
	"nazartaraniuk/alertsProject/internal/domain"

	"github.com/sirupsen/logrus"
)

type AlarmsRepositoryLocal struct {
	Database db.Database
}

func NewAlarmsRepositoryLocal(adminDSN string, appDSN string, dbName string) *AlarmsRepositoryLocal {
	database, err := db.NewDatabase(adminDSN, appDSN, dbName)
	if err != nil {
		logrus.Fatal(err)
	}
	return &AlarmsRepositoryLocal{
		Database: *database,
	}
}

func (r *AlarmsRepositoryLocal) SaveAlarms(ctx context.Context, alarms []domain.RegionAlarmInfo) error {
	err := r.Database.SaveRegions(ctx, alarms)
	if err != nil {
		return err
	}
	return nil
}
