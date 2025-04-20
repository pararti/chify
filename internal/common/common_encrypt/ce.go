package common_encrypt

import (
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/widget"
)

func GetActionButton() (*widget.Check, *widget.Button) {
	modeToggle := widget.NewCheck(lang.L("Decrypt"), func(checked bool) {})
	modeToggle.SetChecked(false)

	actionButton := widget.NewButton(lang.L("Encrypt"), func() {
	})

	modeToggle.OnChanged = func(checked bool) {
		if checked {
			actionButton.SetText(lang.L("Decrypt"))
		} else {
			actionButton.SetText(lang.L("Encrypt"))
		}
	}

	return modeToggle, actionButton
}
