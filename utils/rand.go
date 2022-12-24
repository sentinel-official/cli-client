package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func RandomHexString(l int) string {
	v := make([]byte, l)
	_, _ = rand.Read(v)

	return base64.StdEncoding.EncodeToString(v)
}
