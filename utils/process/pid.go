package process

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/natefinch/atomic"
)

func ReadPID(path string) (int, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}

	id, err := strconv.Atoi(string(bytes))
	if err != nil {
		return 0, err
	}

	return id, nil
}

func WritePID(path string) error {
	return atomic.WriteFile(
		path,
		strings.NewReader(
			fmt.Sprintf("%d", os.Getpid()),
		),
	)
}
