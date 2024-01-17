package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/tcaty/update-watcher/internal/config"
)

type Repository struct {
	slog *slog.Logger
	conn *pgx.Conn
}

func New(cfg config.Postgresql) (*Repository, error) {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("could not connect to database %s: %v", cfg.Database, err)
	}

	r := &Repository{
		slog: slog.Default().With("database", cfg.Database),
		conn: conn,
	}

	return r, nil
}

func (r *Repository) Slog() *slog.Logger {
	return r.slog
}

func (r *Repository) Ping() error {
	return r.conn.Ping(context.Background())
}

func (r *Repository) InitializeTables() error {
	_, err := r.conn.Exec(context.Background(), createVersionsTableQuery)
	return err
}

func (r *Repository) UpdateVersionRecord(target string, version string) (bool, error) {
	doesVersionRecordExist, err := r.doesVersionRecordExist(target)
	if err != nil {
		return false, fmt.Errorf("update version record err: %v", err)
	}
	if doesVersionRecordExist {
		doesVersionRecordNeedUpdate, err := r.doesVersionRecordNeedUpdate(target, version)
		if err != nil {
			return false, fmt.Errorf("update version record err: %v", err)
		}
		if doesVersionRecordNeedUpdate {
			err = r.setVersion(target, version)
			if err != nil {
				return false, fmt.Errorf("update version record err: %v", err)
			}
			return true, nil
		}
	} else {
		if err := r.insertVersionRecord(target, version); err != nil {
			return false, fmt.Errorf("update version record err: %v", err)
		}
	}
	return false, nil
}

func (r *Repository) Close() {
	r.conn.Close(context.Background())
}

func (r *Repository) doesVersionRecordExist(target string) (bool, error) {
	var id int
	row := r.conn.QueryRow(context.Background(), selectIdByTargetQuery, target)
	err := row.Scan(&id)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return false, err
	}
	return id > 0, nil
}

func (r *Repository) insertVersionRecord(target string, version string) error {
	row := r.conn.QueryRow(context.Background(), insertVersionRecordsQuery, target, version)
	err := row.Scan()
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("could not insert version record: %v", err)
	}
	return nil
}

func (r *Repository) doesVersionRecordNeedUpdate(target string, newVersion string) (bool, error) {
	var version string
	row := r.conn.QueryRow(context.Background(), selectVersionByTargetQuery, target)
	err := row.Scan(&version)
	if err != nil {
		return false, fmt.Errorf("unable to check version record update status: %v", err)
	}
	return version != newVersion, nil
}

func (r *Repository) setVersion(target string, newVersion string) error {
	row := r.conn.QueryRow(context.Background(), setVersionByTargetQuery, newVersion, target)
	err := row.Scan()
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("could not update version record: %v", err)
	}
	return nil
}
