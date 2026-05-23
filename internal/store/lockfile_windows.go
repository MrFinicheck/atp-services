//go:build windows

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

const (
	processQueryLimitedInformation = 0x1000
	processStillActive             = 259
)

func acquireDirLock(dir string, maxWait time.Duration) (*os.File, error) {
	lockPath := filepath.Join(dir, ".atp.lock")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}

	deadline := time.Now().Add(maxWait)
	for {
		f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0o600)
		if err == nil {
			_, _ = fmt.Fprintf(f, "%d\n", os.Getpid())
			return f, nil
		}

		if !os.IsExist(err) {
			return nil, err
		}
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
	lockPath := f.Name()
	_ = f.Close()
	_ = os.Remove(lockPath)
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
	handle, err := syscall.OpenProcess(processQueryLimitedInformation, false, uint32(pid))
	if err != nil {
		return false
	}
	defer syscall.CloseHandle(handle)

	var code uint32
	if err := syscall.GetExitCodeProcess(handle, &code); err != nil {
		return false
	}
	return code == processStillActive
}

func tryClearStaleLevelDBLock(dbPath string) {
	_ = os.Remove(filepath.Join(dbPath, "LOCK"))
}
