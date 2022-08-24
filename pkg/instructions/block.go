package instructions

type Block interface {
	RenderTerminal() string
	RenderString() string
	RenderMarkdown() string
}
