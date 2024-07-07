package types

import "image/color"

type Hexable interface {
	Hex() string
}

type Gray interface {
	Gray() color.Gray
}
