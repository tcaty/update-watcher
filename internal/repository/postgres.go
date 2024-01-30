package repository

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/tcaty/update-watcher/internal/config"
	"github.com/tcaty/update-watcher/internal/entities"
	"github.com/tcaty/update-watcher/pkg/postgres"
)

type Postgres struct {
	*postgres.Postgres
}

func New(cfg config.Postgresql) (*Postgres, error) {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	pstgrs, err := postgres.New(connString)
	if err != nil {
		return nil, err
	}

	return &Postgres{pstgrs}, nil
}

func (p *Postgres) UpdateVersionRecord(vr entities.VersionRecord) (bool, error) {
	doesVersionRecordExist, err := p.doesVersionRecordExist(vr)
	if err != nil {
		return false, err
	}

	if doesVersionRecordExist {
		doesVersionRecordNeedUpdate, err := p.doesVersionRecordNeedUpdate(vr)
		if err != nil {
			return false, err
		}

		if doesVersionRecordNeedUpdate {
			err = p.setVersion(vr)
			if err != nil {
				return false, err
			}
			return true, nil
		}
	} else {
		if err := p.insertVersionRecord(vr); err != nil {
			return false, err
		}
	}

	return false, nil
}

func (p *Postgres) doesVersionRecordExist(vr entities.VersionRecord) (bool, error) {
	sql, args, err := p.Builder.
		Select(id).
		From(versionRecords).
		Where(sq.Eq{target: vr.Target}).
		ToSql()

	if err != nil {
		return false, err
	}

	var id int
	row := p.Pool.QueryRow(context.Background(), sql, args...)
	err = row.Scan(&id)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return false, fmt.Errorf("unable to check version record for existing: %v", err)
	}

	exists := id > 0

	return exists, nil
}

func (p *Postgres) insertVersionRecord(vr entities.VersionRecord) error {
	sql, args, err := p.Builder.
		Insert(versionRecords).
		Columns(target, version).
		Values(vr.Target, vr.Version).
		ToSql()

	if err != nil {
		return err
	}

	row := p.Pool.QueryRow(context.Background(), sql, args...)
	err = row.Scan()

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("could not insert version record: %v", err)
	}

	return nil
}

func (p *Postgres) doesVersionRecordNeedUpdate(vr entities.VersionRecord) (bool, error) {
	sql, args, err := p.Builder.
		Select(version).
		From(versionRecords).
		Where(sq.Eq{target: vr.Target}).
		ToSql()

	if err != nil {
		return false, err
	}

	var version string
	row := p.Pool.QueryRow(context.Background(), sql, args...)
	err = row.Scan(&version)

	if err != nil {
		return false, fmt.Errorf("unable to check version record update status: %v", err)
	}

	needsUpdate := version != vr.Version

	return needsUpdate, nil
}

func (p *Postgres) setVersion(vr entities.VersionRecord) error {
	sql, args, err := p.Builder.
		Update(versionRecords).
		Set(version, vr.Version).
		Where(sq.Eq{target: vr.Target}).
		ToSql()

	if err != nil {
		return err
	}

	row := p.Pool.QueryRow(context.Background(), sql, args...)
	err = row.Scan()

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("could not update version record: %v", err)
	}

	return nil
}
