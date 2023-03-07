package target

import "github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"

// TestAccessRule returns an AccessRule fixture to be used in tests.
func TestGroup(opt ...func(*Group)) Group {

	ar := Group{
		ID:   "test-target-group",
		Icon: "aws-sso",
		TargetSchema: GroupTargetSchema{
			From:   "test/test/v1.1.1",
			Schema: providerregistrysdk.Target{Properties: map[string]providerregistrysdk.TargetField{}},
		},
	}

	for _, o := range opt {
		o(&ar)
	}

	return ar
}
