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
    running_total INT,
    payout_date TEXT,
    is_paid_out BOOLEAN);`

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

	getOrdersByMerchantReferenceID = `SELECT * FROM ORDERS WHERE merchant_reference=:merchRef`

	getMerchantByReferenceID = `SELECT * FROM MERCHANTS WHERE reference=:referenceID`

	insertOrder = `INSERT INTO ORDERS(id, merchant_reference, amount, created_at) VALUES(?,?,?,?);`

	insertDisbursement = `INSERT INTO DISBURSEMENT(id, disbursement_group_id, merchReference, order_id, order_fee, running_total, payout_date, is_paid_out)
	VALUES (?,?,?,?,?,?,?);`

	getDisbursementGroupID = `SELECT (disbursement_group_id) FROM DISBURSEMENT WHERE payout_date=:today AND merchReference=:merchRef`
)

type DisburserRepoRepository interface {
	GetOrdersByMerchantUUID(merchantUUID uuid.UUID) ([]disburse.Order, error)
	GetOrdersByMerchantReferenceID(ctx context.Context, merchRef string) ([]disburse.Order, error)
	GetMerchantDisbursementsByRange(logger slog.Logger, merchantUUID uuid.UUID, start time.Time, end time.Time) (disburse.Report, error)
	GetMerchant(merchantUUID uuid.UUID) (disburse.Merchant, error)
	GetMerchantByReferenceID(merchantReferenceID string) (disburse.Merchant, error)
	GetDisbursementGroupID(ctx context.Context, today string, merchRef string) (string, error)
	InsertOrder(order disburse.Order) error
	InsertDisbursement(disbursement Disbursement) (lastInsertID int64, err error)
}

type DisburserRepo struct {
	db                             *sql.DB
	ctx                            *context.Context
	logger                         *slog.Logger
	insertOrder                    *sql.Stmt
	insertDisbursement             *sql.Stmt
	getOrdersByMerchantReferenceID *sql.Stmt
	getMerchantByRefID             *sql.Stmt
	getDisbursementGroupID         *sql.Stmt
}

type Disbursement struct {
	ID                  string `DB:"ID"`
	DisbursementGroupID string `DB:"disbursement_group_id"`
	MerchReference      string `DB:"merchReference"`
	OrderID             string `DB:"order_id"`
	OrderFee            int64  `DB:"order_fee"`
	RunningTotal        int64  `DB:"running_total"`
	PayoutDate          string `DB:"payout_date"`
	IsPaidOut           bool   `DB:"is_paid_out"`
}

func NewDisburserRepo(l *slog.Logger, ctx *context.Context, db *sql.DB) (*DisburserRepo, error) {
	insOrderStmt, err := db.Prepare(insertOrder)
	if err != nil {
		return &DisburserRepo{}, err
	}

	insDisbursementStmt, err := db.Prepare(insertDisbursement)
	if err != nil {
		return &DisburserRepo{}, err
	}

	getOrdersByMerchRefID, err := db.Prepare(getOrdersByMerchantReferenceID)
	if err != nil {
		return &DisburserRepo{}, err
	}

	getMerchantByRefID, err := db.Prepare(getMerchantByReferenceID)
	if err != nil {
		return &DisburserRepo{}, err
	}

	getDisburseGroupID, err := db.Prepare(getOrdersByMerchantReferenceID)
	if err != nil {
		return &DisburserRepo{}, err
	}

	return &DisburserRepo{
		db:                             db,
		ctx:                            ctx,
		logger:                         l,
		insertOrder:                    insOrderStmt,
		insertDisbursement:             insDisbursementStmt,
		getOrdersByMerchantReferenceID: getOrdersByMerchRefID,
		getMerchantByRefID:             getMerchantByRefID,
		getDisbursementGroupID:         getDisburseGroupID,
	}, nil
}

func (dr *DisburserRepo) GetOrdersByMerchantUUID(merchantUUID uuid.UUID) ([]disburse.Order, error) {
	//TODO Implement
	return []disburse.Order{}, errors.New("not implemented")
}

func (dr *DisburserRepo) GetOrdersByMerchantReferenceID(ctx context.Context, merchRef string) ([]disburse.Order, error) {
	var orders []disburse.Order
	rows, err := dr.getOrdersByMerchantReferenceID.QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	counter := 0
	for rows.Next() {
		err = rows.Scan(orders[counter])
		if err != nil {
			return nil, err
		}
		counter++
	}
	return orders, nil
}
func (dr *DisburserRepo) GetMerchantDisbursementsByRange(logger slog.Logger, merchantUUID uuid.UUID, start time.Time, end time.Time) (disburse.Report, error) {
	//TODO Implement
	return disburse.Report{}, errors.New("not implemented")
}

func (dr *DisburserRepo) GetMerchant(merchantUUID uuid.UUID) (disburse.Merchant, error) {
	//TODO Implemement
	return disburse.Merchant{}, errors.New("not implemented")
}

func (dr *DisburserRepo) GetMerchantByReferenceID(merchantReferenceID string) (disburse.Merchant, error) {
	m := disburse.Merchant{}
	err := dr.getMerchantByRefID.QueryRow(merchantReferenceID).Scan(m)
	if err != nil {
		return m, err
	}
	return m, nil
}

// GetDisbursementGroupID returns the row with groupID if exists or err which should be ErrNoRows which tells us we need to create the groupID
func (dr *DisburserRepo) GetDisbursementGroupID(ctx context.Context, today string, merchRef string) (string, error) {
	var refId string
	row := dr.getDisbursementGroupID.QueryRowContext(ctx, today, merchRef)
	err := row.Err()
	if err != nil {
		return "", err
	}

	err = row.Scan(refId)
	if err != nil {
		return "", err
	}
	return refId, nil
}

func (dr *DisburserRepo) InsertOrder(o disburse.Order) error {
	_, err := dr.insertOrder.Exec(o.ID, o.MerchantReference, o.Amount, o.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (dr *DisburserRepo) InsertDisbursement(disbursement Disbursement) (lastInsertID int64, err error) {
	res, err := dr.insertDisbursement.Exec(disbursement)
	if err != nil {
		return 0, err
	}

	lID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lID, nil
}
