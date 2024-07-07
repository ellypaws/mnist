package utils

import (
	"errors"
	"github.com/patrikeh/go-deep/examples/mnist/server/types"
	"image"
	"image/color"
)

func ImageToTensor(img image.Image) ([]types.Tensor, error) {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	if width != 28 || height != 28 {
		return nil, errors.New("image must be 28x28 pixels")
	}

	var tensor []types.Tensor
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := img.At(x, y)
			grayColor := color.GrayModel.Convert(c).(color.Gray)
			tensor = append(tensor, types.ToTensor(grayColor.Y))
		}
	}

	return tensor, nil
}
