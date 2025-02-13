package store

import (
	"database/sql"
	"fmt"
	"main/internal/config"

	_ "github.com/lib/pq"
)

type Store struct {
	DB *sql.DB
}

func NewStore(cfg config.Config) (*Store, error) {
	var err error
	store := Store{}
	store.DB, err = initPostgres(cfg)
	if err != nil {
		return nil, err
	}
	return &store, err
}

func initPostgres(cfg config.Config) (*sql.DB, error) {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Dbname)
	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("The database is connected")
	return db, err
}
