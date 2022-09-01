package main

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/common-fate/granted-approvals/cmd/gdeploy/middleware"
	"github.com/common-fate/granted-approvals/internal/build"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/stretchr/testify/assert"
)

func TestIsReleaseVersionDifferent(t *testing.T) {

	type testcase struct {
		name     string
		gVersion string
		dConfig  deploy.Deployment
		want     string
	}

	testCases := []testcase{{
		name:     "Ok",
		gVersion: "v2.10.11",
		want:     "false",
		dConfig: deploy.Deployment{
			Release: "v2.10.11",
		},
	},
		{
			name:     "Invalid URL",
			gVersion: build.Version,
			dConfig: deploy.Deployment{
				Release: "httpgmail.com",
			},
			want: fmt.Sprintf("invalid URL. Please update your release version in 'granted-deployment.yml' to %s", build.Version),
		},

		{
			name:     "Valid URL",
			gVersion: build.Version,
			dConfig: deploy.Deployment{
				Release: "https://gmail.com",
			},
			want: "false",
		},

		{
			name:     "gdeploy and granted-approval version match",
			gVersion: "v1.02.02",
			dConfig: deploy.Deployment{
				Release: "v1.02.02",
			},
			want: "false",
		},

		{
			name:     "gdeploy and granted-approval version different",
			gVersion: "v1.02.02",
			dConfig: deploy.Deployment{
				Release: "v1.02.022",
			},
			want: "true",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			var got string
			isDifferent, err := middleware.IsReleaseVersionDifferent(tc.dConfig, tc.gVersion)

			if err != nil {
				got = err.Error()
			} else {
				got = strconv.FormatBool(isDifferent)
			}

			assert.Equal(t, tc.want, got)
		})
	}
}
