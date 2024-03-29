package disburse

import (
	"context"
	"github.com/google/uuid"
	"github.com/levtk/sequra/repo"
	"github.com/levtk/sequra/types"
	"log/slog"
	"net/http"
	"time"
)

type Disburser interface {
	ProcessOrder(op OrderProcessor) error
	ImportOrders(o Importer) error
	GenerateReports(r Reporter) ([]Report, error)
}

type Importer interface {
	ImportOrders() ([]types.Disbursement, map[string]types.Merchant, []types.Monthly, error)
	Import(w http.ResponseWriter, r *http.Request)
}
type OrderProcessor interface {
	ProcessOrder(logger *slog.Logger, ctx context.Context, repo repo.DisburserRepoRepository, o *Order) error
	ProcessBatchDistributions([]types.Disbursement) error
	ProcessBatchMonthly([]types.Monthly) error
}

type Seller interface {
	GetMinMonthlyFee() (int64, error)
	GetMinMonthlyFeeRemaining() (int64, error)
	GetNextPayoutDate() (time.Time, error)
	CalculateDailyTotalOrders() (int64, error)
	CalculateWeeklyTotalOrders() (int64, error)
}
type Reporter interface {
	DisbursementsByYear(logger *slog.Logger, repo repo.DisburserRepoRepository, YYYY string) (types.DisbursementReport, error)
	DisbursementsByRange(logger *slog.Logger, ctx context.Context, repo repo.DisburserRepoRepository, start time.Time, end time.Time) (Report, error)
	MerchantDisbursements(logger *slog.Logger, ctx context.Context, repo repo.DisburserRepoRepository, merchantUUID uuid.UUID, start time.Time, end time.Time) (Report, error)
	NumberMonthlyPaymentsByYear(logger *slog.Logger, YYYY string, disbursements []types.Disbursement) (Report, error)
	DisbursementReport(logger *slog.Logger, repo repo.DisburserRepoRepository, YYYY string) (types.DisbursementReport, error)
	GetDisbursementReport(w http.ResponseWriter, r *http.Request)
}
