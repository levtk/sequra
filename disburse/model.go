package disburse

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/levtk/sequra/repo"
	"log/slog"
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

type Disburser interface {
	ProcessOrder(op OrderProcessor) error
	ImportOrders(o Importer) error
	GenerateReports(r Reporter) ([]Report, error)
}

type DisburserService struct {
	ProcessOrder OrderProcessor
	Importer     Importer
	Reporter     Reporter
	logger       slog.Logger
}

func NewDisburserService(logger slog.Logger, ctx context.Context, db *sql.DB) (*DisburserService, error) {
	repo, err := repo.NewDisburserRepo(logger, ctx, db)
}

type Importer interface {
	ImportOrders() error
}
type OrderProcessor interface {
	CalculateOrderFee(o Order) (int64, error)
	GetMinMonthlyFeeRemaining(o Order, s Seller) (int64, error)
}

type Seller interface {
	GetMinMonthlyFee() (int64, error)
	GetRemainingMonthlyFee() (int64, error)
}
type Reporter interface {
	DisbursementsByYear(logger slog.Logger, ctx context.Context) (Report, error)
	DisbursementsByRange(logger slog.Logger, ctx context.Context, start time.Time, end time.Time) (Report, error)
	MerchantDisbursements(logger slog.Logger, ctx context.Context, merchantUUID uuid.UUID, start time.Time, end time.Time) (Report, error)
}

type Import struct {
	logger            slog.Logger
	ctx               context.Context
	ordersFileName    string
	merchantsFileName string
}

func NewImport(logger slog.Logger, ctx context.Context) *Import {
	return &Import{
		logger:            logger,
		ctx:               ctx,
		ordersFileName:    OREDERS_FILENAME,
		merchantsFileName: MERCHANTS_FILENAME,
	}
}

func (i *Import) ImportOrders() error {
	orders, err := parseDataFromOrders(i.ordersFileName)
	if err != nil {
		i.logger.Error("failed to parse data from orders", err.Error())
	}
}

type Report struct {
	Name     string
	Merchant Merchant
	Start    time.Time
	End      time.Time
	data     []byte
}

// DisbursementsByYear meets the requiements outlined in the system requirement for calculating the total number of disbursements,
// amount disbursed to merchants, amount of order fees, number of minimum monthly fees charged, and total amount in monthly fees charged.
func (r *Report) DisbursementsByYear(logger slog.Logger, ctx context.Context) (Report, error) {
	//TODO Implement
	return Report{}, errors.New("not implemented")
}

type Order struct {
	ID                string    `json:"id,omitempty"`
	MerchantReference string    `json:"merchant_reference,omitempty"`
	Amount            int64     `json:"amount,omitempty"`
	CreatedAt         time.Time `json:"created_at,omitempty"`
}

type Merchant struct {
	ID                    uuid.UUID `json:"id,omitempty"`
	Reference             string    `json:"reference,omitempty"`
	Email                 string    `json:"email,omitempty"`
	LiveOn                time.Time `json:"live_on,omitempty"`
	DisbursementFrequency string    `json:"disbursement_frequency,omitempty"`
	MinMonthlyFee         string    `json:"minimum_monthly_fee,omitempty"`
}

type DisbursementReport struct {
	Year                          time.Time `json:"year,omitempty"`
	NumberOfDisbursements         int64     `json:"number_of_disbursements,omitempty"`
	AmountDisbursedToMerchants    int64     `json:"amount_disbursed_to_merchants,omitempty"`
	AmountOfOrderFees             int64     `json:"amount_of_order_fees,omitempty"`
	NumberOfMinMonthlyFeesCharged int32     `json:"number_of_min_monthly_fees_charged,omitempty"`
	AmountOfMonthlyFeeCharged     int64     `json:"amount_of_monthly_fee_charged,omitempty"`
}

func (o Order) ProcessOrder() error {
	_, err := calculateOrderFee(o.Amount)
	//TODO add func to determine total orders per interval freq and account for min monthly fee
	return err
}

func calculateOrderFee(order int64) (orderFee int64, err error) {
	if order > 0 && order < 5000 {
		orderFee = RATE_LESS_THAN_50 * order / 100
		return orderFee, nil
	}

	if order > 5000 && order < 30000 {
		orderFee = RATE_BETWEEN_50_AND_300 * order / 100
		return orderFee, nil
	}

	if order > 30000 {
		orderFee = RATE_ABOVE_300 * order / 1000
		return orderFee, nil
	}

	if order > MAX_ORDER {
		return -1, errors.New("order submitted above max order value permitted")
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
