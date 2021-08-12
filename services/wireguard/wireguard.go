package wireguard

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alessio/shellescape"

	"github.com/sentinel-official/cli-client/services/wireguard/types"
	clienttypes "github.com/sentinel-official/cli-client/types"
)

var (
	_ clienttypes.Service = (*WireGuard)(nil)
)

type WireGuard struct {
	cfg  *types.Config
	info []byte
	home string
}

func NewWireGuard() *WireGuard {
	return &WireGuard{}
}

func (w *WireGuard) WithConfig(v *types.Config) *WireGuard { w.cfg = v; return w }
func (w *WireGuard) WithInfo(v []byte) *WireGuard          { w.info = v; return w }
func (w *WireGuard) WithHome(v string) *WireGuard          { w.home = v; return w }

func (w *WireGuard) Info() []byte { return w.info }

func (w *WireGuard) IsUp() bool {
	iFace, err := w.RealInterface()
	if err != nil {
		return false
	}

	output, err := exec.Command("wg", strings.Split(
		fmt.Sprintf("show %s", shellescape.Quote(iFace)), " ")...).CombinedOutput()
	if err != nil {
		return false
	}
	if strings.Contains(string(output), "No such device") {
		return false
	}

	return true
}

func (w *WireGuard) PreUp() error {
	return w.cfg.WriteToFile(w.home)
}

func (w *WireGuard) Up() error {
	var (
		path = filepath.Join(w.home, fmt.Sprintf("%s.conf", w.cfg.Name))
		cmd  = exec.Command("wg-quick", strings.Split(
			fmt.Sprintf("up %s", shellescape.Quote(path)), " ")...)
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (w *WireGuard) PostUp() error  { return nil }
func (w *WireGuard) PreDown() error { return nil }

func (w *WireGuard) Down() error {
	var (
		path = filepath.Join(w.home, fmt.Sprintf("%s.conf", w.cfg.Name))
		cmd  = exec.Command("wg-quick", strings.Split(
			fmt.Sprintf("down %s", shellescape.Quote(path)), " ")...)
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (w *WireGuard) PostDown() error {
	path := filepath.Join(w.home, fmt.Sprintf("%s.conf", w.cfg.Name))
	if _, err := os.Stat(path); err == nil {
		return os.Remove(path)
	}

	return nil
}

func (w *WireGuard) Transfer() (u int64, d int64, err error) {
	iFace, err := w.RealInterface()
	if err != nil {
		return 0, 0, err
	}

	output, err := exec.Command("wg", strings.Split(
		fmt.Sprintf("show %s transfer", shellescape.Quote(iFace)), " ")...).Output()
	if err != nil {
		return 0, 0, err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		columns := strings.Split(line, "\t")
		if len(columns) != 3 {
			continue
		}

		d, err := strconv.ParseInt(columns[1], 10, 64)
		if err != nil {
			return 0, 0, err
		}

		u, err := strconv.ParseInt(columns[2], 10, 64)
		if err != nil {
			return 0, 0, err
		}

		return d, u, nil
	}

	return 0, 0, nil
}
