package store

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/storage"
)

func Open(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	dbPath := filepath.Join(dir, "data")
	db, err := openWithRetry(dbPath)
	if err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

func openWithRetry(path string) (*leveldb.DB, error) {
	opts := &opt.Options{ErrorIfMissing: false}
	var lastErr error
	for attempt := 0; attempt < 20; attempt++ {
		db, err := leveldb.OpenFile(path, opts)
		if err == nil {
			return db, nil
		}
		lastErr = err
		if isLockError(err) {
			time.Sleep(time.Duration(150*(attempt+1)) * time.Millisecond)
			continue
		}
		if strings.Contains(strings.ToLower(err.Error()), "corrupt") {
			return leveldb.RecoverFile(path, opts)
		}
		return nil, err
	}
	return nil, lastErr
}

func isLockError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, storage.ErrLocked) {
		return true
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "locked") ||
		strings.Contains(msg, "resource temporarily unavailable") ||
		strings.Contains(msg, "being used by another process")
}
