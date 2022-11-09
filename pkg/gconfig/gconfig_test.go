package gconfig

import (
	"context"
	"encoding/json"
	"errors"
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
	emptyField := SecretStringField("test", &emptySecret, "testing", WithNoArgs(""))
	var nonEmptySecret SecretStringValue
	nonEmptyField := SecretStringField("test", &nonEmptySecret, "testing", WithNoArgs(""))
	err := nonEmptyField.Set("value")
	if err != nil {
		t.Fatal(err)
	}
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
	secretAfterSetting := SecretStringValue{"some value"}

	var value StringValue
	valueSetting := StringValue{"some value"}

	var optionalValue OptionalStringValue
	ov := "some value"
	optionalValueSetting := OptionalStringValue{Value: &ov}

	fn := WithNoArgs("granted/path")
	// The following tests ensure that secrets stay secret in logs and prints
	testcases := []testcase{
		{name: "secretString", giveField: SecretStringField("test", &secret, "testing", fn), giveValue: "some value", wantField: &Field{key: "test", description: "testing", value: &secretAfterSetting, secret: true, optional: false, secretPathFunc: fn, hasChanged: true, cliPrompt: CLIPromptTypePassword}},
		{name: "string", giveField: StringField("test", &value, "testing"), giveValue: "some value", wantField: &Field{key: "test", description: "testing", value: &valueSetting, secret: false, optional: false, hasChanged: true}},
		{name: "optionalString", giveField: OptionalStringField("test", &optionalValue, "testing"), giveValue: "some value", wantField: &Field{key: "test", description: "testing", value: &optionalValueSetting, secret: false, optional: true, hasChanged: true}},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			assert.NoError(t, tc.giveField.Set(tc.giveValue))
			assert.Equal(t, tc.wantField.secretPathFunc == nil, tc.giveField.secretPathFunc == nil)
			tc.giveField.secretPathFunc = nil
			tc.wantField.secretPathFunc = nil
			assert.Equal(t, tc.wantField, tc.giveField)
			assert.Equal(t, tc.wantField.key, tc.giveField.Key())
			assert.Equal(t, tc.wantField.description, tc.giveField.Description())
			assert.Equal(t, tc.wantField.optional, tc.giveField.IsOptional())
			assert.Equal(t, tc.wantField.secret, tc.giveField.IsSecret())
			assert.Equal(t, tc.wantField.value.Get(), tc.giveField.Get())
			assert.Equal(t, tc.wantField.hasChanged, tc.giveField.HasChanged())

			if f, ok := tc.giveField.value.(*OptionalStringValue); ok {
				assert.Equal(t, tc.giveField.Get() == "", !f.IsSet())

			}
		})
	}

}

func TestLoad(t *testing.T) {
	type testStruct struct {
		a StringValue
		b SecretStringValue
		// optional
		c OptionalStringValue
		// optional not present
		d OptionalStringValue
	}
	type testcase struct {
		name       string
		giveStruct *testStruct
		giveConfig Config
		giveLoader Loader
		wantStruct *testStruct
		wantError  error
	}

	var test1 testStruct
	oc := "testvaluec"
	c := OptionalStringValue{Value: &oc}
	test1Expected := testStruct{
		a: StringValue{"testvaluea"},
		b: SecretStringValue{"testvalueb"},
		c: c,
	}

	var test2 testStruct

	testcases := []testcase{
		{name: "loading config works as expected when values are non nil", giveStruct: &test1, giveConfig: Config{
			StringField("a", &test1.a, "usage"),
			SecretStringField("b", &test1.b, "usage", WithNoArgs("test-path")),
			OptionalStringField("c", &test1.c, "usage"),
			OptionalStringField("d", &test1.d, "usage"),
		}, giveLoader: &MapLoader{Values: map[string]string{
			"a": "testvaluea",
			"b": "testvalueb",
			"c": "testvaluec",
		}}, wantStruct: &test1Expected},
		{name: "not found in map returns error", giveStruct: &test2, giveConfig: Config{
			StringField("a", &test2.a, "usage"),
		}, giveLoader: &MapLoader{Values: map[string]string{}}, wantError: errors.New("could not find a in map")},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.giveConfig.Load(context.Background(), tc.giveLoader)
			if tc.wantError != nil {
				assert.EqualError(t, err, tc.wantError.Error())
			} else if err != nil {
				t.Fatal(err)
			} else {
				assert.Equal(t, tc.wantStruct, tc.giveStruct)
			}
		})
	}

}
func TestDump(t *testing.T) {
	type testcase struct {
		name       string
		giveConfig Config
		giveDumper Dumper
		wantMap    map[string]string
		wantError  error
	}
	a := StringValue{"testing"}
	b := SecretStringValue{"password"}
	testcases := []testcase{
		{name: "ok", giveConfig: Config{}, giveDumper: SafeDumper{}, wantMap: map[string]string{}, wantError: nil},
		{name: "with values, redacted secret", giveConfig: Config{StringField("a", &a, ""), SecretStringField("b", &b, "", WithNoArgs(""))}, giveDumper: SafeDumper{}, wantMap: map[string]string{"a": "testing", "b": "*****"}, wantError: nil},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.giveConfig.Dump(context.Background(), tc.giveDumper)
			if tc.wantError != nil {
				assert.EqualError(t, err, tc.wantError.Error())
			} else if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tc.wantMap, res)
		})
	}
}

func TestNilValuesPanic(t *testing.T) {
	type testcase struct {
		name      string
		callback  func()
		wantPanic error
	}
	// These tests cases test that making a field with a nil value causes a panic because it is not the supported useage
	testcases := []testcase{
		{name: "StringField with nil value panics", callback: func() { StringField("", nil, "") }, wantPanic: ErrFieldValueMustNotBeNil},
		{name: "SecretStringField with nil value panics", callback: func() { SecretStringField("", nil, "", WithNoArgs("")) }, wantPanic: ErrFieldValueMustNotBeNil},
		{name: "OptionalStringField with nil value panics", callback: func() { OptionalStringField("", nil, "") }, wantPanic: ErrFieldValueMustNotBeNil},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantPanic != nil {
				defer func() {
					err := recover()
					if err != tc.wantPanic {
						t.Fatalf("Wrong panic message: %s", err)
					}
				}()
			}
			tc.callback()
		})
	}

}
