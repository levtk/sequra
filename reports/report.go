package reports

import (
	"context"
	"github.com/levtk/sequra/types"
	"log/slog"
	"time"
)

type Report struct {
	Logger   *slog.Logger
	Ctx      context.Context
	Name     string
	Merchant types.Merchant
	Start    time.Time
	End      time.Time
	Data     []byte
}
