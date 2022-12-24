package utils

import (
	"bufio"

	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
)

func ReadPassword(backend string, r *bufio.Reader) (s string, err error) {
	if backend == keyring.BackendFile {
		s, err = input.GetPassword("Enter keyring passphrase: ", r)
		if err != nil {
			return "", err
		}
	}

	return s, nil
}
