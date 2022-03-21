package pgx

import (
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
}

func NewPGXDB(cfg Config) (*bun.DB, error) {
	config, err := pgx.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=prefer", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName))

	if err != nil {
		return nil, err
	}

	sqldb := stdlib.OpenDB(*config)
	db := bun.NewDB(sqldb, pgdialect.New())

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}
