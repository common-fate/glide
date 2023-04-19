package preflightsvc

import (
	"context"
	"reflect"
	"testing"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/common-fate/pkg/cache"
	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb/ddbmock"
)

func TestValidateNoDuplicates(t *testing.T) {

	tests := []struct {
		name    string
		args    types.CreatePreflightRequest
		wantErr bool
	}{
		{
			name: "No duplicates",
			args: types.CreatePreflightRequest{
				Targets: []string{"target1", "target2", "target3"},
			},
			wantErr: false,
		},
		{
			name: "Duplicate targets",
			args: types.CreatePreflightRequest{
				Targets: []string{"target1", "target2", "target1"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateNoDuplicates(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("ValidateNoDuplicates() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_ValidateAccessToAllTargets(t *testing.T) {
	clk := clock.NewMock()
	type args struct {
		user             identity.User
		preflightRequest types.CreatePreflightRequest
	}
	tests := []struct {
		name          string
		args          args
		mockGetTarget cache.Target
		want          []cache.Target
		wantErr       bool
	}{
		{
			name: "ok",
			args: args{
				user: identity.User{Groups: []string{"group_1"}},
				preflightRequest: types.CreatePreflightRequest{
					Targets: []string{"tg_1"},
				},
			},
			mockGetTarget: cache.Target{
				ID:     "tg_1",
				Groups: map[string]struct{}{"group_1": {}},
			},
			want: []cache.Target{{
				ID:     "tg_1",
				Groups: map[string]struct{}{"group_1": {}},
			}},
		},
		{
			name: "fail",
			args: args{
				user: identity.User{Groups: []string{"group_1"}},
				preflightRequest: types.CreatePreflightRequest{
					Targets: []string{"tg_1"},
				},
			},
			mockGetTarget: cache.Target{
				ID:     "tg_1",
				Groups: map[string]struct{}{"group_2": {}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := ddbmock.New(t)

			db.MockQuery(&storage.GetCachedTarget{
				Result: &tt.mockGetTarget,
			})
			s := &Service{
				DB:    db,
				Clock: clk,
			}
			got, err := s.ValidateAccessToAllTargets(context.Background(), tt.args.user, tt.args.preflightRequest)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.ValidateAccessToAllTargets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.ValidateAccessToAllTargets() = %v, want %v", got, tt.want)
			}
		})
	}
}
