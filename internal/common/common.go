package common

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/widget"
)

func GetInput() (*widget.Label, *widget.Entry, *widget.Button) {
	inputLabel := widget.NewLabel(lang.L("Input"))
	inputEntry := widget.NewMultiLineEntry()
	inputEntry.SetMinRowsVisible(6)
	inputEntry.Wrapping = fyne.TextWrapBreak

	resetButton := widget.NewButton(lang.L("Reset"), func() {
		inputEntry.SetText("")
	})

	return inputLabel, inputEntry, resetButton
}

func GetOutput() (*widget.Label, *widget.Entry, *widget.Button) {
	outputLabel := widget.NewLabel(lang.L("Output"))
	outputEntry := widget.NewMultiLineEntry()
	outputEntry.SetMinRowsVisible(6)
	outputEntry.Wrapping = fyne.TextWrapBreak
	copyButton := widget.NewButton(lang.L("Copy"), func() {
		if outputEntry.Text != "" {
			fyne.CurrentApp().Clipboard().SetContent(outputEntry.Text)
		}
	})

	return outputLabel, outputEntry, copyButton
}

func GetHeader(text string) *widget.Label {
	header := widget.NewLabel(text)
	header.TextStyle.Bold = true
	header.Alignment = fyne.TextAlignCenter

	return header
}
