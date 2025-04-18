package common_encoding

import (
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/widget"
)

func GetActionButton() (*widget.Check, *widget.Button) {
	modeToggle := widget.NewCheck(lang.L("Decode"), func(checked bool) {})
	modeToggle.SetChecked(false)

	actionButton := widget.NewButton(lang.L("Encode"), func() {
	})

	modeToggle.OnChanged = func(checked bool) {
		if checked {
			actionButton.SetText(lang.L("Decode"))
		} else {
			actionButton.SetText(lang.L("Encode"))
		}
	}

	return modeToggle, actionButton
}
