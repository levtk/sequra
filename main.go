package main

import (
	"context"
	"database/sql"
	d "github.com/levtk/sequra/disburse"
	_ "github.com/mattn/go-sqlite3"
	"log/slog"
	"net/http"
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

	DisburserService, err := d.NewDisburserService(logger, ctx, db)
	if err != nil {
		logger.Error("failed to instantiate the disburser service on ", "hostname", hostname, "error", err.Error())
	}

	r := http.NewServeMux()

	r.HandleFunc("/disbursement", DisburserService.Reporter.GetDisbursementReport)
	r.HandleFunc("/import", DisburserService.Importer.Import)

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		logger.Error("failed to launch http server on port 8080", "error", err)
	}
}
