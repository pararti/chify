package encoding

import (
	"encoding/hex"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"log"
	"pararti/chify/internal/common"
	"pararti/chify/internal/common/common_encoding"
)

type Hex struct {
	Name string
}

func NewHex() *Hex {
	return &Hex{Name: "Hex"}
}

func (h *Hex) BuildForm() *fyne.Container {
	header := common.GetHeader(h.Name)
	inputLabel, inputEntry, resetButton := common.GetInput()
	modeToggle, actionButton := common_encoding.GetActionButton()

	outputLabel, outputEntry, copyButton := common.GetOutput()
	outputLabel.SetText(outputLabel.Text)

	actionButton.OnTapped = func() {
		if inputEntry.Text == "" {
			return
		}
		if modeToggle.Checked {
			go func() {
				fyne.Do(func() {
					actionButton.Disable()
					defer actionButton.Enable()

					result := make([]byte, len([]byte(inputEntry.Text)))
					_, err := hex.Decode(result, []byte(inputEntry.Text))

					if err != nil {
						log.Println("Decoding error: ", err)
						outputEntry.SetText("Error: " + err.Error())
						return
					}

					outputEntry.SetText(string(result))
				})
			}()
		} else {
			go func() {
				fyne.Do(func() {
					actionButton.Disable()
					defer actionButton.Enable()

					result := hex.EncodeToString([]byte(inputEntry.Text))

					outputEntry.SetText(result)
				})
			}()
		}
	}

	return container.NewVBox(
		header,
		inputLabel,
		container.NewBorder(nil, nil, nil, resetButton, inputEntry),
		container.NewVBox(modeToggle, actionButton),
		outputLabel,
		container.NewBorder(nil, nil, nil, copyButton, outputEntry),
	)
}
