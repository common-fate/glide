package psetup

import (
	"embed"
	"errors"
	"testing"

	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/ok
var okFS embed.FS

//go:embed testdata/invalid_template
var invalidTemplateFS embed.FS

func TestParseDocsFS(t *testing.T) {
	exampleField := gconfig.StringField("example", &gconfig.StringValue{}, "")

	cfg := gconfig.Config{
		exampleField,
	}

	type testcase struct {
		name    string
		want    []Step
		fs      embed.FS
		td      TemplateData
		cfg     gconfig.Config
		wantErr error
	}

	testcases := []testcase{
		{
			name: "ok",
			fs:   okFS,
			cfg:  cfg,
			td: TemplateData{
				AccessHandlerRoleARN: "<example-role-arn>",
			},
			want: []Step{
				{
					Title: "Test setup step",
					Instructions: `Test body

Data: <example-role-arn>

Test second line`,
					ConfigFields: []*gconfig.Field{exampleField},
				},
			},
		},
		{
			name: "invalid config value",
			fs:   okFS,
			cfg: gconfig.Config{
				gconfig.StringField("otherField", &gconfig.StringValue{}, ""),
			},
			wantErr: errors.New("parsing configFields in testdata/ok/example.md (do the configField values match the provider config?): field with key example not found"),
		},
		{
			name:    "invalid template variable",
			fs:      invalidTemplateFS,
			cfg:     cfg,
			wantErr: errors.New("rendering instructions template for testdata/invalid_template/invalid.md: template: instructions:1:9: executing \"instructions\" at <.VariableWhichDoesntExist>: can't evaluate field VariableWhichDoesntExist in type psetup.TemplateData"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseDocsFS(tc.fs, tc.cfg, tc.td)
			if err != nil && tc.wantErr == nil {
				t.Fatal(err)
			}
			if tc.wantErr != nil {
				assert.EqualError(t, err, tc.wantErr.Error())
			}
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestStepCountFS(t *testing.T) {
	type testcase struct {
		name    string
		want    int
		fs      embed.FS
		wantErr error
	}

	testcases := []testcase{
		{
			name: "ok",
			fs:   okFS,
			// only 1 step in the example folder
			want: 1,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := StepCountFS(tc.fs)
			if err != nil && tc.wantErr == nil {
				t.Fatal(err)
			}
			if tc.wantErr != nil {
				assert.EqualError(t, err, tc.wantErr.Error())
			}
			assert.Equal(t, tc.want, got)
		})
	}
}
