package encoding

import (
	"encoding/ascii85"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"log"
	"math"
	"pararti/chify/internal/common"
	"pararti/chify/internal/common/common_encoding"
)

type Ascii85 struct {
	Name string
}

func NewAscii85() *Ascii85 {
	return &Ascii85{Name: "Ascii85"}
}

func (a *Ascii85) BuildForm() *fyne.Container {
	header := common.GetHeader(a.Name)
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

					dstSize := int(math.Ceil(float64(len([]byte(inputEntry.Text))) * 0.8))
					println(dstSize)

					dst := make([]byte, dstSize)
					_, _, err := ascii85.Decode(dst, []byte(inputEntry.Text), true)

					if err != nil {
						log.Println("Decoding error: ", err)
						outputEntry.SetText("Error: " + err.Error())
						return
					}

					outputEntry.SetText(string(dst))
				})
			}()
		} else {
			go func() {
				fyne.Do(func() {
					actionButton.Disable()
					defer actionButton.Enable()

					dstSize := int(math.Ceil(float64(len([]byte(inputEntry.Text))) * 1.25))
					dstSize += (5 - dstSize%5) % 5 //the size must be a multiple of 5
					dst := make([]byte, dstSize)
					ascii85.Encode(dst, []byte(inputEntry.Text))

					outputEntry.SetText(string(dst))
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
