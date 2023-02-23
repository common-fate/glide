package targetsvc

import (
	"context"

	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/target"

	"github.com/common-fate/common-fate/pkg/types"
)

func (s *Service) CreateRoute(ctx context.Context, group string, req types.CreateTargetGroupLink) (*target.Route, error) {
	q := &storage.GetTargetGroup{
		ID: group,
	}
	_, err := s.DB.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	qh := &storage.GetHandler{
		ID: req.DeploymentId,
	}
	_, err = s.DB.Query(ctx, qh)
	if err != nil {
		return nil, err
	}
	route := target.Route{
		Group:   group,
		Handler: req.DeploymentId,
		// hardcoded mode until multi mode is supported
		Mode:     "Default",
		Priority: req.Priority,
		// invalid initially, healthcheck service will update this async
		Valid: false,
	}

	err = s.DB.Put(ctx, &route)
	if err != nil {
		return nil, err
	}
	return &route, nil
}
