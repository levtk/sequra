package repo

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/levtk/sequra/reports"
	"github.com/levtk/sequra/types"
	"log/slog"
	"time"
)

const (
	insertMerchant = `INSERT INTO MERCHANTS (id, reference, email, live_on, disbursement_frequency, minimum_monthly_fee) VALUES (
                    ?,?,?,?,?,?);`

	getOrdersByMerchantUUID = `SELECT * FROM ORDERS WHERE id=:merchantUUID;`

	getOrdersByMerchantReferenceID = `SELECT * FROM ORDERS WHERE merchant_reference=:merchRef;`

	getMerchantByReferenceID = `SELECT * FROM MERCHANTS WHERE reference=:referenceID;`

	insertOrder = `INSERT INTO ORDERS(id, merchant_reference, amount, created_at) VALUES(?,?,?,?);`

	insertDisbursement = `INSERT INTO DISBURSEMENT(record_uuid, disbursement_group_id, merchReference, order_id, order_fee, order_fee_running_total, payout_date, payout_running_total, payout_total, is_paid_out)
	VALUES (?,?,?,?,?,?,?,?,?,?);`

	getDisbursementGroupID = `SELECT (disbursement_group_id) FROM DISBURSEMENT WHERE payout_date=:today AND merchReference=:merchRef;`

	getNumberOfDisbursementsByYear = `SELECT COUNT() FROM DISBURSEMENT WHERE is_paid_out=1 AND payout_date LIKE :YYYY;`

	getTotalCommissionAndTotalPayoutByYear = `SELECT  COUNT() AS number_of_disbursements, SUM(DISBURSEMENT.payout_total) AS amt_disbursed_to_merchants, SUM(DISBURSEMENT.order_fee_running_total) AS amount_of_order_fees FROM DISBURSEMENT WHERE is_paid_out = 1 AND payout_date LIKE @YYYY || '%';`

	insertMonthly = `INSERT INTO MONTHLY(id, merchant_id, merchant_reference, monthly_fee_date, did_pay_fee, 
                    monthly_fee, total_order_amt, order_fee_total, createdAt, updatedAt) VALUES (?,?,?,?,?,?,?,?,?,?);`

	getMonthlyFeeTotalsByYear = `SELECT COUNT(*) as count, SUM(monthly_fee) AS total_monthly_fees, SUM(order_fee_total) AS total_order_fees, SUM(amt_monthly_fee_paid) AS total_monthly_fees_paid FROM MONTHLY
										WHERE createdAt LIKE @YYYY || '%'  AND did_pay_fee = 1;`
)

type DisburserRepoRepository interface {
	GetOrdersByMerchantUUID(merchantUUID uuid.UUID) ([]types.Order, error)
	GetOrdersByMerchantReferenceID(ctx context.Context, merchRef string) ([]types.Order, error)
	GetMerchantDisbursementsByRange(logger slog.Logger, merchantUUID uuid.UUID, start time.Time, end time.Time) (reports.Report, error)
	GetMerchant(merchantUUID uuid.UUID) (types.Merchant, error)
	GetMerchantByReferenceID(merchantReferenceID string) (types.Merchant, error)
	GetDisbursementGroupID(ctx context.Context, today time.Time, merchRef string) (string, error)
	InsertOrder(order types.Order) error
	InsertDisbursement(disbursement types.Disbursement) (lastInsertID int64, err error)
	InsertMerchant(m types.Merchant) error
	GetNumberOfDisbursementsByYear(yyyy string) (int64, error)
	GetTotalCommissionsAndPayoutByYear(yyyy string) (types.DisbursementReport, error)
	InsertMonthly(m types.Monthly) error
	GetMonthlyFeesPaidByYear(YYYY string) (count, totalMonthlyFees, totalOrderFees sql.NullInt64, err error)
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
	insMonthly                             *sql.Stmt
	getMonthlyFeesPaidByYear               *sql.Stmt
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

	insertMonthlyStmt, err := db.Prepare(insertMonthly)
	if err != nil {
		return &DisburserRepo{}, err
	}

	getMonthlyFeesPaidByYearStmt, err := db.Prepare(getMonthlyFeeTotalsByYear)
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
		insMonthly:                             insertMonthlyStmt,
		getMonthlyFeesPaidByYear:               getMonthlyFeesPaidByYearStmt,
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
func (dr *DisburserRepo) GetDisbursementGroupID(ctx context.Context, today time.Time, merchRef string) (string, error) {
	var refId string
	t := today.Format(time.DateOnly)
	row := dr.getDisbursementGroupID.QueryRowContext(ctx, t, merchRef)
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
	disrpt := types.DisbursementReport{}
	row := dr.getTotalCommissionAndTotalPayoutByYear.QueryRow(yyyy)
	err := row.Scan(&disrpt.NumberOfDisbursements, &disrpt.AmountDisbursedToMerchants, &disrpt.AmountOfOrderFees)
	if err != nil {
		return disrpt, err
	}

	return disrpt, nil
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

func (dr *DisburserRepo) InsertMonthly(m types.Monthly) error {
	id := m.ID.String()
	merchID := m.MerchantID.String()
	monDate := m.MonthlyFeeDate.String()
	createdAt := m.CreatedAt.String()
	_, err := dr.insMonthly.Exec(id, merchID, m.MerchantReference, monDate, m.DidPayFee, m.MonthlyFee, m.TotalOrderAmt, m.OrderFeeTotal, createdAt, time.Now().UTC().Format(time.DateTime))
	if err != nil {
		dr.logger.Info("failed to insert", "monthly", m)
		return err
	}
	return nil
}

func (dr *DisburserRepo) GetMonthlyFeesPaidByYear(YYYY string) (count, totalMonthlyFees, totalOrderFees sql.NullInt64, err error) {
	dest := &struct {
		count                sql.NullInt64
		totalMonthlyFees     sql.NullInt64
		totalOrderFees       sql.NullInt64
		totalMonthlyFeesPaid sql.NullInt64
	}{}
	row := dr.getMonthlyFeesPaidByYear.QueryRow(YYYY)
	err = row.Scan(&dest.count, &dest.totalMonthlyFees, &dest.totalOrderFees, &dest.totalMonthlyFeesPaid)
	if err != nil {
		dr.logger.Error("failed to get monthly fees paid by year")
		return sql.NullInt64{}, sql.NullInt64{}, sql.NullInt64{}, err
	}
	return dest.count, dest.totalMonthlyFees, dest.totalOrderFees, nil
}
