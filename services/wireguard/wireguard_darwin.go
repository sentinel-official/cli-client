package wireguard

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alessio/shellescape"
)

func (w *WireGuard) RealInterface() (string, error) {
	nameFile, err := os.Open(
		fmt.Sprintf("/var/run/wireguard/%s.name", shellescape.Quote(w.cfg.Name)))
	if err != nil {
		return "", err
	}

	scanner := bufio.NewReader(nameFile)

	line, err := scanner.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.Trim(line, "\n"), nil
}

func (w *WireGuard) ExecFile(name string) string {
	return name
}

func (w *WireGuard) Up() error {
	var (
		cfgFilePath = filepath.Join(w.Home(), fmt.Sprintf("%s.conf", w.cfg.Name))
		cmd         = exec.Command(w.ExecFile("wg-quick"), strings.Split(
			fmt.Sprintf("up %s", shellescape.Quote(cfgFilePath)), " ")...)
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (w *WireGuard) Down() error {
	var (
		cfgFilePath = filepath.Join(w.Home(), fmt.Sprintf("%s.conf", w.cfg.Name))
		cmd         = exec.Command(w.ExecFile("wg-quick"), strings.Split(
			fmt.Sprintf("down %s", shellescape.Quote(cfgFilePath)), " ")...)
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
