package main

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log/slog"
	"os"

	d "github.com/levtk/sequra/disburse"
)

const (
	driverName = "sqlite3"
	DSN        = "disbursement.DSN"
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
		logger.Error("failed to instantiate the disburser service on host: ", hostname)
	}
	logger.Info("creating tables for disbursement database...")
	err = disburserService.Repo.CreateTables()
	if err != nil {
		logger.Error("failed to create tables for disbursement app. Does db exist? ", err.Error())
	}

	//TODO write func to check if orders were already imported. Store in db table hash of file
	orders, merchants, err := disburserService.Importer.ImportOrders()
	if err != nil {
		logger.Error("failed to import orders or merchants", err.Error())
	}

	for _, v := range merchants {
		err := disburserService.Repo.InsertMerchant(v) //TODO fix imports by making models module
		if err != nil {
			logger.Error(err.Error())
		}

	}
	for _, v := range orders {
		err := disburserService.ProcessOrder.ProcessOrder(logger, ctx, disburserService.Repo, &v)
		if err != nil {
			logger.Error(err.Error())
		}
	}

	if err != nil {
		logger.Error(err.Error())
		return
	}
}
