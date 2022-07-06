package genv

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// SSMLoader loads configuration from a serialized JSON payload.
// If any values begin with 'awsssm://', it tries to look up the
// AWS SSM value.
type SSMLoader struct {
	Data []byte
}

func (l SSMLoader) Load(ctx context.Context) (map[string]string, error) {
	var res map[string]string
	var mu sync.Mutex

	err := json.Unmarshal(l.Data, &res)
	if err != nil {
		return nil, err
	}

	// use an errgroup so we can look up parameter values in parallel.
	g, gctx := errgroup.WithContext(ctx)

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	client := ssm.NewFromConfig(cfg)

	for k, v := range res {
		if strings.HasPrefix(v, "awsssm://") {
			// important: we need to copy the key and value in this closure,
			// otherwise 'k' and 'v' will change to the next loop iteration
			// while we're loading the SSM parameters.
			name := strings.TrimPrefix(v, "awsssm://")
			key := k
			g.Go(func() error {
				output, err := client.GetParameter(gctx, &ssm.GetParameterInput{
					Name:           &name,
					WithDecryption: true,
				})
				if err != nil {
					return errors.Wrapf(err, "looking up %s", name)
				}
				if output.Parameter.Value == nil {
					return fmt.Errorf("looking up %s: parameter value was nil", name)
				}
				// lock the mutex to ensure we're safe to write to the map
				// without other Goroutines writing over us.
				mu.Lock()
				defer mu.Unlock()
				res[key] = *output.Parameter.Value
				return nil
			})
		}
	}

	err = g.Wait()
	if err != nil {
		return nil, err
	}

	return res, nil
}
