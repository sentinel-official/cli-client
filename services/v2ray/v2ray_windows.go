package v2ray

import (
	"os"
	"os/exec"
	"path/filepath"
)

func (s *V2Ray) execFile(name string) string {
	return ".\\" + filepath.Join("V2Ray", name+".exe")
}

func (s *V2Ray) Up() error {
	cmd := exec.Command(
		s.execFile("v2ray"),
		"run", "--config", s.configFilePath(),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Start()
}
