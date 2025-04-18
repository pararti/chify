package encoding

import (
	"encoding/base32"
	"encoding/base64"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/widget"
	"log"
	"pararti/chify/internal/common"
	"pararti/chify/internal/common/common_encoding"
)

type Base struct {
	Name string
}

type baseMode int

const (
	BASE32 baseMode = iota
	BASE64
)

func (b baseMode) String() string {
	return [...]string{"base32", "base64"}[b]
}

func NewBase() *Base {
	return &Base{Name: "Base(32,64)"}
}

func (b *Base) BuildForm() *fyne.Container {
	header := common.GetHeader(b.Name)
	inputLabel, inputEntry, resetButton := common.GetInput()
	modeToggle, actionButton := common_encoding.GetActionButton()

	// Coding selector
	baseModeLabel := widget.NewLabel(lang.L("Mode"))
	baseModeSelector := widget.NewSelect([]string{"base32", "base64"}, nil)
	baseModeSelector.SetSelected("base32")
	var currentBase = BASE32

	baseModeSelector.OnChanged = func(selected string) {
		switch selected {
		case "base32":
			currentBase = BASE32
		case "base64":
			currentBase = BASE64
		}
	}

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

					var result []byte
					var err error
					switch currentBase {
					case BASE32:
						result, err = base32.StdEncoding.DecodeString(inputEntry.Text)
					case BASE64:
						result, err = base64.StdEncoding.DecodeString(inputEntry.Text)
					}

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

					result := ""
					switch currentBase {
					case BASE32:
						result = base32.StdEncoding.EncodeToString([]byte(inputEntry.Text))
					case BASE64:
						result = base64.StdEncoding.EncodeToString([]byte(inputEntry.Text))
					}

					outputEntry.SetText(result)
				})
			}()
		}
	}

	return container.NewVBox(
		header,
		container.NewHBox(baseModeLabel, baseModeSelector),
		inputLabel,
		container.NewBorder(nil, nil, nil, resetButton, inputEntry),
		container.NewVBox(modeToggle, actionButton),
		outputLabel,
		container.NewBorder(nil, nil, nil, copyButton, outputEntry),
	)
}
