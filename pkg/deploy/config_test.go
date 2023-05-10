package deploy

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

var exampleConfig = `
deployment:
  stackName: "test"
  account: "123456789012"
  region: "us-east-1"
  release: "v0.1.0"
  parameters:
    CognitoDomainPrefix: "a"
`

func TestParseConfig(t *testing.T) {
	var c Config
	err := yaml.Unmarshal([]byte(exampleConfig), &c)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "a", c.Deployment.Parameters.CognitoDomainPrefix)
}

func TestTestCfnParams(t *testing.T) {
	type testcase struct {
		name string
		give Config
		want string
	}

	testcases := []testcase{
		{
			name: "ok",
			give: Config{
				Deployment: Deployment{
					Parameters: Parameters{
						CognitoDomainPrefix: "test",
					},
				},
			},
			want: `[{"ParameterKey":"CognitoDomainPrefix","ParameterValue":"test","ResolvedValue":null,"UsePreviousValue":null}]`,
		},
		{
			name: "provider config",
			give: Config{
				Deployment: Deployment{
					Parameters: Parameters{},
				},
			},
			want: `[{"ParameterKey":"CognitoDomainPrefix","ParameterValue":"","ResolvedValue":null,"UsePreviousValue":null}]`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.give.CfnParams()
			if err != nil {
				t.Fatal(err)
			}
			gotJSON, err := json.Marshal(got)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tc.want, string(gotJSON))
		})
	}
}

func TestCfnTemplateURL(t *testing.T) {
	type testcase struct {
		name string
		give Config
		want string
	}

	testcases := []testcase{
		{
			name: "version tag",
			give: Config{
				Deployment: Deployment{
					Region:  "ap-southeast-2",
					Release: "v0.1.0",
				},
			},
			want: "https://granted-releases-ap-southeast-2.s3.amazonaws.com/v0.1.0/Granted.template.json",
		},
		{
			name: "custom URL",
			give: Config{
				Deployment: Deployment{
					Region:  "ap-southeast-2",
					Release: "https://custom-release.s3.amazonaws.com/template.json",
				},
			},
			want: "https://custom-release.s3.amazonaws.com/template.json",
		},
		{
			// note: this currently won't return an error, even though CloudFormation will refuse to deploy it.
			// this test captures this behaviour - in future we can add more validation around the URL.
			name: "custom URL not in S3",
			give: Config{
				Deployment: Deployment{
					Release: "https://some-other-url.com",
				},
			},
			want: "https://some-other-url.com",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.give.CfnTemplateURL()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestFeatureMap(t *testing.T) {
	// Tests that the Add method works as expected
	var p Parameters
	p.IdentityConfiguration.Upsert("test", map[string]string{})
	p.IdentityConfiguration.Upsert("test2", map[string]string{})
	p.IdentityConfiguration.Upsert("test", map[string]string{})

	// assert that the map is as expected
	assert.Equal(t, FeatureMap{"test": map[string]string{}, "test2": map[string]string{}}, p.IdentityConfiguration)
}
