package disburse

import (
	"context"
	"github.com/levtk/sequra/repo"
	"log/slog"
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
