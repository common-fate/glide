package instructions

import "fmt"

type Instructions struct {
	Introduction []Block
	Steps        []Step
	Conclusion   []Block
}

// renders for the terminal, including colours
func (i Instructions) RenderTerminal() string {
	out := ""
	for _, block := range i.Introduction {
		out += fmt.Sprintln(block.RenderTerminal())
	}
	for _, step := range i.Steps {
		out += fmt.Sprintln(step.RenderTerminal())
	}
	for _, block := range i.Conclusion {
		out += fmt.Sprintln(block.RenderTerminal())
	}
	return out
}

// renders without colours
func (i Instructions) RenderString() string {
	out := ""
	for _, step := range i.Steps {
		out += fmt.Sprintln(step.RenderString())
	}
	return out
}

// render markdown
func (i Instructions) RenderMarkdown() string {
	out := ""
	for _, step := range i.Steps {
		out += fmt.Sprintln(step.RenderMarkdown())
	}
	return out
}

type Step struct {
	Title  string
	Blocks []Block
}

// renders for the terminal, including colours
func (s Step) RenderTerminal() string {
	out := s.Title
	if len(s.Blocks) > 0 {
		out += "\n"
	}
	for _, block := range s.Blocks {
		out += fmt.Sprintln(block.RenderTerminal())
	}
	return out
}

// renders without colours
func (s Step) RenderString() string {
	out := s.Title
	if len(s.Blocks) > 0 {
		out += "\n"
	}
	for _, block := range s.Blocks {
		out += fmt.Sprintln(block.RenderTerminal())
	}
	return out
}

// render markdown
func (s Step) RenderMarkdown() string {
	out := s.Title
	if len(s.Blocks) > 0 {
		out += "\n"
	}
	for _, block := range s.Blocks {
		out += fmt.Sprintln(block.RenderTerminal())
	}
	return out
}

type Block interface {
	RenderTerminal() string
	RenderString() string
	RenderMarkdown() string
}

type TextBlock string

// renders for the terminal, including colours
func (tb TextBlock) RenderTerminal() string {
	return string(tb)
}

// renders without colours
func (tb TextBlock) RenderString() string {
	return string(tb)
}

// render markdown
func (tb TextBlock) RenderMarkdown() string {
	return string(tb)
}

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
