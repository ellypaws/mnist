package types

import (
	"fmt"
	"image/color"
)

type Tensor float64

func (t Tensor) Byte() byte {
	return byte(t * 255)
}

func (t Tensor) Float64() float64 {
	return float64(t)
}

func (t Tensor) Hex() string {
	return fmt.Sprintf("%02x", t.Byte())
}

func TensorToBytes(tensors []Tensor) []Byte {
	if tensors == nil {
		return nil
	}
	var bytes = make([]Byte, len(tensors))
	for i, v := range tensors {
		bytes[i] = Byte(v.Byte())
	}
	return bytes
}

func ToTensor[number interface{ ~uint8 | ~float64 }](t number) Tensor {
	return Tensor(t) / 255.0
}

func Coerce[from interface{ ~uint8 | ~float64 }, to interface{ ~uint8 | ~float64 }](in []from) []to {
	if in == nil {
		return nil
	}
	var out = make([]to, len(in))
	for i, v := range in {
		out[i] = to(v)
	}
	return out
}

func (t Tensor) RGBA() color.RGBA {
	b := t.Byte()
	return color.RGBA{R: b, G: b, B: b, A: 255}
}

func (t Tensor) Gray() color.Gray {
	return color.Gray{Y: t.Byte()}
}
