package pdk

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"

	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
	"github.com/mitchellh/mapstructure"
)

type LocalRuntime struct {
	Path string
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

	var data map[string]interface{}
	err = json.Unmarshal(out, &data)
	if err != nil {
		return LoadResourceResponse{}, err
	}

	err = mapstructure.Decode(data, &resources)
	if err != nil {
		return LoadResourceResponse{}, err
	}
	return
}

func (l LocalRuntime) Describe(ctx context.Context) (*providerregistrysdk.DescribeResponse, error) {
	cmd := exec.Command("pdk-cli", "test", "describe")
	cmd.Dir = l.Path
	cmd.Env = append(cmd.Env, os.Environ()...)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var describe providerregistrysdk.DescribeResponse
	err = json.Unmarshal(out, &describe)
	if err != nil {
		return nil, err
	}

	return &describe, nil
}
func (l LocalRuntime) Grant(ctx context.Context, subject string, target Target) (err error) {
	// @TODO this is untested/ not implemented in the local CLI
	cmd := exec.Command("pdk-cli", "test", "grant")
	cmd.Dir = l.Path
	cmd.Env = append(cmd.Env, os.Environ()...)
	_, err = cmd.Output()
	return err

}
func (l LocalRuntime) Revoke(ctx context.Context, subject string, target Target) (err error) {
	// @TODO this is untested/ not implemented in the local CLI
	cmd := exec.Command("pdk-cli", "test", "revoke")
	cmd.Dir = l.Path
	cmd.Env = append(cmd.Env, os.Environ()...)
	_, err = cmd.Output()
	return err

}
