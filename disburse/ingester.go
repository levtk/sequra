package disburse

import (
	"bufio"
	"encoding/csv"
	"github.com/google/uuid"
	"github.com/levtk/sequra/types"
	"io"
	"log/slog"
	"math"
	"os"
	"strconv"
	"sync"
	"time"
)

// parseDataFromOrders parses the order data that was exported to a semicolon separated file formatted
// per the legacy design specification prior to the new requirements documented in [link to jira story]
func parseDataFromOrders(fileName string) (Orders, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	o := make([]*Order, 1500000)
	var counter = 0
	ofd, err := os.Open(fileName)
	if err != nil {
		return o, err
	}
	defer func(ofd *os.File) {
		err := ofd.Close()
		if err != nil {

		}
	}(ofd)

	r := csv.NewReader(bufio.NewReader(ofd))
	r.Comma = ';'

	for {
		rec, err := r.Read()
		if err != nil && err != io.EOF {
			logger.Error("error while reading orders file", "error", err)
		}

		if err == nil {
			line, _ := r.FieldPos(1)
			if line == 1 { //skipping header line
				continue
			}

			amount, err := strconv.ParseFloat(rec[2], 64)
			if err != nil {
				return o, err
			}

			a := int64(math.Round(amount * 100))
			if err != nil {
				return o, err
			}
			createdAt, err := time.Parse("2006-01-02", rec[3])
			if err != nil {
				return o, err
			}

			order := Order{
				ID:                rec[0],
				MerchantReference: rec[1],
				MerchantID:        uuid.UUID{},
				Amount:            a,
				CreatedAt:         createdAt,
				RWMutex:           sync.RWMutex{},
			}
			o[counter] = &order
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
// which returns a map[string]types.Merchant where the key is Merchant.reference
func parseDataFromMerchants(fileName string) (map[string]types.Merchant, error) {
	var m = map[string]types.Merchant{}
	mfd, err := os.Open(fileName)

	if err != nil {
		return map[string]types.Merchant{}, err
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
			merchant := types.Merchant{ID: uuid, Reference: rec[1], Email: rec[2], LiveOn: liveon, DisbursementFrequency: rec[4], MinMonthlyFee: rec[5]}
			m[rec[1]] = merchant
		}
	}
}
