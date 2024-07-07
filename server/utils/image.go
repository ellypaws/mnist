package utils

import (
	"github.com/patrikeh/go-deep/server/types"
	"github.com/patrikeh/go-deep/training"
	"image"
	"image/color"
	"math"
)

func ImageToTensor(img image.Image) []types.Tensor {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var tensors []types.Byte
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := img.At(x, y)
			grayColor := color.GrayModel.Convert(c).(color.Gray)
			tensors = append(tensors, types.Byte(grayColor.Y))
		}
	}

	return types.BytesToTensor(tensors)
}

func ImageToBytes(img image.Image) []types.Byte {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var bytes []types.Byte
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := img.At(x, y)
			grayColor := color.GrayModel.Convert(c).(color.Gray)
			bytes = append(bytes, types.Byte(grayColor.Y))
		}
	}
	return bytes
}

func TensorToImage(tensor []types.Tensor) image.Image {
	img := image.NewGray(image.Rect(0, 0, 28, 28))
	for i, v := range tensor {
		x, y := i%28, i/28
		img.SetGray(x, y, v.Gray())
	}
	return img
}

func RotateImage(img image.Image, degrees float64) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	newImg := image.NewRGBA(bounds)
	angle := degrees * (math.Pi / 180)

	cx, cy := width/2, height/2

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			nx := int(math.Cos(angle)*float64(x-cx) - math.Sin(angle)*float64(y-cy) + float64(cx))
			ny := int(math.Sin(angle)*float64(x-cx) + math.Cos(angle)*float64(y-cy) + float64(cy))
			if nx >= 0 && nx < width && ny >= 0 && ny < height {
				newImg.Set(nx, ny, img.At(x, y))
			}
		}
	}

	return newImg
}

func ZoomImage(img image.Image, factor float64) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	centerX, centerY := width/2, height/2

	newImg := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			srcX := centerX + int(float64(x-centerX)/factor)
			srcY := centerY + int(float64(y-centerY)/factor)

			if srcX >= 0 && srcX < width && srcY >= 0 && srcY < height {
				newImg.Set(x, y, img.At(srcX, srcY))
			} else {
				newImg.Set(x, y, color.RGBA{A: 255})
			}
		}
	}

	return newImg
}

func Translate(img image.Image, deltaX, deltaY int) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	newImg := image.NewRGBA(bounds)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			srcX, srcY := x-deltaX, y-deltaY
			if srcX >= 0 && srcX < width && srcY >= 0 && srcY < height {
				newImg.Set(x, y, img.At(srcX, srcY))
			}
		}
	}

	return newImg
}

type SyntheticConfig struct {
	Rotate    *Values
	Translate *Values
	Zoom      *Values
}

type Values struct {
	Min float64
	Max float64
}

type SyntheticData struct {
	Rotated    []training.Example
	Translated []training.Example
	Zoomed     []training.Example
}

func (c SyntheticConfig) Synthesize(data []training.Example) SyntheticData {
	var syntheticData SyntheticData

	var images = make([]image.Image, len(data))
	for i, example := range data {
		tensors := DataToTensor(example.Input)
		images[i] = TensorToImage(tensors)
	}

	for i, img := range images {
		if c.Rotate != nil {
			rotated := RotateImage(img, RandBetween(c.Rotate.Min, c.Rotate.Max))
			syntheticData.Rotated = append(syntheticData.Rotated,
				training.Example{
					Input:    types.Coerce[types.Byte, float64](ImageToBytes(rotated)),
					Response: data[i].Response,
				},
			)

		}

		if c.Translate != nil {
			translated := Translate(img, int(RandBetween(c.Translate.Min, c.Translate.Max)), int(RandBetween(c.Translate.Min, c.Translate.Max)))
			syntheticData.Translated = append(syntheticData.Translated,
				training.Example{
					Input:    types.Coerce[types.Byte, float64](ImageToBytes(translated)),
					Response: data[i].Response,
				},
			)
		}

		if c.Zoom != nil {
			zoomed := ZoomImage(img, RandBetween(c.Zoom.Min, c.Zoom.Max))
			syntheticData.Zoomed = append(syntheticData.Zoomed,
				training.Example{
					Input:    types.Coerce[types.Byte, float64](ImageToBytes(zoomed)),
					Response: data[i].Response,
				},
			)
		}
	}

	return syntheticData
}

func normalize(value uint8) float64 {
	return float64(value) / 255.0
}
