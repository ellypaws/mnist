package utils

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	"image/color"
	"image/png"
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

func ImageToTensor(img image.Image) ([]float64, error) {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	if width != 28 || height != 28 {
		return nil, errors.New("image must be 28x28 pixels")
	}

	var tensor []float64
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := img.At(x, y)
			grayColor := color.GrayModel.Convert(c).(color.Gray)
			tensor = append(tensor, Normalize(grayColor.Y))
		}
	}

	return tensor, nil
}

func Normalize(color uint8) float64 {
	return float64(color) / 255.0
}
