package targetdeploymentsvc

import (
	"context"

	"github.com/common-fate/common-fate/pkg/targetgroup"
	"github.com/common-fate/common-fate/pkg/types"
)

func (s *Service) CreateTargetGroupDeployment(ctx context.Context, req types.CreateTargetGroupDeploymentRequest) (*targetgroup.Deployment, error) {

	// TODO: run pre-lim checks to ensure aws account/arn are valid
	dbInput := targetgroup.Deployment{
		ID:           req.Id,
		FunctionARN:  req.FunctionArn,
		Runtime:      req.Runtime,
		AWSAccount:   req.AwsAccount,
		Healthy:      false,
		Diagnostics:  []targetgroup.Diagnostic{},
		ActiveConfig: map[string]targetgroup.Config{},
		Provider:     targetgroup.Provider{},
	}

	dbInput.Provider.Name = req.Provider.Name
	dbInput.Provider.Version = req.Provider.Version
	dbInput.Provider.Version = req.Provider.Version

	/**

	  TODO
	  - determine the specific spec for active config,
	  - what are the value input requirements,
	  - what are the value output requirements? i.e. what additional processing is needed

	  ...

	  Below is a rough idea of how to extract values taken from:
	  pkg/service/psetupsvc/create.go:89

	*/

	// initialise the config values if the provider supports it.
	// if configer, ok := b.ActiveConfig.(targetgroup.Config); ok {
	// 	for _, field := range configer.Config() {
	// 		ps.ConfigValues[field.Key()] = ""
	// 	}
	// }

	// @TODO: run a check here to ensure no overwrites occur ...

	err := s.DB.Put(ctx, &dbInput)
	if err != nil {
		return nil, err
	}

	return &dbInput, nil
}
