package utils

import (
	"bufio"

	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
)

func GetPassword(backend string, r *bufio.Reader) (string, error) {
	if backend == keyring.BackendFile {
		password, err := input.GetPassword("Enter keyring passphrase: ", r)
		if err != nil {
			return "", err
		}

		return password, nil
	}

	return "", nil
}
