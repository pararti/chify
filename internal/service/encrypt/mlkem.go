package encrypt

import (
	"crypto/mlkem"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"log"
	"pararti/chify/internal/common"
	"pararti/chify/internal/common/common_encrypt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/widget"
)

type MLKEM struct {
	Name string
}

type keySizeMode int

const (
	mode768 keySizeMode = iota
	mode1024
)

func (m keySizeMode) String() string {
	return [...]string{"ML-KEM-768", "ML-KEM-1024"}[m]
}

func NewMLKEM() *MLKEM {
	return &MLKEM{Name: "ML-KEM(Kyber)"}
}

func (m *MLKEM) BuildForm() *fyne.Container {
	header := common.GetHeader(m.Name)
	inputLabel, inputEntry, resetButton := common.GetInput()
	modeToggle, actionButton := common_encrypt.GetActionButton()

	actionButton.Text = lang.L("Encapsulate")
	modeToggle.Text = lang.L("Decapsulate")
	modeToggle.OnChanged = func(checked bool) {
		if checked {
			actionButton.SetText(lang.L("Decapsulate"))
		} else {
			actionButton.SetText(lang.L("Encapsulate"))
		}
	}

	// KeySize selector
	keySizeLabel := widget.NewLabel(lang.L("KeySize"))
	keySizeSelect := widget.NewSelect([]string{"ML-KEM-768", "ML-KEM-1024"}, nil)
	keySizeSelect.SetSelected("ML-KEM-768")
	var currentKeySize keySizeMode = mode768

	keySizeDescription := widget.NewLabel("ML-KEM-768 - Recommended security level (NIST Level 3)")
	keySizeDescription.TextStyle.Italic = true

	keySizeSelect.OnChanged = func(selected string) {
		switch selected {
		case "ML-KEM-768":
			currentKeySize = mode768
			keySizeDescription.SetText("ML-KEM-768 - Recommended security level (NIST Level 3)")
		case "ML-KEM-1024":
			currentKeySize = mode1024
			keySizeDescription.SetText("ML-KEM-1024 - Higher security level (NIST Level 5)")
		}
	}

	publicKeyLabel := widget.NewLabel(lang.L("PublicKey"))
	publicKeyEntry := widget.NewMultiLineEntry()
	publicKeyEntry.Wrapping = fyne.TextWrapBreak
	publicKeyEntry.SetMinRowsVisible(3)
	publicKeyCopyButton := widget.NewButton(lang.L("Copy"), func() {
		if publicKeyEntry.Text != "" {
			fyne.CurrentApp().Clipboard().SetContent(publicKeyEntry.Text)
		}
	})

	privateKeyLabel := widget.NewLabel(lang.L("PrivateKey"))
	privateKeyEntry := widget.NewMultiLineEntry()
	privateKeyEntry.Wrapping = fyne.TextWrapBreak
	privateKeyEntry.SetMinRowsVisible(3)
	privateKeyCopyButton := widget.NewButton(lang.L("Copy"), func() {
		if privateKeyEntry.Text != "" {
			fyne.CurrentApp().Clipboard().SetContent(privateKeyEntry.Text)
		}
	})

	generateKeyButton := widget.NewButton(lang.L("GenerateKeys"), func() {})

	// Output display
	sharedKeyLabel := widget.NewLabel(lang.L("SharedKey") + "(hex)")
	sharedKeyEntry := widget.NewEntry()
	sharedKeyCopyButton := widget.NewButton(lang.L("Copy"), func() {
		if publicKeyEntry.Text != "" {
			fyne.CurrentApp().Clipboard().SetContent(sharedKeyEntry.Text)
		}
	})
	outputLabel, outputEntry, copyButton := common.GetOutput()
	outputLabel.SetText(outputLabel.Text + "(base64)")

	// Main action states
	var decapsulationKey *mlkem.DecapsulationKey768
	var decapsulationKey1024 *mlkem.DecapsulationKey1024
	var encapsulationKey *mlkem.EncapsulationKey768
	var encapsulationKey1024 *mlkem.EncapsulationKey1024

	// Generate keys button action
	generateKeyButton.OnTapped = func() {
		var err error
		var publicKeyText, privateKeyText string

		switch currentKeySize {
		case mode768:
			decapsulationKey, err = mlkem.GenerateKey768()
			if err != nil {
				log.Println("Error generating 768 key:", err)
				outputEntry.SetText("Error: " + err.Error())
				return
			}
			encapsulationKey = decapsulationKey.EncapsulationKey()
			publicKeyText = base64.StdEncoding.EncodeToString(encapsulationKey.Bytes())
			privateKeyText = base64.StdEncoding.EncodeToString(decapsulationKey.Bytes())
		case mode1024:
			decapsulationKey1024, err = mlkem.GenerateKey1024()
			if err != nil {
				log.Println("Error generating 1024 key:", err)
				outputEntry.SetText("Error: " + err.Error())
				return
			}
			encapsulationKey1024 = decapsulationKey1024.EncapsulationKey()
			publicKeyText = base64.StdEncoding.EncodeToString(encapsulationKey1024.Bytes())
			privateKeyText = base64.StdEncoding.EncodeToString(decapsulationKey1024.Bytes())
		}

		// Display keys
		publicKeyEntry.SetText(publicKeyText)
		privateKeyEntry.SetText(privateKeyText)
	}

	// Set up validators and actions based on encryption/decryption mode
	actionButton.OnTapped = func() {
		if modeToggle.Checked {
			if decapsulationKey == nil && decapsulationKey1024 == nil {
				outputEntry.SetText("Error: You need to generate or insert a private key first")
				return
			}

			inputEntry.Validator = func(s string) error {
				if len(s) == 0 {
					return errors.New(lang.L("Required"))
				}
				_, err := base64.StdEncoding.DecodeString(s)
				if err != nil {
					log.Println("Invalid base64 input for ciphertext:", err)
					return errors.New(lang.L("InvalidBase64"))
				}
				return nil
			}

			if err := inputEntry.Validate(); err != nil {
				inputEntry.SetValidationError(err)
				return
			}

			go func() {
				fyne.Do(func() {
					actionButton.Disable()
					defer actionButton.Enable()

					ciphertext, err := base64.StdEncoding.DecodeString(inputEntry.Text)
					if err != nil {
						log.Println("Base64 decode error:", err)
						outputEntry.SetText("Error: Invalid Base64 input")
						return
					}

					var sharedKey []byte
					switch currentKeySize {
					case mode768:
						if decapsulationKey == nil {
							outputEntry.SetText("Error: No ML-KEM-768 private key available")
							return
						}
						sharedKey, err = decapsulationKey.Decapsulate(ciphertext)
					case mode1024:
						if decapsulationKey1024 == nil {
							outputEntry.SetText("Error: No ML-KEM-1024 private key available")
							return
						}
						sharedKey, err = decapsulationKey1024.Decapsulate(ciphertext)
					}

					if err != nil {
						log.Println("Decapsulation error:", err)
						outputEntry.SetText("Error: " + err.Error())
						return
					}

					sharedKeyEntry.SetText(hex.EncodeToString(sharedKey))
				})
			}()
		} else {
			// Encapsulation mode - uses public key to generate shared key and ciphertext
			if encapsulationKey == nil && encapsulationKey1024 == nil {
				outputEntry.SetText("Error: You need to generate or insert a public key first")
				return
			}

			go func() {
				fyne.Do(func() {
					actionButton.Disable()
					defer actionButton.Enable()

					var sharedKey, ciphertext []byte
					switch currentKeySize {
					case mode768:
						if encapsulationKey == nil {
							outputEntry.SetText("Error: No ML-KEM-768 public key available")
							return
						}
						sharedKey, ciphertext = encapsulationKey.Encapsulate()
					case mode1024:
						if encapsulationKey1024 == nil {
							outputEntry.SetText("Error: No ML-KEM-1024 public key available")
							return
						}
						sharedKey, ciphertext = encapsulationKey1024.Encapsulate()
					}

					sharedKeyEntry.SetText(hex.EncodeToString(sharedKey))
					outputEntry.SetText(base64.StdEncoding.EncodeToString(ciphertext))
				})
			}()
		}
	}

	return container.NewVBox(
		header,
		container.NewHBox(keySizeLabel, keySizeSelect),
		keySizeDescription,
		generateKeyButton,
		publicKeyLabel,
		container.NewBorder(nil, nil, nil, publicKeyCopyButton, publicKeyEntry),
		privateKeyLabel,
		container.NewBorder(nil, nil, nil, privateKeyCopyButton, privateKeyEntry),
		inputLabel,
		container.NewBorder(nil, nil, nil, resetButton, inputEntry),
		container.NewVBox(modeToggle, actionButton),
		sharedKeyLabel,
		container.NewBorder(nil, nil, nil, sharedKeyCopyButton, sharedKeyEntry),
		outputLabel,
		container.NewBorder(nil, nil, nil, copyButton, outputEntry),
	)
}
