package main

import (
	"fmt"
	"gorgonia.org/tensor/native"
	"image"
	"image/png"
	"log"
	"os"

	"gorgonia.org/tensor"
	"mnist-go/mnist"
)

func main() {
	inputs, targets, err := mnist.Load("train", "./dataset", tensor.Float64)
	if err != nil {
		log.Fatal(err)
	}
	cols := inputs.Shape()[1]
	imageBackend := make([]uint8, cols)
	for i := 0; i < cols; i++ {
		v, _ := inputs.At(0, i)
		imageBackend[i] = uint8((v.(float64) - 0.1) * 0.9 * 255)
	}
	img := &image.Gray{
		Pix:    imageBackend,
		Stride: 28,
		Rect:   image.Rect(0, 0, 28, 28),
	}
	w, _ := os.Create("output.png")
	vals, _ := native.MatrixF64(targets.(*tensor.Dense))
	fmt.Println(vals[0])
	err = png.Encode(w, img)
}
