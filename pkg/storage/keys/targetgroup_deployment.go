package keys

import "fmt"

const TargetGroupDeploymentKey = "TARGET_GROUP_DEPLOYMENT#"

type targetGroupDeploymentKeys struct {
	PK1                string
	SK1                func(targetGroupDeploymentId string) string
	GSIPK1             func(targetGroupId string) string
	GSISK1             func(valid bool, healthy bool, priority int) string
	GSISK1ValidHealthy func(valid bool, healthy bool) string
}

var TargetGroupDeployment = targetGroupDeploymentKeys{
	PK1:    TargetGroupDeploymentKey,
	SK1:    func(targetGroupDeploymentId string) string { return targetGroupDeploymentId },
	GSIPK1: func(targetGroupId string) string { return targetGroupId },
	GSISK1: func(valid bool, healthy bool, priority int) string {
		return fmt.Sprintf("%v#%v#%d#", valid, healthy, priority)
	},
	GSISK1ValidHealthy: func(valid bool, healthy bool) string {
		return fmt.Sprintf("%v#%v#", valid, healthy)
	},
}
