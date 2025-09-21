package db

import (
	"context"
	"database/sql"
	"fmt"
	"nazartaraniuk/alertsProject/internal/domain"
	"time"

	"github.com/lib/pq"
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
	if !exists {
		q := fmt.Sprintf(`CREATE DATABASE %s WITH OWNER %s TEMPLATE template0 ENCODING 'UTF8'`,
			pq.QuoteIdentifier(dbName), pq.QuoteIdentifier(dbOwnerFromDSN(appDSN)))
		if _, err := admin.Exec(q); err != nil {
			return nil, err
		}
	}

	app, err := sql.Open("postgres", appDSN) // appDSN: postgres://user:pass@host:5432/<dbName>?sslmode=disable
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

	stmts := []string{
		`CREATE TABLE IF NOT EXISTS regions (
		  region_id        TEXT PRIMARY KEY,
		  region_type      TEXT NOT NULL,
		  region_name      TEXT NOT NULL,
		  region_eng_name  TEXT NOT NULL,
		  last_update      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS region_active_alarms (
		  id           BIGSERIAL PRIMARY KEY,
		  region_id    TEXT NOT NULL REFERENCES regions(region_id) ON DELETE CASCADE,
		  region_type  TEXT NOT NULL,
		  type         TEXT NOT NULL,
		  last_update  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		  UNIQUE (region_id, type)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_region_active_alarms_region_id ON region_active_alarms(region_id)`,
	}
	for _, q := range stmts {
		if _, err := app.Exec(q); err != nil {
			_ = app.Close()
			return nil, err
		}
	}

	return &Database{
		connectionString: appDSN,
		DB:               app,
	}, nil
}

func dbOwnerFromDSN(_ string) string {
	return "myuser"
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
