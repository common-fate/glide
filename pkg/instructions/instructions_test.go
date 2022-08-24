package instructions

import (
	"fmt"
	"testing"
)

func TestInstructions(t *testing.T) {
	i := Instructions{
		Introduction: []Block{TextBlock("hello world")},
		Steps: []Step{
			{
				Title: "step 1",
				Blocks: []Block{
					TextBlock("someting like this"),
					CodeBlock{
						Language: "go",
						Code:     `fmt.Println("hello world")`,
					},
				},
			},
		},
	}
	s := i.RenderTerminal()
	_ = s
	fmt.Println(s)
}
