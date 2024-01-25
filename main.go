package main

import (
	"context"
	"database/sql"
	d "github.com/levtk/sequra/disburse"
	_ "github.com/mattn/go-sqlite3"
	"log/slog"
	"os"
)

const (
	driverName = "sqlite3"
	DSN        = "./disbursement.sqlite"
)

func main() {
	ctx := context.Background()
	hostname, err := os.Hostname()
	if err != nil {
		slog.Error("error getting hostname for local system", err.Error())
		return
	}

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	logger.Info("starting disbursement service on", "hostname", hostname)
	logger.Info("staring database...")
	db, err := sql.Open(driverName, DSN)
	if err != nil {
		logger.Error("failed to connect to db", err.Error())
	}

	disburserService, err := d.NewDisburserService(logger, ctx, db)
	if err != nil {
		logger.Error("failed to instantiate the disburser service on ", "hostname", hostname, "error", err.Error())
	}

	//TODO write func to check if orders were already imported. Store in db table hash of file
	distributions, merchants, monthly, err := disburserService.Importer.ImportOrders()
	if err != nil {
		logger.Error("failed to import orders or merchants", err.Error())
	}

	for _, v := range merchants {
		err := disburserService.Repo.InsertMerchant(v) //TODO fix imports by making types module
		if err != nil {
			logger.Error(err.Error())
		}

	}

	err = disburserService.ProcessOrder.ProcessBatchMonthly(monthly)
	if err != nil {
		logger.Error("failed to process batch monthly records", "error", err)
	}

	err = disburserService.ProcessOrder.ProcessBatchDistributions(distributions)
	if err != nil {
		logger.Error(err.Error())
	}
}
