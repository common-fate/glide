package storage

import (
	"context"
	"testing"

	"github.com/common-fate/common-fate/pkg/requestsv2.go"

	"github.com/stretchr/testify/assert"
)

func TestListEntitlementResources(t *testing.T) {

	type testcase struct {
		name string

		//filter options
		accessRules  []string
		filters      []string
		resourceName string

		insertBefore []requestsv2.ResourceOption
		want         []requestsv2.ResourceOption
		notWant      []requestsv2.ResourceOption
		wantErr      error
	}
	accessRule1 := "test"
	accessRule2 := "different"

	testData := CreateSeedData(accessRule1, accessRule2)

	testcases := []testcase{
		{

			name:         "get top level resources works with correct access rule",
			insertBefore: testData,
			accessRules:  []string{accessRule1},
			resourceName: "accountId",

			want:    CreateAccountIdOptions(accessRule1, accessRule2),
			notWant: CreatePermissionSetOptions(accessRule2),
		},
		{

			name:         "resource with multiple access rule relations only returns its own",
			insertBefore: testData,
			accessRules:  []string{accessRule2},
			resourceName: "accountId",

			want: []requestsv2.ResourceOption{
				requestsv2.ResourceOption{

					Label: "accountId",
					Value: "123456789012",
					Provider: requestsv2.TargetFrom{
						Kind:      "Account",
						Name:      "AWS",
						Publisher: "common-fate",
						Version:   "v0.1.0",
					},
					Type:        "Account",
					TargetGroup: "test",
					AccessRules: []string{
						accessRule1,
						accessRule2,
					},
				},
			},
			notWant: CreatePermissionSetOptions(accessRule2),
		},
		{

			name:         "get top level resources returns empty with wrong access rule",
			insertBefore: testData,
			accessRules:  []string{accessRule2},
			resourceName: "accountId",

			want:    []requestsv2.ResourceOption{},
			notWant: CreatePermissionSetOptions(accessRule2),
		},
		{

			name:         "get filtered options returns correct results",
			insertBefore: testData,
			accessRules:  []string{accessRule1},
			resourceName: "permissionSetArn",
			filters:      []string{"123456789012"},

			want: []requestsv2.ResourceOption{
				{
					Label: "permissionSetArn",
					Value: "123-abc",
					Provider: requestsv2.TargetFrom{
						Kind:      "Account",
						Name:      "AWS",
						Publisher: "common-fate",
						Version:   "v0.1.0",
					},
					Type:        "Account",
					TargetGroup: "test",
					AccessRules: []string{
						accessRule1,
					},
					RelatedTo: []string{"123456789012"},
				},
				{
					Label: "permissionSetArn",
					Value: "bar",
					Provider: requestsv2.TargetFrom{
						Kind:      "Account",
						Name:      "AWS",
						Publisher: "common-fate",
						Version:   "v0.1.0",
					},
					Type:        "Account",
					TargetGroup: "test",
					AccessRules: []string{
						accessRule1,
					},
					RelatedTo: []string{"123456789012"},
				},
			},
			notWant: CreateAccountIdOptions(accessRule1, accessRule2),
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			s := newTestingStorage(t)

			// insert any required fixture data
			for _, r := range tc.insertBefore {
				err := s.Put(ctx, &r)
				if err != nil {
					t.Fatal(err)
				}
			}

			q := ListEntitlementResources{
				Provider: requestsv2.TargetFrom{
					Kind:      "Account",
					Name:      "AWS",
					Publisher: "common-fate",
					Version:   "v0.1.0",
				},
				Argument:        tc.resourceName,
				UserAccessRules: tc.accessRules,
				FilterValues:    tc.filters,
			}
			_, err := s.Query(ctx, &q)
			if err != nil && tc.wantErr == nil {
				t.Fatal(err)
			}
			got := q.Result

			if tc.wantErr != nil {
				assert.Equal(t, tc.wantErr, err)
			}
			for _, item := range tc.want {
				assert.Contains(t, got, item)
			}
			for _, item := range tc.notWant {
				assert.NotContains(t, got, item, "expected item to not be in results")
			}
		})
	}
}

