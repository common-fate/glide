package instructions

type TextBlock string

// renders for the terminal, including colours
func (tb TextBlock) RenderTerminal() string {
	return string(tb)
}
