package disburse

import (
	"cmp"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/levtk/sequra/repo"
	"github.com/levtk/sequra/types"
	"log/slog"
	"net/http"
	"os"
	"slices"
	"time"
)

func NewImport(logger *slog.Logger, ctx context.Context, repo *repo.DisburserRepo) *Import {
	return &Import{
		Logger:            logger,
		Ctx:               ctx,
		OrdersFileName:    types.OREDERS_FILENAME,
		MerchantsFileName: types.MERCHANTS_FILENAME,
		Repo:              repo,
	}
}

func (i *Import) ImportOrders() ([]types.Disbursement, map[string]types.Merchant, []types.Monthly, error) {
	var orders Orders
	var disbursements []types.Disbursement
	var merchants map[string]types.Merchant
	var monthly []types.Monthly

	orders, err := parseDataFromOrders(i.OrdersFileName)
	if err != nil {
		i.Logger.Error("failed to parse data from orders", "error", err.Error())
		return disbursements, merchants, monthly, err
	}

	sortOrdersByMerchant(orders)

	merchants, err = parseDataFromMerchants(i.MerchantsFileName)
	if err != nil {
		i.Logger.Error("failed to parse data from merchants", "error", err.Error())
		return disbursements, merchants, monthly, err
	}

	disbursements, monthly, err = buildDisbursementRecordsFromImport(orders, merchants)
	return disbursements, merchants, monthly, err
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

// TODO add monthly fees charged logic and to disbursement or another table.
// calculatePayout takes a sorted list of type Orders and calculates their distribution payouts and creates the distribution id.
func buildDisbursementRecordsFromImport(o Orders, m map[string]types.Merchant) ([]types.Disbursement, []types.Monthly, error) {
	var merchant types.Merchant
	disbursements := make([]types.Disbursement, 1500000)
	monthly := make([]types.Monthly, 2000000)
	monthlyCount := 0
	var disbursementGroupID, frequency string
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	for i := 0; i < len(o); i++ {
		if o[i] != nil {
			merchant = m[o[i].MerchantReference]
			frequency = merchant.DisbursementFrequency

			if i == 0 { //For the first order, create the disbursement for index 0
				disbursements[i].RecordUUID = uuid.NewString()
				disbursementGroupID = uuid.NewString()
				orderFee, err := o[i].CalculateOrderFee()
				if err != nil {
					return disbursements, monthly, err
				}

				disbursements[i].OrderFee = orderFee
				disbursements[i].OrderFeeRunningTotal = orderFee
				disbursements[i].PayoutRunningTotal = o[i].Amount - orderFee
				disbursements[i].DisbursementGroupID = disbursementGroupID
				disbursements[i].MerchReference = m[o[i].MerchantReference].Reference
				disbursements[i].OrderID = o[i].ID

				if frequency == types.DAILY {
					disbursements[i].PayoutDate = o[i].CreatedAt
				}

				if m[o[i].MerchantReference].DisbursementFrequency == types.WEEKLY {
					disbursements[i].PayoutDate, err = merchant.CalculatePastPayoutDate(o[i].CreatedAt)
					if err != nil {
						return disbursements, monthly, err
					}
				}
				continue
			}
			if i > 0 && o[i] != nil { //after first disbursement record is created we can calculate values based on the past records
				newPayoutPeriod, err := isNewPayoutPeriod(o[i-1], o[i], m[o[i].MerchantReference]) //checking to see if this is a new payout period
				if err != nil {
					return nil, monthly, err
				}

				if o[i] != nil {
					frequency = m[o[i].MerchantReference].DisbursementFrequency
					switch frequency {
					case types.DAILY:
						if !newPayoutPeriod {
							orderFee, err := o[i].CalculateOrderFee()
							if err != nil {
								return disbursements, monthly, err
							}
							disbursements[i].RecordUUID = uuid.NewString()
							disbursements[i].MerchReference = o[i].MerchantReference

							disbursements[i].OrderFeeRunningTotal = orderFee + disbursements[i-1].OrderFeeRunningTotal
							disbursements[i].PayoutRunningTotal = disbursements[i-1].PayoutRunningTotal + (o[i].Amount - orderFee)
							disbursements[i].DisbursementGroupID = disbursements[i-1].DisbursementGroupID //use the previous disbursement group RecordUUID

							disbursements[i].OrderID = o[i].ID
							disbursements[i].OrderFee = orderFee
							disbursements[i].PayoutDate = o[i].CreatedAt
							continue

						} else {
							orderFee, err := o[i].CalculateOrderFee()
							if err != nil {
								return disbursements, monthly, err
							}

							disbursements[i].RecordUUID = uuid.NewString()
							disbursements[i-1].PayoutTotal = disbursements[i-1].PayoutRunningTotal //The last running total record within the frequency period becomes the PayoutTotal
							disbursements[i].DisbursementGroupID = uuid.NewString()                //create new disbursement group RecordUUID
							disbursements[i].MerchReference = o[i].MerchantReference               //set the merchant var to the new value
							disbursements[i].OrderFeeRunningTotal = orderFee
							disbursements[i].PayoutRunningTotal = o[i].Amount - orderFee
							disbursements[i].OrderID = o[i].ID
							disbursements[i].OrderFee = orderFee
							disbursements[i].PayoutDate = o[i].CreatedAt

							disbursements[i-1].PayoutTotal = disbursements[i-1].PayoutRunningTotal
							disbursements[i-1].IsPaidOut = true

							if types.IsNewMonth(disbursements[i-1].PayoutDate, disbursements[i].PayoutDate) {
								monthlyFee, err := types.StrToInt64(m[o[i].MerchantReference].MinMonthlyFee)
								if err != nil {
									logger.Error("failed to parse monthly fee to int64", "error", err)
								}

								var didPayFee = 1
								if (disbursements[i-1].OrderFeeRunningTotal-monthlyFee > 0) || m[o[i].MerchantReference].MinMonthlyFee == "0.0" {
									didPayFee = 0
								}

								monthly[monthlyCount] = types.Monthly{
									ID:                uuid.New(),
									MerchantReference: m[o[i].MerchantReference].Reference,
									MerchantID:        merchant.ID,
									MonthlyFeeDate:    disbursements[i].PayoutDate,
									DidPayFee:         didPayFee,
									MonthlyFee:        monthlyFee,
									TotalOrderAmt:     disbursements[i-1].PayoutRunningTotal,
									OrderFeeTotal:     disbursements[i-1].OrderFeeRunningTotal,
									CreatedAt:         disbursements[i].PayoutDate,
									UpdatedAt:         time.Now().UTC(),
								}
								monthlyCount++
							}
							continue
						}
					case types.WEEKLY:
						orderFee, err := o[i].CalculateOrderFee()
						if err != nil {
							return disbursements, monthly, err
						}

						if !newPayoutPeriod {
							previousRecordPayoutDate, err := merchant.CalculatePastPayoutDate(o[i-1].CreatedAt)
							if err != nil {
								return disbursements, monthly, err
							}

							disbursements[i].RecordUUID = uuid.NewString()
							disbursements[i].MerchReference = o[i].MerchantReference
							disbursements[i].OrderFeeRunningTotal = orderFee + disbursements[i-1].OrderFeeRunningTotal
							disbursements[i].PayoutRunningTotal = disbursements[i-1].PayoutRunningTotal + (o[i].Amount - orderFee)
							disbursements[i].DisbursementGroupID = disbursements[i-1].DisbursementGroupID
							disbursements[i].OrderID = o[i].ID
							disbursements[i].OrderFee = orderFee
							disbursements[i].PayoutDate = previousRecordPayoutDate.UTC()
							continue
						} else {
							currentRecordsPayoutDate, err := merchant.CalculatePastPayoutDate(o[i].CreatedAt)
							if err != nil {
								return disbursements, monthly, err
							}

							disbursements[i].RecordUUID = uuid.NewString()
							disbursements[i].MerchReference = o[i].MerchantReference
							disbursements[i].OrderFeeRunningTotal = orderFee
							disbursements[i].PayoutRunningTotal = o[i].Amount - orderFee
							disbursements[i].DisbursementGroupID = uuid.NewString()
							disbursements[i].OrderID = o[i].ID
							disbursements[i].OrderFee = orderFee
							disbursements[i].PayoutDate = currentRecordsPayoutDate.UTC()

							disbursements[i-1].PayoutTotal = disbursements[i-1].PayoutRunningTotal
							disbursements[i-1].IsPaidOut = true

							if types.IsNewMonth(disbursements[i-1].PayoutDate, disbursements[i].PayoutDate) {
								monthlyFee, err := types.StrToInt64(m[o[i].MerchantReference].MinMonthlyFee)
								if err != nil {
									logger.Error("failed to parse monthly fee to int64", "error", err)
								}

								didPayFee := 1
								if (disbursements[i-1].OrderFeeRunningTotal-monthlyFee > 0) || m[o[i].MerchantReference].MinMonthlyFee == "0.0" {
									didPayFee = 0
								}

								monthly[monthlyCount] = types.Monthly{
									ID:                uuid.New(),
									MerchantReference: m[o[i].MerchantReference].Reference,
									MerchantID:        merchant.ID,
									MonthlyFeeDate:    disbursements[i].PayoutDate,
									DidPayFee:         didPayFee,
									MonthlyFee:        monthlyFee,
									TotalOrderAmt:     disbursements[i-1].PayoutRunningTotal,
									OrderFeeTotal:     disbursements[i-1].OrderFeeRunningTotal,
									CreatedAt:         disbursements[i].PayoutDate,
									UpdatedAt:         time.Now().UTC(),
								}
								monthlyCount++
							}
							continue
						}
					default:
						logger.Info("failed to match select case for payout period", "got", frequency)

					}
				}
			}
		}
	}
	return disbursements, monthly, nil
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
			disbursements[i].RecordUUID = uuid.NewString()
			disbursements[i].OrderFeeRunningTotal = orderFee + disbursements[i-1].OrderFeeRunningTotal
			disbursements[i].DisbursementGroupID = disbursements[i-1].DisbursementGroupID
			disbursements[i].OrderID = o[i].ID
			disbursements[i].OrderFee = orderFee
			disbursements[i].PayoutDate = disbursements[i-1].PayoutDate
			return disbursements, nil
		} else {
			newPayoutPeriod = true
			disbursements[i].DisbursementGroupID = uuid.NewString()
			disbursements[i].RecordUUID = uuid.NewString()
			disbursements[i].OrderFeeRunningTotal = orderFee
			disbursements[i].OrderID = o[i].ID
			disbursements[i].OrderFee = orderFee
			merchant := m[o[i].MerchantReference]
			pastPayoutDate, err := merchant.CalculatePastPayoutDate(o[i].CreatedAt)
			if err != nil {
				return disbursements, err
			}
			disbursements[i].PayoutDate = pastPayoutDate
			return disbursements, nil
		}
	}
	return disbursements, nil
}

func calculateDisbursementPeriodPayout(disbursements []types.Disbursement, i int) (int64, error) {
	payoutTotal := disbursements[i].PayoutRunningTotal - disbursements[i].OrderFeeRunningTotal
	if payoutTotal < 0 {
		return 0, errors.New("payout total less than zero and is not an accounting adjustment")
	}
	return payoutTotal, nil
}

func (i *Import) Import(w http.ResponseWriter, r *http.Request) {

	op := NewOrderProcessor(i.Logger, i.Ctx, i.Repo)
	distributions, merchants, monthly, err := i.ImportOrders()
	if err != nil {
		i.Logger.Error("failed to import orders or merchants", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, v := range merchants {
		err := i.Repo.InsertMerchant(v)
		if err != nil {
			i.Logger.Error(err.Error())
		}

	}

	err = op.ProcessBatchMonthly(monthly)
	if err != nil {
		i.Logger.Error("failed to process batch monthly records", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = op.ProcessBatchDistributions(distributions)
	if err != nil {
		i.Logger.Error("failed to process batch distributions", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
