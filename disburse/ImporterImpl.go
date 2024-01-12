package disburse

import (
	"context"
	"github.com/google/uuid"
	"github.com/levtk/sequra/repo"
	"github.com/levtk/sequra/types"
	"log/slog"
	"sort"
	"time"
)

func NewImport(logger *slog.Logger, ctx context.Context, repo repo.DisburserRepoRepository) *Import {
	return &Import{
		Logger:            logger,
		Ctx:               ctx,
		OrdersFileName:    types.OREDERS_FILENAME,
		MerchantsFileName: types.MERCHANTS_FILENAME,
		Repo:              repo,
	}
}

func (i *Import) ImportOrders() ([]types.Disbursement, map[string]types.Merchant, error) {
	var orders Orders
	var disbursements []types.Disbursement
	var merchants map[string]types.Merchant

	orders, err := parseDataFromOrders(i.OrdersFileName)
	if err != nil {
		i.Logger.Error("failed to parse data from orders", "error", err.Error())
		return disbursements, merchants, err
	}

	sortOrdersByMerchant(orders)
	sortOrdersByOrderDate(orders)

	merchants, err = parseDataFromMerchants(i.MerchantsFileName)
	if err != nil {
		i.Logger.Error("failed to parse data from merchants", "error", err.Error())
		return disbursements, merchants, err
	}

	disbursements, err = buildDisbursementRecordsFromImport(orders, merchants)
	return disbursements, merchants, nil
}

func sortOrdersByMerchant(orders Orders) {
	sort.Sort(ByMerchRef{orders})
}

func sortOrdersByOrderDate(orders Orders) {
	sort.Sort(ByOrderDate{orders})
}

// calculatePayout takes a sorted list of type Orders and calculates their distribution payouts and creates the distribution id.
func buildDisbursementRecordsFromImport(o Orders, m map[string]types.Merchant) ([]types.Disbursement, error) {
	var merchant types.Merchant
	var disbursements []types.Disbursement
	var disbursementGroupID, frequency string
	var count int
	newPayoutPeriod := true

	for i := 0; i < len(o); i++ {
		if newPayoutPeriod {
			disbursements[count].PayoutTotal = disbursements[count].RunningTotal //The last running total record within the frequency period becomes the PayoutTotal
			frequency = m[o[i].MerchantReference].DisbursementFrequency
			disbursementGroupID = uuid.NewString()
			disbursements[count].DisbursementGroupID = disbursementGroupID

			merchant = m[o[i].MerchantReference]
		}
		if i < len(o)-1 && !newPayoutPeriod {
			switch frequency {

			case types.DAILY:
				if o[i].MerchantReference == o[i+1].MerchantReference && o[i].CreatedAt == o[i+1].CreatedAt {
					orderFee, err := o[i].CalculateOrderFee()
					if err != nil {
						return disbursements, err
					}

					disbursements[count].RunningTotal = orderFee + disbursements[count-1].RunningTotal
					disbursements[count].DisbursementGroupID = disbursementGroupID
					disbursements[count].OrderID = o[i].ID
					disbursements[count].OrderFee = orderFee
					disbursements[count].PayoutDate = o[i].CreatedAt.Format(time.DateOnly)

					count++

				} else {
					count++
					newPayoutPeriod = true
				}

			case types.WEEKLY:
				orderFee, err := o[i].CalculateOrderFee()
				if err != nil {
					return disbursements, err
				}

				currentRecordPayoutDate, err := merchant.CalculatePastPayoutDate(o[i].CreatedAt)
				if err != nil {
					return disbursements, err
				}
				nextRecordsPayoutDate, err := merchant.CalculatePastPayoutDate(o[i+1].CreatedAt)
				if err != nil {
					return disbursements, err
				}

				if o[i].MerchantReference == o[i+1].MerchantReference && currentRecordPayoutDate == nextRecordsPayoutDate {
					disbursements[count].RunningTotal = orderFee + disbursements[count-1].RunningTotal
					disbursements[count].DisbursementGroupID = disbursementGroupID
					disbursements[count].OrderID = o[i].ID
					disbursements[count].OrderFee = orderFee
					disbursements[count].PayoutDate = currentRecordPayoutDate.UTC().Format(time.DateOnly)

					count++

				} else {
					newPayoutPeriod = true
					count++
				}
			}
		}
	}
	return disbursements, nil
}
