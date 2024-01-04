package main

import (
	"fmt"
	"log/slog"
	"os"
)

const (
	RATE_LESS_THAN_50       int64 = 10
	RATE_BETWEEN_50_AND_300 int64 = 5
	RATE_ABOVE_300          int64 = 25
)

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		slog.Error("error getting hostname for local system", err.Error())
		return
	}
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	logger.Info("starting disbursement service on", "hostname", hostname)
	totalPurchaseFees, err := calculateTotalPurchaseFee(30100)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	fmt.Printf("The total fee for the single purchase is: €%.2f", float32(totalPurchaseFees)/100)

}

func calculateTotalPurchaseFee(totalPurchase int64) (totalPurchaseFee int64, err error) {
	if totalPurchase > 0 && totalPurchase < 5000 {
		totalPurchaseFee = RATE_LESS_THAN_50 * totalPurchase / 100
	}

	if totalPurchase > 5000 && totalPurchase < 30000 {
		totalPurchaseFee = RATE_BETWEEN_50_AND_300 * totalPurchase / 100
	}

	if totalPurchase > 30000 {
		totalPurchaseFee = RATE_ABOVE_300 * totalPurchase / 1000
	}

	return totalPurchaseFee, nil
}
