package linker

// Helper package to clean up the links to the docs in the gdeploy cli
type Linker struct {
	BaseURL string
}

func (l *Linker) ReturnBaseURL() string {
	return l.BaseURL
}

func (l *Linker) MakeURL(path string) string {
	return l.BaseURL + "/" + path
}
