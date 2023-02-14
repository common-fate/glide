package pdk

import (
	"context"
	"fmt"

	"github.com/common-fate/common-fate/pkg/targetgroup"
)

// uselocal enables development mode using alocal cli instead of calling out to deployed lambdas
// set this to true to enable local handler
var UseLocal bool

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
	Describe(ctx context.Context) (describeResponse targetgroup.ProviderDescribe, err error)
	Grant(ctx context.Context, subject string, target Target) (err error)
	Revoke(ctx context.Context, subject string, target Target) (err error)
}

func GetRuntime(ctx context.Context, deployment targetgroup.Deployment) (ProviderRuntime, error) {
	var pr ProviderRuntime
	if UseLocal {
		path, ok := LocalDeploymentMap[deployment.ID]
		if !ok {
			return nil, fmt.Errorf("local runtime path not configured for deployment %s", deployment.ID)
		}
		pr = LocalRuntime{
			Path: path,
		}
	} else {
		p, err := NewLambdaRuntime(ctx, deployment.FunctionARN)
		if err != nil {
			return nil, err
		}
		pr = p
	}
	return pr, nil
}
