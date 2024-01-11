package types

import (
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
	ID                  string `json:"ID" DB:"ID"`
	DisbursementGroupID string `json:"DisbursementGroupID" DB:"disbursement_group_id"`
	TransactionID       string `json:"TransactionID" DB:"transaction_id"`
	MerchReference      string `json:"MerchReference" DB:"merchReference"`
	OrderID             string `json:"OrderID" DB:"order_id"`
	OrderFee            int64  `json:"OrderFee" DB:"order_fee"`
	RunningTotal        int64  `json:"RunningTotal" DB:"running_total"`
	PayoutDate          string `json:"PayoutDate" DB:"payout_date"`
	PayoutTotal         int64  `json:"PayoutTotal" DB:"payout_total"`
	IsPaidOut           bool   `json:"IsPaidOut" DB:"is_paid_out"`
}
