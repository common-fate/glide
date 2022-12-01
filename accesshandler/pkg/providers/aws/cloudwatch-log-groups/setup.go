package cloudwatchloggroups

import "embed"

//go:embed setup
var setupDocs embed.FS

// SetupDocs returns the embedded filesystem containing setup documentation.
func (p *Provider) SetupDocs() embed.FS {
	return setupDocs
}
