package deploy

import (
	"reflect"
	"testing"
)

func Test_parseHeadersFromKVPairs(t *testing.T) {
	type args struct {
		headersString string
	}
	tests := []struct {
		name    string
		args    args
		want    []header
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				headersString: "TESTVAL=example,SECONDVAL=two",
			},
			want: []header{
				{
					Key:   "TESTVAL",
					Value: "example",
				},
				{
					Key:   "SECONDVAL",
					Value: "two",
				},
			},
		},
		{
			name: "empty string",
			args: args{
				headersString: "",
			},
			want: nil,
		},
		{
			name: "with spaces",
			args: args{
				headersString: " TESTVAL=example, SECONDVAL=two ",
			},
			want: []header{
				{
					Key:   "TESTVAL",
					Value: "example",
				},
				{
					Key:   "SECONDVAL",
					Value: "two",
				},
			},
		},
		{
			name: "invalid format",
			args: args{
				headersString: "badformat",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseHeadersFromKVPairs(tt.args.headersString)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseHeadersFromKVPairs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseHeadersFromKVPairs() = %v, want %v", got, tt.want)
			}
		})
	}
}
