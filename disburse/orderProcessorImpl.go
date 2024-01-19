package disburse

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/levtk/sequra/repo"
	"github.com/levtk/sequra/types"
	"log/slog"
	"time"
)

func NewOrderProcessor(l *slog.Logger, ctx context.Context, disburserRepo *repo.DisburserRepo) *OProcessor {
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
func (op *OProcessor) ProcessOrder(logger *slog.Logger, ctx context.Context, disburserRepo repo.DisburserRepoRepository, o *Order) error {
	op.Order = o
	of, err := op.Order.CalculateOrderFee()
	if err != nil {
		return err
	}

	ok, err := op.Order.IsBeforeTimeCutOff()
	if ok && err == nil {
		o.Lock()
		merch, err := disburserRepo.GetMerchantByReferenceID(o.MerchantReference)
		o.Unlock()
		if err != nil {
			logger.Error("failed to get merchant by reference id", "error", err.Error())
			return err
		}
		o.Lock()
		disbursement, err := buildDisbursement(logger, ctx, disburserRepo, o, merch, of)
		o.Unlock()
		if err != nil {
			logger.Error("could not build disbursement", "error", err.Error())
			return err
		}

		_, err = disburserRepo.InsertDisbursement(disbursement)
		if err != nil {
			logger.Error("failed to insert disbursement", "error", err.Error())
			return err
		}
	}
	return nil
}

// buildDisbursement contains the logic to determine if the order is before the cutoff time and whether the merchant is disbursed daily or weekly. It then
// builds the Disbursement struct filling the required fields.
func buildDisbursement(logger *slog.Logger, ctx context.Context, disburserRepo repo.DisburserRepoRepository, o *Order, merch types.Merchant, orderFee int64) (types.Disbursement, error) {
	disbursementID := uuid.NewString()
	var pd time.Time
	var payoutDate string
	disbursementFreq := merch.DisbursementFrequency
	switch disbursementFreq {
	case types.DAILY:
		ok, err := o.IsBeforeTimeCutOff()
		if ok && err == nil {
			pd = time.Now().UTC()
		}
		if !ok && err == nil {
			pd = pd.AddDate(0, 0, 1)
		}
		payoutDate = pd.Format(time.DateOnly)

	case types.WEEKLY:
		pd, err := merch.GetNextPayoutDate()
		if err != nil {
			return types.Disbursement{}, err
		}
		payoutDate = pd.Format(time.DateOnly)

	default:
		return types.Disbursement{}, errors.New("merchants disbursement frequency is not supported")
	}

	var disbursementGroupID string
	disbGrpID, err := disburserRepo.GetDisbursementGroupID(ctx, payoutDate, merch.Reference)
	if err == nil {
		disbursementGroupID = disbGrpID
	} else {
		logger.Error("could not get disbursement group id or create it from disburserRepoRepository", "error", err.Error())
		return types.Disbursement{}, err
	}

	return types.Disbursement{
		RecordUUID:           disbursementID,
		DisbursementGroupID:  disbursementGroupID,
		MerchReference:       merch.Reference,
		OrderID:              o.ID,
		OrderFee:             orderFee,
		OrderFeeRunningTotal: 0,
		PayoutDate:           payoutDate,
		IsPaidOut:            false,
	}, err
}

func (op *OProcessor) ProcessBatchDistributions(disbursements []types.Disbursement) error {
	count := 0
	for i := 0; i < len(disbursements); i++ {
		//if disbursements[i].DisbursementGroupID == "" {
		//	break
		//}
		if disbursements[i].RecordUUID == "" {
			continue
		}
		_, err := op.disburserRepoRepository.InsertDisbursement(disbursements[i])
		if err != nil {
			op.logger.Error("error inserting disbursement record", err.Error())
			return err
		}
		count++

	}
	op.logger.Info("number of disbursement records inserted: ", "count", count)
	return nil
}
