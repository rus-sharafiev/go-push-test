package db

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	pool *pgxpool.Pool
}

// -- Create instance -------------------------------------------------------------

var (
	connectOnce sync.Once
	Instance    Postgres
	ConnErr     error
)

func Connect(dbConnString string) error {
	connectOnce.Do(func() {

		pool, err := pgxpool.New(context.Background(), dbConnString)
		if err != nil {
			log.Fatalf("\x1b[2mPostgreSQL:\x1b[0m\x1b[31m unable to create database connection: %s \x1b[0m\n\n", err.Error())
			ConnErr = err
		}

		fmt.Println("\x1b[32mСonnection to the database has been established\x1b[0m")
		Instance = Postgres{pool}
	})
	return ConnErr
}

// -- Methods ---------------------------------------------------------------------

func (p *Postgres) Query(query *string, args ...any) (pgx.Rows, error) {
	rows, err := p.pool.Query(context.Background(), *query, args...)
	return rows, err
}

func (p *Postgres) QueryRow(query *string, args ...any) pgx.Row {
	return p.pool.QueryRow(context.Background(), *query, args...)
}

func (p *Postgres) PgxPoolClose() {
	p.pool.Close()
}
