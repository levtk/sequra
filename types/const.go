package types

const (
	RATE_LESS_THAN_50       int64  = 10
	RATE_BETWEEN_50_AND_300 int64  = 5
	RATE_ABOVE_300          int64  = 25
	MAX_ORDER               int64  = 1000000 //Should be configured per Merchant during onboarding
	TIME_CUT_OFF            string = "23:00:00"
	OREDERS_FILENAME               = "orders.csv"
	MERCHANTS_FILENAME             = "merchants.csv"
	WEEKLY                         = "WEEKLY"
	DAILY                          = "DAILY"
	SQL_NO_ROWS                    = "sql: no rows in result set"
)
