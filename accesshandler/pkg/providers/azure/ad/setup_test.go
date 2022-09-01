package ad

import (
	"testing"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/psetup"
)

func TestSetup(t *testing.T) {
	p := Provider{}
	_, err := psetup.ParseDocsFS(p.SetupDocs(), p.Config(), psetup.TemplateData{})
	if err != nil {
		t.Fatal(err)
	}
}
