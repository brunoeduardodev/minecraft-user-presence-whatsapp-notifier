package loggerstore

import (
	"context"
	"database/sql"

	"github.com/brunoeduardodev/minecraft-user-presence-whatsapp-notifier/internal/logger"
)

type LoggerStore struct {
	db *sql.DB
}

func Connect(file string) (*LoggerStore, error) {
	db, err := sql.Open("sqlite3", file)

	if err != nil {
		return nil, err
	}

	return &LoggerStore{
		db: db,
	}, nil
}

func (l *LoggerStore) Close() error {
	return l.db.Close()
}

func (l *LoggerStore) Prepare(ctx context.Context) error {
	logger.Info(ctx, "preparing database")
	stmt, err := l.db.PrepareContext(ctx, "CREATE TABLE IF NOT EXISTS logs (id INTEGER PRIMARY KEY AUTOINCREMENT, log TEXT, Timestamp DATETIME DEFAULT CURRENT_TIMESTAMP)")
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return err
	}

	logger.Info(ctx, "database prepared")

	return nil
}

func (l *LoggerStore) AddLog(ctx context.Context, log string) error {
	stmt, err := l.db.PrepareContext(ctx, "INSERT INTO logs (log) VALUES (?)")
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, log)
	if err != nil {
		return err
	}

	return nil
}

func (l *LoggerStore) CheckIfExists(ctx context.Context, log string) (bool, error) {
	// check if log exists in the database, and was created in the last 2 days
	stmt, err := l.db.PrepareContext(ctx, "SELECT * FROM logs WHERE log = ? AND timestamp > datetime('now', '-2 day')")
	if err != nil {
		return false, err
	}

	rows, err := stmt.QueryContext(ctx, log)
	if err != nil {
		return false, err
	}

	defer rows.Close()

	for rows.Next() {
		return true, nil
	}

	return false, nil
}
