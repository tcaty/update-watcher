package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	Pool    *pgxpool.Pool
	Builder squirrel.StatementBuilderType
}

func New(connString string) (*Postgres, error) {
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, err
	}
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	return &Postgres{
		Pool:    pool,
		Builder: builder,
	}, err
}

func (p *Postgres) ConnString() string {
	return fmt.Sprintf(
		"postgres://%s:xxx@%s:%d/%s",
		p.Pool.Config().ConnConfig.User,
		p.Pool.Config().ConnConfig.Host,
		p.Pool.Config().ConnConfig.Port,
		p.Pool.Config().ConnConfig.Database,
	)
}

func (p *Postgres) Close() {
	p.Pool.Close()
}

func (p *Postgres) Ping() error {
	return p.Pool.Ping(context.Background())
}
