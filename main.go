package main

import (
	"log/slog"
	"os"
)

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		slog.Error("error getting hostname for local system", err.Error())
		return
	}

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	logger.Info("starting disbursement service on", "hostname", hostname)

	if err != nil {
		logger.Error(err.Error())
		return
	}
}
