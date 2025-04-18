package encrypt

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"log"
	"pararti/chify/internal/common"
	"pararti/chify/internal/common/common_encrypt"
	"strconv"

	"golang.org/x/crypto/chacha20"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/widget"
)

type ChaCha20 struct {
	Name string
}

func NewChaCha20() *ChaCha20 {
	return &ChaCha20{Name: "ChaCha20"}
}

func (c *ChaCha20) BuildForm() *fyne.Container {
	header := common.GetHeader(c.Name)
	inputLabel, inputEntry, resetButton := common.GetInput()
	modeToggle, actionButton := common_encrypt.GetActionButton()

	keyLabel := widget.NewLabel(lang.L("Key"))
	keyEntry := widget.NewEntry()
	keyEntry.PlaceHolder = lang.L("KeyMustBe32Bytes")

	nonceLabel := widget.NewLabel(lang.L("Nonce"))
	nonceEntry := widget.NewEntry()
	nonceEntry.PlaceHolder = lang.L("NonceMustBe12Bytes")

	counterLabel := widget.NewLabel(lang.L("Counter"))
	counterEntry := widget.NewEntry()
	counterEntry.Text = "1"

	outputLabel, outputEntry, copyButton := common.GetOutput()

	// Generate random key and nonce buttons
	generateKeyButton := widget.NewButton(lang.L("Generate"), func() {
		key := make([]byte, 32)
		_, err := rand.Read(key)
		if err != nil {
			log.Println("Error generating random key:", err)
			outputEntry.SetText("Error generating random key: " + err.Error())
			return
		}

		keyEntry.SetText(hex.EncodeToString(key)[:32])
	})

	generateNonceButton := widget.NewButton(lang.L("Generate"), func() {
		nonce := make([]byte, 12)
		_, err := rand.Read(nonce)
		if err != nil {
			log.Println("Error generating random nonce:", err)
			outputEntry.SetText("Error generating random nonce: " + err.Error())
			return
		}

		nonceEntry.SetText(hex.EncodeToString(nonce)[:12])
	})

	keyEntry.Validator = func(s string) error {
		if s == "" {
			return errors.New(lang.L("Required"))
		}

		if len([]byte(keyEntry.Text)) != 32 {
			return errors.New(lang.L("KeyMustBe32Bytes"))
		}

		return nil
	}

	nonceEntry.Validator = func(s string) error {
		if s == "" {
			return errors.New(lang.L("Required"))
		}

		if len([]byte(nonceEntry.Text)) != 12 {
			return errors.New(lang.L("NonceMustBe12Bytes"))
		}

		return nil
	}

	counterEntry.Validator = func(s string) error {
		if s == "" {
			return errors.New("Error: " + lang.L("Counter") + " " + lang.L("Required"))
		}

		return nil
	}

	actionButton.OnTapped = func() {
		// Validate inputs
		err := keyEntry.Validate()
		if err != nil {
			keyEntry.SetValidationError(err)
			return
		}

		err = nonceEntry.Validate()
		if err != nil {
			nonceEntry.SetValidationError(err)
			return
		}

		err = counterEntry.Validate()
		if err != nil {
			counterEntry.SetValidationError(err)
			return
		}

		if inputEntry.Text == "" {
			return
		}

		// Process the encryption/decryption
		go func() {
			fyne.Do(func() {
				actionButton.Disable()
				defer actionButton.Enable()

				keyBytes := []byte(keyEntry.Text)
				nonceBytes := []byte(nonceEntry.Text)

				var counter uint32 = 1
				if counterVal, err := strconv.ParseUint(counterEntry.Text, 10, 32); err == nil {
					counter = uint32(counterVal)
				} else {
					log.Println("Invalid counter value, using 1:", err)
				}

				if modeToggle.Checked {
					ciphertext, err := base64.StdEncoding.DecodeString(inputEntry.Text)
					if err != nil {
						outputEntry.SetText("Error: Invalid Base64 input")
						return
					}

					cipher, err := chacha20.NewUnauthenticatedCipher(keyBytes, nonceBytes)
					if err != nil {
						outputEntry.SetText("Error: " + err.Error())
						return
					}

					cipher.SetCounter(counter)

					plaintext := make([]byte, len(ciphertext))
					cipher.XORKeyStream(plaintext, ciphertext)

					outputEntry.SetText(string(plaintext))
				} else {
					plaintext := []byte(inputEntry.Text)

					cipher, err := chacha20.NewUnauthenticatedCipher(keyBytes, nonceBytes)
					if err != nil {
						outputEntry.SetText("Error: " + err.Error())
						return
					}

					cipher.SetCounter(counter)

					ciphertext := make([]byte, len(plaintext))
					cipher.XORKeyStream(ciphertext, plaintext)

					outputEntry.SetText(base64.StdEncoding.EncodeToString(ciphertext))
				}
			})
		}()
	}

	return container.NewVBox(
		header,
		inputLabel,
		container.NewBorder(nil, nil, nil, resetButton, inputEntry),
		keyLabel,
		container.NewBorder(nil, nil, nil, generateKeyButton, keyEntry),
		nonceLabel,
		container.NewBorder(nil, nil, nil, generateNonceButton, nonceEntry),
		counterLabel,
		counterEntry,
		container.NewVBox(modeToggle, actionButton),
		outputLabel,
		container.NewBorder(nil, nil, nil, copyButton, outputEntry),
	)
}
