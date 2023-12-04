package okta

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOktaOrganizationURLValidation(t *testing.T) {
	type testcase struct {
		name      string
		giveURL   string
		wantError error
	}

	testcases := []testcase{

		{name: "ok", giveURL: "https://josh.okta.com"},
		{name: "http not allowed", giveURL: "http://josh.okta.com", wantError: errors.New("okta Organization URL must use https scheme")},
		{name: "non okta host not allowed", giveURL: "https://bad.hacker.com", wantError: errors.New("okta Organization URL must use the okta.com host. For security, if you use a custom domain for your Okta instance you need to configure the okta provider directly via the gdeploy CLI.")},
	}
	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {
			err := validateOktaURL(tc.giveURL)
			if tc.wantError == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.wantError.Error())
			}
		})
	}

}
