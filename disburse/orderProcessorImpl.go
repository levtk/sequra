package disburse

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/levtk/sequra/repo"
	"log/slog"
	"time"
)

func NewOrderProcessor(l *slog.Logger, ctx *context.Context, disburserRepo *repo.DisburserRepo) *OProcessor {
	op := &OProcessor{
		logger:                  l,
		ctx:                     ctx,
		disburserRepoRepository: disburserRepo,
		Order:                   nil,
	}
	return op
}

// ProcessOrder processes an order by performing calculations on fees, order cutoff time, and disbursement frequencies. It then
// // inserts the resulting disbursement object into the disbursement table. This does not include disbursing payments which is another process.
func (op *OProcessor) ProcessOrder(logger *slog.Logger, ctx context.Context, disburserRepo *repo.DisburserRepo, o *Order) error {
	op.Order = o
	of, err := op.Order.CalculateOrderFee()
	if err != nil {
		return err
	}

	ok, err := op.Order.IsBeforeTimeCutOff()
	if ok && err == nil {
		merch, err := disburserRepo.GetMerchantByReferenceID(o.MerchantReference)
		if err != nil {
			logger.Error("failed to get merchant by reference id", err.Error())
			return err
		}

		disbursement, err := buildDisbursement(logger, ctx, disburserRepo, o, merch, of)
		if err != nil {
			logger.Error("could not build disbursement", err.Error())
			return err
		}

		_, err = disburserRepo.InsertDisbursement(disbursement)
		if err != nil {
			logger.Error("failed to insert disbursement", err.Error())
			return err
		}
	}
	return nil
}

// buildDisbursement contains the logic to determine if the order is before the cutoff time and whether the merchant is disbursed daily or weekly. It then
// builds the repo.Disbursement struct filling the required fields.
func buildDisbursement(logger *slog.Logger, ctx context.Context, disburserRepo *repo.DisburserRepo, o *Order, merch Merchant, orderFee int64) (repo.Disbursement, error) {
	disbursementID := uuid.NewString()
	var pd time.Time
	var payoutDate string
	disbursementFreq := merch.DisbursementFrequency
	switch disbursementFreq {
	case DAILY:
		ok, err := o.IsBeforeTimeCutOff()
		if ok && err == nil {
			pd = time.Now().UTC()
		}
		if !ok && err == nil {
			pd = pd.AddDate(0, 0, 1)
		}
		payoutDate = pd.Format(time.DateOnly)

	case WEEKLY:
		pd, err := merch.GetNextPayoutDate()
		if err != nil {
			return repo.Disbursement{}, err
		}
		payoutDate = pd.Format(time.DateOnly)

	default:
		return repo.Disbursement{}, errors.New("merchants disbursement frequency is not supported")
	}

	var disbursementGroupID string
	disbGrpID, err := disburserRepo.GetDisbursementGroupID(ctx, payoutDate, merch.Reference)
	if err == nil {
		disbursementGroupID = disbGrpID
	} else if errors.Is(sql.ErrNoRows, err) {
		disbursementGroupID = uuid.NewString()
	} else {
		logger.Error("could not get disbursement group id from disburserRepoRepository", err.Error())
		return repo.Disbursement{}, err
	}

	return repo.Disbursement{
		ID:                  disbursementID,
		DisbursementGroupID: disbursementGroupID,
		MerchReference:      merch.Reference,
		OrderID:             o.ID,
		OrderFee:            orderFee,
		RunningTotal:        0,
		PayoutDate:          payoutDate,
		IsPaidOut:           false,
	}, err
}
