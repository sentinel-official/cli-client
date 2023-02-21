package v2ray

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/spf13/viper"

	"github.com/sentinel-official/cli-client/services/v2ray/types"
	clienttypes "github.com/sentinel-official/cli-client/types"
)

var (
	_ clienttypes.Service = (*V2Ray)(nil)
)

type V2Ray struct {
	cfg *types.Config
}

func NewV2Ray(cfg *types.Config) *V2Ray {
	return &V2Ray{
		cfg: cfg,
	}
}

func (s *V2Ray) home() string           { return viper.GetString(flags.FlagHome) }
func (s *V2Ray) configFilePath() string { return filepath.Join(s.home(), types.DefaultConfigFileName) }
func (s *V2Ray) pid() int32             { return s.cfg.PID }

func (s *V2Ray) Info() []byte {
	buf, err := json.Marshal(s.cfg)
	if err != nil {
		panic(err)
	}

	return buf
}

func (s *V2Ray) PreUp() error {
	cfgFilePath := s.configFilePath()
	return s.cfg.WriteToFile(cfgFilePath)
}

func (s *V2Ray) IsUp() bool {
	ok, err := process.PidExists(s.pid())
	if err != nil {
		return false
	}
	if !ok {
		return false
	}

	proc, err := process.NewProcess(s.pid())
	if err != nil {
		return false
	}

	ok, err = proc.IsRunning()
	if err != nil {
		return false
	}
	if !ok {
		return false
	}

	name, err := proc.Name()
	if err != nil {
		return false
	}
	if name != v2ray {
		return false
	}

	return true
}

func (s *V2Ray) PostUp() error  { return nil }
func (s *V2Ray) PreDown() error { return nil }

func (s *V2Ray) Down() error {
	proc, err := process.NewProcess(s.pid())
	if err != nil {
		return err
	}

	return proc.Kill()
}

func (s *V2Ray) PostDown() error {
	cfgFilePath := s.configFilePath()
	if _, err := os.Stat(cfgFilePath); err != nil {
		return nil
	}

	return os.Remove(cfgFilePath)
}

func (s *V2Ray) Transfer() (int64, int64, error) {
	return 0, 0, nil
}
