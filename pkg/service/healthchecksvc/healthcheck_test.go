package healthchecksvc

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheckValidity(t *testing.T) {

	type testcase struct {
		name       string
		schema1    map[string]providerregistrysdk.TargetArgument
		schema2    map[string]providerregistrysdk.TargetArgument
		valid_want bool
	}

	testcases := []testcase{
		{
			name: "identical-valid",
			schema1: map[string]providerregistrysdk.TargetArgument{
				"1": providerregistrysdk.TargetArgument{
					Id:           "abc",
					ResourceName: aws.String("abc"),
				},
			},
			schema2: map[string]providerregistrysdk.TargetArgument{
				"1": providerregistrysdk.TargetArgument{
					Id:           "abc",
					ResourceName: aws.String("abc"),
				},
			},
			valid_want: true,
		},
		{
			name: "different-invalid",
			schema1: map[string]providerregistrysdk.TargetArgument{
				"1": providerregistrysdk.TargetArgument{
					Id:           "123",
					ResourceName: aws.String("123"),
				},
			},
			schema2: map[string]providerregistrysdk.TargetArgument{
				"1": providerregistrysdk.TargetArgument{
					Id:           "abc",
					ResourceName: aws.String("abc"),
				},
			},
			valid_want: false,
		},
		{
			name: "different-invalid",
			schema1: map[string]providerregistrysdk.TargetArgument{
				"1": providerregistrysdk.TargetArgument{
					Id:           "123",
					ResourceName: aws.String("123"),
				},
			},
			schema2: map[string]providerregistrysdk.TargetArgument{
				"1": providerregistrysdk.TargetArgument{
					Id:           "abc",
					ResourceName: aws.String("abc"),
				},
			},
			valid_want: false,
		},
		{
			name: "resource-name-nil-valid",
			schema1: map[string]providerregistrysdk.TargetArgument{
				"1": providerregistrysdk.TargetArgument{
					Id:           "abc",
					ResourceName: nil,
				},
			},
			schema2: map[string]providerregistrysdk.TargetArgument{
				"1": providerregistrysdk.TargetArgument{
					Id:           "abc",
					ResourceName: nil,
				},
			},
			valid_want: true,
		},
	}

	for _, tc := range testcases {

		tc := tc

		t.Run(tc.name, func(t *testing.T) {

			s := Service{}

			validity := s.validateProviderSchema(tc.schema1, tc.schema2)

			assert.Equal(t, tc.valid_want, validity)

		})
	}
}
