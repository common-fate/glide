package identitysync

import (
	"context"
	"reflect"
	"testing"

	"github.com/common-fate/common-fate/pkg/gconfig"
	"github.com/common-fate/common-fate/pkg/identity"
)

func TestAzureSync_idpUserFromAzureUser(t *testing.T) {
	mail := "mail"
	invalidIdentifer := "invalid"

	type fields struct {
		token           gconfig.SecretStringValue
		tenantID        gconfig.StringValue
		clientID        gconfig.StringValue
		clientSecret    gconfig.SecretStringValue
		emailIdentifier gconfig.OptionalStringValue
	}
	type args struct {
		ctx       context.Context
		azureUser map[string]interface{}
		groups    []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    identity.IDPUser
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				ctx: context.Background(),
				azureUser: map[string]interface{}{
					"businessPhones":    []string{"425-555-0100"},
					"displayName":       "MOD Administrator",
					"givenName":         "MOD",
					"jobTitle":          nil,
					"mail":              nil,
					"mobilePhone":       "425-555-0101",
					"officeLocation":    nil,
					"preferredLanguage": "en-US",
					"surname":           "Administrator",
					"userPrincipalName": "admin@contoso.com",
					"id":                "4562bcc8-c436-4f95-b7c0-4f8ce89dca5e",
				},
			},
			want: identity.IDPUser{
				ID:        "4562bcc8-c436-4f95-b7c0-4f8ce89dca5e",
				FirstName: "MOD",
				LastName:  "Administrator",
				Email:     "admin@contoso.com",
			},
		},
		{
			name: "assign groups",
			args: args{
				ctx: context.Background(),
				azureUser: map[string]interface{}{
					"givenName":         "MOD",
					"surname":           "Administrator",
					"userPrincipalName": "admin@contoso.com",
					"id":                "4562bcc8-c436-4f95-b7c0-4f8ce89dca5e",
				},
				groups: []string{"a", "b"},
			},
			want: identity.IDPUser{
				ID:        "4562bcc8-c436-4f95-b7c0-4f8ce89dca5e",
				FirstName: "MOD",
				LastName:  "Administrator",
				Email:     "admin@contoso.com",
				Groups:    []string{"a", "b"},
			},
		},
		{
			name: "handle custom email identifier",
			fields: fields{
				// use the 'mail' attribute from the Azure AD API response.
				emailIdentifier: gconfig.OptionalStringValue{Value: &mail},
			},
			args: args{
				ctx: context.Background(),
				azureUser: map[string]interface{}{
					"givenName": "MOD",
					"surname":   "Administrator",
					// in this case, the UPN is an external user.
					"userPrincipalName": "admin#EXT@contoso.com",
					"mail":              "admin@contoso.com",
					"id":                "4562bcc8-c436-4f95-b7c0-4f8ce89dca5e",
				},
				groups: []string{"a", "b"},
			},
			want: identity.IDPUser{
				ID:        "4562bcc8-c436-4f95-b7c0-4f8ce89dca5e",
				FirstName: "MOD",
				LastName:  "Administrator",
				Email:     "admin@contoso.com",
				Groups:    []string{"a", "b"},
			},
		},
		{
			name: "invalid email identifier",
			fields: fields{
				emailIdentifier: gconfig.OptionalStringValue{Value: &invalidIdentifer},
			},
			args: args{
				ctx: context.Background(),
				azureUser: map[string]interface{}{
					"givenName":         "MOD",
					"surname":           "Administrator",
					"userPrincipalName": "admin#EXT@contoso.com",
					"mail":              "admin@contoso.com",
					"id":                "4562bcc8-c436-4f95-b7c0-4f8ce89dca5e",
				},
				groups: []string{"a", "b"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AzureSync{
				token:           tt.fields.token,
				tenantID:        tt.fields.tenantID,
				clientID:        tt.fields.clientID,
				clientSecret:    tt.fields.clientSecret,
				emailIdentifier: tt.fields.emailIdentifier,
			}
			got, err := a.idpUserFromAzureUser(tt.args.ctx, tt.args.azureUser, tt.args.groups)
			if (err != nil) != tt.wantErr {
				t.Errorf("AzureSync.idpUserFromAzureUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AzureSync.idpUserFromAzureUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
