package disburse

import (
	"context"
	"github.com/levtk/sequra/repo"
	"log/slog"
)

func NewImport(logger *slog.Logger, ctx context.Context, repo repo.DisburserRepoRepository) *Import {
	return &Import{
		logger:            logger,
		ctx:               ctx,
		ordersFileName:    OREDERS_FILENAME,
		merchantsFileName: MERCHANTS_FILENAME,
		repo:              repo,
	}
}

func (i *Import) ImportOrders() ([]Order, map[string]Merchant, error) {
	var orders []Order
	var merchants map[string]Merchant

	orders, err := parseDataFromOrders(i.ordersFileName)
	if err != nil {
		i.logger.Error("failed to parse data from orders", err.Error())
		return orders, merchants, err
	}

	merchants, err = parseDataFromMerchants(i.merchantsFileName)
	if err != nil {
		i.logger.Error("failed to parse data from merchants", err.Error())
		return orders, merchants, err
	}

	return orders, merchants, nil
}
