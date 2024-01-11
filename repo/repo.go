package repo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/levtk/sequra/reports"
	"github.com/levtk/sequra/types"
	_ "github.com/mattn/go-sqlite3"
	"log/slog"
	"time"
)

const (
	insertMerchant = `INSERT INTO MERCHANTS (id, reference, email, live_on, disbursement_frequency, minimum_monthly_fee) VALUES (
                    ?,?,?,?,?,?);`

	getOrdersByMerchantUUID = `SELECT * FROM ORDERS WHERE id=:merchantUUID;`

	getOrdersByMerchantReferenceID = `SELECT * FROM ORDERS WHERE merchant_reference=:merchRef`

	getMerchantByReferenceID = `SELECT * FROM MERCHANTS WHERE reference=:referenceID`

	insertOrder = `INSERT INTO ORDERS(id, merchant_reference, amount, created_at) VALUES(?,?,?,?);`

	insertDisbursement = `INSERT INTO DISBURSEMENT(id, disbursement_group_id, merchReference, order_id, order_fee, running_total, payout_date, is_paid_out)
	VALUES (?,?,?,?,?,?,?,?);`

	getDisbursementGroupID = `SELECT (disbursement_group_id) FROM DISBURSEMENT WHERE payout_date=:today AND merchReference=:merchRef`
)

type DisburserRepoRepository interface {
	GetOrdersByMerchantUUID(merchantUUID uuid.UUID) ([]types.Order, error)
	GetOrdersByMerchantReferenceID(ctx context.Context, merchRef string) ([]types.Order, error)
	GetMerchantDisbursementsByRange(logger slog.Logger, merchantUUID uuid.UUID, start time.Time, end time.Time) (reports.Report, error)
	GetMerchant(merchantUUID uuid.UUID) (types.Merchant, error)
	GetMerchantByReferenceID(merchantReferenceID string) (types.Merchant, error)
	GetDisbursementGroupID(ctx context.Context, today string, merchRef string) (string, error)
	InsertOrder(order types.Order) error
	InsertDisbursement(disbursement types.Disbursement) (lastInsertID int64, err error)
	InsertMerchant(m types.Merchant) error
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
	}, nil
}

func (dr *DisburserRepo) GetOrdersByMerchantUUID(merchantUUID uuid.UUID) ([]types.Order, error) {
	//TODO Implement
	return []types.Order{}, errors.New("not implemented")
}

func (dr *DisburserRepo) GetOrdersByMerchantReferenceID(ctx context.Context, merchRef string) ([]types.Order, error) {
	var orders []types.Order
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
func (dr *DisburserRepo) GetMerchantDisbursementsByRange(logger slog.Logger, merchantUUID uuid.UUID, start time.Time, end time.Time) (reports.Report, error) {
	//TODO Implement
	return reports.Report{}, errors.New("not implemented")
}

func (dr *DisburserRepo) GetMerchant(merchantUUID uuid.UUID) (types.Merchant, error) {
	//TODO Implemement
	return types.Merchant{}, errors.New("not implemented")
}

func (dr *DisburserRepo) GetMerchantByReferenceID(merchantReferenceID string) (types.Merchant, error) {
	var liveOn string
	m := &types.Merchant{}

	err := dr.getMerchantByRefID.QueryRow(merchantReferenceID).Scan(&m.ID, &m.Reference, &m.Email, &liveOn, &m.DisbursementFrequency, &m.MinMonthlyFee)
	if err != nil {
		return *m, err
	}
	m.LiveOn, err = time.Parse("2006-01-02 15:04:05+00:00", liveOn)
	if err != nil {
		return *m, err
	}
	return *m, nil
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
	if errors.Is(err, sql.ErrNoRows) {
		dgID := uuid.NewString()
		return dgID, nil
	}
	return refId, nil
}

func (dr *DisburserRepo) InsertOrder(o types.Order) error {
	_, err := dr.insertOrder.Exec(o.ID, o.MerchantReference, o.Amount, o.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (dr *DisburserRepo) InsertDisbursement(d types.Disbursement) (lastInsertID int64, err error) {
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

func (dr *DisburserRepo) InsertMerchant(m types.Merchant) error {
	_, err := dr.insertMerchant.Exec(m.ID, m.Reference, m.Email, m.LiveOn, m.DisbursementFrequency, m.MinMonthlyFee)
	if err != nil {
		return err
	}
	return nil
}
