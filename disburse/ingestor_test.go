package disburse

import (
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
		name    string
		want    []Order
		wantErr bool
	}{

		{name: "success", want: o, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := importDataFromOrders("../orders_test.csv")
			if (err != nil) != tt.wantErr {
				t.Errorf("importDataFromOrders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("importDataFromOrders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_importDataFromMerchants(t *testing.T) {
	tests := []struct {
		name    string
		want    map[string]Merchant
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := importDataFromMerchants()
			if (err != nil) != tt.wantErr {
				t.Errorf("importDataFromMerchants() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("importDataFromMerchants() = %v, want %v", got, tt.want)
			}
		})
	}
}
