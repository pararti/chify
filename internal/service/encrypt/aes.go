package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"pararti/chify/internal/common"
	"pararti/chify/internal/common/common_encrypt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/widget"
)

type AES struct {
	Name string
}

type encryptionMode int

const (
	modeCBC encryptionMode = iota
	modeGCM
	modeCTR
)

func (m encryptionMode) String() string {
	return [...]string{"CBC", "GCM", "CTR"}[m]
}

func NewAES() *AES {
	return &AES{Name: "AES"}
}

func (a *AES) BuildForm() *fyne.Container {
	header := common.GetHeader(a.Name)
	inputLabel, inputEntry, resetButton := common.GetInput()
	modeToggle, actionButton := common_encrypt.GetActionButton()

	modeLabel := widget.NewLabel(lang.L("Mode"))
	modeSelect := widget.NewSelect([]string{"CBC", "GCM", "CTR"}, nil)
	modeSelect.SetSelected("CBC")
	var currentMode encryptionMode = modeCBC

	modeDescription := widget.NewLabel("CBC - Cipher Block Chaining")
	modeDescription.TextStyle.Italic = true

	modeSelect.OnChanged = func(selected string) {
		switch selected {
		case "CBC":
			currentMode = modeCBC
			modeDescription.SetText("CBC - Cipher Block Chaining")
		case "GCM":
			currentMode = modeGCM
			modeDescription.SetText("GCM - Galois/Counter Mode (Authenticated)")
		case "CTR":
			currentMode = modeCTR
			modeDescription.SetText("CTR - Counter Mode")
		}
	}

	keyLabel := widget.NewLabel(lang.L("Key"))
	keyEntry := widget.NewEntry()
	keyEntry.PlaceHolder = lang.L("KeyAesError")
	keyEntry.OnChanged = func(key string) {
		kLen := len([]byte(key))
		switch kLen {
		case 16:
			keyLabel.SetText(lang.L("Key") + " aes128")
		case 24:
			keyLabel.SetText(lang.L("Key") + " aes192")
		case 32:
			keyLabel.SetText(lang.L("Key") + " aes256")
		default:
			keyLabel.SetText(lang.L("Key") + " " + lang.L("IncorrectKeyCount") + " " + strconv.Itoa(kLen))
		}
	}
	keyEntry.Validator = keyValidator

	outputLabel, outputEntry, copyButton := common.GetOutput()
	outputLabel.SetText(outputLabel.Text + " (base64)")

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

	actionButton.OnTapped = func() {
		err := keyEntry.Validate()
		if err != nil {
			keyEntry.SetValidationError(err)
			return
		}
		c, err := createCipher([]byte(keyEntry.Text))
		if err != nil {
			log.Println(err)
			keyEntry.SetValidationError(err)
			return
		}

		if inputEntry.Text == "" {
			return
		}

		go func() {
			fyne.Do(func() {
				actionButton.Disable()
				defer actionButton.Enable()
				output := ""

				if modeToggle.Checked {
					decodedInput, err := base64.StdEncoding.DecodeString(inputEntry.Text)
					if err != nil {
						log.Println("Base64 decode error:", err)
						outputEntry.SetText("Error: Invalid Base64 input")
						return
					}

					var decryptedData []byte
					switch currentMode {
					case modeCBC:
						decryptedData = decryptCBC(c, decodedInput)
					case modeGCM:
						decryptedData, err = decryptGCM(c, decodedInput)
					case modeCTR:
						decryptedData, err = decryptCTR(c, decodedInput)
					}

					if err != nil {
						log.Println("Decryption error:", err)
						outputEntry.SetText("Error: " + err.Error())
						return
					}

					output = string(decryptedData)
				} else {
					var encryptedData []byte
					var err error

					switch currentMode {
					case modeCBC:
						encryptedData = encryptCBC(c, []byte(inputEntry.Text))
					case modeGCM:
						encryptedData, err = encryptGCM(c, []byte(inputEntry.Text))
					case modeCTR:
						encryptedData, err = encryptCTR(c, []byte(inputEntry.Text))
					}

					if err != nil {
						log.Println("Encryption error:", err)
						outputEntry.SetText("Error: " + err.Error())
						return
					}

					output = base64.StdEncoding.EncodeToString(encryptedData)
				}

				outputEntry.SetText(output)
			})
		}()
	}

	return container.NewVBox(
		header,
		container.NewHBox(modeLabel, modeSelect),
		modeDescription,
		inputLabel,
		container.NewBorder(nil, nil, nil, resetButton, inputEntry),
		keyLabel,
		container.NewBorder(nil, nil, nil, generateKeyButton, keyEntry),
		container.NewVBox(modeToggle, actionButton),
		outputLabel,
		container.NewBorder(nil, nil, nil, copyButton, outputEntry),
	)
}

func encryptCBC(c cipher.Block, data []byte) []byte {
	data = pKCS7Padding(data, c.BlockSize())
	iv := make([]byte, c.BlockSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Println("Error generating IV:", err)
		return []byte{}
	}

	mode := cipher.NewCBCEncrypter(c, iv)
	result := make([]byte, len(data))
	mode.CryptBlocks(result, data)

	result = append(iv, result...)

	return result
}

func decryptCBC(c cipher.Block, data []byte) []byte {
	blockSize := c.BlockSize()

	if len(data) < blockSize {
		log.Println("Data too short for decryption")
		return []byte("Error: data too short for decryption")
	}

	iv := data[:blockSize]
	data = data[blockSize:]

	mode := cipher.NewCBCDecrypter(c, iv)
	result := make([]byte, len(data))
	mode.CryptBlocks(result, data)

	result = pKCS7UnPadding(result)

	return result
}

func encryptGCM(c cipher.Block, data []byte) ([]byte, error) {
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	result := gcm.Seal(nonce, nonce, data, nil)
	return result, nil
}

func decryptGCM(c cipher.Block, data []byte) ([]byte, error) {
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	if len(data) < gcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func encryptCTR(c cipher.Block, data []byte) ([]byte, error) {
	iv := make([]byte, c.BlockSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	ctr := cipher.NewCTR(c, iv)

	ciphertext := make([]byte, len(data))
	ctr.XORKeyStream(ciphertext, data)

	result := append(iv, ciphertext...)

	return result, nil
}

func decryptCTR(c cipher.Block, data []byte) ([]byte, error) {
	blockSize := c.BlockSize()

	if len(data) < blockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := data[:blockSize]
	ciphertext := data[blockSize:]

	ctr := cipher.NewCTR(c, iv)

	plaintext := make([]byte, len(ciphertext))
	ctr.XORKeyStream(plaintext, ciphertext)

	return plaintext, nil
}

func createCipher(key []byte) (cipher.Block, error) {
	return aes.NewCipher(key)
}

func keyValidator(key string) error {
	keyLength := len([]byte(key))
	switch keyLength {
	case 16, 24, 32:
		return nil
	default:
		return errors.New(lang.L("KeyAesError"))
	}
}

func pKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func pKCS7UnPadding(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}
