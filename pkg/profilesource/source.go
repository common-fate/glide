package profilesource

import (
	"context"

	"github.com/common-fate/awsconfigfile"
	"github.com/common-fate/common-fate/pkg/client"
)

// Source reads available AWS SSO profiles from the Common Fate API.
// It implements the awsconfigfile.Source interface
type Source struct {
	SSORegion    string
	StartURL     string
	Client       *client.Client
	DashboardURL string
}

func (s Source) GetProfiles(ctx context.Context) ([]awsconfigfile.SSOProfile, error) {
	cf := s.Client

	rules, err := cf.UserListAccessRulesWithResponse(ctx)
	if err != nil {
		return nil, err
	}
	var profiles []awsconfigfile.SSOProfile

	for _, r := range rules.JSON200.AccessRules {
		ruleDetail, err := cf.UserGetAccessRuleWithResponse(ctx, r.ID)
		if err != nil {
			return nil, err
		}

		// only rules for the aws-sso Access Provider are relevant here
		if ruleDetail.JSON200.Target.Provider.Type != "aws-sso" {
			continue
		}

		accountId := ruleDetail.JSON200.Target.Arguments.AdditionalProperties["accountId"]
		permissionSetArn := ruleDetail.JSON200.Target.Arguments.AdditionalProperties["permissionSetArn"]

		// add all options to our profile map
		for _, acc := range accountId.Options {
			for _, ps := range permissionSetArn.Options {
				p := awsconfigfile.SSOProfile{
					AccountID:     acc.Value,
					AccountName:   acc.Label,
					RoleName:      ps.Label,
					SSOStartURL:   s.StartURL,
					SSORegion:     s.SSORegion,
					GeneratedFrom: "commonfate",
					CommonFateURL: s.DashboardURL,
				}
				profiles = append(profiles, p)
			}
		}
	}
	return profiles, nil
}
