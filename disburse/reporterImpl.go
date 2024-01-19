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

// TODO refactor this Reporter interface. A report does not need a DB connection. A reporter does.
// DisbursementsByYear meets the requirements outlined in the system requirement for calculating the total number of disbursements,
// amount disbursed to merchants, amount of order fees, number of minimum monthly fees charged, and total amount in monthly fees charged.
func (r *Report) DisbursementsByYear(logger *slog.Logger, ctx context.Context, repo repo.DisburserRepoRepository) (Report, error) {
	var dispReport2022, dispReport2023 types.DisbursementReport
	dispReport2022, err := repo.GetTotalCommissionsAndPayoutByYear("2022")
	if err != nil {
		logger.Error("failed to get number of disbursements for 2022", "error", err)
		return Report{}, err
	}

	dispReport2023, err := repo.GetTotalCommissionsAndPayoutByYear("2023")
	if err != nil {
		logger.Error("failed to get number of disbursements for 2023", "error", err)
		return Report{}, err
	}

	//TODO add querry for monthly fees.
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
