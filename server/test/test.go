package main

import (
	"fmt"
	"github.com/patrikeh/go-deep/server/mnist"
	"github.com/patrikeh/go-deep/server/utils"
)

const dataPath = "dist/drawing_data.csv"

func main() {
	dataSet, err := mnist.Examples(dataPath)
	if err != nil {
		panic(err)
	}
	dataSet.Shuffle()

	for i, data := range dataSet {
		if i >= 2 {
			break
		}

		tensors := utils.DataToTensor(data.Input)
		fmt.Printf("%s\n", utils.String(tensors))

		image := utils.TensorToImage(tensors)
		rotatedImage := utils.RotateImage(image, utils.RandBetween(-25., 10.))
		rotatedBytes := utils.ImageToBytes(rotatedImage)

		fmt.Printf("rotated:\n%s\n", utils.String(rotatedBytes))

		translatedImage := utils.Translate(image, utils.RandBetween(-10, 10), utils.RandBetween(-10, 10))
		translatedBytes := utils.ImageToBytes(translatedImage)

		fmt.Printf("translated:\n%s\n", utils.String(translatedBytes))

		zoomImage := utils.ZoomImage(image, utils.RandBetween(1.05, 1.15))
		zoomBytes := utils.ImageToBytes(zoomImage)

		fmt.Printf("zoomed:\n%s\n", utils.String(zoomBytes))
	}

	synthesizer := utils.SyntheticConfig{
		Rotate:    &utils.Values{Min: -25., Max: 10.},
		Translate: &utils.Values{Min: -10., Max: 10.},
		Zoom:      &utils.Values{Min: 1.05, Max: 1.15},
	}

	syntheticData := synthesizer.Synthesize(dataSet)

	fmt.Printf("synthetic data:\nrotated: %d, translated: %d, zoomed: %d\n",
		len(syntheticData.Rotated), len(syntheticData.Translated), len(syntheticData.Zoomed))
}
