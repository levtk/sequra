package disburse

import (
	"context"
	"github.com/google/uuid"
	"github.com/levtk/sequra/repo"
	"log/slog"
	"time"
)

func NewOrderProcessor(l *slog.Logger, ctx *context.Context, repo *repo.DisburserRepo) *OProcessor {
	op := &OProcessor{
		logger: l,
		ctx:    ctx,
		repo:   repo,
		Order:  nil,
	}
	return op
}

func (op *OProcessor) ProcessOrder(logger *slog.Logger, ctx *context.Context, repo *repo.DisburserRepo, o *Order) error {
	op.Order = o
	of, err := op.Order.CalculateOrderFee()
	if err != nil {
		return err
	}

	ok, err := op.Order.IsBeforeTimeCutOff()
	if ok && err == nil {
		//TODO create save to disbursement table with appropriate payout frequency date tagged
		merch, err := repo.GetMerchantByReferenceID(o.MerchantReference)
		if err != nil {
			logger.Error("failed to get merchant by reference id", err.Error())
			return err
		}

		if merch.DisbursementFrequency == DAILY {
			disbursementID := uuid.NewString()
			today := time.Now().UTC().Format(time.DateOnly)

			disbusrement := Disbursement{
				ID:                  disbursementID,
				DisbursementGroupID: "", //TODO write func to find existing group ID for merchant else create new one
				MerchReference:      o.MerchantReference,
				OrderID:             o.ID,
				OrderFee:            of,
				RunningTotal:        0,
				PayoutDate:          today,
				IsPaidOut:           false,
			}
		}

		if merch.DisbursementFrequency == WEEKLY {
			pd, err := merch.GetNextPayoutDate()
			if err != nil {
				logger.Error("could not get next payout date", err.Error())
				return err
			}
		}

	}

	if !ok && err == nil {
		//TODO create GetMerchPayoutFrequency and if daily add to tomorrows payout.
		//If weekly and liveOn day of week is not today, add to next payout date for merchant
	}

	if err != nil {
		op.logger.Error("failed to process order", "order-id", op.Order.ID, "merchant-reference", op.Order.MerchantReference)
		return err
	}

	return nil
}
