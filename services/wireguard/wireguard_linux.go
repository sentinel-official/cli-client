package wireguard

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alessio/shellescape"
)

func (w *WireGuard) RealInterface() (string, error) {
	return w.cfg.Name, nil
}

func (w *WireGuard) ExecFile(name string) string {
	return name
}

func (w *WireGuard) Up() error {
	var (
		cfgFilePath = filepath.Join(w.home, fmt.Sprintf("%s.conf", w.cfg.Name))
		cmd         = exec.Command(w.ExecFile("wg-quick"), strings.Split(
			fmt.Sprintf("up %s", shellescape.Quote(cfgFilePath)), " ")...)
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (w *WireGuard) Down() error {
	iFace, err := w.RealInterface()
	if err != nil {
		return err
	}

	cmd := exec.Command(w.ExecFile("wg-quick"), strings.Split(
		fmt.Sprintf("down %s", shellescape.Quote(iFace)), " ")...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
