package disburse

import (
	"github.com/google/uuid"
	"reflect"
	"testing"
	"time"
)

func Test_importDataFromOrders(t *testing.T) {
	o := make([]Order, 1310000)
	created, _ := time.Parse(time.DateOnly, "2023-02-01")
	o1 := Order{ID: "e653f3e14bc4", MerchantReference: "padberg_group", Amount: 10229, CreatedAt: created}
	o2 := Order{ID: "20b674c93ea6", MerchantReference: "padberg_group", Amount: 43321, CreatedAt: created}
	o3 := Order{ID: "0b73fb1d3332", MerchantReference: "padberg_group", Amount: 19437, CreatedAt: created}
	o[0] = o1
	o[1] = o2
	o[2] = o3

	tests := []struct {
		name     string
		fileName string
		want     []Order
		wantErr  bool
	}{

		{name: "success", fileName: "../orders_test.csv", want: o, wantErr: false},
		{name: "fileNotFound", fileName: "../order.csv", want: []Order{}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDataFromOrders(tt.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDataFromOrders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseDataFromOrders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_importDataFromMerchants(t *testing.T) {
	uu1, _ := uuid.Parse("86312006-4d7e-45c4-9c28-788f4aa68a62")
	uu2, _ := uuid.Parse("d1649242-a612-46ba-82d8-225542bb9576")
	uu3, _ := uuid.Parse("a616488f-c8b2-45dd-b29f-364d12a20238")
	uu4, _ := uuid.Parse("9b6d2b8a-f06c-4298-8f27-f33545eb5899")

	lo1, _ := time.Parse(time.DateOnly, "2023-02-01")
	lo2, _ := time.Parse(time.DateOnly, "2022-12-14")
	lo3, _ := time.Parse(time.DateOnly, "2022-12-10")
	lo4, _ := time.Parse(time.DateOnly, "2022-11-09")

	refs := []string{"padberg_group", "deckow_gibson", "romaguera_and_sons", "rosenbaum_parisian"}

	merchants := map[string]Merchant{
		"padberg_group": {
			ID:                    uu1,
			Reference:             "padberg_group",
			Email:                 "info@padberg-group.com",
			LiveOn:                lo1,
			DisbursementFrequency: "DAILY",
			MinMonthlyFee:         "0.0",
		},
		"deckow_gibson": {
			ID:                    uu2,
			Reference:             "deckow_gibson",
			Email:                 "info@deckow-gibson.com",
			LiveOn:                lo2,
			DisbursementFrequency: "DAILY",
			MinMonthlyFee:         "30.0",
		},
		"romaguera_and_sons": {
			ID:                    uu3,
			Reference:             "romaguera_and_sons",
			Email:                 "info@romaguera-and-sons.com",
			LiveOn:                lo3,
			DisbursementFrequency: "DAILY",
			MinMonthlyFee:         "15.0",
		},
		"rosenbaum_parisian": {
			ID:                    uu4,
			Reference:             "rosenbaum_parisian",
			Email:                 "info@rosenbaum-parisian.com",
			LiveOn:                lo4,
			DisbursementFrequency: "WEEKLY",
			MinMonthlyFee:         "15.0",
		},
	}
	tests := []struct {
		name     string
		fileName string
		want     map[string]Merchant
		wantErr  bool
	}{
		{name: "success", fileName: "../merchants_test.csv", want: merchants, wantErr: false},
		{name: "fileNotFound", fileName: "../merch.csv", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDataFromMerchants(tt.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDataFromMerchants() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for _, v := range refs {
				if !reflect.DeepEqual(got[v], tt.want[v]) {
					t.Errorf("parseDataFromMerchants() = %v, want %v", got[v], tt.want[v])
				}
			}
		})
	}
}
