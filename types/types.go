package types

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type Order struct {
	ID                string    `json:"id,omitempty" DB:"id"`
	MerchantReference string    `json:"merchant_reference,omitempty" DB:"merchant_reference"`
	MerchantID        uuid.UUID `json:"merchant_id,omitempty" DB:"merchant_id""`
	Amount            int64     `json:"amount,omitempty" DB:"amount"`
	CreatedAt         time.Time `json:"created_at,omitempty" DB:"created_at"`
}

type Merchant struct {
	ID                    uuid.UUID `json:"id,omitempty" DB:"id"`
	Reference             string    `json:"reference,omitempty" DB:"reference"`
	Email                 string    `json:"email,omitempty" DB:"email"`
	LiveOn                time.Time `json:"live_on,omitempty" DB:"live_on"`
	DisbursementFrequency string    `json:"disbursement_frequency,omitempty" DB:"disbursement_frequency"`
	MinMonthlyFee         string    `json:"minimum_monthly_fee,omitempty" DB:"minimum_monthly_fee"`
}

type Disbursement struct {
	RecordUUID           uuid.UUID `json:"RecordUUID" DB:"record_uuid"`
	DisbursementGroupID  uuid.UUID `json:"DisbursementGroupID" DB:"disbursement_group_id"`
	TransactionID        string    `json:"TransactionID" DB:"transaction_id"`
	MerchReference       string    `json:"MerchReference" DB:"merchReference"`
	OrderID              string    `json:"OrderID" DB:"order_id"`
	OrderFee             int64     `json:"OrderFee" DB:"order_fee"`
	OrderFeeRunningTotal int64     `json:"OrderFeeRunningTotal" DB:"order_fee_running_total"`
	PayoutDate           time.Time `json:"PayoutDate" DB:"payout_date"`
	PayoutRunningTotal   int64     `json:"PayoutRunningTotal" DB:"payout_running_total"`
	PayoutTotal          int64     `json:"PayoutTotal" DB:"payout_total"`
	IsPaidOut            bool      `json:"IsPaidOut" DB:"is_paid_out"`
}

type DisbursementReport struct {
	Year                          time.Time     `json:"year,omitempty"`
	NumberOfDisbursements         sql.NullInt64 `json:"number_of_disbursements,omitempty" DB:"number_of_disbursements"`
	AmountDisbursedToMerchants    sql.NullInt64 `json:"amount_disbursed_to_merchants,omitempty" DB:"amt_disbursed_to_merchants"`
	AmountOfOrderFees             sql.NullInt64 `json:"amount_of_order_fees,omitempty" DB:"amount_of_order_fees"`
	NumberOfMinMonthlyFeesCharged sql.NullInt64 `json:"number_of_min_monthly_fees_charged,omitempty"`
	AmountOfMonthlyFeeCharged     sql.NullInt64 `json:"amount_of_monthly_fee_charged,omitempty"`
}

type Monthly struct {
	ID                uuid.UUID `json:"id,omitempty" DB:"id"`
	MerchantReference string    `json:"merchant_reference,omitempty" DB:"merchant_reference"`
	MerchantID        uuid.UUID `json:"merchant_id,omitempty" DB:"merchant_id"`
	MonthlyFeeDate    time.Time `json:"monthly_fee_date" DB:"monthly_fee_date"`
	DidPayFee         int       `json:"did_pay_fee,omitempty" DB:"did_pay_fee"`
	MonthlyFee        int64     `json:"monthly_fee,omitempty" DB:"monthly_fee"`
	TotalOrderAmt     int64     `json:"total_order_amt,omitempty" DB:"total_order_amt"`
	OrderFeeTotal     int64     `json:"order_fee_total" DB:"order_fee_total"`
	CreatedAt         time.Time `json:"created_at" DB:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" DB:"updated_at"`
}
