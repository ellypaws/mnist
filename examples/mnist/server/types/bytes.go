package types

import (
	"fmt"
	"image/color"
)

type Byte byte

func (b Byte) Tensor() Tensor {
	return Tensor(b) / 255.0
}

func (b Byte) Byte() byte {
	return byte(b)
}

func (b Byte) Float64() float64 {
	return float64(b) / 255.0
}

func (b Byte) Hex() string {
	return fmt.Sprintf("%02x", b)
}

func BytesToTensor(bytes []Byte) []Tensor {
	if bytes == nil {
		return nil
	}
	var tensors = make([]Tensor, len(bytes))
	for i, b := range bytes {
		tensors[i] = ToTensor(b)
	}
	return tensors
}

func ToByte[number interface{ ~uint8 | ~float64 }](t number) Byte {
	return Byte(t * 255)
}

func (b Byte) RGBA() color.RGBA {
	return color.RGBA{R: uint8(b), G: uint8(b), B: uint8(b), A: 255}
}
