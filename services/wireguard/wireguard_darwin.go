package wireguard

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func (s *WireGuard) realInterface() (string, error) {
	nameFile, err := os.Open(fmt.Sprintf("/var/run/wireguard/%s.name", s.cfg.Name))
	if err != nil {
		return "", err
	}

	reader := bufio.NewReader(nameFile)

	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.Trim(line, "\n"), nil
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
