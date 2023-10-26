package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
)

type InputData struct {
	InputString string `json:"input_string"`
}

type OutputData struct {
	OutputField string `json:"output_field"`
}

func HandleRequest(in InputData) (OutputData, error) {

	fmt.Println("Hello from handler")
	var output OutputData

	if in.InputString == "YES" {
		output = OutputData{OutputField: "YES"}
	} else {
		output = OutputData{OutputField: "NO"}
	}

	fmt.Println("Bye from Handler")

	return output, nil
}

func main() {
	lambda.Start(HandleRequest)
}
