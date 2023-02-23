package handler

func TestHandler(id string, opt ...func(*Handler)) Handler {

	ar := Handler{
		ID:          id,
		Runtime:     "aws-lambda",
		AWSAccount:  "123456789012",
		Healthy:     false,
		Diagnostics: []Diagnostic{},
	}

	for _, o := range opt {
		o(&ar)
	}

	return ar
}
