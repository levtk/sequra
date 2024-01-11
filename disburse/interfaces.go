package disburse

import (
	"context"
	"github.com/google/uuid"
	"github.com/levtk/sequra/repo"
	"github.com/levtk/sequra/types"
	"log/slog"
	"time"
)

type Disburser interface {
	ProcessOrder(op OrderProcessor) error
	ImportOrders(o Importer) error
	GenerateReports(r Reporter) ([]Report, error)
}

type Importer interface {
	ImportOrders() (Orders, map[string]types.Merchant, error)
}
type OrderProcessor interface {
	ProcessOrder(logger *slog.Logger, ctx context.Context, repo repo.DisburserRepoRepository, o *Order) error
}

type Seller interface {
	GetMinMonthlyFee() (int64, error)
	GetMinMonthlyFeeRemaining() (int64, error)
	GetNextPayoutDate() (time.Time, error)
}
type Reporter interface {
	DisbursementsByYear(logger *slog.Logger, ctx context.Context, repo repo.DisburserRepoRepository) (Report, error)
	DisbursementsByRange(logger *slog.Logger, ctx context.Context, repo repo.DisburserRepoRepository, start time.Time, end time.Time) (Report, error)
	MerchantDisbursements(logger *slog.Logger, ctx context.Context, repo repo.DisburserRepoRepository, merchantUUID uuid.UUID, start time.Time, end time.Time) (Report, error)
}
