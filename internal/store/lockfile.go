package store

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func acquireDirLock(dir string, maxWait time.Duration) (*os.File, error) {
	lockPath := filepath.Join(dir, ".atp.lock")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}

	deadline := time.Now().Add(maxWait)
	for {
		f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o600)
		if err != nil {
			return nil, err
		}
		if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err == nil {
			_, _ = f.Seek(0, 0)
			_ = f.Truncate(0)
			_, _ = fmt.Fprintf(f, "%d\n", os.Getpid())
			return f, nil
		}
		_ = f.Close()

		if pid, ok := readLockPID(dir); ok && !isProcessAlive(pid) {
			_ = os.Remove(lockPath)
			continue
		}

		if time.Now().After(deadline) {
			return nil, fmt.Errorf("каталог данных занят другим процессом — закройте другой экземпляр (wails dev / веб-сервер)")
		}
		time.Sleep(250 * time.Millisecond)
	}
}

func releaseDirLock(f *os.File) {
	if f == nil {
		return
	}
	_ = syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	_ = f.Close()
}

func readLockPID(dir string) (int, bool) {
	b, err := os.ReadFile(filepath.Join(dir, ".atp.lock"))
	if err != nil {
		return 0, false
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(b)))
	if err != nil {
		return 0, false
	}
	return pid, true
}

func isProcessAlive(pid int) bool {
	if pid <= 0 {
		return false
	}
	return syscall.Kill(pid, 0) == nil
}

func tryClearStaleLevelDBLock(dbPath string) {
	lockPath := filepath.Join(dbPath, "LOCK")
	f, err := os.OpenFile(lockPath, os.O_RDWR, 0o644)
	if err != nil {
		return
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err == nil {
		_ = syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
		_ = f.Close()
		_ = os.Remove(lockPath)
		return
	}
	_ = f.Close()
}
