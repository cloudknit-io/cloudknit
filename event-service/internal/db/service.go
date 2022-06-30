package db

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/compuzest/zlifecycle-event-service/internal/env"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type MigrationCommand string

const (
	DriverSQLMock                  = "sqlmock"
	DriverMySQL                    = "mysql"
	MigrateUp     MigrationCommand = "up"
	MigrateDown   MigrationCommand = "down"
)

type config struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	Schema   string
	Driver   string
}

func newConfig() *config {
	c := config{}

	c.Host = os.Getenv("DB_HOST")
	if c.Host == "" {
		c.Host = "localhost"
	}
	c.Port = os.Getenv("DB_PORT")
	if c.Port == "" {
		c.Port = "3306"
	}
	c.Username = os.Getenv("DB_USERNAME")
	if c.Username == "" {
		c.Username = "root"
	}
	c.Password = os.Getenv("DB_PASSWORD")
	if c.Password == "" {
		c.Password = "zlifecycle"
	}
	c.Database = os.Getenv("DB_NAME")
	if c.Database == "" {
		c.Database = "event"
	}
	c.Driver = os.Getenv("DB_DRIVER")
	if c.Driver == "" {
		c.Driver = "mysql"
	}

	return &c
}

func newConnectionURL(cfg *config, withProtocol bool) (string, error) {
	switch cfg.Driver {
	case DriverSQLMock:
		return "sqlmock", nil
	case DriverMySQL:
		if withProtocol {
			escapedPW := url.QueryEscape(cfg.Password)
			return fmt.Sprintf("mysql://%s:%s@tcp(%s:%s)/%s", cfg.Username, escapedPW, cfg.Host, cfg.Port, cfg.Database), nil
		}
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database), nil
	default:
		return "", errors.Errorf("unsupported database driver: %s", cfg.Driver)
	}
}

func NewDatabase(ctx context.Context) (*sqlx.DB, error) {
	cfg := newConfig()

	connURL, err := newConnectionURL(newConfig(), false)
	if err != nil {
		return nil, err
	}

	db, err := sqlx.ConnectContext(ctx, cfg.Driver, connURL)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Migrate(command MigrationCommand) (change bool, err error) {
	connURL, err := newConnectionURL(newConfig(), true)
	if err != nil {
		return false, err
	}

	m, err := migrate.New(env.Config().MigrationsDir, connURL)
	if err != nil {
		return false, errors.Wrap(err, "error creating migrations runner")
	}
	switch command {
	case MigrateUp:
		err = m.Up()
	case MigrateDown:
		err = m.Down()
	default:
		return false, errors.Errorf("invalid migration command: %s", command)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		return false, nil
	}
	if err != nil {
		return false, errors.Wrapf(err, "error executing migrate %s command", command)
	}
	return true, nil
}
