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
