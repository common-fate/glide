package gconfig

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecretStringValue(t *testing.T) {
	type testcase struct {
		name string
		give string
		want string
	}
	var emptySecret SecretStringValue
	emptyField := SecretString("test", &emptySecret, "testing")
	var nonEmptySecret SecretStringValue
	nonEmptyField := SecretString("test", &nonEmptySecret, "testing")
	nonEmptyField.Set("value")

	sBytes, err := json.Marshal(map[string]interface{}{"test": nonEmptySecret})
	if err != nil {
		t.Fatal(err)
	}
	fBytes, err := json.Marshal(map[string]interface{}{"test": nonEmptyField})
	if err != nil {
		t.Fatal(err)
	}
	// The following tests ensure that secrets stay secret in logs and prints
	testcases := []testcase{
		{name: "empty SecretStringValue.String() renders redacted value", give: emptySecret.String(), want: "*****"},
		{name: "empty SecretStringValue.Get() renders raw value", give: emptySecret.Get(), want: ""},
		{name: "empty Field.String() renders redacted value", give: emptyField.String(), want: "*****"},
		{name: "empty Field.Get() renders raw value", give: emptyField.Get(), want: ""},
		{name: "non empty SecretStringValue.String() renders redacted value", give: nonEmptySecret.String(), want: "*****"},
		{name: "non empty SecretStringValue.Get() renders raw value", give: nonEmptySecret.Get(), want: "value"},
		{name: "non empty Field.String() renders redacted value", give: nonEmptyField.String(), want: "*****"},
		{name: "non empty Field.Get() renders raw value", give: nonEmptyField.Get(), want: "value"},
		{name: "fmt.Sprint() on a SecretStringValue renders redacted value", give: fmt.Sprint(nonEmptySecret), want: "*****"},
		{name: "fmt.Sprint() on a Field renders redacted value", give: fmt.Sprint(nonEmptyField), want: "*****"},
		{name: "json.Marshal() on a SecretStringValue renders redacted value", give: string(sBytes), want: `{"test":"*****"}`},
		{name: "json.Marshal() on a Field renders nothing", give: string(fBytes), want: `{"test":{}}`},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.give)
		})
	}

}

func TestGeneralUsage(t *testing.T) {
	type testcase struct {
		name      string
		giveField *Field
		giveValue string
		wantField *Field
	}
	var secret SecretStringValue
	secretAfterSetting := SecretStringValue("some value")

	var value StringValue
	valueSetting := StringValue("some value")

	var optionalValue StringValue
	optionalValueSetting := StringValue("some value")

	// The following tests ensure that secrets stay secret in logs and prints
	testcases := []testcase{
		{name: "secretString", giveField: SecretString("test", &secret, "testing"), giveValue: "some value", wantField: &Field{key: "test", usage: "testing", value: &secretAfterSetting, secret: true, optional: false}},
		{name: "string", giveField: String("test", &value, "testing"), giveValue: "some value", wantField: &Field{key: "test", usage: "testing", value: &valueSetting, secret: false, optional: false}},
		{name: "optionalString", giveField: OptionalString("test", &optionalValue, "testing"), giveValue: "some value", wantField: &Field{key: "test", usage: "testing", value: &optionalValueSetting, secret: false, optional: true}},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert.NoError(t, tc.giveField.Set(tc.giveValue))
			assert.Equal(t, tc.wantField, tc.giveField)
			assert.Equal(t, tc.wantField.key, tc.giveField.Key())
			assert.Equal(t, tc.wantField.usage, tc.giveField.Usage())
			assert.Equal(t, tc.wantField.optional, tc.giveField.IsOptional())
			assert.Equal(t, tc.wantField.secret, tc.giveField.IsSecret())
			assert.Equal(t, tc.wantField.value.Get(), tc.giveField.Get())
		})
	}

}
