package config

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// AWSSSMParamToken returns a token which can be fetched when using LoadAndReplaceSSMValues
func AWSSSMParamToken(path string, version string) string {
	return fmt.Sprintf("awsssm://%s:%s", path, version)
}

// LoadAndReplaceSSMValues will replace ssm descriptors with their real values from SSM
// If any values begin with 'awsssm://', it tries to look up the
// AWS SSM value.
// ssm parameters must have a version number specified by a trailing :versionnumber, e.g "awsssm:///granted/example:1"
// LoadAndReplaceSSMValues expects dst to be a non nil struct.
// Structs should be flat with all string values.
func LoadAndReplaceSSMValues(ctx context.Context, dst interface{}) error {
	if dst == nil {
		return fmt.Errorf("dst cannot be nil")
	}
	v := reflect.ValueOf(dst).Elem()
	// expect dst to be a struct
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("expected dst to be a struct, found: %s", v.Kind().String())
	}
	if !v.CanAddr() {
		return fmt.Errorf("cannot assign to the item passed, item must be a pointer in order to assign")
	}

	// use an errgroup so we can look up parameter values in parallel.
	g, gctx := errgroup.WithContext(ctx)

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	client := ssm.NewFromConfig(cfg)

	for i := 0; i < v.NumField(); i++ {
		f := v.Type().Field(i)
		if v.Field(i).Kind() != reflect.String {
			return fmt.Errorf("expected all struct fields to be of kind String, found: %s:%s", f.Name, v.Field(i).Kind().String())
		}
		value := v.Field(i).Interface().(string)
		if strings.HasPrefix(value, "awsssm://") {
			// important: we need to copy the i in this closure,
			// otherwise 'i' will change to the next loop iteration
			// while we're loading the SSM parameters.
			name := strings.TrimPrefix(value, "awsssm://")
			// check to ensure this parameter is versioned with a trailing ":versionnumber"
			parts := strings.Split(name, ":")
			if len(parts) != 2 {
				return fmt.Errorf("invalid format for awsssm parameter, expected versioned parameter. awsssm:///granted/example:1")
			}
			fieldIndex := i
			g.Go(func() error {
				output, err := client.GetParameter(gctx, &ssm.GetParameterInput{
					// This API fetched a version of a param if the path ends in ":versionnumber"
					Name:           aws.String(name),
					WithDecryption: true,
				})
				if err != nil {
					return errors.Wrapf(err, "looking up %s", name)
				}
				if output.Parameter.Value == nil {
					return fmt.Errorf("looking up %s: parameter value was nil", name)
				}
				// set the new value on the input struct
				v.Field(fieldIndex).SetString(aws.ToString(output.Parameter.Value))
				return nil
			})
		}
	}
	return g.Wait()
}
