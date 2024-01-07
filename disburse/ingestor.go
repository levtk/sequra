package disburse

import (
	"bufio"
	"encoding/csv"
	"io"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// parseDataFromOrders parses the order data that was exported to a semicolon separated file formatted
// per the legacy design specification prior to the new requirements documented in [link to jira story]
func parseDataFromOrders(fileName string) ([]Order, error) {
	var o = make([]Order, 1310000)
	var counter = 0
	ofd, err := os.Open(fileName)
	if err != nil {
		return []Order{}, err
	}
	defer ofd.Close()

	r := csv.NewReader(bufio.NewReader(ofd))
	r.Comma = ';'

	for {

		rec, err := r.Read()
		if err == nil {
			line, _ := r.FieldPos(0)
			if line == 1 { //skipping header line
				continue
			}
			o[counter].ID = rec[0]
			o[counter].MerchantReference = rec[1]
			amount, err := strconv.ParseFloat(rec[2], 64)
			if err != nil {
				return o, err
			}

			o[counter].Amount = int64(math.Round(amount * 100))
			if err != nil {
				return o, err
			}
			o[counter].CreatedAt, err = time.Parse("2006-01-02", rec[3])
			if err != nil {
				return o, err
			}
			counter++
		}

		if err == io.EOF {
			return o, nil
		}

		if err != nil && err != io.EOF {
			return o, err
		}
	}
}

// parseDataFromMerchants parses the order data that was exported to a semicolon separated file formatted
// per the legacy design specification prior to the new requirements documented in [link to jira story]
func parseDataFromMerchants(fileName string) (map[string]Merchant, error) {
	var m = map[string]Merchant{}
	mfd, err := os.Open(fileName)

	if err != nil {
		return map[string]Merchant{}, err
	}

	defer mfd.Close()

	r := csv.NewReader(bufio.NewReader(mfd))
	r.Comma = ';'

	for {
		rec, err := r.Read()
		if err == io.EOF {
			return m, nil
		}

		if err != nil {
			return m, err
		}

		line, _ := r.FieldPos(0)
		if line == 1 { //skipping header line
			continue
		}

		uuid, err := uuid.Parse(rec[0])
		if err != nil {
			return m, err
		}

		liveon, err := time.Parse(time.DateOnly, rec[3])
		if err != nil {
			return m, err
		}

		if err == nil {
			merchant := Merchant{ID: uuid, Reference: rec[1], Email: rec[2], LiveOn: liveon, DisbursementFrequency: rec[4], MinMonthlyFee: rec[5]}
			m[rec[1]] = merchant
		}
	}
}
