package types

import (
	"os"
	"path/filepath"
)

const (
	APIPathPrefix   = "/api/v1"
	BuildFolderName = "build"
	ConfigFilename  = "config.toml"
	StatusFilename  = "status.json"
	TokenLength     = 32
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
