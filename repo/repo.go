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

	insertDisbursement = `INSERT INTO DISBURSEMENT(record_uuid, disbursement_group_id, merchReference, order_id, order_fee, order_fee_running_total, payout_date, payout_running_total, payout_total, is_paid_out)
	VALUES (?,?,?,?,?,?,?,?,?,?);`

	getDisbursementGroupID = `SELECT (disbursement_group_id) FROM DISBURSEMENT WHERE payout_date=:today AND merchReference=:merchRef`

	getNumberOfDisbursementsByYear = `SELECT COUNT() FROM DISBURSEMENT WHERE is_paid_out=1 AND payout_date LIKE '%' + :year + '%'`

	getTotalCommissionAndTotalPayoutByYear = `SELECT  count(), SUM(DISBURSEMENT.payout_total), sum(DISBURSEMENT.order_fee_running_total) from DISBURSEMENT WHERE is_paid_out = 1 AND payout_date LIKE '%' + :year + '%'`
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
	GetNumberOfDisbursementsByYear(yyyy string) (int64, error)
	GetTotalCommissionsAndPayoutByYear(yyyy string) (types.DisbursementReport, error)
}

type DisburserRepo struct {
	db                                     *sql.DB
	ctx                                    context.Context
	logger                                 *slog.Logger
	insertOrder                            *sql.Stmt
	insertDisbursement                     *sql.Stmt
	insertMerchant                         *sql.Stmt
	getOrdersByMerchantReferenceID         *sql.Stmt
	getMerchantByRefID                     *sql.Stmt
	getDisbursementGroupID                 *sql.Stmt
	getNumberOfDisbursementsByYear         *sql.Stmt
	getTotalCommissionAndTotalPayoutByYear *sql.Stmt
	createDisbursementsTable               *sql.Stmt
	createMerchantsTable                   *sql.Stmt
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

	getNumDisbursementsByYear, err := db.Prepare(getNumberOfDisbursementsByYear)
	if err != nil {
		return &DisburserRepo{}, err
	}

	getTotalCommAndPayoutByYear, err := db.Prepare(getTotalCommissionAndTotalPayoutByYear)
	if err != nil {
		return &DisburserRepo{}, err
	}

	return &DisburserRepo{
		db:                                     db,
		ctx:                                    ctx,
		logger:                                 l,
		insertOrder:                            insOrderStmt,
		insertDisbursement:                     insDisbursementStmt,
		insertMerchant:                         insertMerchantStmt,
		getOrdersByMerchantReferenceID:         getOrdersByMerchRefID,
		getMerchantByRefID:                     getMerchantByRefID,
		getDisbursementGroupID:                 getDisburseGroupID,
		getNumberOfDisbursementsByYear:         getNumDisbursementsByYear,
		getTotalCommissionAndTotalPayoutByYear: getTotalCommAndPayoutByYear,
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

// GetNumberOfDisbursementsByYear takes the year format of YYYY as a string and returns the number of disbursements for that year or an error.
func (dr *DisburserRepo) GetNumberOfDisbursementsByYear(yyyy string) (int64, error) {
	var n int64
	row := dr.getNumberOfDisbursementsByYear.QueryRow(yyyy)
	err := row.Scan(n)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (dr *DisburserRepo) GetTotalCommissionsAndPayoutByYear(yyyy string) (types.DisbursementReport, error) {
	var dispReport types.DisbursementReport
	row := dr.getTotalCommissionAndTotalPayoutByYear.QueryRow(yyyy)
	err := row.Scan(dispReport)
	if err != nil {
		return dispReport, err
	}

	return dispReport, nil
}

func (dr *DisburserRepo) InsertOrder(o types.Order) error {
	_, err := dr.insertOrder.Exec(o.ID, o.MerchantReference, o.Amount, o.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (dr *DisburserRepo) InsertDisbursement(d types.Disbursement) (lastInsertID int64, err error) {
	res, err := dr.insertDisbursement.Exec(d.RecordUUID, d.DisbursementGroupID, d.MerchReference, d.OrderID, d.OrderFee, d.OrderFeeRunningTotal, d.PayoutDate, d.PayoutRunningTotal, d.PayoutTotal, d.IsPaidOut)
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
