package pdk

import (
	"context"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/pkg/targetgroup"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
)

// make of deploymentID to relative path
// example: ../../testvault-provider/provider
var LocalDeploymentMap map[string]string

func init() {
	LocalDeploymentMap = make(map[string]string)
}

type Data struct {
	ID    string                 `mapstructure:"id"`
	Name  string                 `mapstructure:"name"`
	Other map[string]interface{} `mapstructure:",remain"`
}

type Resource struct {
	Type string `mapstructure:"type"`
	Data Data   `mapstructure:"data"`
}

type LoadResourceResponse struct {
	Resources []Resource `mapstructure:"resources"`

	PendingTasks []struct {
		Name string      `mapstructure:"name"`
		Ctx  interface{} `mapstructure:"ctx"`
	} `mapstructure:"pendingTasks"`
}

type ProviderRuntime interface {
	FetchResources(ctx context.Context, name string, contx interface{}) (resources LoadResourceResponse, err error)
	Describe(ctx context.Context) (describeResponse *providerregistrysdk.DescribeResponse, err error)
	Grant(ctx context.Context, subject string, target Target) (err error)
	Revoke(ctx context.Context, subject string, target Target) (err error)
}

func GetRuntime(ctx context.Context, deployment targetgroup.Deployment) (ProviderRuntime, error) {
	log := logger.Get(ctx)
	var pr ProviderRuntime
	path, ok := LocalDeploymentMap[deployment.ID]
	if ok {
		log.Debugw("found local runtime configuration for deployment", "deployment", deployment, "path", path)
		pr = LocalRuntime{
			Path: path,
		}

	} else {
		log.Debugw("no local runtime configuration for deployment, using lambda runtime", "deployment", deployment)
		p, err := NewLambdaRuntime(ctx, deployment.FunctionARN())
		if err != nil {
			return nil, err
		}
		pr = p
	}
	return pr, nil
}
