package disburse

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/levtk/sequra/repo"
	"github.com/levtk/sequra/types"
	"log/slog"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestMerchant_CalculatePastPayoutDate(t *testing.T) {
	liveOn, err := time.Parse(time.DateOnly, "2022-11-09")
	orderDate, err := time.Parse(time.DateOnly, "2022-12-08")
	wantDate, err := time.Parse(time.DateOnly, "2022-12-14")
	if err != nil {
		t.Log(err.Error())
	}
	type fields struct {
		ID                    uuid.UUID
		Reference             string
		Email                 string
		LiveOn                time.Time
		DisbursementFrequency string
		MinMonthlyFee         string
	}
	type args struct {
		t time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   time.Time
	}{
		{name: "weekly liveOn Date of 2022-11-09",
			fields: fields{
				ID:                    uuid.UUID{},
				Reference:             "deckow_gibson",
				Email:                 "info@deckow-gibson.com",
				LiveOn:                liveOn,
				DisbursementFrequency: "WEEKLY",
				MinMonthlyFee:         "15.0",
			},
			args: args{orderDate},
			want: wantDate},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &types.Merchant{
				ID:                    tt.fields.ID,
				Reference:             tt.fields.Reference,
				Email:                 tt.fields.Email,
				LiveOn:                tt.fields.LiveOn,
				DisbursementFrequency: tt.fields.DisbursementFrequency,
				MinMonthlyFee:         tt.fields.MinMonthlyFee,
			}
			if got, _ := m.CalculatePastPayoutDate(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalculatePastPayoutDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMerchant_GetMinMonthlyFee(t *testing.T) {
	type fields struct {
		ID                    uuid.UUID
		Reference             string
		Email                 string
		LiveOn                time.Time
		DisbursementFrequency string
		MinMonthlyFee         string
	}
	tests := []struct {
		name    string
		fields  fields
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &types.Merchant{
				ID:                    tt.fields.ID,
				Reference:             tt.fields.Reference,
				Email:                 tt.fields.Email,
				LiveOn:                tt.fields.LiveOn,
				DisbursementFrequency: tt.fields.DisbursementFrequency,
				MinMonthlyFee:         tt.fields.MinMonthlyFee,
			}
			got, err := m.GetMinMonthlyFee()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMinMonthlyFee() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetMinMonthlyFee() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMerchant_GetMinMonthlyFeeRemaining(t *testing.T) {
	type fields struct {
		ID                    uuid.UUID
		Reference             string
		Email                 string
		LiveOn                time.Time
		DisbursementFrequency string
		MinMonthlyFee         string
	}
	tests := []struct {
		name    string
		fields  fields
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &types.Merchant{
				ID:                    tt.fields.ID,
				Reference:             tt.fields.Reference,
				Email:                 tt.fields.Email,
				LiveOn:                tt.fields.LiveOn,
				DisbursementFrequency: tt.fields.DisbursementFrequency,
				MinMonthlyFee:         tt.fields.MinMonthlyFee,
			}
			got, err := m.GetMinMonthlyFeeRemaining()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMinMonthlyFeeRemaining() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetMinMonthlyFeeRemaining() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMerchant_GetNextPayoutDate(t *testing.T) {
	type fields struct {
		ID                    uuid.UUID
		Reference             string
		Email                 string
		LiveOn                time.Time
		DisbursementFrequency string
		MinMonthlyFee         string
	}
	tests := []struct {
		name    string
		fields  fields
		want    time.Time
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &types.Merchant{
				ID:                    tt.fields.ID,
				Reference:             tt.fields.Reference,
				Email:                 tt.fields.Email,
				LiveOn:                tt.fields.LiveOn,
				DisbursementFrequency: tt.fields.DisbursementFrequency,
				MinMonthlyFee:         tt.fields.MinMonthlyFee,
			}
			got, err := m.GetNextPayoutDate()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNextPayoutDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNextPayoutDate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMerchant_CalculateDailyTotalOrders(t *testing.T) {
	type fields struct {
		ID                    uuid.UUID
		Reference             string
		Email                 string
		LiveOn                time.Time
		DisbursementFrequency string
		MinMonthlyFee         string
	}
	tests := []struct {
		name    string
		fields  fields
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := types.Merchant{
				ID:                    tt.fields.ID,
				Reference:             tt.fields.Reference,
				Email:                 tt.fields.Email,
				LiveOn:                tt.fields.LiveOn,
				DisbursementFrequency: tt.fields.DisbursementFrequency,
				MinMonthlyFee:         tt.fields.MinMonthlyFee,
			}
			got, err := m.CalculateDailyTotalOrders()
			if (err != nil) != tt.wantErr {
				t.Errorf("calculateDailyTotalOrders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("calculateDailyTotalOrders() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMerchant_CalculateWeeklyTotalOrders(t *testing.T) {
	type fields struct {
		ID                    uuid.UUID
		Reference             string
		Email                 string
		LiveOn                time.Time
		DisbursementFrequency string
		MinMonthlyFee         string
	}
	tests := []struct {
		name    string
		fields  fields
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := types.Merchant{
				ID:                    tt.fields.ID,
				Reference:             tt.fields.Reference,
				Email:                 tt.fields.Email,
				LiveOn:                tt.fields.LiveOn,
				DisbursementFrequency: tt.fields.DisbursementFrequency,
				MinMonthlyFee:         tt.fields.MinMonthlyFee,
			}
			got, err := m.CalculateWeeklyTotalOrders()
			if (err != nil) != tt.wantErr {
				t.Errorf("calculateWeeklyTotalOrders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("calculateWeeklyTotalOrders() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDisburserService(t *testing.T) {
	type args struct {
		logger *slog.Logger
		ctx    context.Context
		db     *sql.DB
	}
	tests := []struct {
		name    string
		args    args
		want    *DisburserService
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDisburserService(tt.args.logger, tt.args.ctx, tt.args.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDisburserService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDisburserService() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewOrder(t *testing.T) {
	type args struct {
		id                string
		merchantReference string
		amount            int64
	}
	tests := []struct {
		name string
		args args
		want *Order
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOrder(tt.args.id, tt.args.merchantReference, tt.args.amount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewReporter(t *testing.T) {
	type args struct {
		logger *slog.Logger
		ctx    context.Context
		repo   *repo.DisburserRepo
	}
	tests := []struct {
		name string
		args args
		want *Report
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewReporter(tt.args.logger, tt.args.ctx, tt.args.repo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewReporter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrder_CalculateOrderFee(t *testing.T) {
	type fields struct {
		ID                string
		MerchantReference string
		MerchantID        uuid.UUID
		Amount            int64
		CreatedAt         time.Time
		RWMutex           sync.RWMutex
	}
	tests := []struct {
		name    string
		fields  fields
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Order{
				ID:                tt.fields.ID,
				MerchantReference: tt.fields.MerchantReference,
				MerchantID:        tt.fields.MerchantID,
				Amount:            tt.fields.Amount,
				CreatedAt:         tt.fields.CreatedAt,
				RWMutex:           tt.fields.RWMutex,
			}
			got, err := o.CalculateOrderFee()
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateOrderFee() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CalculateOrderFee() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrder_IsBeforeTimeCutOff(t *testing.T) {
	type fields struct {
		ID                string
		MerchantReference string
		MerchantID        uuid.UUID
		Amount            int64
		CreatedAt         time.Time
		RWMutex           sync.RWMutex
	}
	tests := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Order{
				ID:                tt.fields.ID,
				MerchantReference: tt.fields.MerchantReference,
				MerchantID:        tt.fields.MerchantID,
				Amount:            tt.fields.Amount,
				CreatedAt:         tt.fields.CreatedAt,
				RWMutex:           tt.fields.RWMutex,
			}
			got, err := o.IsBeforeTimeCutOff()
			if (err != nil) != tt.wantErr {
				t.Errorf("IsBeforeTimeCutOff() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsBeforeTimeCutOff() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrder_ProcessOrder(t *testing.T) {
	type fields struct {
		ID                string
		MerchantReference string
		MerchantID        uuid.UUID
		Amount            int64
		CreatedAt         time.Time
		RWMutex           sync.RWMutex
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Order{
				ID:                tt.fields.ID,
				MerchantReference: tt.fields.MerchantReference,
				MerchantID:        tt.fields.MerchantID,
				Amount:            tt.fields.Amount,
				CreatedAt:         tt.fields.CreatedAt,
				RWMutex:           tt.fields.RWMutex,
			}
			if err := o.ProcessOrder(); (err != nil) != tt.wantErr {
				t.Errorf("ProcessOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOrders_Len(t *testing.T) {
	tests := []struct {
		name string
		s    Orders
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrders_Swap(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		s    Orders
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Swap(tt.args.i, tt.args.j)
		})
	}
}

func Test_calculateOrderFee(t *testing.T) {
	type args struct {
		orderAmt int64
	}
	tests := []struct {
		name         string
		args         args
		wantOrderFee int64
		wantErr      bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOrderFee, err := calculateOrderFee(tt.args.orderAmt)
			if (err != nil) != tt.wantErr {
				t.Errorf("calculateOrderFee() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOrderFee != tt.wantOrderFee {
				t.Errorf("calculateOrderFee() gotOrderFee = %v, want %v", gotOrderFee, tt.wantOrderFee)
			}
		})
	}
}

func Test_getMerchant(t *testing.T) {
	type args struct {
		merchRef string
	}
	tests := []struct {
		name    string
		args    args
		want    types.Merchant
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getMerchant(tt.args.merchRef)
			if (err != nil) != tt.wantErr {
				t.Errorf("getMerchant() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getMerchant() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getMerchantReferenceFromOrder(t *testing.T) {
	type args struct {
		o Order
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getMerchantReferenceFromOrder(tt.args.o)
			if (err != nil) != tt.wantErr {
				t.Errorf("getMerchantReferenceFromOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getMerchantReferenceFromOrder() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getOrdersByMerchRef(t *testing.T) {
	type args struct {
		merchRef string
	}
	tests := []struct {
		name    string
		args    args
		want    []Order
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getOrdersByMerchRef(tt.args.merchRef)
			if (err != nil) != tt.wantErr {
				t.Errorf("getOrdersByMerchRef() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getOrdersByMerchRef() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isBeforeCutOffTime(t *testing.T) {
	tests := []struct {
		name    string
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := isBeforeCutOffTime()
			if (err != nil) != tt.wantErr {
				t.Errorf("isBeforeCutOffTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("isBeforeCutOffTime() got = %v, want %v", got, tt.want)
			}
		})
	}
}
