package instructions

import "fmt"

type CodeBlock struct {
	Language string
	Code     string
}

// renders for the terminal, including colours
func (cb CodeBlock) RenderTerminal() string {
	return cb.Code
}

// renders without colours
func (cb CodeBlock) RenderString() string {
	return cb.Code
}

// render markdown
// as an example this can render a markdown codeblock with a language selector
func (cb CodeBlock) RenderMarkdown() string {
	return fmt.Sprintf("``` %s\n%s\n```", cb.Language, cb.Code)
}
