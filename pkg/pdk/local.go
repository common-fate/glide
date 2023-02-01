package pdk

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"

	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
)

type LocalRuntime struct {
	Path string
}

func (l LocalRuntime) Schema(ctx context.Context) (schema providerregistrysdk.ProviderSchema, err error) {
	cmd := exec.Command("pdk-cli", "test", "schema")
	cmd.Dir = l.Path
	cmd.Env = append(cmd.Env, os.Environ()...)
	out, err := cmd.Output()
	if err != nil {
		return providerregistrysdk.ProviderSchema{}, err
	}
	err = json.Unmarshal(out, &schema)
	if err != nil {
		return providerregistrysdk.ProviderSchema{}, err
	}
	return
}

func (l LocalRuntime) FetchResources(ctx context.Context, name string, contx interface{}) (resources LoadResourceResponse, err error) {
	b, err := json.Marshal(contx)
	if err != nil {
		return LoadResourceResponse{}, err
	}
	cmd := exec.Command("pdk-cli", "test", "fetch-resources", "--name="+name, "--ctx="+string(b))
	cmd.Dir = l.Path
	cmd.Env = append(cmd.Env, os.Environ()...)
	out, err := cmd.Output()
	if err != nil {
		return LoadResourceResponse{}, err
	}
	err = json.Unmarshal(out, &resources)
	if err != nil {
		return LoadResourceResponse{}, err
	}
	return
}
