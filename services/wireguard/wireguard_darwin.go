package wireguard

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/alessio/shellescape"
)

func (w *WireGuard) PreUp() error {
	return w.cfg.WriteToFile(w.Home())
}

func (w *WireGuard) RealInterface() (string, error) {
	nameFile, err := os.Open(
		fmt.Sprintf("/var/run/wireguard/%s.name", shellescape.Quote(w.cfg.Name)))
	if err != nil {
		return "", err
	}

	var (
		scanner = bufio.NewReader(nameFile)
	)

	line, err := scanner.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.Trim(line, "\n"), nil
}
