package main

import (
	"testing"

	"github.com/common-fate/granted-approvals/cmd/gdeploy/middleware"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/stretchr/testify/assert"
)

func TestIsReleaseVersionDifferent(t *testing.T) {

	type testcase struct {
		name           string
		gVersion       string
		dConfig        deploy.Deployment
		ignoreMismatch bool
		want           bool
		wantError      error
		wantIsGreater  bool
	}

	testCases := []testcase{
		{
			name:     "Ok",
			gVersion: "v2.10.11",
			want:     false,
			dConfig: deploy.Deployment{
				Release: "v2.10.11",
			},
		},
		{
			name:     "ignore check in dev",
			gVersion: "dev",
			dConfig: deploy.Deployment{
				Release: "httpgmail.com",
			},
			wantIsGreater: true,
			want:          false,
		},
		{
			name:     "Valid URL",
			gVersion: "0.1.1",
			dConfig: deploy.Deployment{
				Release: "https://gmail.com",
			},
			want:          false,
			wantIsGreater: true,
		},

		{
			name:     "gdeploy and common-fate version match",
			gVersion: "v1.02.02",
			dConfig: deploy.Deployment{
				Release: "v1.02.02",
			},
			want: false,
		},
		{
			name:     "gdeploy and common-fate version different",
			gVersion: "v1.02.02",
			dConfig: deploy.Deployment{
				Release: "v1.02.022",
			},
			want: true,
		},
		{
			name:     "gdeploy version missing v prefix",
			gVersion: "1.02.02",
			dConfig: deploy.Deployment{
				Release: "v1.02.02",
			},
			want: false,
		},
		{
			name:     "gdeploy version dev, release version number for UAT",
			gVersion: "dev",
			dConfig: deploy.Deployment{
				Release: "v1.02.02",
			},
			want:          false,
			wantIsGreater: true,
		},
		{
			name:     "ignore mismatch",
			gVersion: "0.1.1",
			dConfig: deploy.Deployment{
				Release: "v1.02.02",
			},
			ignoreMismatch: true,
			want:           false,
			wantIsGreater:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isDifferent, isGreater, err := middleware.IsReleaseVersionDifferent(tc.dConfig, tc.gVersion, tc.ignoreMismatch)
			if tc.wantError != nil {
				assert.EqualError(t, err, tc.wantError.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.want, isDifferent)
			assert.Equal(t, tc.wantIsGreater, isGreater)
		})
	}
}
