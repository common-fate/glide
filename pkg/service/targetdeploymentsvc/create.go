package targetdeploymentsvc

import (
	"context"

	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/targetgroup"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

func (s *Service) CreateTargetGroupDeployment(ctx context.Context, req types.CreateTargetGroupDeploymentRequest) (*targetgroup.Deployment, error) {

	// run pre-lim checks to ensure input data is valid
	if !IsValidAwsAccountNumber(req.AwsAccount) {
		return nil, ErrInvalidAwsAccountNumber
	}

	// fetch existing deployment to ensure no overlap
	q := storage.GetTargetGroupDeployment{ID: req.Id}
	_, err := s.DB.Query(ctx, &q)
	// database error unrelated to no items
	if err != nil && err != ddb.ErrNoItems {
		return nil, err
	}
	// we've found a pre-existing deployment
	if err == nil {
		return nil, ErrTargetGroupDeploymentIdAlreadyExists
	}

	// create deployment
	dbInput := targetgroup.Deployment{
		ID:          req.Id,
		FunctionARN: req.FunctionArn,
		Runtime:     req.Runtime,
		AWSAccount:  req.AwsAccount,
		AwsRegion:   req.AwsRegion,
		Healthy:     false,
		Diagnostics: []targetgroup.Diagnostic{
			{
				Level:   string(types.ProviderSetupDiagnosticLogLevelINFO),
				Message: "offline: lambda cannot be reached/invoked",
			},
		},
	}

	err = s.DB.Put(ctx, &dbInput)
	if err != nil {
		return nil, err
	}

	return &dbInput, nil
}

func allDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func IsValidAwsAccountNumber(s string) bool {
	return len(s) == 12 && allDigits(s)
}
