package disburse

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/levtk/sequra/repo"
	"log/slog"
	"time"
)

// DisbursementsByYear meets the requirements outlined in the system requirement for calculating the total number of disbursements,
// amount disbursed to merchants, amount of order fees, number of minimum monthly fees charged, and total amount in monthly fees charged.
func (r *Report) DisbursementsByYear(logger *slog.Logger, ctx context.Context, repo repo.DisburserRepoRepository) (Report, error) {
	//TODO Implement
	return Report{}, errors.New("not implemented")
}

func (r *Report) DisbursementsByRange(logger *slog.Logger, ctx context.Context, repo repo.DisburserRepoRepository, start time.Time, end time.Time) (Report, error) {
	//TODO Implement
	return Report{}, errors.New("not implemented")
}

func (r *Report) MerchantDisbursements(logger *slog.Logger, ctx context.Context, repo repo.DisburserRepoRepository, merchantUUID uuid.UUID, start time.Time, end time.Time) (Report, error) {
	//TODO Implement
	return Report{}, errors.New("not implemented")
}
