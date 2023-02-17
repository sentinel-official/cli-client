package v2ray

import (
	"os"
	"os/exec"
	"path/filepath"
)

const (
	v2ray = "v2ray.exe"
)

func (s *V2Ray) execFile(name string) string {
	return ".\\" + filepath.Join("V2Ray", name)
}

func (s *V2Ray) Up() error {
	cmd := exec.Command(
		s.execFile(v2ray),
		"run", "--config", s.configFilePath(),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	s.cfg.PID = int32(cmd.Process.Pid)
	return nil
}
