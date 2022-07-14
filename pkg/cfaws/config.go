package cfaws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type contextkey struct{}

var AWSConfigContextKey contextkey

func ConfigFromContextOrDefault(ctx context.Context) (aws.Config, error) {
	if cfg := ctx.Value(AWSConfigContextKey); cfg != nil {
		return cfg.(aws.Config), nil
	} else {
		return config.LoadDefaultConfig(ctx)
	}
}

func SetConfigInContext(ctx context.Context, cfg aws.Config) context.Context {
	return context.WithValue(ctx, AWSConfigContextKey, cfg)
}
