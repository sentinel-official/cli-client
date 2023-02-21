package wireguard

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func (s *WireGuard) realInterface() (string, error) {
	return s.cfg.Name, nil
}

func (s *WireGuard) execFile(name string) string {
	return name
}

func (s *WireGuard) Up() error {
	cmd := exec.Command(
		s.execFile("wg-quick"),
		strings.Split(
			fmt.Sprintf("up %s", s.configFilePath()),
			" ",
		)...,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (s *WireGuard) Down() error {
	cmd := exec.Command(
		s.execFile("wg-quick"),
		strings.Split(
			fmt.Sprintf("down %s", s.configFilePath()),
			" ",
		)...,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
