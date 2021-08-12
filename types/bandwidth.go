package types

import (
	"fmt"

	hubtypes "github.com/sentinel-official/hub/types"

	netutils "github.com/sentinel-official/cli-client/utils/net"
)

type Bandwidth struct {
	Upload   int64 `json:"upload"`
	Download int64 `json:"download"`
}

func (b Bandwidth) String() string {
	return fmt.Sprintf(
		"%s+%s",
		netutils.ToReadable(b.Upload, 2),
		netutils.ToReadable(b.Download, 2),
	)
}

func NewBandwidthFromRaw(v hubtypes.Bandwidth) Bandwidth {
	return Bandwidth{
		Upload:   v.Upload.Int64(),
		Download: v.Download.Int64(),
	}
}
