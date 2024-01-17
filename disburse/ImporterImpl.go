package disburse

import (
	"cmp"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/levtk/sequra/repo"
	"github.com/levtk/sequra/types"
	"log/slog"
	"os"
	"slices"
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
	//sortOrdersByOrderDate(orders)

	merchants, err = parseDataFromMerchants(i.MerchantsFileName)
	if err != nil {
		i.Logger.Error("failed to parse data from merchants", "error", err.Error())
		return disbursements, merchants, err
	}

	disbursements, err = buildDisbursementRecordsFromImport(orders, merchants)
	return disbursements, merchants, err
}

func sortOrdersByMerchant(orders Orders) {
	slices.SortFunc(orders, func(a, b *Order) int {
		if a != nil && b != nil {
			if n := cmp.Compare(a.MerchantReference, b.MerchantReference); n != 0 {
				return n
			}
			return a.CreatedAt.Compare(b.CreatedAt)
		}
		return 0
	})
}

//func sortOrdersByOrderDate(orders Orders) {
//	sort.Sort(ByOrderDate{orders})
//}

// calculatePayout takes a sorted list of type Orders and calculates their distribution payouts and creates the distribution id.
func buildDisbursementRecordsFromImport(o Orders, m map[string]types.Merchant) ([]types.Disbursement, error) {
	var merchant types.Merchant
	disbursements := make([]types.Disbursement, 1500000)
	var disbursementGroupID, frequency string
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	for i := 0; i < len(o); i++ {
		if o[i] != nil {

			merchant = m[o[i].MerchantReference]
			frequency = merchant.DisbursementFrequency
			if o[i] == nil {
				return disbursements, nil
			} //Loop through all orders
			if i == 0 { //For the first order, create the disbursement for index 0
				disbursements[i].ID = uuid.NewString()
				disbursementGroupID = uuid.NewString()
				orderFee, err := o[i].CalculateOrderFee()
				if err != nil {
					return disbursements, err
				}

				disbursements[i].OrderFee = orderFee
				disbursements[i].RunningTotal = orderFee
				disbursements[i].DisbursementGroupID = disbursementGroupID
				disbursements[i].MerchReference = m[o[i].MerchantReference].Reference
				disbursements[i].OrderID = o[i].ID

				if frequency == types.DAILY {
					createdAt := o[i].CreatedAt.Format(time.DateOnly)
					disbursements[i].PayoutDate = createdAt
				}

				if m[o[i].MerchantReference].DisbursementFrequency == types.WEEKLY {
					payoutDate, err := merchant.CalculatePastPayoutDate(o[i].CreatedAt)
					if err != nil {
						return disbursements, err
					}

					disbursements[i].PayoutDate = payoutDate.Format(time.DateOnly)
				}
				continue
			}
			if i > 0 && o[i] != nil {
				newPayoutPeriod, err := isNewPayoutPeriod(o[i-1], o[i], m[o[i].MerchantReference]) //checking to see if this is a new payout period
				if err != nil {
					return nil, err
				}
				if o[i] != nil {
					frequency = m[o[i].MerchantReference].DisbursementFrequency
					switch frequency {
					case types.DAILY:
						if !newPayoutPeriod {
							orderFee, err := o[i].CalculateOrderFee()
							if err != nil {
								return disbursements, err
							}
							disbursements[i].ID = uuid.NewString()
							disbursements[i].MerchReference = o[i].MerchantReference

							disbursements[i].RunningTotal = orderFee + disbursements[i-1].RunningTotal
							disbursements[i].DisbursementGroupID = disbursements[i-1].DisbursementGroupID //use the previous disbursement group ID

							disbursements[i].OrderID = o[i].ID
							disbursements[i].OrderFee = orderFee
							disbursements[i].PayoutDate = o[i].CreatedAt.Format(time.DateOnly)
							continue
						}
						if newPayoutPeriod {
							orderFee, err := o[i].CalculateOrderFee()
							if err != nil {
								return disbursements, err
							}
							disbursements[i].ID = uuid.NewString()
							disbursements[i-1].PayoutTotal = disbursements[i-1].RunningTotal //The last running total record within the frequency period becomes the PayoutTotal
							disbursements[i].DisbursementGroupID = uuid.NewString()          //create new disbursement group ID
							disbursements[i].MerchReference = o[i].MerchantReference         //set the merchant var to the new value
							disbursements[i].RunningTotal = orderFee
							disbursements[i].OrderID = o[i].ID
							disbursements[i].OrderFee = orderFee
							disbursements[i].PayoutDate = o[i].CreatedAt.Format(time.DateOnly)
							continue
						}
					case types.WEEKLY:
						orderFee, err := o[i].CalculateOrderFee()
						if err != nil {
							return disbursements, err
						}

						if !newPayoutPeriod {
							previousRecordPayoutDate, err := merchant.CalculatePastPayoutDate(o[i-1].CreatedAt)
							if err != nil {
								return disbursements, err
							}

							disbursements[i].ID = uuid.NewString()
							disbursements[i].MerchReference = o[i].MerchantReference
							disbursements[i].RunningTotal = orderFee + disbursements[i-1].RunningTotal
							disbursements[i].DisbursementGroupID = disbursements[i-1].DisbursementGroupID
							disbursements[i].OrderID = o[i].ID
							disbursements[i].OrderFee = orderFee
							disbursements[i].PayoutDate = previousRecordPayoutDate.UTC().Format(time.DateOnly)
							continue
						} else {
							currentRecordsPayoutDate, err := merchant.CalculatePastPayoutDate(o[i].CreatedAt)
							if err != nil {
								return disbursements, err
							}

							disbursements[i].ID = uuid.NewString()
							disbursements[i].MerchReference = o[i].MerchantReference
							disbursements[i].RunningTotal = orderFee
							disbursements[i].DisbursementGroupID = uuid.NewString()
							disbursements[i].OrderID = o[i].ID
							disbursements[i].OrderFee = orderFee
							disbursements[i].PayoutDate = currentRecordsPayoutDate.UTC().Format(time.DateOnly)
							continue
						}
					default:
						logger.Info("failed to match select case for payout period", "got", frequency)

					}
				}
			}
		}
	}
	return disbursements, nil
}

