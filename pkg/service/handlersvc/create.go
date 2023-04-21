package handlersvc

import (
	"context"

	"github.com/common-fate/common-fate/pkg/handler"
	"github.com/common-fate/common-fate/pkg/storage"

	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

// Input data validation is handled by the API layer
func (s *Service) RegisterHandler(ctx context.Context, req types.RegisterHandlerRequest) (*handler.Handler, error) {
	// fetch existing deployment to ensure no overlap
	q := storage.GetHandler{ID: req.Id}
	_, err := s.DB.Query(ctx, &q)
	// database error unrelated to no items
	if err != nil && err != ddb.ErrNoItems {
		return nil, err
	}
	// we've found a pre-existing deployment
	if err == nil {
		return nil, ErrHandlerIdAlreadyExists
	}

	// create deployment
	dbInput := handler.Handler{
		ID:         req.Id,
		Runtime:    string(req.Runtime),
		AWSAccount: req.AwsAccount,
		AWSRegion:  req.AwsRegion,
		Healthy:    false,
		Diagnostics: []handler.Diagnostic{
			{
				Level:   types.LogLevelINFO,
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
