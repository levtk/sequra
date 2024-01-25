package types

import (
	"math"
	"strconv"
	"time"
)

// StrToInt64 parses strings to int64 for use as money
func StrToInt64(s string) (int64, error) {
	fltFromStr, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}

	return int64(math.Round(fltFromStr * 100)), nil
}

// IsNewMonth takes two arguments of time.Time and returns true if they do not fall
// in the same month. It should be noted it does not evaluate the year.
func IsNewMonth(lastPeriod time.Time, thisPeriod time.Time) bool {
	if lastPeriod.Month() != thisPeriod.Month() {
		return true
	} else {
		return false
	}
}
