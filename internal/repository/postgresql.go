package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/tcaty/update-watcher/internal/config"
)

type Repository struct {
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
	fmt.Printf("Connecting to database %s ...\n", connString)
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Connection to database %s completed successfully!\n", connString)
	return &Repository{conn: conn}, nil
}

func (r *Repository) Close() {
	r.conn.Close(context.Background())
}
