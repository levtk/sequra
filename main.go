package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
)

const (
	RATE_LESS_THAN_50       int64 = 10
	RATE_BETWEEN_50_AND_300 int64 = 5
	RATE_ABOVE_300          int64 = 25
	MAX_ORDER               int64 = 1000000
)

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		slog.Error("error getting hostname for local system", err.Error())
		return
	}
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	logger.Info("starting disbursement service on", "hostname", hostname)
	orderFee, err := calculateOrderFee(30100)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	fmt.Printf("The total fee for the single purchase is: €%.2f", float32(orderFee)/100)

}

func calculateOrderFee(order int64) (orderFee int64, err error) {
	if order > 0 && order < 5000 {
		orderFee = RATE_LESS_THAN_50 * order / 100
		return orderFee, nil
	}

	if order > 5000 && order < 30000 {
		orderFee = RATE_BETWEEN_50_AND_300 * order / 100
		return orderFee, nil
	}

	if order > 30000 {
		orderFee = RATE_ABOVE_300 * order / 1000
		return orderFee, nil
	}

	if order > MAX_ORDER {
		slog.Error("order submitted above max order value permitted")
		return -1, errors.New("order submitted above max order value permitted")
	}

	return orderFee, nil
}
