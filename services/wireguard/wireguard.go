package wireguard

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/viper"

	"github.com/sentinel-official/cli-client/services/wireguard/types"
	clienttypes "github.com/sentinel-official/cli-client/types"
)

var (
	_ clienttypes.Service = (*WireGuard)(nil)
)

type WireGuard struct {
	cfg  *types.Config
	info []byte
}

func NewWireGuard(cfg *types.Config, info []byte) *WireGuard {
	return &WireGuard{
		cfg:  cfg,
		info: info,
	}
}

func (s *WireGuard) home() string { return viper.GetString(flags.FlagHome) }

func (s *WireGuard) configFilePath() string {
	return filepath.Join(s.home(), types.DefaultConfigFileName)
}

func (s *WireGuard) Info() []byte { return s.info }

func (s *WireGuard) IsUp() bool {
	iFace, err := s.realInterface()
	if err != nil {
		return false
	}

	cmd := exec.Command(
		s.execFile("wg"),
		strings.Split(
			fmt.Sprintf("show %s", iFace),
			" ",
		)...,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	if strings.Contains(string(output), "No such device") {
		return false
	}

	return true
}

func (s *WireGuard) PreUp() error {
	cfgFilePath := s.configFilePath()
	return s.cfg.WriteToFile(cfgFilePath)
}

func (s *WireGuard) PostUp() error  { return nil }
func (s *WireGuard) PreDown() error { return nil }

func (s *WireGuard) PostDown() error {
	cfgFilePath := s.configFilePath()
	if _, err := os.Stat(cfgFilePath); err != nil {
		return nil
	}

	return os.Remove(cfgFilePath)
}

func (s *WireGuard) Transfer() (u int64, d int64, err error) {
	iFace, err := s.realInterface()
	if err != nil {
		return 0, 0, err
	}

	cmd := exec.Command(
		s.execFile("wg"),
		strings.Split(
			fmt.Sprintf("show %s transfer", iFace),
			" ",
		)...,
	)

	output, err := cmd.Output()
	if err != nil {
		return 0, 0, err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		columns := strings.Split(line, "\t")
		if len(columns) != 3 {
			continue
		}

		d, err = strconv.ParseInt(columns[1], 10, 64)
		if err != nil {
			return 0, 0, err
		}

		u, err = strconv.ParseInt(columns[2], 10, 64)
		if err != nil {
			return 0, 0, err
		}

		return d, u, nil
	}

	return 0, 0, nil
}
