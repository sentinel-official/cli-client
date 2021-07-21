package rand

import (
	"crypto/rand"
	"encoding/base64"
)

func RandomStringHex(l int) string {
	v := make([]byte, l)
	_, _ = rand.Read(v)

	return base64.StdEncoding.EncodeToString(v)
}
