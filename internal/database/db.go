package database

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"

	"gitlab-code-review-notifier/internal"
	"gitlab-code-review-notifier/pkg/envutil"
	"gitlab-code-review-notifier/pkg/log"
)

type db struct {
	migrationsPath string
	*sqlx.DB
	log.Loggable
}

func NewDb(host string, port string, user string, password string, dbname string) (*db, error) {
	connectionStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)
	return NewDbFromUrl(connectionStr)
}

func NewDbFromUrl(url string) (*db, error) {
	sqlxdb, err := sqlx.Connect("postgres", url)
	if err != nil {
		return nil, err
	}
	return &db{
		migrationsPath: "file://internal/database/migrations",
		DB:             sqlxdb,
	}, nil
}

func NewDbFromEnv() (*db, error) {
	var db *db
	var err error

	if dbUrl := envutil.GetEnvStr(internal.EnvDbUrl); len(dbUrl) > 0 {
		db, err = NewDbFromUrl(dbUrl)
	} else {
		dbHost := envutil.MustGetEnvStr(internal.EnvDbHost)
		dbPort := envutil.MustGetEnvStr(internal.EnvDbPort)
		dbUser := envutil.MustGetEnvStr(internal.EnvDbUser)
		dbPassword := envutil.MustGetEnvStr(internal.EnvDbPassword)
		dbName := envutil.MustGetEnvStr(internal.EnvDbName)
		db, err = NewDb(dbHost, dbPort, dbUser, dbPassword, dbName)
	}

	return db, err
}

func (d *db) Migrate() error {
	driver, err := postgres.WithInstance(d.DB.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("create new migration driver: %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance(d.migrationsPath, "postgres", driver)
	if err != nil {
		return fmt.Errorf("create new migration instance: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("run migrations: %v", err)
	}
	return nil
}
