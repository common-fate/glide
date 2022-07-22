package genv

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/joho/godotenv"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

func TestSSMLoader(t *testing.T) {
	_ = godotenv.Load()

	if os.Getenv("TESTING_SSM") == "" {
		t.Skip("TESTING_SSM env var not set")
	}

	ctx := context.Background()

	// random ID for testing
	id := ksuid.New().String()
	val := "ssmContents"

	key := "/granted/test/" + id
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		t.Fatal(err)
	}

	client := ssm.NewFromConfig(cfg)
	_, err = client.PutParameter(ctx, &ssm.PutParameterInput{
		Name:  &key,
		Value: &val,
		Type:  types.ParameterTypeSecureString,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		_, err = client.DeleteParameter(ctx, &ssm.DeleteParameterInput{
			Name: &key,
		})
		if err != nil {
			t.Logf("failed deleting SSM parameter: %s", err.Error())
		}
	})

	vals := map[string]string{
		"myTestValue":    "awsssm://" + key,
		"plainTextValue": "value",
	}
	valsJSON, err := json.Marshal(vals)
	if err != nil {
		t.Fatal(err)
	}

	l := SSMLoader{
		Data: valsJSON,
	}

	time.Sleep(time.Second * 2)

	got, err := l.Load(ctx)
	if err != nil {
		t.Fatal(err)
	}

	want := map[string]string{
		"myTestValue":    "ssmContents", // this one should have been resolved to the actual contents of the SSM parameter.
		"plainTextValue": "value",
	}
	assert.Equal(t, want, got)
}
