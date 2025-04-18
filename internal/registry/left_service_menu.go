package registry

import (
	"pararti/chify/internal/service"
	encoding2 "pararti/chify/internal/service/encode"
	encrypt2 "pararti/chify/internal/service/encrypt"
	hash2 "pararti/chify/internal/service/hash"
)

type SubMenuElement struct {
	Name    string
	Service service.FormBuilder
}

type MenuElement struct {
	Category string
	Elements []*SubMenuElement
}

var DefaultService = encrypt2.NewAES()

var LeftServiceMenu = []*MenuElement{
	{
		Category: "crypto",
		Elements: []*SubMenuElement{
			{
				Name:    "aes",
				Service: encrypt2.NewAES(),
			},
			{
				Name:    "chacha20",
				Service: encrypt2.NewChaCha20(),
			},
			{
				Name:    "ml-kem",
				Service: encrypt2.NewMLKEM(),
			},
		},
	},
	{
		Category: "encode",
		Elements: []*SubMenuElement{
			{
				Name:    "ascii85",
				Service: encoding2.NewAscii85(),
			},
			{
				Name:    "base",
				Service: encoding2.NewBase(),
			},
			{
				Name:    "hex",
				Service: encoding2.NewHex(),
			},
		},
	},
	{
		Category: "hash",
		Elements: []*SubMenuElement{
			{
				Name:    "md5",
				Service: hash2.NewMd5(),
			},
			{
				Name:    "sha",
				Service: hash2.NewSha(),
			},
		},
	},
}
