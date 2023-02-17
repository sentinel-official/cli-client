package wireguard

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/sentinel-official/cli-client/services/wireguard/types"
)

func (s *WireGuard) realInterface() (string, error) {
	return types.DefaultInterface, nil
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
	iFace, err := s.realInterface()
	if err != nil {
		return err
	}

	cmd := exec.Command(
		s.execFile("wg-quick"),
		strings.Split(
			fmt.Sprintf("down %s", iFace),
			" ",
		)...,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
