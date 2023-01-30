package pdk

import (
	"context"
	"encoding/json"
	"os/exec"

	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
)

type LocalRuntime struct {
	Path string
}

func (l LocalRuntime) Schema(ctx context.Context) (schema providerregistrysdk.ProviderSchema, err error) {
	cmd := exec.Command("pdk-cli", "test", "arg-schema")
	cmd.Dir = l.Path
	out, err := cmd.Output()
	if err != nil {
		return providerregistrysdk.ProviderSchema{}, err
	}
	err = json.Unmarshal(out, &schema.Target)
	if err != nil {
		return providerregistrysdk.ProviderSchema{}, err
	}
	return
}
