package repo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/levtk/sequra/disburse"
	_ "github.com/mattn/go-sqlite3"
	"log/slog"
	"time"
)

const (
	createDisbursementTable = `CREATE TABLE IF NOT EXISTS DISBURSEMENT (
    id TEXT NOT NULL PRIMARY KEY,
    disbursement_group_id TEXT,
    merchReference TEXT NOT NULL,
    order_id TEXT NOT NULL,
    order_fee INT NOT NULL,
    running_total INT);`

	createOrdersTable = `CREATE TABLE IF NOT EXISTS ORDERS (
    id TEXT NOT NULL PRIMARY KEY,
    merchant_reference TEXT NOT NULL,
    amount INT NOT NULL,
    created_at TEXT);`

	createMerchantsTable = `CREATE TABLE IF NOT EXISTS MERCHANTS (
    id TEXT PRIMARY KEY,
    reference TEXT,
    email TEXT,
    live_on TEXT,
    disbursement_frequency TEXT,
    minimum_monthly_fee TEXT);`

	getOrdersByMerchantUUID = `SELECT * FROM ORDERS WHERE id=:merchantUUID;`

	insertOrder = `INSERT INTO ORDERS(id, merchant_reference, amount, created_at) VALUES(?,?,?,?);`
)

type DisburserRepoRepository interface {
	GetOrdersByMerchantUUID(merchantUUID uuid.UUID) ([]disburse.Order, error)
	GetMerchantDisbursementsByRange(logger slog.Logger, merchant uuid.UUID, start time.Time, end time.Time)
	GetMerchant(merchantUUID uuid.UUID) (disburse.Merchant, error)
	InsertOrder(order disburse.Order) error
}

type DisburserRepo struct {
	db          *sql.DB
	ctx         context.Context
	logger      slog.Logger
	insertOrder *sql.Stmt
}

func NewDisburserRepo(l slog.Logger, ctx context.Context, db *sql.DB) (*DisburserRepo, error) {
	insOrderStmt, err := db.Prepare(insertOrder)
	if err != nil {
		return &DisburserRepo{}, err
	}
	return &DisburserRepo{
		db:          db,
		ctx:         ctx,
		logger:      l,
		insertOrder: insOrderStmt,
	}, nil
}

func (dr *DisburserRepo) GetOrdersByMerchantUUID(merchantUUID uuid.UUID) ([]disburse.Order, error) {
	//TODO Implement
	return []disburse.Order{}, errors.New("not implemented")
}

func (dr *DisburserRepo) GetMerchantDisbursementsByRange(logger slog.Logger, merchantUUID uuid.UUID, start time.Time, end time.Time) (disburse.Report, error) {
	//TODO Implement
	return disburse.Report{}, errors.New("not implemented")
}

func (dr *DisburserRepo) GetMerchant(merchantUUID uuid.UUID) (disburse.Merchant, error) {
	//TODO Implemement
	return disburse.Merchant{}, errors.New("not implemented")
}

func (dr *DisburserRepo) InsertOrder(o disburse.Order) error {
	_, err := dr.insertOrder.Exec(o.ID, o.MerchantReference, o.Amount, o.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
