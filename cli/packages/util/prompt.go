package util

import (
	"errors"

	"github.com/manifoldco/promptui"
)

type Prompt struct {
	errorMsg string
	label    string
}

func GetPromptInput(p Prompt) string {
	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New(p.errorMsg)
		}
		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }}",
		Valid:   "{{ . | green }}",
		Invalid: "{{ . | red }}",
		Success: "{{ . | bold }}",
	}

	prompt := promptui.Prompt{
		Label:     p.label,
		Templates: templates,
		Validate:  validate,
	}

	result, err := prompt.Run()
	HandleError(err, true)

	return result
}
