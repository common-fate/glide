package instructions

import "fmt"

type Instructions struct {
	Introduction []Block
	Steps        []Step
	Conclusion   []Block
}

// renders for the terminal
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
