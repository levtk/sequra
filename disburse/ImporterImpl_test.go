package disburse

import (
	"context"
	"github.com/levtk/sequra/repo"
	"github.com/levtk/sequra/types"
	"log/slog"
	"reflect"
	"testing"
)

func TestImport_ImportOrders(t *testing.T) {
	type fields struct {
		Logger            *slog.Logger
		Ctx               context.Context
		Repo              repo.DisburserRepoRepository
		OrdersFileName    string
		MerchantsFileName string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []types.Disbursement
		want1   map[string]types.Merchant
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Import{
				Logger:            tt.fields.Logger,
				Ctx:               tt.fields.Ctx,
				Repo:              tt.fields.Repo,
				OrdersFileName:    tt.fields.OrdersFileName,
				MerchantsFileName: tt.fields.MerchantsFileName,
			}
			got, got1, err := i.ImportOrders()
			if (err != nil) != tt.wantErr {
				t.Errorf("ImportOrders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ImportOrders() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ImportOrders() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNewImport(t *testing.T) {
	type args struct {
		logger *slog.Logger
		ctx    context.Context
		repo   repo.DisburserRepoRepository
	}
	tests := []struct {
		name string
		args args
		want *Import
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewImport(tt.args.logger, tt.args.ctx, tt.args.repo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewImport() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildDisbursementRecordsFromImport(t *testing.T) {
	var orders = make([]*Order, 5)
	o1, _ := newOrder("e653f3e14bc4", "padberg_group", 10229, "2023-02-01")
	o2, _ := newOrder("20b674c93ea6", "padberg_group", 43321, "2023-02-01")
	o3, _ := newOrder("adaf77dffa91", "padberg_group", 724, "2023-02-02")
	o4, _ := newOrder("f1d9ec2b3d51", "rosenbaum_parisian", 8286, "2022-11-09")
	o5, _ := newOrder("858df04cb2b7", "rosenbaum_parisian", 5959, "2022-11-17")
	orders[0] = o1
	orders[1] = o2
	orders[2] = o3
	orders[3] = o4
	orders[4] = o5
	type args struct {
		o Orders
		m map[string]types.Merchant
	}
	tests := []struct {
		name    string
		args    args
		want    []types.Disbursement
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				o: orders,
				m: nil,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildDisbursementRecordsFromImport(tt.args.o, tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildDisbursementRecordsFromImport() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildDisbursementRecordsFromImport() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isNewPayoutPeriod(t *testing.T) {
	type args struct {
		o1 *Order
		o2 *Order
		m  types.Merchant
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := isNewPayoutPeriod(tt.args.o1, tt.args.o2, tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("isNewPayoutPeriod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("isNewPayoutPeriod() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sortOrdersByMerchant(t *testing.T) {
	type args struct {
		orders Orders
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortOrdersByMerchant(tt.args.orders)
		})
	}
}

func Test_buildWeeklyRecord(t *testing.T) {
	type args struct {
		o             Orders
		m             map[string]types.Merchant
		disbursements []types.Disbursement
		i             int
	}
	tests := []struct {
		name    string
		args    args
		want    []types.Disbursement
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildWeeklyRecord(tt.args.o, tt.args.m, tt.args.disbursements, tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildWeeklyRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildWeeklyRecord() got = %v, want %v", got, tt.want)
			}
		})
	}
}
