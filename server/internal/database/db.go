package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/shlmvgleb/mtg-task/server/internal/config"
	log "github.com/sirupsen/logrus"
)

func New(ctx context.Context, config *config.PostgresConfig) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DbName,
	)

	db, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("error while creating a connection to database: %w", err)
	}

	if err = db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("error while ping a database: %w", err)
	}

	log.Infoln("successfully connected to postgres")

	err = migrateDb(connStr)
	if err != nil {
		return nil, fmt.Errorf("error while migrating a database: %w", err)
	}

	log.Infoln("postgres migrations successfully executed")
	return db, nil
}

func migrateDb(connStr string) error {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Errorln("close error:", err)
		}
	}()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil {
		log.Infof("migration message: %v", err)
	}

	return nil
}
