package utils

import (
	"bufio"
	"os"
)

func ReadLineFromFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}

	reader := bufio.NewReader(file)

	buf, _, err := reader.ReadLine()
	if err != nil {
		return "", err
	}

	return string(buf), nil
}
