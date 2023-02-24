package identitysync

import (
	"reflect"
	"testing"

	"github.com/common-fate/common-fate/pkg/identity"
)

func Test_filterGroups(t *testing.T) {
	type args struct {
		groups       []identity.IDPGroup
		filterString string
	}
	tests := []struct {
		name    string
		args    args
		want    []identity.IDPGroup
		wantErr bool
	}{
		{
			name: "regex starts with: single match",
			args: args{
				groups: []identity.IDPGroup{
					{Name: "foo"},
					{Name: "bar"},
				},
				filterString: "^foo",
			},
			want: []identity.IDPGroup{
				{Name: "foo"},
			},
			wantErr: false,
		},
		{
			name: "regex ends with: double match",
			args: args{
				groups: []identity.IDPGroup{
					{Name: "foo"},
					{Name: "oooooo"},
					{Name: "bar"},
				},
				filterString: "oo$",
			},
			want: []identity.IDPGroup{
				{Name: "foo"},
				{Name: "oooooo"},
			},
			wantErr: false,
		},
		{
			name: "regex starts with: empty reponse",
			args: args{
				groups: []identity.IDPGroup{
					{Name: "bar"},
				},
				filterString: "^foo",
			},
			want:    []identity.IDPGroup{},
			wantErr: false,
		},
		{
			name: "invalid regex throws err",
			args: args{
				groups: []identity.IDPGroup{
					{Name: "bar"},
				},
				filterString: "^(foo",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FilterGroups(tt.args.groups, tt.args.filterString)
			if (err != nil) != tt.wantErr {
				t.Errorf("filterGroups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterGroups() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func TestOneLoginSync_listGroupsWithFilter(t *testing.T) {
// 	type fields struct {
// 		clientID       gconfig.StringValue
// 		clientSecret   gconfig.SecretStringValue
// 		token          gconfig.SecretStringValue
// 		idpGroupFilter gconfig.StringValue
// 		baseURL        gconfig.StringValue
// 	}
// 	type args struct {
// 		ctx          context.Context
// 		filterString *string
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    []identity.IDPGroup
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := &OneLoginSync{
// 				clientID:     tt.fields.clientID,
// 				clientSecret: tt.fields.clientSecret,
// 				token:        tt.fields.token,
// 				// idpGroupFilter: tt.fields.idpGroupFilter,
// 				baseURL: tt.fields.baseURL,
// 			}

// 			// @TODO: figure out how we can:
// 			// mock the response from GET "/api/1/roles" http request for a http.NewRequest method

// 			got, err := s.listGroupsWithFilter(tt.args.ctx, tt.args.filterString)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("OneLoginSync.listGroupsWithFilter() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("OneLoginSync.listGroupsWithFilter() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
