package repo

import (
	"context"
	"database/sql"
	"errors"
	
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"log/slog"
	"time"
)

const (
	createDisbursementTable = `CREATE TABLE IF NOT EXISTS DISBURSEMENT (
    id TEXT NOT NULL PRIMARY KEY,
    disbursement_group_id TEXT,
    transaction_id TEXT, -- not implemented. when payout is confirmed by payment method/system the id for the payment transaction should be saved
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

	insertMerchant = `INSERT INTO MERCHANTS (id, reference, email, live_on, disbursement_frequency, minimum_monthly_fee) VALUES (
                    ?,?,?,?,?,?);`

	getOrdersByMerchantUUID = `SELECT * FROM ORDERS WHERE id=:merchantUUID;`

	getOrdersByMerchantReferenceID = `SELECT * FROM ORDERS WHERE merchant_reference=:merchRef`

	getMerchantByReferenceID = `SELECT * FROM MERCHANTS WHERE reference=:referenceID`

	insertOrder = `INSERT INTO ORDERS(id, merchant_reference, amount, created_at) VALUES(?,?,?,?);`

	insertDisbursement = `INSERT INTO DISBURSEMENT(id, disbursement_group_id, merchReference, order_id, order_fee, running_total, payout_date, is_paid_out)
	VALUES (?,?,?,?,?,?,?);`

	getDisbursementGroupID = `SELECT (disbursement_group_id) FROM DISBURSEMENT WHERE payout_date=:today AND merchReference=:merchRef`
)

type DisburserRepoRepository interface {
	GetOrdersByMerchantUUID(merchantUUID uuid.UUID) ([]Order, error)
	GetOrdersByMerchantReferenceID(ctx context.Context, merchRef string) ([]Order, error)
	GetMerchantDisbursementsByRange(logger slog.Logger, merchantUUID uuid.UUID, start time.Time, end time.Time) (Report, error)
	GetMerchant(merchantUUID uuid.UUID) (Merchant, error)
	GetMerchantByReferenceID(merchantReferenceID string) (Merchant, error)
	GetDisbursementGroupID(ctx context.Context, today string, merchRef string) (string, error)
	InsertOrder(order Order) error
	InsertDisbursement(disbursement Disbursement) (lastInsertID int64, err error)
	InsertMerchant(m Merchant) error
	CreateTables() error
}

type DisburserRepo struct {
	db                             *sql.DB
	ctx                            context.Context
	logger                         *slog.Logger
	insertOrder                    *sql.Stmt
	insertDisbursement             *sql.Stmt
	insertMerchant                 *sql.Stmt
	getOrdersByMerchantReferenceID *sql.Stmt
	getMerchantByRefID             *sql.Stmt
	getDisbursementGroupID         *sql.Stmt
	createDisbursementsTable       *sql.Stmt
	createOrdersTable              *sql.Stmt
	createMerchantsTable           *sql.Stmt
}

type DBDisbursement struct {
	ID                  string `DB:"ID"`
	DisbursementGroupID string `DB:"disbursement_group_id"`
	MerchReference      string `DB:"merchReference"`
	OrderID             string `DB:"order_id"`
	OrderFee            int64  `DB:"order_fee"`
	RunningTotal        int64  `DB:"running_total"`
	PayoutDate          string `DB:"payout_date"`
	IsPaidOut           bool   `DB:"is_paid_out"`
}

type DBOrder struct {
	ID                string    `json:"id,omitempty"`
	MerchantReference string    `json:"merchant_reference,omitempty"`
	MerchantID        uuid.UUID `json:"merchant_id,omitempty"`
	Amount            int64     `json:"amount,omitempty"`
	CreatedAt         time.Time `json:"created_at,omitempty"`
}

type DBMerchant struct {
	ID                    uuid.UUID `json:"id,omitempty"`
	Reference             string    `json:"reference,omitempty"`
	Email                 string    `json:"email,omitempty"`
	LiveOn                time.Time `json:"live_on,omitempty"`
	DisbursementFrequency string    `json:"disbursement_frequency,omitempty"`
	MinMonthlyFee         string    `json:"minimum_monthly_fee,omitempty"`
}

type DBReport struct {
	logger   *slog.Logger
	ctx      context.Context
	Name     string
	Merchant DBMerchant
	repo     DisburserRepoRepository
	Start    time.Time
	End      time.Time
	data     []byte
}

func NewDisburserRepo(l *slog.Logger, ctx context.Context, db *sql.DB) (*DisburserRepo, error) {
	insOrderStmt, err := db.Prepare(insertOrder)
	if err != nil {
		return &DisburserRepo{}, err
	}

	insDisbursementStmt, err := db.Prepare(insertDisbursement)
	if err != nil {
		return &DisburserRepo{}, err
	}

	insertMerchantStmt, err := db.Prepare(insertMerchant)
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

	getDisburseGroupID, err := db.Prepare(getDisbursementGroupID)
	if err != nil {
		return &DisburserRepo{}, err
	}

	createDisTable, err := db.Prepare(createDisbursementTable)
	if err != nil {
		return &DisburserRepo{}, err
	}

	createOrdersTabel, err := db.Prepare(createOrdersTable)
	if err != nil {
		return &DisburserRepo{}, err
	}

	createMerchTable, err := db.Prepare(createMerchantsTable)
	if err != nil {
		return &DisburserRepo{}, err
	}

	return &DisburserRepo{
		db:                             db,
		ctx:                            ctx,
		logger:                         l,
		insertOrder:                    insOrderStmt,
		insertDisbursement:             insDisbursementStmt,
		insertMerchant:                 insertMerchantStmt,
		getOrdersByMerchantReferenceID: getOrdersByMerchRefID,
		getMerchantByRefID:             getMerchantByRefID,
		getDisbursementGroupID:         getDisburseGroupID,
		createDisbursementsTable:       createDisTable,
		createMerchantsTable:           createMerchTable,
		createOrdersTable:              createOrdersTabel,
	}, nil
}

func (dr *DisburserRepo) GetOrdersByMerchantUUID(merchantUUID uuid.UUID) ([]Order, error) {
	//TODO Implement
	return []Order{}, errors.New("not implemented")
}

func (dr *DisburserRepo) GetOrdersByMerchantReferenceID(ctx context.Context, merchRef string) ([]Order, error) {
	var orders []Order
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
func (dr *DisburserRepo) GetMerchantDisbursementsByRange(logger slog.Logger, merchantUUID uuid.UUID, start time.Time, end time.Time) (Report, error) {
	//TODO Implement
	return Report{}, errors.New("not implemented")
}

func (dr *DisburserRepo) GetMerchant(merchantUUID uuid.UUID) (Merchant, error) {
	//TODO Implemement
	return Merchant{}, errors.New("not implemented")
}

func (dr *DisburserRepo) GetMerchantByReferenceID(merchantReferenceID string) (Merchant, error) {
	m := Merchant{}
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

func (dr *DisburserRepo) InsertOrder(o Order) error {
	_, err := dr.insertOrder.Exec(o.ID, o.MerchantReference, o.Amount, o.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (dr *DisburserRepo) InsertDisbursement(d Disbursement) (lastInsertID int64, err error) {
	res, err := dr.insertDisbursement.Exec(d.ID, d.DisbursementGroupID, d.MerchReference, d.OrderID, d.OrderFee, d.RunningTotal, d.PayoutDate, d.IsPaidOut)
	if err != nil {
		return 0, err
	}

	lID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lID, nil
}

func (dr *DisburserRepo) InsertMerchant(m Merchant) error {
	_, err := dr.insertMerchant.Exec(m.ID, m.Reference, m.Email, m.LiveOn, m.DisbursementFrequency, m.MinMonthlyFee)
	if err != nil {
		return err
	}
	return nil
}

func (dr *DisburserRepo) CreateTables() error {
	_, err := dr.createDisbursementsTable.Exec()
	if err != nil {
		return err
	}

	_, err = dr.createMerchantsTable.Exec()
	if err != nil {
		return err
	}

	_, err = dr.createOrdersTable.Exec()
	if err != nil {
		return err
	}

	return nil
}
