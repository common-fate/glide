package providertest

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _  = runtime.Caller(0)
	basepath    = filepath.Dir(b)
	projectPath = filepath.Join(basepath, "../..")
)

func LoadFixture(ctx context.Context, name string, f interface{}) error {
	fixtureFile := filepath.Join(projectPath, "fixtures", fmt.Sprintf("%s.json", name))

	bytes, err := ioutil.ReadFile(fixtureFile)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, f)
}
