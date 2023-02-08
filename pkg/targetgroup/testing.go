package targetgroup

import "github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"

// TestAccessRule returns an AccessRule fixture to be used in tests.
func TestTargetGroup(opt ...func(*TargetGroup)) TargetGroup {

	ar := TargetGroup{
		ID:   "test-target-group",
		Icon: "aws-sso",
		TargetSchema: GroupTargetSchema{
			From:   "test/test/v1.1.1",
			Schema: providerregistrysdk.TargetSchema{AdditionalProperties: map[string]providerregistrysdk.TargetArgument{}},
		},
		TargetDeployments: []DeploymentRegistration{},
	}

	for _, o := range opt {
		o(&ar)
	}

	return ar
}
