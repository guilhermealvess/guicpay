package logger

import (
	"log"

	"go.uber.org/zap"
)

var Logger *zap.Logger

func init() {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	Logger = zapLogger
}
