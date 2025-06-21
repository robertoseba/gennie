package ui

import "github.com/charmbracelet/huh"

type promptModel struct {
	form     *huh.Form
	confirm  *huh.Confirm
	isActive bool
	answer   bool
}

func newPromptModel() *promptModel {
	prompt := huh.NewConfirm()
	prompt.Negative("Cancel").Affirmative("Continue")
	prompt.WithTheme(huh.ThemeDracula())
	model := &promptModel{
		confirm:  prompt,
		isActive: false,
	}
	model.confirm.Value(&model.answer)
	model.form = huh.NewForm(
		huh.NewGroup(
			prompt,
		),
	)

	return model
}
