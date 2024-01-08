package disburse

import (
	"context"
	"github.com/google/uuid"
	"github.com/levtk/sequra/repo"
	"log/slog"
	"time"
)

type Disburser interface {
	ProcessOrder(op OrderProcessor) error
	ImportOrders(o Importer) error
	GenerateReports(r Reporter) ([]Report, error)
}

type Importer interface {
	ImportOrders() ([]Order, map[string]Merchant, error)
}
type OrderProcessor interface {
	ProcessOrder(logger *slog.Logger, ctx *context.Context, repo *repo.DisburserRepo, o *Order) error
}

type Seller interface {
	GetMinMonthlyFee() (int64, error)
	GetMinMonthlyFeeRemaining() (int64, error)
}
type Reporter interface {
	DisbursementsByYear(logger *slog.Logger, ctx *context.Context, repo *repo.DisburserRepo) (Report, error)
	DisbursementsByRange(logger *slog.Logger, ctx *context.Context, repo *repo.DisburserRepo, start time.Time, end time.Time) (Report, error)
	MerchantDisbursements(logger *slog.Logger, ctx *context.Context, repo *repo.DisburserRepo, merchantUUID uuid.UUID, start time.Time, end time.Time) (Report, error)
}
