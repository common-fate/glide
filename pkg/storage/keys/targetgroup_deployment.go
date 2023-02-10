package keys

const TargetGroupDeploymentKey = "TARGET_GROUP_DEPLOYMENT#"

type targetGroupDeploymentKeys struct {
	PK1                string
	SK1                func(targetGroupDeploymentId string) string
	GSIPK1             func(targetGroupId string) string
	GSISK1             func(validity string, health string, priority string) string
	GSIPK1ValidHealthy func(targetGroupId string) string
	GSISK1ValidHealthy string
}

var TargetGroupDeployment = targetGroupDeploymentKeys{
	PK1:    TargetGroupDeploymentKey,
	SK1:    func(targetGroupDeploymentId string) string { return targetGroupDeploymentId },
	GSIPK1: func(targetGroupId string) string { return targetGroupId },
	GSISK1: func(validity string, health string, priority string) string {
		return validity + "#" + health + "#" + priority
	},
	GSIPK1ValidHealthy: func(targetGroupId string) string { return targetGroupId },
	GSISK1ValidHealthy: "true#true#",
}
