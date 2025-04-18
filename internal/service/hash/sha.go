package hash

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha3"
	"crypto/sha512"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/widget"
	"pararti/chify/internal/common"
	"pararti/chify/internal/common/common_hash"
)

type Sha struct {
	Name string
}

type hashMode int

const (
	SHA1 hashMode = iota
	SHA224
	SHA256
	SHA3_224
	SHA3_256
	SHA3_384
	SHA3_512
	SHA512_224
	SHA512_256
	SHA384
	SHA512
)

var shas = []string{"sha1", "sha224", "sha256", "sha3-224", "sha3-256", "sha3-384", "sha3-512", "sha512-224", "sha512-256", "sha384", "sha512"}

func (h hashMode) String() string {
	return shas[h]
}

func NewSha() *Sha {
	return &Sha{Name: "SHA"}
}

func (s *Sha) BuildForm() *fyne.Container {
	header := common.GetHeader(s.Name)
	inputLabel, inputEntry, resetButton := common.GetInput()
	actionButton := common_hash.GetActionButton()

	// Sha hash selector
	baseModeLabel := widget.NewLabel(lang.L("Mode"))
	baseModeSelector := widget.NewSelect(shas, nil)
	baseModeSelector.SetSelected("sha1")
	var currentSha = SHA1

	baseModeSelector.OnChanged = func(selected string) {
		switch selected {
		case "sha1":
			currentSha = SHA1
		case "sha224":
			currentSha = SHA224
		case "sha256":
			currentSha = SHA256
		case "sha3-224":
			currentSha = SHA3_224
		case "sha3-256":
			currentSha = SHA3_256
		case "sha3-384":
			currentSha = SHA3_384
		case "sha3-512":
			currentSha = SHA3_512
		case "sha512-224":
			currentSha = SHA512_224
		case "sha512-256":
			currentSha = SHA512_256
		case "sha384":
			currentSha = SHA384
		case "sha512":
			currentSha = SHA512

		}
	}

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

				var h any
				switch currentSha {
				case SHA1:
					h = sha1.Sum([]byte(inputEntry.Text))
				case SHA224:
					h = sha256.Sum224([]byte(inputEntry.Text))
				case SHA256:
					h = sha256.Sum256([]byte(inputEntry.Text))
				case SHA3_224:
					h = sha3.Sum224([]byte(inputEntry.Text))
				case SHA3_256:
					h = sha3.Sum256([]byte(inputEntry.Text))
				case SHA3_384:
					h = sha3.Sum384([]byte(inputEntry.Text))
				case SHA3_512:
					h = sha3.Sum512([]byte(inputEntry.Text))
				case SHA512_224:
					h = sha512.Sum512_224([]byte(inputEntry.Text))
				case SHA512_256:
					h = sha512.Sum512_256([]byte(inputEntry.Text))
				case SHA384:
					h = sha512.Sum384([]byte(inputEntry.Text))
				case SHA512:
					h = sha512.Sum512([]byte(inputEntry.Text))
				}

				outputEntry.SetText(fmt.Sprintf("%x", h))
			})
		}()
	}

	return container.NewVBox(
		header,
		container.NewHBox(baseModeLabel, baseModeSelector),
		inputLabel,
		container.NewBorder(nil, nil, nil, resetButton, inputEntry),
		actionButton,
		outputLabel,
		container.NewBorder(nil, nil, nil, copyButton, outputEntry),
	)
}
