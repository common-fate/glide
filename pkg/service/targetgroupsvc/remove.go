package targetgroupsvc

import (
	"context"

	"github.com/common-fate/common-fate/pkg/storage"
)

// Unlink a target group deployment from its target group
func (s *Service) RemoveTargetGroupLink(ctx context.Context, deploymentID string, groupID string) error {

	// Get the target group
	q := storage.GetTargetGroup{ID: groupID}
	_, err := s.DB.Query(ctx, &q)
	if err != nil {
		return err
	}

	// Get the target group deployment
	q2 := storage.GetTargetGroupDeployment{ID: deploymentID}
	_, err = s.DB.Query(ctx, &q2)
	if err != nil {
		return err
	}

	// Remove the TargetGroupAssignment
	q2.Result.TargetGroupAssignment = nil
	err = s.DB.Put(ctx, &q2.Result)
	if err != nil {
		return err
	}

	return nil
}
