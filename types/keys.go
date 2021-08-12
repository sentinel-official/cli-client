package types

import (
	"os"
	"path/filepath"
)

const (
	APIPathPrefix  = "/api/v1"
	StatusFilename = "status.json"
	Listen         = "127.0.0.1:11112"
)

var (
	Home = func() string {
		home, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}

		return filepath.Join(home, ".sentinelcli")
	}()
)
