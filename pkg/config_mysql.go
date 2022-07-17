package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql" // mysql driver
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const ConnectionMaxLifeTimeMinutes = 5

type Executor interface {
	Get(dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}

type DB interface {
	Executor
}

type MainConnection sqlx.DB

type Connection struct {
	Main *sqlx.DB
}

func NewDB(conn Connection) DB {
	return &executorImpl{
		db: conn,
	}
}

func InitConnection(ctx context.Context, l *zap.Logger) (Connection, func(), error) {
	master, err := InitMainConnection(ctx, l)
	if err != nil {
		l.Error("Can't initialize master database connection", zap.Error(err))
	}
	con := Connection{
		Main: (*sqlx.DB)(master),
	}
	cleanup := func() {
		l.Info("shutting down master and slave db connections")
		if err := master.Close(); err != nil {
			l.Error("Can't close master db connection", zap.Error(err))
		}
	}
	return con, cleanup, err
}

func InitMainConnection(ctx context.Context, logger *zap.Logger) (*MainConnection, error) {
	db, err := connect(ctx)
	if err != nil {
		return nil, err
	}
	return (*MainConnection)(db), nil
}

func connect(ctx context.Context) (*sqlx.DB, error) {
	uri := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?parseTime=true&charset=utf8mb4",
		"root",
		"root",
		"localhost",
		3310,
		"volt")
	db, err := sqlx.ConnectContext(ctx, "mysql", uri)
	if err != nil {
		return nil, errors.Wrap(err, "error connect to db: ")
	}
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(ConnectionMaxLifeTimeMinutes * time.Minute)
	err = db.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "error ping db: ")
	}
	return db, err
}

func (db Connection) CLose() error {
	var wrappedErr []string
	if err := db.Main.Close(); err != nil {
		wrappedErr = append(wrappedErr, err.Error())
	}
	if len(wrappedErr) > 0 {
		return fmt.Errorf("cant close db: %s", strings.Join(wrappedErr, "; "))
	}
	return nil
}
