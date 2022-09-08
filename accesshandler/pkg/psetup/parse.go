package psetup

import (
	"fmt"
	"io/fs"
	"strings"
	"text/template"

	"github.com/common-fate/frontmatter"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/pkg/errors"
)

// TemplateData are the particular values to be passed to the instructions template
type TemplateData struct {
	// AccessHandlerExecutionRoleARN is the ARN of the role that the Access Handler runs as
	AccessHandlerExecutionRoleARN string
	// WebhookURL is the Granted Approval's webhook URL
	WebhookURL string
}

type Step struct {
	Title        string
	Instructions string
	// ConfigFields which the user needs to enter as part of this step
	ConfigFields []*gconfig.Field
}

// FrontMatter is the frontmatter in the provider setup docs.
type FrontMatter struct {
	Title        string   `yaml:"title"`
	ConfigFields []string `yaml:"configFields"`
}

// StepCountFS reads the filesystem to determine how many steps there are in a provider setup process.
func StepCountFS(fsys fs.FS) (int, error) {
	var count int
	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, readErr error) error {
		if readErr != nil {
			return readErr
		}

		// skip directories, we just want to read markdown files.
		if d.IsDir() {
			return nil
		}

		count += 1
		return nil
	})
	return count, err
}

func ParseDocsFS(fsys fs.FS, cfg gconfig.Config, td TemplateData) ([]Step, error) {
	var steps []Step
	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, readErr error) error {
		if readErr != nil {
			return readErr
		}

		// skip directories, we just want to read markdown files.
		if d.IsDir() {
			return nil
		}

		f, err := fsys.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		var m FrontMatter
		body, err := frontmatter.MustParse(f, &m)
		if err != nil {
			return err
		}

		bodyStr := string(body)
		bodyStr = strings.TrimSpace(bodyStr)

		step := Step{
			Title: m.Title,
		}

		// find the relevant config fields for this setup step.
		for _, cf := range m.ConfigFields {
			field, err := cfg.FindFieldByKey(cf)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("parsing configFields in %s (do the configField values match the provider config?)", path))
			}
			step.ConfigFields = append(step.ConfigFields, field)
		}

		// render any template data into the body of the instructions
		t, err := template.New("instructions").Parse(bodyStr)
		if err != nil {
			return err
		}

		b := new(strings.Builder)
		err = t.Execute(b, td)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("rendering instructions template for %s", path))
		}
		step.Instructions = b.String()

		steps = append(steps, step)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return steps, nil
}
