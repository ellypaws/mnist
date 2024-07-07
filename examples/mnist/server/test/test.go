package main

import (
	"fmt"
	"github.com/patrikeh/go-deep/examples/mnist/server/mnist"
	"github.com/patrikeh/go-deep/examples/mnist/server/utils"
)

const dataPath = "dist/drawing_data.csv"

func main() {
	dataSet, err := mnist.Examples(dataPath)
	if err != nil {
		panic(err)
	}
	dataSet.Shuffle()

	for i, data := range dataSet {
		if i >= 10 {
			break
		}

		tensors := utils.DataToTensor(data.Input)
		fmt.Printf("%s\n", utils.String(tensors))
	}
}
