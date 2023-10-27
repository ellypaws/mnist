package main

import (
	"fmt"
	"log"

	"gorgonia.org/tensor"
	"mnist-go/mnist"
)

func main() {
	for _, typ := range []string{"test", "train"} {
		inputs, targets, err := mnist.Load(typ, "./dataset", tensor.Float64)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(typ+" inputs:", inputs.Shape())
		fmt.Println(typ+" data:", targets.Shape())
	}
}
