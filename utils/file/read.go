package file

import (
	"bufio"
	"os"
)

func ReadLine(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}

	var (
		reader = bufio.NewReader(file)
	)

	buf, _, err := reader.ReadLine()
	if err != nil {
		return "", err
	}

	return string(buf), nil
}
