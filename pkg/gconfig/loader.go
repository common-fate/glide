package gconfig

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/common-fate/granted-approvals/pkg/cfaws"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type SecretGetter interface {
	GetSecret(ctx context.Context, path string) (string, error)
}

// these are our secret backends
// only SSM for now
var secretGetterRegistry = map[string]SecretGetter{
	"awsssm://": SSMGetter{},
}

// MapLoader looks up values in it's Values map
// when loading configuration.
//
// It's useful for writing tests which use genv to configure things.
type MapLoader struct {
	Values map[string]string
}

// Under the hood, this just uses the json loader so we get all the SSM loading capability
func (l *MapLoader) Load(ctx context.Context) (map[string]string, error) {
	b, err := json.Marshal(l.Values)
	if err != nil {
		return nil, err
	}
	return JSONLoader{Data: b}.Load(ctx)
}

// JSONLoader loads configuration from a serialized JSON payload
// set in the 'Data' field.
// if any values are prefixed with one of teh known prefixes, there are further processed
// e.g values prefixed with "awsssm://" will be treated as an ssm parameter and will be fetched via the aws SDK
type JSONLoader struct {
	Data []byte
}

func (l JSONLoader) Load(ctx context.Context) (map[string]string, error) {
	var res map[string]string

	err := json.Unmarshal(l.Data, &res)
	if err != nil {
		return nil, err
	}
	var mu sync.Mutex

	// use an errgroup so we can look up parameter values in parallel.
	g, gctx := errgroup.WithContext(ctx)
	// After unmarshalling the json, check for any value which match a secret getter
	// if it does, get the secret value
	for k, v := range res {
		for getterKey, getter := range secretGetterRegistry {
			if strings.HasPrefix(v, getterKey) {
				// important: we need to copy the key and value in this closure,
				// otherwise 'k' and 'v' will change to the next loop iteration
				// while we're loading the value
				name := strings.TrimPrefix(v, getterKey)
				key := k
				g.Go(func() error {
					value, err := getter.GetSecret(gctx, name)
					if err != nil {
						return err
					}
					// lock the mutex to ensure we're safe to write to the map
					// without other Goroutines writing over us.
					mu.Lock()
					defer mu.Unlock()
					res[key] = value
					return nil
				})
				continue
			}
		}
	}
	err = g.Wait()
	if err != nil {
		return nil, err
	}
	return res, nil
}

type SSMGetter struct{}

func (g SSMGetter) GetSecret(ctx context.Context, path string) (string, error) {
	cfg, err := cfaws.ConfigFromContextOrDefault(ctx)
	if err != nil {
		return "", err
	}
	client := ssm.NewFromConfig(cfg)
	output, err := client.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           &path,
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "", errors.Wrapf(err, "looking up %s in ssm", path)
	}
	if output.Parameter.Value == nil {
		return "", fmt.Errorf("looking up %s in ssm: parameter value was nil", path)
	}
	return *output.Parameter.Value, nil
}
