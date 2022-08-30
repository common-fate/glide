package ssov2

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArgSchema(t *testing.T) {
	p := Provider{}

	res := p.ArgSchema()
	out, err := json.Marshal(res)
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile("./testdata/argschema.json")
	if err != nil {
		t.Fatal(err)
	}
	buffer := new(bytes.Buffer)
	err = json.Compact(buffer, want)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, buffer.String(), string(out))
}
