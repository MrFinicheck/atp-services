package app

import (
	"os"
	"path/filepath"
)

func resolveDataDir(override string) string {
	if override != "" {
		return override
	}
	if v := os.Getenv("ATP_DATA_DIR"); v != "" {
		return v
	}
	// Local ./data avoids lock conflicts with a stale web server during wails dev.
	if cwd, err := os.Getwd(); err == nil {
		local := filepath.Join(cwd, "data")
		if err := os.MkdirAll(local, 0o755); err == nil {
			return local
		}
	}
	home, err := os.UserConfigDir()
	if err != nil {
		return "data"
	}
	return filepath.Join(home, "atp-services")
}
