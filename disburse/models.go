package disburse

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/levtk/sequra/repo"
	"log/slog"
	"strconv"
	"time"

	"github.com/google/uuid"
)

const (
	RATE_LESS_THAN_50       int64  = 10
	RATE_BETWEEN_50_AND_300 int64  = 5
	RATE_ABOVE_300          int64  = 25
	MAX_ORDER               int64  = 1000000 //Should be configured per Merchant during onboarding
	TIME_CUT_OFF            string = "08:00"
	OREDERS_FILENAME               = "orders.csv"
	MERCHANTS_FILENAME             = "merchants.csv"
)

type DisburserService struct {
	logger       *slog.Logger
	ctx          *context.Context
	ProcessOrder OrderProcessor
	Importer     Importer
	Reporter     Reporter
	Repo         repo.DisburserRepoRepository
}

func NewDisburserService(logger *slog.Logger, ctx *context.Context, db *sql.DB) (*DisburserService, error) {
	repo, err := repo.NewDisburserRepo(logger, ctx, db)
	if err != nil {
		return &DisburserService{}, err
	}

	importer := NewImport(logger, ctx, repo)
	orderProcessor := NewOrderProcessor(logger, ctx, repo)
	reporter := NewReporter(logger, ctx, repo)
	return &DisburserService{
		logger:       logger,
		ctx:          ctx,
		ProcessOrder: orderProcessor,
		Importer:     importer,
		Reporter:     reporter,
		Repo:         repo,
	}, nil

}

func NewReporter(logger *slog.Logger, ctx *context.Context, repo *repo.DisburserRepo) *Report {
	return &Report{
		logger:   logger,
		ctx:      ctx,
		repo:     repo,
		Name:     "",
		Merchant: Merchant{},
		Start:    time.Time{},
		End:      time.Time{},
		data:     nil,
	}
}

// Report implements the Reporter interface
type Report struct {
	logger   *slog.Logger
	ctx      *context.Context
	Name     string
	Merchant Merchant
	repo     repo.DisburserRepoRepository
	Start    time.Time
	End      time.Time
	data     []byte
}

type Import struct {
	logger            *slog.Logger
	ctx               *context.Context
	ordersFileName    string
	merchantsFileName string
	repo              repo.DisburserRepoRepository
}

type OProcessor struct {
	Order  *Order
	repo   repo.DisburserRepoRepository
	logger *slog.Logger
	ctx    *context.Context
}

type Order struct {
	ID                string    `json:"id,omitempty"`
	MerchantReference string    `json:"merchant_reference,omitempty"`
	MerchantID        uuid.UUID `json:"merchant_id,omitempty"`
	Amount            int64     `json:"amount,omitempty"`
	CreatedAt         time.Time `json:"created_at,omitempty"`
}

func NewOrder(id string, merchantReference string, amount int64) *Order {
	t := time.Now().UTC()
	return &Order{
		ID:                id,
		MerchantReference: merchantReference,
		Amount:            amount,
		CreatedAt:         t,
	}
}

func (o *Order) IsBeforeTimeCutOff() (bool, error) {
	cutoff, err := time.Parse(time.TimeOnly, TIME_CUT_OFF)
	if err != nil {
		return false, err
	}

	if time.Now().UTC().Before(cutoff) {
		return true, nil
	}

	return false, nil
}

func (o *Order) CalculateOrderFee() (int64, error) {
	fee, err := calculateOrderFee(o.Amount)
	if err != nil {
		return 0, err
	}
	return fee, nil
}

type Merchant struct {
	ID                    uuid.UUID `json:"id,omitempty"`
	Reference             string    `json:"reference,omitempty"`
	Email                 string    `json:"email,omitempty"`
	LiveOn                time.Time `json:"live_on,omitempty"`
	DisbursementFrequency string    `json:"disbursement_frequency,omitempty"`
	MinMonthlyFee         string    `json:"minimum_monthly_fee,omitempty"`
}

func (m *Merchant) GetMinMonthlyFee() (int64, error) {
	mmf, err := strconv.ParseFloat(m.MinMonthlyFee, 64)
	if err != nil {
		return 0, err
	}
	return int64(mmf * 100), nil
}

func (m *Merchant) GetMinMonthlyFeeRemaining() (int64, error) {
	//TODO Implement
	return 0, errors.New("not implemented.")
}

type DisbursementReport struct {
	Year                          time.Time `json:"year,omitempty"`
	NumberOfDisbursements         int64     `json:"number_of_disbursements,omitempty"`
	AmountDisbursedToMerchants    int64     `json:"amount_disbursed_to_merchants,omitempty"`
	AmountOfOrderFees             int64     `json:"amount_of_order_fees,omitempty"`
	NumberOfMinMonthlyFeesCharged int32     `json:"number_of_min_monthly_fees_charged,omitempty"`
	AmountOfMonthlyFeeCharged     int64     `json:"amount_of_monthly_fee_charged,omitempty"`
}

func (o *Order) ProcessOrder() error {
	_, err := calculateOrderFee(o.Amount)
	//TODO add func to determine total orders per interval freq and account for min monthly fee
	return err
}

type Disbursement struct {
	ID                  string `json:"ID" DB:"ID"`
	DisbursementGroupID string `json:"DisbursementGroupID" DB:"disbursement_group_id"`
	MerchReference      string `json:"MerchReference" DB:"merchReference"`
	OrderID             string `json:"OrderID" DB:"order_id"`
	OrderFee            int64  `json:"OrderFee" DB:"order_fee"`
	RunningTotal        int64  `json:"RunningTotal" DB:"running_total"`
	PayoutDate          string `json:"PayoutDate" DB:"payout_date"`
	IsPaidOut           bool   `json:"IsPaidOut" DB:"is_paid_out"`
}

func calculateOrderFee(orderAmt int64) (orderFee int64, err error) {
	if orderAmt > 0 && orderAmt < 5000 {
		orderFee = RATE_LESS_THAN_50 * orderAmt / 100
		return orderFee, nil
	}

	if orderAmt > 5000 && orderAmt < 30000 {
		orderFee = RATE_BETWEEN_50_AND_300 * orderAmt / 100
		return orderFee, nil
	}

	if orderAmt > 30000 {
		orderFee = RATE_ABOVE_300 * orderAmt / 1000
		return orderFee, nil
	}

	if orderAmt > MAX_ORDER {
		return -1, errors.New("orderamt submitted above max orderamt value permitted")
	}

	return orderFee, nil
}

func getMerchantReferenceFromOrder(o Order) (string, error) {
	merchRef := o.MerchantReference

	if merchRef == "" {
		return "", fmt.Errorf("merchant reference not found for order %s", o.ID)
	}

	return merchRef, nil
}

func getMerchant(merchRef string) (Merchant, error) {
	// TODO implement
	return Merchant{}, nil
}

func (m Merchant) calculateDailyTotalOrders() (int64, error) {
	//TODO implement
	return -1, nil
}

func (m Merchant) calculateWeeklyTotalOrders() (int64, error) {
	//TODO implement
	return -1, nil
}

func getOrdersByMerchRef(merchRef string) ([]Order, error) {
	//TODO implement
	return []Order{}, nil
}

func isBeforeCutOffTime() (bool, error) {
	now := time.Now().UTC()
	cutOff, err := time.Parse("15:04", TIME_CUT_OFF)
	if err != nil {
		return false, err
	}
	if now.Sub(cutOff) < 0 {
		return true, nil
	} else {
		return false, nil
	}
}