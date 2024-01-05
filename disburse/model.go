package disburse

import (
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

const (
	RATE_LESS_THAN_50       int64 = 10
	RATE_BETWEEN_50_AND_300 int64 = 5
	RATE_ABOVE_300          int64 = 25
	MAX_ORDER               int64 = 1000000
)

type Disburser interface {
	ProcessOrder(order Order) error
}

type Order struct {
	ID                string    `json:"id,omitempty"`
	MerchantReference string    `json:"merchant_reference,omitempty"`
	Amount            string    `json:"amount,omitempty"`
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
		slog.Error("order submitted above max order value permitted")
		return -1, errors.New("order submitted above max order value permitted")
	}

	return orderFee, nil
}
