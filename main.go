package main

import (
	"fmt"
	"log/slog"
	"os"

	d "github.com/levtk/sequra/disburse"
)

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		slog.Error("error getting hostname for local system", err.Error())
		return
	}

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	logger.Info("starting disbursement service on", "hostname", hostname)

	err = d.Disburser.ProcessOrder()

	if err != nil {
		logger.Error(err.Error())
		return
	}

	fmt.Printf("The total fee for the single purchase is: €%.2f", float32(orderFee)/100)
}
