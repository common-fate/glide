package instructions

import "fmt"

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
