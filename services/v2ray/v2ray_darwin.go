package v2ray

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	v2ray = "v2ray"
)

func (s *V2Ray) execFile(name string) string {
	return name
}

func (s *V2Ray) Up() error {
	cmd := exec.Command(
		s.execFile(v2ray),
		strings.Split(
			fmt.Sprintf("run --config %s", s.configFilePath()),
			" ",
		)...,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	s.cfg.PID = int32(cmd.Process.Pid)
	return nil
}
