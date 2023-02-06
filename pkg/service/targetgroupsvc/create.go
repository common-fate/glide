package targetgroupsvc

import (
	"context"

	"github.com/common-fate/common-fate/pkg/targetgroup"
	"github.com/common-fate/common-fate/pkg/types"
	"go.uber.org/zap"
)

func (s *Service) CreateTargetGroup(ctx context.Context, req types.CreateTargetGroupRequest) (*targetgroup.TargetGroup, error) {
	log := zap.S()

	//look up target schema for the provider version

	group := targetgroup.TargetGroup{ID: req.ID}
	//based on the target schema provider type set the Icon

	log.Debugw("saving target group", "group", group)
	// save the request.
	err := s.DB.Put(ctx, &group)
	if err != nil {
		return nil, err
	}
	return &group, nil
}
