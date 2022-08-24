package instructions

type CodeBlock struct {
	Language string
	Code     string
}

// renders for the terminal, including colours
func (cb CodeBlock) RenderTerminal() string {
	return cb.Code
}
