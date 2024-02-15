package disburse

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/levtk/sequra/repo"
	"github.com/levtk/sequra/types"
	"log/slog"
	"net/http"
	"time"
)

// DisbursementsByYear meets the requirements outlined in the system requirement for calculating the total number of disbursements,
// amount disbursed to merchants, amount of order fees, number of minimum monthly fees charged, and total amount in monthly fees charged.
func (r *Report) DisbursementsByYear(logger *slog.Logger, repo repo.DisburserRepoRepository, YYYY string) (types.DisbursementReport, error) {
	disbursementReport := types.DisbursementReport{}
	numMonthlyFeesCharged, amtOfMonthlyFeeCharged, amtOrderFees, err := repo.GetMonthlyFeesPaidByYear(YYYY)
	if err != nil {
		logger.Error("failed to get monthly fees paid by year", "error", err)
		return disbursementReport, err
	}

	disprpt, err := repo.GetTotalCommissionsAndPayoutByYear(YYYY)
	if err != nil {
		logger.Error("failed to get total commissions and payouts by year", "error", err)
		return types.DisbursementReport{}, err
	}

	disprpt.NumberOfMinMonthlyFeesCharged = numMonthlyFeesCharged
	disprpt.AmountOfMonthlyFeeCharged = amtOfMonthlyFeeCharged
	disprpt.AmountOfOrderFees = amtOrderFees

	return disprpt, nil
}

func (r *Report) DisbursementsByRange(logger *slog.Logger, ctx context.Context, repo repo.DisburserRepoRepository, start time.Time, end time.Time) (Report, error) {
	//TODO Implement
	return Report{}, errors.New("not implemented")
}

func (r *Report) MerchantDisbursements(logger *slog.Logger, ctx context.Context, repo repo.DisburserRepoRepository, merchantUUID uuid.UUID, start time.Time, end time.Time) (Report, error) {
	//TODO Implement
	return Report{}, errors.New("not implemented")
}

func (r *Report) NumberMonthlyPaymentsByYear(logger *slog.Logger, YYYY string, disbursements []types.Disbursement) (Report, error) {
	for i := 0; i < len(disbursements); i++ {
		if i > 0 && disbursements[i-1].MerchReference == disbursements[i].MerchReference {

		}
	}
	return Report{}, errors.New("not implemented")
}

func (r *Report) DisbursementReport(logger *slog.Logger, repo repo.DisburserRepoRepository, YYYY string) (types.DisbursementReport, error) {
	disbursementReport := types.DisbursementReport{}
	numMonthlyFeesCharged, amtOfMonthlyFeeCharged, _, err := repo.GetMonthlyFeesPaidByYear(YYYY)
	if err != nil {
		logger.Error("failed to get monthly fees paid by year", "error", err)
		return disbursementReport, err
	}

	disprpt, err := repo.GetTotalCommissionsAndPayoutByYear(YYYY)
	if err != nil {
		logger.Error("failed to get total commissions and payouts by year", "error", err)
		return types.DisbursementReport{}, err
	}

	disprpt.NumberOfMinMonthlyFeesCharged = numMonthlyFeesCharged
	disprpt.AmountOfMonthlyFeeCharged = amtOfMonthlyFeeCharged

	return disprpt, nil
}

func (r *Report) GetDisbursementReport(w http.ResponseWriter, req *http.Request) {
	reportRequest := struct {
		Name string
		YYYY string
	}{}

	err := json.NewDecoder(req.Body).Decode(&reportRequest)
	if err != nil {
		r.Logger.Error("failed to decode report request from http request", "error", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	report, err := r.DisbursementReport(r.Logger, r.Repo, reportRequest.YYYY)
	if err != nil {
		r.Logger.Error("failed to get disbursement report from repo", "error", err)
	}

	rpt, err := json.Marshal(report)
	if err != nil {
		r.Logger.Error("failed to encode report", "error", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(rpt)
}
