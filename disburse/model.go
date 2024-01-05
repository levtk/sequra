package disburse

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	RATE_LESS_THAN_50       int64  = 10
	RATE_BETWEEN_50_AND_300 int64  = 5
	RATE_ABOVE_300          int64  = 25
	MAX_ORDER               int64  = 1000000
	TIME_CUT_OFF            string = "08:00"
)

type Disburser interface {
	ProcessOrder() error
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
