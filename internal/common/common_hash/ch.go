package common_hash

import (
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/widget"
)

func GetActionButton() *widget.Button {

	return widget.NewButton(lang.L("HashName"), nil) // Hash - is reserved

}
