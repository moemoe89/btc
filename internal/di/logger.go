package di

import (
	"log"

	"github.com/moemoe89/btc/pkg/logging"
)

func GetLogger() logging.Logger {
	l, err := logging.NewLogger()
	if err != nil {
		log.Fatal(err)
	}

	return l
}
