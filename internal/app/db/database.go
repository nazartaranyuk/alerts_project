package db

import (
	"context"
	"database/sql"
	"nazartaraniuk/alertsProject/internal/domain"
	"time"

	_ "github.com/lib/pq"
)

type Database struct {
	connectionString string
	DB               *sql.DB
}

func NewDatabase(adminDSN, appDSN, dbName string) (*Database, error) {
	admin, err := sql.Open("postgres", adminDSN) // adminDSN: postgres://user:pass@host:5432/postgres?sslmode=disable
	if err != nil {
		return nil, err
	}
	defer admin.Close()

	if err := admin.Ping(); err != nil {
		return nil, err
	}

	var exists bool
	if err := admin.QueryRow(`SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname=$1)`, dbName).Scan(&exists); err != nil {
		return nil, err
	}

	app, err := sql.Open("postgres", appDSN)
	if err != nil {
		return nil, err
	}
	if err := app.Ping(); err != nil {
		_ = app.Close()
		return nil, err
	}

	app.SetMaxOpenConns(10)
	app.SetMaxIdleConns(5)
	app.SetConnMaxLifetime(30 * time.Minute)

	return &Database{
		connectionString: appDSN,
		DB:               app,
	}, nil
}

func (r *Database) SaveRegions(ctx context.Context, regions []domain.RegionAlarmInfo) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	upsertRegion, err := tx.PrepareContext(ctx, `
INSERT INTO regions (region_id, region_type, region_name, region_eng_name, last_update)
VALUES ($1,$2,$3,$4,$5)
ON CONFLICT (region_id) DO UPDATE
SET region_type=EXCLUDED.region_type,
    region_name=EXCLUDED.region_name,
    region_eng_name=EXCLUDED.region_eng_name,
    last_update=EXCLUDED.last_update`)
	if err != nil {
		return err
	}
	defer upsertRegion.Close()

	insertAlarm, err := tx.PrepareContext(ctx, `
INSERT INTO region_active_alarms (region_id, region_type, type, last_update)
VALUES ($1,$2,$3,$4)
ON CONFLICT (region_id, type) DO UPDATE
SET region_type=EXCLUDED.region_type,
    last_update=EXCLUDED.last_update`)
	if err != nil {
		return err
	}
	defer insertAlarm.Close()

	for _, region := range regions {
		if _, err = upsertRegion.ExecContext(ctx,
			region.RegionID, region.RegionType, region.RegionName,
			region.RegionEngName, region.LastUpdate); err != nil {
			return err
		}
		if _, err = tx.ExecContext(ctx, `DELETE FROM region_active_alarms WHERE region_id=$1`, region.RegionID); err != nil {
			return err
		}
		for _, alert := range region.ActiveAlerts {
			rtype := alert.RegionType
			if rtype == "" {
				rtype = region.RegionType
			}
			if _, err = insertAlarm.ExecContext(ctx,
				region.RegionID, rtype, alert.Type, alert.LastUpdate); err != nil {
				return err
			}
		}
	}

	err = tx.Commit()
	return err
}
