package deploy

import "fmt"

type Release struct {
	ProductionReleasesBucket      string
	ProductionReleaseBucketPrefix string
}

// CDKContextArgs returns the CDK context arguments
// in the form "-c" "ArgName=ArgValue"
func (s Release) CDKContextArgs() []string {
	var args []string
	// pass context variables through as CLI arguments. This will eventually allow them to be
	// overridden in automated deployment workflows like in CI pipelines.
	args = append(args, "-c", fmt.Sprintf("productionReleasesBucket=%s", s.ProductionReleasesBucket))
	args = append(args, "-c", fmt.Sprintf("productionReleaseBucketPrefix=%s", s.ProductionReleaseBucketPrefix))
	return args
}
