package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"url-shortener/internal/storage"

	sl "url-shortener/pkg/logger/slog"

	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

var once sync.Once
var dbConn *Storage
var log *slog.Logger = sl.GetLogger()

func Close(ctx context.Context, storage *Storage) error {
	var err error
	shutdown := make(chan struct{}, 1)
	go func() {
		err = storage.db.Close()
		shutdown <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("failed to clsose db: %w", ctx.Err())
	case <-shutdown:
	}

	return err
}

func New(storagePath string) (*Storage, error) {
	once.Do(func() {
		conn, err := InitDB(storagePath)
		if err != nil {
			log.Error("failed to init storage", sl.Err(err))
			return
		}
		dbConn = conn
	})
	return dbConn, nil
}

func InitDB(storagePath string) (*Storage, error) {
	const caller = "storage.sqlite.New"

	log = log.With(slog.String("caller", caller))

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", caller, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", caller, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", caller, err)
	}

	log.Info("initiated DB")
	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) error {
	const caller = "storage.sqlite.SaveURL"
	log = log.With(slog.String("caller", caller))

	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES(?, ?)")
	if err != nil {
		return fmt.Errorf("%s: %w", caller, err)
	}

	_, err = stmt.Exec(urlToSave, alias)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("%s: failed to insert: %w", caller, storage.ErrURLAlreadyExists)
		}

		return fmt.Errorf("%s: %w", caller, err)
	}

	log.Info("saved url", slog.String("url", urlToSave), slog.String("alias", alias))
	return nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const caller = "storage.sqlite.GetURL"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias=?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", caller, err)
	}

	var urlFromDB string
	err = stmt.QueryRow(alias).Scan(&urlFromDB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: failed to get: %w", caller, storage.ErrURLNotFound)
		}
		return "", fmt.Errorf("%s: %w", caller, err)
	}

	return urlFromDB, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const caller = "storage.sqlite.DeleteURL"
	log = log.With(slog.String("caller", caller))

	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias=? RETURNING id")
	if err != nil {
		return fmt.Errorf("%s: %w", caller, err)
	}

	var deletedID int
	err = stmt.QueryRow(alias).Scan(&deletedID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: failed to delete: %w", caller, storage.ErrURLNotFound)
		}
		return fmt.Errorf("%s: %w", caller, err)
	}

	log.Info("deleted alias", slog.String("alias", alias))
	return nil
}
