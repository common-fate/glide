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
    CognitoDomainPrefix: ""
	ProviderConfiguration:
	  okta:
        uses: "commonfate/okta@v1"
        with:
          orgUrl: "https://test.internal"
          apiToken: "awsssm:///granted/okta/apiToken"
`

func TestParseConfig(t *testing.T) {
	var c Config
	err := yaml.Unmarshal([]byte(exampleConfig), &c)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "commonfate/okta@v1", c.Deployment.Parameters.ProviderConfiguration["okta"].Uses)
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
					Parameters: Parameters{
						ProviderConfiguration: map[string]Feature{
							"okta": {
								Uses: "commonfate/okta@v1",
								With: map[string]string{
									"orgUrl": "test.internal",
								},
							},
						},
					},
				},
			},
			want: `[{"ParameterKey":"CognitoDomainPrefix","ParameterValue":"","ResolvedValue":null,"UsePreviousValue":null},{"ParameterKey":"ProviderConfiguration","ParameterValue":"{\"okta\":{\"uses\":\"commonfate/okta@v1\",\"with\":{\"orgUrl\":\"test.internal\"}}}","ResolvedValue":null,"UsePreviousValue":null}]`,
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

func TestFeatures(t *testing.T) {
	var p Parameters
	f1 := Feature{Uses: "hello", With: map[string]string{"hello": "world"}}
	err := p.IdentityConfiguration.Upsert(f1)
	assert.NoError(t, err)
	f, err := p.IdentityConfiguration.Get("hello")
	assert.NoError(t, err)
	assert.Equal(t, &f1, f)
	f1Updated := Feature{Uses: "hello", With: map[string]string{"hello": "you"}}
	err = p.IdentityConfiguration.Upsert(f1Updated)
	assert.NoError(t, err)
	f, err = p.IdentityConfiguration.Get("hello")
	assert.NoError(t, err)
	assert.Equal(t, &f1Updated, f)

	// Assert that f1 was not mutated unexpectedly but updating the array
	assert.NotEqual(t, f1, f)

}
func TestFeatureMap(t *testing.T) {
	// Tests that the Add method works as expected
	var p Parameters
	err := p.ProviderConfiguration.Add("test", Feature{})
	assert.NoError(t, err)
	err = p.ProviderConfiguration.Add("test2", Feature{})
	assert.NoError(t, err)

	// Expect this to return an error
	err = p.ProviderConfiguration.Add("test", Feature{})
	assert.EqualError(t, err, "provider test already exists in the config")

	// assert that the map is as expected
	assert.Equal(t, FeatureMap{"test": Feature{}, "test2": Feature{}}, p.ProviderConfiguration)
}
