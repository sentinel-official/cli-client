package wireguard

import (
	"os"
	"os/exec"
	"path/filepath"
)

func (s *WireGuard) realInterface() (string, error) {
	return s.cfg.Name, nil
}

func (s *WireGuard) execFile(name string) string {
	return ".\\" + filepath.Join("WireGuard", name+".exe")
}

func (s *WireGuard) Up() error {
	var (
		cmd = exec.Command(
			s.execFile("wireguard"),
			"/installtunnelservice", s.configFilePath(),
		)
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
		s.execFile("wireguard"),
		"/uninstalltunnelservice", iFace,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
