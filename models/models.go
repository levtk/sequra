package models

import (
	"context"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

type Import struct {
	logger            *slog.Logger
	ctx               context.Context
	ordersFileName    string
	merchantsFileName string
}

type Order struct {
	ID                string    `json:"id,omitempty"`
	MerchantReference string    `json:"merchant_reference,omitempty"`
	MerchantID        uuid.UUID `json:"merchant_id,omitempty"`
	Amount            int64     `json:"amount,omitempty"`
	CreatedAt         time.Time `json:"created_at,omitempty"`
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
	Start    time.Time
	End      time.Time
	data     []byte
}