func isNewPayoutPeriod(o1 *Order, o2 *Order, m types.Merchant) (bool, error) {
	if o1 != nil && o2 != nil {
		if m.DisbursementFrequency == types.DAILY {
			if (o1.MerchantReference == o2.MerchantReference) && (o1.CreatedAt != o2.CreatedAt) {
				return true, nil
			}

			if o1.MerchantReference != o2.MerchantReference {
				return true, nil
			}
			return false, nil
		}

		if m.DisbursementFrequency == types.WEEKLY {
			if o1.MerchantReference != o2.MerchantReference {
				return true, nil
			}

			payoutDate1, err := m.CalculatePastPayoutDate(o1.CreatedAt)
			if err != nil {
				return false, err
			}

			payoutDate2, err := m.CalculatePastPayoutDate(o2.CreatedAt)
			if err != nil {
				return false, err
			}

			if o1.MerchantReference == o2.MerchantReference && payoutDate1 == payoutDate2 {
				return false, nil
			}

			if o1.MerchantReference == o2.MerchantReference && payoutDate1 != payoutDate2 {
				return true, nil
			}

		}

		return false, errors.New("no determination of payout same merchant")
	}
	return false, nil
}

func buildWeeklyRecord(o Orders, m map[string]types.Merchant, disbursements []types.Disbursement, i int) ([]types.Disbursement, error) {
	orderFee, err := o[i].CalculateOrderFee()
	if err != nil {
		return disbursements, err
	}

	if i <= len(o)-1 {
		newPayoutPeriod, err := isNewPayoutPeriod(o[i-1], o[i], m[o[i].MerchantReference])
		if err != nil {
			return disbursements, err
		}

		if !newPayoutPeriod {
			disbursements[i].ID = uuid.NewString()
			disbursements[i].RunningTotal = orderFee + disbursements[i-1].RunningTotal
			disbursements[i].DisbursementGroupID = disbursements[i-1].DisbursementGroupID
			disbursements[i].OrderID = o[i].ID
			disbursements[i].OrderFee = orderFee
			disbursements[i].PayoutDate = disbursements[i-1].PayoutDate
			return disbursements, nil
		} else {
			newPayoutPeriod = true
			disbursements[i].DisbursementGroupID = uuid.NewString()
			disbursements[i].ID = uuid.NewString()
			disbursements[i].RunningTotal = orderFee
			disbursements[i].OrderID = o[i].ID
			disbursements[i].OrderFee = orderFee
			merchant := m[o[i].MerchantReference]
			pastPayoutDate, err := merchant.CalculatePastPayoutDate(o[i].CreatedAt)
			if err != nil {
				return disbursements, err
			}
			disbursements[i].PayoutDate = pastPayoutDate.Format(time.DateOnly)
			return disbursements, nil
		}
	}
	return disbursements, nil
}
