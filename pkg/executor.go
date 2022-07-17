package repository

import (
	"context"
	"database/sql"
)

type executorImpl struct {
	db Connection
}

func (i *executorImpl) Get(dest interface{}, query string, args ...interface{}) error {
	return i.db.Main.Get(dest, query, args...)
}

func (i *executorImpl) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	tx := FromContext(ctx)
	if tx != nil && tx.IsActive {
		return tx.Tx.GetContext(ctx, dest, query, args...)
	}
	return i.db.Main.GetContext(ctx, dest, query, args...)
}

func (i *executorImpl) Select(dest interface{}, query string, args ...interface{}) error {
	return i.db.Main.Select(dest, query, args...)
}

func (i *executorImpl) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	tx := FromContext(ctx)
	if tx != nil && tx.IsActive {
		return tx.Tx.SelectContext(ctx, dest, query, args...)
	}
	return i.db.Main.SelectContext(ctx, dest, query, args...)
}

func (i *executorImpl) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return i.db.Main.Query(query, args...)
}

func (i *executorImpl) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	tx := FromContext(ctx)
	if tx != nil && tx.IsActive {
		return tx.Tx.QueryContext(ctx, query, args...)
	}
	return i.db.Main.QueryContext(ctx, query, args...)
}

func (i *executorImpl) QueryRow(query string, args ...interface{}) *sql.Row {
	return i.db.Main.QueryRow(query, args...)
}

func (i *executorImpl) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	tx := FromContext(ctx)
	if tx != nil && tx.IsActive {
		return tx.Tx.QueryRowContext(ctx, query, args...)
	}
	return i.db.Main.QueryRowContext(ctx, query, args...)
}

func (i *executorImpl) Exec(query string, args ...interface{}) (sql.Result, error) {
	return i.db.Main.Exec(query, args...)
}

func (i *executorImpl) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	tx := FromContext(ctx)
	if tx != nil && tx.IsActive {
		return tx.Tx.ExecContext(ctx, query, args...)
	}
	return i.db.Main.ExecContext(ctx, query, args...)
}

func (i *executorImpl) NamedExec(query string, arg interface{}) (sql.Result, error) {
	return i.db.Main.NamedExec(query, arg)
}

func (i *executorImpl) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	tx := FromContext(ctx)
	if tx != nil && tx.IsActive {
		return tx.Tx.NamedExecContext(ctx, query, arg)
	}
	return i.db.Main.NamedExecContext(ctx, query, arg)
}
