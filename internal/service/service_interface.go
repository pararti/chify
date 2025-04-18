package service

import "fyne.io/fyne/v2"

type FormBuilder interface {
	BuildForm() *fyne.Container
}
