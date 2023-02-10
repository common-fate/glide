package pdk

import (
	"context"
	"strings"

	"github.com/common-fate/common-fate/pkg/targetgroup"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
)

// uselocal enables development mode using alocal cli instead of calling out to deployed lambdas
// set this to true to enable local handler
var UseLocal bool

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
	Schema(ctx context.Context) (schema providerregistrysdk.ProviderSchema, err error)
	FetchResources(ctx context.Context, name string, contx interface{}) (resources LoadResourceResponse, err error)
	Describe(ctx context.Context) (describeResponse targetgroup.ProviderDescribe, err error)
}

func GetRuntime(ctx context.Context, arn string) (ProviderRuntime, error) {
	var pr ProviderRuntime
	if UseLocal {
		// bit of a hack to get the local path in here
		path := strings.TrimPrefix(arn, "arn:aws:lambda")
		pr = LocalRuntime{
			Path: path,
		}
	} else {
		p, err := NewLambdaRuntime(ctx, arn)
		if err != nil {
			return nil, err
		}
		pr = p
	}
	return pr, nil
}
