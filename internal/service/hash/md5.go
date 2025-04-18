package hash

import (
	"crypto/md5"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"pararti/chify/internal/common"
	"pararti/chify/internal/common/common_hash"
)

type Md5 struct {
	Name string
}

func NewMd5() *Md5 {
	return &Md5{Name: "MD5"}
}

func (m *Md5) BuildForm() *fyne.Container {
	header := common.GetHeader(m.Name)
	inputLabel, inputEntry, resetButton := common.GetInput()
	actionButton := common_hash.GetActionButton()

	outputLabel, outputEntry, copyButton := common.GetOutput()
	outputLabel.SetText(outputLabel.Text)

	actionButton.OnTapped = func() {
		if inputEntry.Text == "" {
			return
		}
		go func() {
			fyne.Do(func() {
				actionButton.Disable()
				defer actionButton.Enable()

				h := md5.Sum([]byte(inputEntry.Text))

				outputEntry.SetText(fmt.Sprintf("%x", h))
			})
		}()
	}

	return container.NewVBox(
		header,
		inputLabel,
		container.NewBorder(nil, nil, nil, resetButton, inputEntry),
		actionButton,
		outputLabel,
		container.NewBorder(nil, nil, nil, copyButton, outputEntry),
	)
}
