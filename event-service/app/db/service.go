package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	DriverSQLMock = "sqlmock"
	DriverMySQL   = "mysql"
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

func newConnectionURL(cfg *config) (string, error) {
	switch cfg.Driver {
	case DriverSQLMock:
		return "sqlmock", nil
	case DriverMySQL:
		return fmt.Sprintf("%s:%s@(%s:%s)/%s?parseTime=true", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database), nil
	default:
		return "", errors.Errorf("unsupported database driver: %s", cfg.Driver)
	}
}

func NewDatabase(ctx context.Context) (*sqlx.DB, error) {
	cfg := newConfig()

	connURL, err := newConnectionURL(newConfig())
	if err != nil {
		return nil, err
	}

	db, err := sqlx.ConnectContext(ctx, cfg.Driver, connURL)
	if err != nil {
		return nil, err
	}

	return db, nil
}
