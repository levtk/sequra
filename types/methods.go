package types

import (
	"errors"
	"strconv"
	"time"
)

func (m *Merchant) GetMinMonthlyFee() (int64, error) {
	mmf, err := strconv.ParseFloat(m.MinMonthlyFee, 64)
	if err != nil {
		return 0, err
	}
	return int64(mmf * 100), nil
}

func (m *Merchant) GetMinMonthlyFeeRemaining() (int64, error) {
	//TODO Implement
	return 0, errors.New("not implemented.")
}

func (m *Merchant) GetNextPayoutDate() (time.Time, error) {
	wd := m.LiveOn.UTC().Weekday()
	today := time.Now().UTC().Weekday()

	todayDate := time.Now().UTC()

	if int(today) == int(wd) {
		return todayDate, nil
	}

	daysUntil := 7 - int(today)
	return todayDate.AddDate(0, 0, daysUntil), nil
}

// CalculatePastPayoutDate assumes that the imported orders order date was received before 08:00 UTC and thus already recorded
// correctly since we have no way of knowing what time the orders came in based on the data only format. If we can get the full timestamp of the order
// receipt within our systems, we can then determine the true payout date if the system logic for time cutoff was the same when the original order was processed.
func (m *Merchant) CalculatePastPayoutDate(t time.Time) (time.Time, error) {
	wd := m.LiveOn.UTC().Weekday()
	orderDate := t.UTC()
	orderDayOfWeek := orderDate.Weekday()
	if m.DisbursementFrequency == WEEKLY {
		if int(wd) > int(orderDayOfWeek) {
			days := int(wd) - int(orderDayOfWeek)
			return orderDate.AddDate(0, 0, days), nil
		}
		if int(wd) < int(orderDayOfWeek) {
			days := 7 - (int(orderDayOfWeek) - int(wd))
			return orderDate.AddDate(0, 0, days), nil
		}

		if int(wd) == int(orderDayOfWeek) {
			return orderDate, nil
		}
	}
	if m.DisbursementFrequency == DAILY {
		return orderDate.AddDate(0, 0, 1), nil
	}
	return time.Now().UTC(), errors.New("failed to calculate past payout date")
}

func (m Merchant) CalculateDailyTotalOrders() (int64, error) {
	//TODO implement
	return -1, nil
}

func (m Merchant) CalculateWeeklyTotalOrders() (int64, error) {
	//TODO implement
	return -1, nil
}
