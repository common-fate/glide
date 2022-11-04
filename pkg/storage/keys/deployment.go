package keys

const DeploymentKey = "CF_DEPLOYMENT#"

type deploymentKeys struct {
	PK1 string
	SK1 string
}

var Deployment = deploymentKeys{
	PK1: DeploymentKey,
	SK1: DeploymentKey,
}
