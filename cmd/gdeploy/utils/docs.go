package docs

type Linker struct {
	BaseURL string
}

func (l *Linker) ReturnBaseURL() string {
	return l.BaseURL
}

func (l *Linker) MakeURL(path string) string {
	return l.BaseURL + "/" + path
}
