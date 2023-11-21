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

func AssertEqualJson(t *testing.T, wantStr string, got interface{}) {
	gotJSON, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}

	var wantObj interface{}
	err = json.Unmarshal([]byte(wantStr), &wantObj)
	if err != nil {
		t.Fatal(err)
	}
	want, err := json.Marshal(wantObj)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, string(want), string(gotJSON))
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
			want: `
				[
					{
						"ParameterKey": "CognitoDomainPrefix",
						"ParameterValue": "test",
						"ResolvedValue": null,
						"UsePreviousValue": null
					}
				]
			`,
		},
		{
			name: "IdentityGroupFilter",
			give: Config{
				Deployment: Deployment{
					Parameters: Parameters{
						CognitoDomainPrefix: "test",
						IdentityGroupFilter: "test",
					},
				},
			},
			want: `
				[
					{
						"ParameterKey": "CognitoDomainPrefix",
						"ParameterValue": "test",
						"ResolvedValue": null,
						"UsePreviousValue": null
					},
					{
						"ParameterKey": "IdentityGroupFilter",
						"ParameterValue": "test",
						"ResolvedValue": null,
						"UsePreviousValue": null
					}
				]
			`,
		},
		{
			name: "provider config",
			give: Config{
				Deployment: Deployment{
					Parameters: Parameters{
						ProviderConfiguration: map[string]Provider{
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
			want: `
				[
					{
						"ParameterKey": "CognitoDomainPrefix",
						"ParameterValue": "",
						"ResolvedValue": null,
						"UsePreviousValue": null
					},
					{
						"ParameterKey": "ProviderConfiguration",
						"ParameterValue": "{\"okta\":{\"uses\":\"commonfate/okta@v1\",\"with\":{\"orgUrl\":\"test.internal\"}}}",
						"ResolvedValue": null,
						"UsePreviousValue": null
					}
				]
			`,
		},
	}

	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.give.CfnParams()
			if err != nil {
				t.Fatal(err)
			}
			AssertEqualJson(t, tc.want, got)
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

	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {
			got := tc.give.CfnTemplateURL()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestProviderMap(t *testing.T) {
	// Tests that the Add method works as expected
	var p Parameters
	err := p.ProviderConfiguration.Add("test", Provider{})
	assert.NoError(t, err)
	err = p.ProviderConfiguration.Add("test2", Provider{})
	assert.NoError(t, err)

	// Expect this to return an error
	err = p.ProviderConfiguration.Add("test", Provider{})
	assert.EqualError(t, err, "provider test already exists in the config")

	// assert that the map is as expected
	assert.Equal(t, ProviderMap{"test": Provider{}, "test2": Provider{}}, p.ProviderConfiguration)
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

func TestGetIDForNewProvider(t *testing.T) {
	type testcase struct {
		name   string
		giveID string
		give   Config
		want   string
	}

	testcases := []testcase{
		{
			name: "ok",
			give: Config{
				Deployment: Deployment{
					Parameters: Parameters{
						ProviderConfiguration: ProviderMap{},
					},
				},
			},
			giveID: "aws-sso",
			want:   "aws-sso",
		},
		{
			name:   "uninitialised map",
			give:   Config{},
			giveID: "aws-sso",
			want:   "aws-sso",
		},
		{
			name: "duplicate entry",
			give: Config{
				Deployment: Deployment{
					Parameters: Parameters{
						ProviderConfiguration: ProviderMap{
							"aws-sso": Provider{},
						},
					},
				},
			},
			giveID: "aws-sso",
			want:   "aws-sso-2",
		},
		{
			name: "multiple duplicate entries",
			give: Config{
				Deployment: Deployment{
					Parameters: Parameters{
						ProviderConfiguration: ProviderMap{
							"aws-sso":   Provider{},
							"aws-sso-2": Provider{},
							"aws-sso-3": Provider{},
						},
					},
				},
			},
			giveID: "aws-sso",
			want:   "aws-sso-4",
		},
		{
			name: "different ID",
			give: Config{
				Deployment: Deployment{
					Parameters: Parameters{
						ProviderConfiguration: ProviderMap{
							"aws-sso":   Provider{},
							"aws-sso-2": Provider{},
							"aws-sso-3": Provider{},
						},
					},
				},
			},
			giveID: "something-else",
			want:   "something-else",
		},
	}

	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {
			got := tc.give.Deployment.Parameters.ProviderConfiguration.GetIDForNewProvider(tc.giveID)
			assert.Equal(t, tc.want, got)
		})
	}
}
