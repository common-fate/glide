package instructions

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
