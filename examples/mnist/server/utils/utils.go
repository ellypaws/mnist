package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/patrikeh/go-deep/examples/mnist/server/types"
	"image"
	"image/png"
	"math/rand/v2"
	"os"
	"strings"
)

func WrapError(message string, err error) map[string]string {
	if err == nil {
		return nil
	}
	return map[string]string{"error": message, "debug": err.Error()}
}

func Base64ToImage(base64Str string) (image.Image, error) {
	base64Str = strings.TrimPrefix(base64Str, "data:image/png;base64,")
	data, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, err
	}

	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	return img, nil
}

func SaveImage(img image.Image, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

func String[hexable types.Hexable](in []hexable) string {
	var numberPrint strings.Builder
	var column int
	for i, in := range in {
		if i%28 == 0 {
			column++
			numberPrint.WriteString(fmt.Sprintf("\n%2d: ", column))
		}
		numberPrint.WriteString(lipgloss.NewStyle().Foreground(toColor(in)).Render("██"))
	}
	return numberPrint.String()
}

func RandBetween[T float64 | int](min, max T) T {
	return T(float64(min) + (float64(max)-float64(min))*rand.Float64())
}

func toColor(in types.Hexable) lipgloss.Color {
	hex := in.Hex()
	return lipgloss.Color(fmt.Sprintf("#%s%s%s", hex, hex, hex))
}

func DataToTensor[number interface{ ~uint8 | ~float64 }](data []number) []types.Tensor {
	if data == nil {
		return nil
	}
	var tensors = make([]types.Tensor, len(data))
	for i, d := range data {
		tensors[i] = types.ToTensor(d)
	}
	return tensors
}
