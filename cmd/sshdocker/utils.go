package main

import (
	"github.com/manifoldco/promptui"
)

func getValue(label string, validate func(string) error, defaultValue ...string) (string, error) {
	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
	}
	if len(defaultValue) > 0 {
		prompt.Default = defaultValue[0]
	}

	return prompt.Run()
}
