package utils

import (
	"os"
	"log"
)

func GetLogger(prefix string) *log.Logger {
	return log.New(os.Stdout, prefix, log.LstdFlags)
}
