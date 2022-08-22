package flask

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/stretchr/testify/assert"
)

func TestArgSchema(t *testing.T) {
	p := Provider{}
	ctx := context.Background()

	err := p.Config().Load(ctx, &gconfig.MapLoader{
		Values: map[string]interface{}{
			"schema": map[string]interface{}{
				"properties": map[string]interface{}{
					"role": map[string]interface{}{
						"type":        "string",
						"title":       "Role",
						"description": "The Kubernetes Role to grant access to",
					},
				},
				"type":     "object",
				"required": []string{"role"},
			},
			"options": map[string]interface{}{
				"role": []types.Option{
					{
						Label: "admin",
						Value: "admin",
					},
				},
			},
			"type": "eks",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	res := p.ArgSchema()
	out, err := json.Marshal(res)
	if err != nil {
		t.Fatal(err)
	}
	want := `
	{
		"$schema": "http://json-schema.org/draft/2020-12/schema",
		"$id": "https://commonfate.io/demo/eks/args",
		"$ref": "#/$defs/Args",
		"$defs": {
		  "Args": {
			"properties": {
			  "role": {
				"description": "The Kubernetes Role to grant access to",
				"title": "Role",
				"type": "string"
			  }
			},
			"type": "object",
			"required": ["role"]
		  }
		}
	  }	  
	`
	if err != nil {
		t.Fatal(err)
	}
	buffer := new(bytes.Buffer)
	err = json.Compact(buffer, []byte(want))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, buffer.String(), string(out))
}
