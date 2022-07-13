package deploy

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSSMSecretsAreMaskedWhenLogged(t *testing.T) {
	okta := Okta{
		APIToken: "this-should-be-hidden",
	}
	google := Google{
		APIToken: "this-should-be-hidden",
	}
	slack := SlackConfig{
		APIToken: "this-should-be-hidden",
	}
	type testcase struct {
		name string
		give interface{}
		want string
	}
	testcases := []testcase{
		{name: "okta", give: okta, want: "{APIToken: ****, OrgURL: }"},
		{name: "oktaptr", give: &okta, want: "{APIToken: ****, OrgURL: }"},
		{name: "google", give: google, want: "{APIToken: ****, Domain: , AdminEmail }"},
		{name: "google", give: &google, want: "{APIToken: ****, Domain: , AdminEmail }"},
		{name: "slack", give: slack, want: "{APIToken: ****}"},
		{name: "slack", give: &slack, want: "{APIToken: ****}"},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, fmt.Sprintf("%s", tc.give))
		})
	}
}
