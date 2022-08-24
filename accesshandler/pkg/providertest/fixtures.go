package providertest

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _  = runtime.Caller(0)
	basepath    = filepath.Dir(b)
	projectPath = filepath.Join(basepath, "../..")
)

func LoadFixture(ctx context.Context, name string, f interface{}) error {
	switch name {
	case "azure":
		projectPath = filepath.Join(basepath, "../../..")
	case "okta":
		projectPath = filepath.Join(basepath, "../..")
	default:
		projectPath = filepath.Join(basepath, "../..")
	}

	fixtureFile := filepath.Join(projectPath, "fixtures", fmt.Sprintf("%s.json", name))

	bytes, err := os.ReadFile(fixtureFile)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, f)
}
