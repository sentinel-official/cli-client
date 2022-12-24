package types

import (
	"fmt"

	hubtypes "github.com/sentinel-official/hub/types"
)

const (
	B  int64 = 1
	KB       = 1e3 * B
	MB       = 1e3 * KB
	GB       = 1e3 * MB
	TB       = 1e3 * GB
)

func ToReadableBytes(bytes int64, decimals int) (out string) {
	var (
		i    int64
		rem  int64
		unit string
	)

	switch {
	case bytes > TB:
		i = bytes / TB
		rem = bytes - (i * TB)
		unit = "TB"
	case bytes > GB:
		i = bytes / GB
		rem = bytes - (i * GB)
		unit = "GB"
	case bytes > MB:
		i = bytes / MB
		rem = bytes - (i * MB)
		unit = "MB"
	case bytes > KB:
		i = bytes / KB
		rem = bytes - (i * KB)
		unit = "KB"
	default:
		i = bytes / B
		rem = bytes - (i * B)
		unit = "B"
	}

	if decimals == 0 {
		return fmt.Sprintf("%d%s", i, unit)
	}

	width := 0
	switch {
	case rem > GB:
		width = 12
	case rem > MB:
		width = 9
	case rem > KB:
		width = 6
	default:
		width = 3
	}

	remStr := fmt.Sprintf("%d", rem)
	for iter := len(remStr); iter < width; iter++ {
		remStr = "0" + remStr
	}
	if decimals > len(remStr) {
		decimals = len(remStr)
	}

	return fmt.Sprintf("%d.%s%s", i, remStr[:decimals], unit)
}

type Bandwidth struct {
	Upload   int64 `json:"upload"`
	Download int64 `json:"download"`
}

func (b Bandwidth) String() string {
	return fmt.Sprintf(
		"%s,%s",
		ToReadableBytes(b.Upload, 2),
		ToReadableBytes(b.Download, 2),
	)
}

func NewBandwidthFromRaw(v hubtypes.Bandwidth) Bandwidth {
	return Bandwidth{
		Upload:   v.Upload.Int64(),
		Download: v.Download.Int64(),
	}
}
