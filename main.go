package main

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	d "github.com/levtk/sequra/disburse"
	"github.com/spf13/viper"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	ctx := context.Background()
	hostname, err := os.Hostname()
	if err != nil {
		slog.Error("error getting hostname for local system", err.Error())
		return
	}

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	viper.SetConfigFile(".env")
	err = viper.ReadInConfig()
	if err != nil {
		logger.Error("failed to read config file", "error", err.Error())
	}

	logger.Info("starting disbursement service on", "hostname", hostname)
	logger.Info("connecting to database...")
	db, err := sqlx.Connect(viper.GetString("driver"), viper.GetString("DSN"))
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
