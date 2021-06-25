package types

import (
	"os"
	"path/filepath"
)

var (
	DefaultHomeDirectory = func() string {
		home, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}

		return filepath.Join(home, ".sentinelcli")
	}()
)
