package common

import (
	"os"
)

func Success() {
	os.Exit(0)
}

func Failure(exitCode int) {
	os.Exit(exitCode)
}

func HandleError(err error, exitCode int) {
	if err != nil {
		Failure(exitCode)
	}
}
