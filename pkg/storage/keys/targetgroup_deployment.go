package keys

const TargetGroupDeploymentKey = "TARGET_GROUP_DEPLOYMENT#"

type targetGroupDeploymentKeys struct {
	PK1 string
	SK1 func(targetGroupDeploymentId string) string
}

var TargetGroupDeployment = targetGroupDeploymentKeys{
	PK1: TargetGroupDeploymentKey,
	SK1: func(targetGroupDeploymentId string) string { return targetGroupDeploymentId },
}