func CreateAccountIdOptions(accessRule1 string, accessRule2 string) []requestsv2.ResourceOption {
	opt1 := requestsv2.ResourceOption{

		Label: "accountId",
		Value: "123456789012",
		Provider: requestsv2.TargetFrom{
			Kind:      "Account",
			Name:      "AWS",
			Publisher: "common-fate",
			Version:   "v0.1.0",
		},
		Type:        "Account",
		TargetGroup: "test",
		AccessRules: []string{
			accessRule1,
			accessRule2,
		},
	}
	opt1a := requestsv2.ResourceOption{

		Label: "accountId",
		Value: "13579012345",
		Provider: requestsv2.TargetFrom{
			Kind:      "Account",
			Name:      "AWS",
			Publisher: "common-fate",
			Version:   "v0.1.0",
		},
		Type:        "Account",
		TargetGroup: "test",
		AccessRules: []string{
			accessRule1,
		},
	}
	opt1b := requestsv2.ResourceOption{

		Label: "accountId",
		Value: "583847583929",
		Provider: requestsv2.TargetFrom{
			Kind:      "Account",
			Name:      "AWS",
			Publisher: "common-fate",
			Version:   "v0.1.0",
		},
		Type:        "Account",
		TargetGroup: "test",
		AccessRules: []string{
			accessRule1,
		},
	}
	return []requestsv2.ResourceOption{opt1, opt1a, opt1b}
}

func CreatePermissionSetOptions(accessRule1 string) []requestsv2.ResourceOption {
	opt2 := requestsv2.ResourceOption{
		Label: "permissionSetArn",
		Value: "123-abc",
		Provider: requestsv2.TargetFrom{
			Kind:      "Account",
			Name:      "AWS",
			Publisher: "common-fate",
			Version:   "v0.1.0",
		},
		Type:        "Account",
		TargetGroup: "test",
		AccessRules: []string{
			accessRule1,
		},
		RelatedTo: []string{"123456789012"},
	}
	opt2a := requestsv2.ResourceOption{
		Label: "permissionSetArn",
		Value: "bar",
		Provider: requestsv2.TargetFrom{
			Kind:      "Account",
			Name:      "AWS",
			Publisher: "common-fate",
			Version:   "v0.1.0",
		},
		Type:        "Account",
		TargetGroup: "test",
		AccessRules: []string{
			accessRule1,
		},
		RelatedTo: []string{"123456789012"},
	}
	opt2b := requestsv2.ResourceOption{
		Label: "permissionSetArn",
		Value: "foo",
		Provider: requestsv2.TargetFrom{
			Kind:      "Account",
			Name:      "AWS",
			Publisher: "common-fate",
			Version:   "v0.1.0",
		},
		Type:        "Account",
		TargetGroup: "test",
		AccessRules: []string{
			accessRule1,
		},
		RelatedTo: []string{"different"},
	}

	return []requestsv2.ResourceOption{opt2, opt2a, opt2b}
}

func CreateSeedData(accessRule1 string, accessRule2 string) []requestsv2.ResourceOption {
	accountIds := CreateAccountIdOptions(accessRule1, accessRule2)
	permissionSets := CreatePermissionSetOptions(accessRule1)

	opt3 := requestsv2.ResourceOption{
		Label: "groupName",
		Value: "This is a okta group",
		Provider: requestsv2.TargetFrom{
			Kind:      "Group",
			Name:      "Okta",
			Publisher: "common-fate",
			Version:   "v0.1.0",
		},
		Type:        "Account",
		TargetGroup: "test",
		AccessRules: []string{
			accessRule2,
		},
	}

	res := []requestsv2.ResourceOption{}

	res = append(res, accountIds...)
	res = append(res, permissionSets...)
	res = append(res, opt3)

	return res
}
