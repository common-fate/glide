package main

import (
	"fmt"
	"testing"

	"github.com/common-fate/granted-approvals/internal/build"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/stretchr/testify/assert"
)

func TestCheckReleaseVersion(t *testing.T) {

	type testcase struct {
		name     string
		gVersion string
		dConfig  deploy.Deployment
		want     string
	}

	testCases := []testcase{{
		name:     "Ok",
		gVersion: build.Version,
		want:     "",
		dConfig: deploy.Deployment{
			Release: build.Version,
		},
	},
		{
			name:     "Invalid URL",
			gVersion: build.Version,
			dConfig: deploy.Deployment{
				Release: "httpgmail.com",
			},
			want: fmt.Sprintf("Incompatible gdeploy version. Expected %s got %s . ", "httpgmail.com", build.Version)},

		{
			name:     "Valid URL",
			gVersion: build.Version,
			dConfig: deploy.Deployment{
				Release: "https://gmail.com",
			},
			want: ""},

		{
			name:     "gdeploy and granted-approval version match",
			gVersion: "v1.02.02",
			dConfig: deploy.Deployment{
				Release: "v1.02.02",
			},
			want: ""},

		{
			name:     "gdeploy and granted-approval version different",
			gVersion: "v1.02.02",
			dConfig: deploy.Deployment{
				Release: "v1.02.022",
			},
			want: fmt.Sprintf("Incompatible gdeploy version. Expected %s got %s . ", "v1.02.022", "v1.02.02")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			var got string
			res := CheckReleaseVersion(tc.dConfig, tc.gVersion)

			if res == nil {
				got = ""
			} else {
				got = res.Error()
			}

			assert.Equal(t, tc.want, got)
		})
	}
}
