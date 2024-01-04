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
	fmt.Printf("Using connection string: %s\n", connString)
	fmt.Println("Connecting to database...")
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, err
	}
	fmt.Println("Connection to database completed successfully!")
	return &Repository{conn: conn}, nil
}

func (r *Repository) InitializeTables() error {
	_, err := r.conn.Exec(context.Background(), createVersionsTableQuery)
	return err
}

func (r *Repository) Ping() error {
	return r.conn.Ping(context.Background())
}

func (r *Repository) Close() {
	r.conn.Close(context.Background())
}
